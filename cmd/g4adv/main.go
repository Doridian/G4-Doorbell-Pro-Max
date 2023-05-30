package main

import (
	"os"
	"time"

	"github.com/Doridian/G4-Doorbell-Pro-Max/util/bmkt"
	"github.com/rs/zerolog"
)

func main() {
	os.Setenv("TZ", "Etc/UTC")
	consoleWriter := zerolog.NewConsoleWriter()
	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter.TimeFormat = time.RFC3339
	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

	bmktLogger := logger.With().Str("component", "bmkt").Logger()
	bmktCtx, err := bmkt.New(bmktLogger)
	if err != nil {
		panic(err)
	}
	err = bmktCtx.Open()
	if err != nil {
		panic(err)
	}

	logger.Info().Msg("Entering idle main loop...")

	for {
		time.Sleep(time.Second * 1)
	}
}
