package bmkt

// #include "handler.h"
// #cgo LDFLAGS: -lbmkt
import "C"
import (
	"errors"
	"fmt"
	"time"
)

type BMKTContext struct {
	ctx        *C.cb_ctx_t
	MaxRetries int
	RetryDelay time.Duration
}

type runnable func() C.int

func Open() (*BMKTContext, error) {
	ctxPtr := C.bmkt_main_init()
	if ctxPtr == nil {
		return nil, errors.New("unknown error allocating context")
	}

	if ctxPtr.state == C.IF_STATE_INVALID {
		return nil, fmt.Errorf("code %d initializing BMKT", ctxPtr.last_error)
	}

	return &BMKTContext{
		ctx:        ctxPtr,
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 1,
	}, nil
}

func isErrorTransient(err int) bool {
	return err == C.BMKT_SENSOR_NOT_READY || err == C.BMKT_TIMEOUT
}

func (c *BMKTContext) runWithRetry(runfunc runnable) error {
	for curRetry := 0; curRetry < c.MaxRetries; curRetry++ {
		res := int(runfunc())
		if res == C.BMKT_SUCCESS {
			return nil
		}
		if !isErrorTransient(res) {
			return fmt.Errorf("code %d error", res)
		}
		time.Sleep(c.RetryDelay)
	}
	return errors.New("retries exhausted")
}

func (c *BMKTContext) Initialize() error {
	return c.runWithRetry(func() C.int {
		return C.bmkt_init_fps(c.ctx.session)
	})
}

func (c *BMKTContext) Close() {
	if c.ctx == nil {
		return
	}
	ctx := c.ctx
	c.ctx = nil
	C.bmkt_main_close(ctx)
}
