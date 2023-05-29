package main

// #include <libbmkt/bmkt.h>
// #include <libbmkt/custom.h>
// #cgo LDFLAGS: -lnxp-nfc -lbmkt
import "C"

import (
	"github.com/Doridian/G4-Doorbell-Pro-Max/util/bmkt"
	"log"
)

func main() {
	bmktCtx, err := bmkt.Open()
	if err != nil {
		panic(err)
	}
	err = bmktCtx.Initialize()
	if err != nil {
		panic(err)
	}

	log.Printf("%v", bmktCtx)
}
