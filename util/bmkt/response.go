package bmkt

// #include <libbmkt/bmkt_response.h>
import "C"
import "log"

func (c *BMKTContext) handleResponse(response *C.bmkt_response_t) {
	log.Printf("Got C response %v", response)
}
