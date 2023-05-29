package bmkt

// #include "handler.h"
// #cgo LDFLAGS: -lbmkt
import "C"
import (
	"errors"
	"fmt"
)

type BMKTContext struct {
	ctx *C.cb_ctx_t
}

func Open() (*BMKTContext, error) {
	ctxPtr := C.bmkt_main_init()
	if ctxPtr == nil {
		return nil, errors.New("unknown error allocating context")
	}

	if ctxPtr.state == C.IF_STATE_INVALID {
		return nil, fmt.Errorf("code %d initializing BMKT", ctxPtr.last_error)
	}

	return &BMKTContext{
		ctx: ctxPtr,
	}, nil
}

func (c *BMKTContext) Close() {
	if c.ctx == nil {
		return
	}
	ctx := c.ctx
	c.ctx = nil
	C.bmkt_main_close(ctx)
}
