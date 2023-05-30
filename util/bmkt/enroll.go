package bmkt

// #include <libbmkt/custom.h>
import "C"

// TODO: Wait for result
func (ctx *BMKTContext) Enroll(username string, finger_id int) error {
	err := ctx.Cancel()
	if err != nil {
		return err
	}

	c_username, c_username_len := convertStringToCUserID(username)
	c_finger_id := C.uint8_t(finger_id)

	ctx.state = IF_STATE_ENROLLING
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_enroll(ctx.session, c_username, c_username_len, c_finger_id)
	})
}

// TODO: Wait for result
func (ctx *BMKTContext) DeleteEnrollment(username string, finger_id int) error {
	err := ctx.Cancel()
	if err != nil {
		return err
	}

	c_username, c_username_len := convertStringToCUserID(username)
	c_finger_id := C.uint8_t(finger_id)

	ctx.state = IF_STATE_DELETING_USER
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_enroll(ctx.session, c_username, c_username_len, c_finger_id)
	})
}
