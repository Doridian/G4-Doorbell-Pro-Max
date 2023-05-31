package bmkt

// #include <libbmkt/custom.h>
import "C"
import "errors"

func (ctx *Context) autoIdentify() {
	if !ctx.AutoIdentify || ctx.state != IF_STATE_IDLE {
		return
	}

	user, finger, err := ctx.identifyNoCancel()
	if err != nil {
		ctx.logger.Warn().Bool("success", false).Str("op", "auto_identify").Err(err).Send()
		return
	}
	ctx.logger.Info().Bool("success", true).Str("op", "auto_identify").Str("user", user).Int("finger", finger).Send()
}

func (ctx *Context) identifyNoCancel() (string, int, error) {
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()
	if ctx.state != IF_STATE_IDLE {
		return "", 0, errors.New("interrupted by other operation")
	}

	ctx.state = IF_STATE_IDENTIFYING
	ctx.lastIdentifyResult = -1
	ctx.lastIdentifyUser = ""

	err := ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_identify(ctx.session)
	})
	if err != nil {
		return "", 0, err
	}
	ctx.waitForIdle()

	if ctx.lastIdentifyResult == C.BMKT_SUCCESS && ctx.IdentifyCallback != nil {
		go ctx.IdentifyCallback(ctx.lastIdentifyUser, ctx.lastIdentifyFinger)
	}

	return ctx.lastIdentifyUser, ctx.lastIdentifyFinger, wrapBMKTError(ctx.lastIdentifyResult)
}

func (ctx *Context) Identify() (string, int, error) {
	err := ctx.Cancel()
	if err != nil {
		return "", 0, err
	}

	return ctx.identifyNoCancel()
}
