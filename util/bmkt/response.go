package bmkt

// #include <libbmkt/custom.h>
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/rs/zerolog"
)

const (
	IF_STATE_INIT          = 0
	IF_STATE_IDLE          = iota
	IF_STATE_ENROLLING     = iota
	IF_STATE_VERIFYING     = iota
	IF_STATE_IDENTIFYING   = iota
	IF_STATE_CANCELLING    = iota
	IF_STATE_DELETING_ALL  = iota
	IF_STATE_DELETING_USER = iota
	IF_STATE_INVALID       = iota
)

func (ctx *BMKTContext) handleResponseCode(resp *C.bmkt_response_t, op string) {
	ctx.state = IF_STATE_IDLE
	ctx.handleResponseCodeNoIdle(resp, op)
}

func (ctx *BMKTContext) handleResponseCodeNoIdle(resp *C.bmkt_response_t, op string) {
	var evt *zerolog.Event
	if resp.result == C.BMKT_SUCCESS {
		evt = ctx.logger.Info().Bool("success", true)
	} else {
		evt = ctx.logger.Warn().Bool("success", false)
	}
	evt.Str("type", "sensor_response").Str("op", op).Int("result", int(resp.result)).Send()
}

func (ctx *BMKTContext) handleEnrollProgress(progress int) {
	ctx.logger.Info().Str("type", "enroll_progress").Int("progress", progress).Send()
}

func (ctx *BMKTContext) handleDeleteAllProgress(progress int) {
	ctx.logger.Info().Str("type", "delete_all_progress").Int("progress", progress).Send()
}

func (ctx *BMKTContext) handleFingerPresence(present bool, op string) {
	ctx.logger.Info().Str("type", "finger_presence").Str("op", op).Bool("present", present).Send()
	go ctx.autoIdentify()
}

func (ctx *BMKTContext) handleResponse(resp *C.bmkt_response_t) {
	switch resp.response_id {
	// Events
	case C.BMKT_EVT_FINGER_REPORT:
		finger_event := (*C.bmkt_finger_event_t)(unsafe.Pointer(&resp.response))
		switch finger_event.finger_state {
		case C.BMKT_EVT_FINGER_STATE_NOT_ON_SENSOR:
			ctx.handleFingerPresence(false, "finger_report")
		case C.BMKT_EVT_FINGER_STATE_ON_SENSOR:
			ctx.handleFingerPresence(true, "finger_report")
		}

	// Init
	case C.BMKT_RSP_FPS_INIT_OK:
		init_resp := (*C.bmkt_init_resp_t)(unsafe.Pointer(&resp.response))
		is_finger_present := int(init_resp.finger_presence) == 0
		ctx.handleFingerPresence(is_finger_present, "init_fps")
		fallthrough
	case C.BMKT_RSP_FPS_INIT_FAIL:
		ctx.lastInitResult = resp.result
		ctx.handleResponseCode(resp, "init_fps")

	// Enrollment
	case C.BMKT_RSP_ENROLL_READY:
		ctx.state = IF_STATE_ENROLLING
		ctx.lastEnrollResult = -1
		ctx.handleEnrollProgress(0)
	case C.BMKT_RSP_ENROLL_OK:
		fallthrough
	case C.BMKT_RSP_ENROLL_FAIL:
		ctx.handleEnrollProgress(100)
		ctx.lastEnrollResult = resp.result
		ctx.handleResponseCode(resp, "enroll")
	case C.BMKT_RSP_ENROLL_REPORT:
		ctx.state = IF_STATE_ENROLLING
		ctx.lastEnrollResult = -1
		enroll_resp := (*C.bmkt_enroll_resp_t)(unsafe.Pointer(&resp.response))
		ctx.handleEnrollProgress(int(enroll_resp.progress))

	// Verify / verify_resp
	case C.BMKT_RSP_VERIFY_READY:
		ctx.state = IF_STATE_VERIFYING
		ctx.lastVerifyResult = -1
		ctx.lastVerifyFinger = -1
		ctx.lastVerifyUsername = ""
		ctx.logger.Info().Str("type", "sensor_ready").Str("op", "verify").Send()
	case C.BMKT_RSP_VERIFY_OK:
		verify_resp := (*C.bmkt_verify_resp_t)(unsafe.Pointer(&resp.response))

		user_id := convertCUserIDToString(&verify_resp.user_id)
		finger_id := int(verify_resp.finger_id)

		ctx.lastVerifyUsername = user_id
		ctx.lastVerifyFinger = finger_id
		fallthrough
	case C.BMKT_RSP_VERIFY_FAIL:
		ctx.lastVerifyResult = resp.result
		ctx.handleResponseCode(resp, "verify")

	// Identify / id_resp
	case C.BMKT_RSP_ID_READY:
		ctx.state = IF_STATE_IDENTIFYING
		ctx.lastIdentifyResult = -1
		ctx.lastIdentifyFinger = -1
		ctx.lastIdentifyUser = ""
		ctx.logger.Info().Str("type", "sensor_ready").Str("op", "identify").Send()
	case C.BMKT_RSP_ID_OK:
		id_resp := (*C.bmkt_identify_resp_t)(unsafe.Pointer(&resp.response))

		user_id := convertCUserIDToString(&id_resp.user_id)
		finger_id := int(id_resp.finger_id)

		ctx.lastIdentifyUser = user_id
		ctx.lastIdentifyFinger = finger_id
		fallthrough
	case C.BMKT_RSP_ID_FAIL:
		ctx.lastIdentifyResult = resp.result
		ctx.handleResponseCode(resp, "identify")

	// Op cancalltion
	case C.BMKT_RSP_CANCEL_OP_OK:
		fallthrough
	case C.BMKT_RSP_CANCEL_OP_FAIL:
		ctx.lastCancelResult = resp.result
		ctx.handleResponseCode(resp, "cancel")

	// Delete all
	case C.BMKT_RSP_DELETE_PROGRESS:
		ctx.state = IF_STATE_DELETING_ALL
		ctx.lastDeleteAllResult = -1
		del_all_users_resp := (*C.bmkt_del_all_users_resp_t)(unsafe.Pointer(&resp.response))
		ctx.handleDeleteAllProgress(int(del_all_users_resp.progress))
	case C.BMKT_RSP_DEL_FULL_DB_OK:
		fallthrough
	case C.BMKT_RSP_DEL_FULL_DB_FAIL:
		ctx.lastDeleteAllResult = resp.result
		ctx.handleResponseCode(resp, "delete_all")

	// Delete user/finger
	case C.BMKT_RSP_DEL_USER_FP_OK:
		fallthrough
	case C.BMKT_RSP_DEL_USER_FP_FAIL:
		ctx.lastDeleteUserResult = resp.result
		ctx.handleResponseCode(resp, "delete_user")

	case C.BMKT_RSP_CAPTURE_COMPLETE:
		ctx.handleResponseCodeNoIdle(resp, "capture_complete")

	// Unhandled
	default:
		ctx.handleResponseCode(resp, fmt.Sprintf("unknown_0x%02x", resp.response_id))
	}
}
