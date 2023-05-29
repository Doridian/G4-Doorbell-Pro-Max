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
	sensor     C.bmkt_sensor_t
	ctx        *C.bmkt_ctx_t
	MaxRetries int
	RetryDelay time.Duration
}

type runnable func() C.int

func New() (*BMKTContext, error) {
	ctx := &BMKTContext{
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 1,
	}

	// Type 0 means SPI in this library
	ctx.sensor.transport_type = C.SENSOR_TRANSPORT_SPI

	// SPI settings
	ctx.sensor.transport_info.addr = 1
	ctx.sensor.transport_info.subaddr = 1
	ctx.sensor.transport_info.mode = C.SPI_MODE_0
	ctx.sensor.transport_info.speed = 4000000
	ctx.sensor.transport_info.bpw = 8

	// GPIO pin information
	ctx.sensor.transport_info.pin_out.pin = 68
	ctx.sensor.transport_info.pin_out.direction = C.GPIO_DIRECTION_OUT
	ctx.sensor.transport_info.pin_out.edge = C.GPIO_EDGE_NONE
	ctx.sensor.transport_info.pin_out.active_low = 0

	ctx.sensor.transport_info.pin_in.pin = 69
	ctx.sensor.transport_info.pin_in.direction = C.GPIO_DIRECTION_IN
	ctx.sensor.transport_info.pin_in.edge = C.GPIO_EDGE_RISING
	ctx.sensor.transport_info.pin_in.active_low = 0

	// No idea, might just be padding
	ctx.sensor.transport_info.unknown_padding = 0

	return ctx, nil
}

func isErrorTransient(err C.int) bool {
	return err == C.BMKT_SENSOR_NOT_READY || err == C.BMKT_TIMEOUT
}

func wrapBMKTError(err C.int) error {
	if err == C.BMKT_SUCCESS {
		return nil
	}
	return fmt.Errorf("code %d error", int(err))
}

func (c *BMKTContext) runWithRetry(runfunc runnable) error {
	for curRetry := 0; curRetry < c.MaxRetries; curRetry++ {
		res := runfunc()
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

func (c *BMKTContext) Open() error {
	c.ctx = C.bmkt_wrapped_init()
	if c.ctx == nil {
		return errors.New("bmkt_init() failure")
	}
	err := wrapBMKTError(C.bmkt_wrapped_open(c.ctx, &c.sensor))
	if err != nil {
		return err
	}
	return c.runWithRetry(func() C.int {
		return C.bmkt_init_fps(c.ctx)
	})
}

func (c *BMKTContext) Close() {
	if c.ctx == nil {
		return
	}
	ctx := c.ctx
	c.ctx = nil

	C.bmkt_close(ctx)
	C.bmkt_exit(ctx)
}
