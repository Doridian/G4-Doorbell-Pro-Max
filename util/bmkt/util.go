package bmkt

// #include <libbmkt/custom.h>
// #include <string.h>
import "C"
import "unsafe"

func convertCUserIDToString(c_user_id *C.user_id_t) string {
	c_user_id_charptr := (*C.char)(unsafe.Pointer(c_user_id))
	len := C.strnlen(c_user_id_charptr, C.BMKT_MAX_USER_ID_LEN)
	str := C.GoStringN(c_user_id_charptr, C.int(len))
	return str
}

func (ctx *BMKTContext) Cancel() error {
	ctx.state = IF_STATE_CANCELLING
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_cancel_op(ctx.session)
	})
}
