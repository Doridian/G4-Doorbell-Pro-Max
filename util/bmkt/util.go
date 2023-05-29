package bmkt

// #include <libbmkt/custom.h>
import "C"

func (ctx *BMKTContext) Cancel() error {
	ctx.state = IF_STATE_CANCELLING
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_cancel_op(ctx.session)
	})
}
