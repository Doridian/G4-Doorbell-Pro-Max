package bmkt

// #include <libbmkt/custom.h>
import "C"

// TODO: Wait for result
func (ctx *BMKTContext) Identify() error {
	err := ctx.Cancel()
	if err != nil {
		return err
	}
	ctx.state = IF_STATE_IDENTIFYING
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_identify(ctx.session)
	})
}
