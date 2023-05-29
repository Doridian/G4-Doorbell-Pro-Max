package bmkt

// #include "handler.h"
// #cgo LDFLAGS: -lbmkt
import "C"
import (
	"errors"
	"fmt"
)

type BMKTContext struct {
	ctx interface{}
}

const IF_STATE_INVALID = -1
const IF_STATE_IDLE = 0

func Open() (*BMKTContext, error) {
	ctxPtr := C.bmkt_main_init()
	if ctxPtr == nil {
		return nil, errors.New("unknown error allocating context")
	}

	if ctxPtr.state == IF_STATE_INVALID {
		return nil, fmt.Errorf("code %d initializing BMKT", ctxPtr.last_error)
	}

	return &BMKTContext{
		ctx: ctxPtr,
	}, nil
}
