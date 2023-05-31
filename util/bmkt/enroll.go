package bmkt

// #include <libbmkt/custom.h>
import "C"
import "errors"

func (ctx *Context) Enroll(username string, finger_id int) error {
	c_username, c_username_len := convertStringToCUserID(username)
	c_finger_id := C.uint8_t(finger_id)

	err := ctx.Cancel()
	if err != nil {
		return err
	}
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()
	if ctx.state != IF_STATE_IDLE {
		return errors.New("interrupted by other operation")
	}

	ctx.state = IF_STATE_ENROLLING
	ctx.lastEnrollResult = -1

	err = ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_enroll(ctx.session, c_username, c_username_len, c_finger_id)
	})
	if err != nil {
		return err
	}
	ctx.waitForIdle()
	return wrapBMKTError(ctx.lastEnrollResult)
}

func (ctx *Context) DeleteEnrollment(username string, finger_id int) error {
	c_username, c_username_len := convertStringToCUserID(username)
	c_finger_id := C.uint8_t(finger_id)

	err := ctx.Cancel()
	if err != nil {
		return err
	}
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()
	if ctx.state != IF_STATE_IDLE {
		return errors.New("interrupted by other operation")
	}

	ctx.state = IF_STATE_DELETING_USER
	ctx.lastDeleteUserResult = -1

	err = ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_enroll(ctx.session, c_username, c_username_len, c_finger_id)
	})
	if err != nil {
		return err
	}
	ctx.waitForIdle()
	return wrapBMKTError(ctx.lastDeleteUserResult)
}
