package bmkt

// #include <libbmkt/custom.h>
import "C"
import "unsafe"
import "log"

const (
	IF_STATE_INIT         = 0
	IF_STATE_IDLE         = iota
	IF_STATE_ENROLLING    = iota
	IF_STATE_VERIFYING    = iota
	IF_STATE_IDENTIFYING  = iota
	IF_STATE_CANCELLING   = iota
	IF_STATE_DELETING_ALL = iota
	IF_STATE_INVALID      = iota
)

func (c *BMKTContext) handleResponseError(resp *C.bmkt_response_t, op string) {
	c.state = IF_STATE_IDLE
	log.Printf("Got error %d during %s", resp.result, op)
}

func (c *BMKTContext) handleEnrollProgress(progress int) {
	log.Printf("Enroll progress %d %%", progress)
}

func (ctx *BMKTContext) handleResponse(resp *C.bmkt_response_t) {
	switch resp.response_id {
	// Events
	case C.BMKT_EVT_FINGER_REPORT:
		finger_event := (*C.bmkt_finger_event_t)(unsafe.Pointer(&resp.response))
		switch finger_event.finger_state {
		case C.BMKT_EVT_FINGER_STATE_NOT_ON_SENSOR:
			log.Printf("Finger removed from sensor!\n")
			if ctx.state == IF_STATE_IDLE {
				ctx.Identify()
			}
		case C.BMKT_EVT_FINGER_STATE_ON_SENSOR:
			log.Printf("Finger placed on sensor!\n")
		}

	// Init
	case C.BMKT_RSP_FPS_INIT_OK:
		ctx.state = IF_STATE_IDLE
		log.Printf("Init OK!\n")
		init_resp := (*C.bmkt_init_resp_t)(unsafe.Pointer(&resp.response))
		if int(init_resp.finger_presence) == 0 {
			ctx.Identify()
		}
	case C.BMKT_RSP_FPS_INIT_FAIL:
		ctx.handleResponseError(resp, "Init")

	// Enrollment
	case C.BMKT_RSP_ENROLL_READY:
		ctx.state = IF_STATE_ENROLLING
		ctx.handleEnrollProgress(0)
	case C.BMKT_RSP_ENROLL_OK:
		ctx.state = IF_STATE_IDLE
		ctx.handleEnrollProgress(100)
	case C.BMKT_RSP_ENROLL_FAIL:
		ctx.handleResponseError(resp, "Enroll")
		ctx.handleEnrollProgress(-1)
	case C.BMKT_RSP_ENROLL_REPORT:
		ctx.state = IF_STATE_ENROLLING
		enroll_resp := (*C.bmkt_enroll_resp_t)(unsafe.Pointer(&resp.response))
		ctx.handleEnrollProgress(int(enroll_resp.progress))

	// Verify / verify_resp
	case C.BMKT_RSP_VERIFY_READY:
		ctx.state = IF_STATE_VERIFYING
		log.Printf("Verify started!\n")
	case C.BMKT_RSP_VERIFY_OK:
		ctx.state = IF_STATE_IDLE
		verify_resp := (*C.bmkt_verify_resp_t)(unsafe.Pointer(&resp.response))

		user_id := convertCUserIDToString(&verify_resp.user_id)
		finger_id := int(verify_resp.finger_id)

		log.Printf("Verify OK! You are %s finger %d\n", user_id, finger_id)
	case C.BMKT_RSP_VERIFY_FAIL:
		ctx.handleResponseError(resp, "Verify")

	// Identify / id_resp
	case C.BMKT_RSP_ID_READY:
		ctx.state = IF_STATE_IDENTIFYING
		log.Printf("Identify started!\n")
	case C.BMKT_RSP_ID_OK:
		ctx.state = IF_STATE_IDLE
		id_resp := (*C.bmkt_identify_resp_t)(unsafe.Pointer(&resp.response))

		user_id := convertCUserIDToString(&id_resp.user_id)
		finger_id := int(id_resp.finger_id)

		log.Printf("Identify OK! You are %s finger %d\n", user_id, finger_id)
	case C.BMKT_RSP_ID_FAIL:
		ctx.handleResponseError(resp, "Identify")

	// Op cancalltion
	case C.BMKT_RSP_CANCEL_OP_OK:
		ctx.state = IF_STATE_IDLE
		log.Printf("Cancel OK!\n")
	case C.BMKT_RSP_CANCEL_OP_FAIL:
		ctx.handleResponseError(resp, "Cancel")

	// Deletion
	case C.BMKT_RSP_DELETE_PROGRESS:
		ctx.state = IF_STATE_DELETING_ALL
		del_all_users_resp := (*C.bmkt_del_all_users_resp_t)(unsafe.Pointer(&resp.response))
		log.Printf("Delete all progress %d\n", del_all_users_resp.progress)
	case C.BMKT_RSP_DEL_FULL_DB_OK:
		ctx.state = IF_STATE_IDLE
		log.Printf("Delete all OK!\n")
	case C.BMKT_RSP_DEL_FULL_DB_FAIL:
		ctx.handleResponseError(resp, "Delete all")
	case C.BMKT_RSP_DEL_USER_FP_OK:
		ctx.state = IF_STATE_IDLE
		log.Printf("Delete user OK!\n")
	case C.BMKT_RSP_DEL_USER_FP_FAIL:
		ctx.handleResponseError(resp, "Delete user")

	// Unhandled
	default:
		log.Printf("on_response(%d / 0x%02x)\n", resp.response_id, resp.response_id)
	}
}
