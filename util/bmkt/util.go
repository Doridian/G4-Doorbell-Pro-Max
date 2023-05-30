package bmkt

// #include <string.h>
// #include <libbmkt/custom.h>
import "C"
import (
	"time"
	"unsafe"
)

func convertCUserIDToString(c_user_id *C.user_id_t) string {
	c_user_id_charptr := (*C.char)(unsafe.Pointer(c_user_id))
	len := C.strnlen(c_user_id_charptr, C.BMKT_MAX_USER_ID_LEN)
	str := C.GoStringN(c_user_id_charptr, C.int(len))
	return str
}

func convertStringToCUserID(username string) (*C.uint8_t, C.uint32_t) {
	c_username := (*C.uint8_t)(unsafe.Pointer(C.CString(username)))
	c_username_len := C.uint32_t(len(username))
	return c_username, c_username_len
}

func (ctx *BMKTContext) Cancel() error {
	if ctx.state == IF_STATE_IDLE || ctx.state == IF_STATE_INIT {
		return nil
	}

	ctx.state = IF_STATE_CANCELLING
	ctx.lastCancelResult = -1
	err := ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_cancel_op(ctx.session)
	})
	if err != nil {
		return err
	}
	ctx.waitForIdle()
	return wrapBMKTError(ctx.lastCancelResult)
}

func (ctx *BMKTContext) waitForIdle() {
	for ctx.state != IF_STATE_IDLE {
		time.Sleep(time.Microsecond * 100)
	}
}
