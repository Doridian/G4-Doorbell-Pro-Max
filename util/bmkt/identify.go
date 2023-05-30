package bmkt

// #include <libbmkt/custom.h>
import "C"

// TODO: Wait for result
func (ctx *BMKTContext) Identify() error {
	err := ctx.Cancel()
	if err != nil {
		return err
	}
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()

	ctx.state = IF_STATE_IDENTIFYING
	err = ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_identify(ctx.session)
	})
	if err != nil {
		return err
	}
	ctx.waitForIdle()
	return nil
}
