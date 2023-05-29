package bmkt

import "log"

func (ctx *BMKTContext) handleError(code int) {
	log.Printf("Got C error %d", code)
}
