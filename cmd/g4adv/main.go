package main

import (
	"log"
	"time"

	"github.com/Doridian/G4-Doorbell-Pro-Max/util/bmkt"
)

func main() {
	bmktCtx, err := bmkt.New()
	if err != nil {
		panic(err)
	}
	err = bmktCtx.Open()
	if err != nil {
		panic(err)
	}

	log.Printf("Entering idle main loop...")

	for {
		time.Sleep(time.Second * 1)
	}
}
