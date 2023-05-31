package bmkt

// #include <libbmkt/custom.h>
import "C"

func (ctx *BMKTContext) Identify() (string, error) {
	err := ctx.Cancel()
	if err != nil {
		return "", err
	}
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()

	ctx.state = IF_STATE_IDENTIFYING
	ctx.lastIdentifyResult = -1
	ctx.lastIdentifyUser = ""

	err = ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_identify(ctx.session)
	})
	if err != nil {
		return "", err
	}
	ctx.waitForIdle()

	if ctx.lastIdentifyResult == C.BMKT_SUCCESS && ctx.IdentifyCallback != nil {
		go ctx.IdentifyCallback(ctx.lastIdentifyUser, ctx.lastVerifyFinger)
	}

	return ctx.lastIdentifyUser, wrapBMKTError(ctx.lastIdentifyResult)
}
