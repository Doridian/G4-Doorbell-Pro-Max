package main

// #include "../../libbmkt/bmkt.h"
// #include "../../libbmkt/custom.h"
// #cgo LDFLAGS: -lnxp-nfc -lbmkt
import "C"

import (
	"log"
)

func main() {
	log.Printf("%v", C.bmkt_init)
}
