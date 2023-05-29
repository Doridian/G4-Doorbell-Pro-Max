package bmkt

import "log"

func (c *BMKTContext) handleError(code int) {
	log.Printf("Got C error %d", code)
}
