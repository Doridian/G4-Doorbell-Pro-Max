package bmkt

// #include <libbmkt/custom.h>
import "C"

func (ctx *BMKTContext) Identify() error {
	err := ctx.Cancel()
	if err != nil {
		return err
	}
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_identify(ctx.session)
	})
}
