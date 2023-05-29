package bmkt

// #include <libbmkt/custom.h>
// #include "wrapper.h"
import "C"

import (
	"errors"
	"fmt"
	"time"
)

type runnable func() C.int

//export c_on_error
func c_on_error(cid C.uint64_t, code C.uint16_t) {
	ctx := bmktContexts[uint64(cid)]
	if ctx == nil {
		return
	}
	ctx.handleError(int(code))
}

//export c_on_response
func c_on_response(cid C.uint64_t, resp *C.bmkt_response_t) {
	ctx := bmktContexts[uint64(cid)]
	if ctx == nil {
		return
	}
	ctx.handleResponse(resp)
}

func isErrorTransient(err C.int) bool {
	return err == C.BMKT_SENSOR_NOT_READY || err == C.BMKT_TIMEOUT || err == C.BMKT_OP_TIME_OUT
}

func wrapBMKTError(err C.int) error {
	if err == C.BMKT_SUCCESS {
		return nil
	}
	return fmt.Errorf("code %d error", int(err))
}

func (ctx *BMKTContext) wrapAndRunWithRetry(runfunc runnable) error {
	for curRetry := 0; curRetry < ctx.MaxRetries; curRetry++ {
		res := runfunc()
		if res == C.BMKT_SUCCESS {
			return nil
		}
		if !isErrorTransient(res) {
			return fmt.Errorf("code %d error", res)
		}
		time.Sleep(ctx.RetryDelay)
	}
	return errors.New("retries exhausted")
}

func (ctx *BMKTContext) wrappedOpen() error {
	bmktContexts[ctx.id] = ctx
	ctx.session = C.bmkt_wrapped_init()
	if ctx.session == nil {
		return errors.New("bmkt_init() failure")
	}
	err := wrapBMKTError(C.bmkt_wrapped_open(ctx.session, &ctx.sensor, &ctx.cid))
	if err != nil {
		return err
	}
	return ctx.wrapAndRunWithRetry(func() C.int {
		return C.bmkt_init_fps(ctx.session)
	})
}
