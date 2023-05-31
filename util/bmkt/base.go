package bmkt

// #cgo LDFLAGS: -lbmkt
// #include <libbmkt/custom.h>
import "C"
import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

var maxID uint64
var bmktContexts = make(map[uint64]*Context)

type IdentifyCallback func(user string, finger_id int)

type Context struct {
	MaxRetries       int
	RetryDelay       time.Duration
	AutoIdentify     bool
	IdentifyCallback IdentifyCallback

	id  uint64
	cid C.uint64_t

	sensor      C.bmkt_sensor_t
	session     *C.bmkt_ctx_t
	sessionLock sync.Mutex

	state  int
	logger zerolog.Logger

	lastInitResult       C.int
	lastEnrollResult     C.int
	lastCancelResult     C.int
	lastDeleteUserResult C.int
	lastDeleteAllResult  C.int

	lastVerifyResult   C.int
	lastVerifyUsername string
	lastVerifyFinger   int

	lastIdentifyResult C.int
	lastIdentifyUser   string
	lastIdentifyFinger int
}

func New(logger zerolog.Logger) (*Context, error) {
	id := atomic.AddUint64(&maxID, 1)
	ctx := &Context{
		MaxRetries:   3,
		RetryDelay:   time.Millisecond * 1,
		AutoIdentify: true,

		id:  id,
		cid: C.uint64_t(id),

		state:  IF_STATE_INVALID,
		logger: logger,
	}
	bmktContexts[ctx.id] = ctx

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

func (ctx *Context) Open() error {
	ctx.sessionLock.Lock()
	defer ctx.sessionLock.Unlock()

	err := ctx.wrappedOpen()
	if err != nil {
		ctx.Close()
	}
	return err
}

func (ctx *Context) Close() {
	delete(bmktContexts, ctx.id)
	if ctx.session == nil {
		return
	}
	session := ctx.session
	ctx.session = nil

	C.bmkt_close(session)
	C.bmkt_exit(session)
}
