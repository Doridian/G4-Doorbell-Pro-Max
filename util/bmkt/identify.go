package bmkt

// #include <libbmkt/custom.h>
import "C"
import "errors"

func (ctx *BMKTContext) autoIdentify() {
	if !ctx.AutoIdentify || ctx.state != IF_STATE_IDLE {
		return
	}

	_, err := ctx.identifyNoCancel()
	if err != nil {
		ctx.logger.Warn().Bool("success", false).Str("op", "auto_identify").Err(err).Send()
		return
	}
	ctx.logger.Info().Bool("success", true).Str("op", "auto_identify").Send()
}

func (ctx *BMKTContext) identifyNoCancel() (string, error) {
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()
	if ctx.state != IF_STATE_IDLE {
		return "", errors.New("interrupted by other operation")
	}

	ctx.state = IF_STATE_IDENTIFYING
	ctx.lastIdentifyResult = -1
	ctx.lastIdentifyUser = ""

	err := ctx.wrapAndRunWithRetry(func() C.int {
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

func (ctx *BMKTContext) Identify() (string, error) {
	err := ctx.Cancel()
	if err != nil {
		return "", err
	}

	return ctx.identifyNoCancel()
}
