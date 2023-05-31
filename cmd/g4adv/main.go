package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/Doridian/G4-Doorbell-Pro-Max/util/bmkt"
	"github.com/Doridian/G4-Doorbell-Pro-Max/util/mqtt"
	"github.com/rs/zerolog"
)

var mqttClient *mqtt.Client
var logger zerolog.Logger
var bmktCtx *bmkt.Context

func main() {
	os.Setenv("TZ", "Etc/UTC")

	var err error

	consoleWriter := zerolog.NewConsoleWriter()
	zerolog.TimeFieldFormat = time.RFC3339
	consoleWriter.TimeFormat = time.RFC3339
	logger = zerolog.New(consoleWriter).With().Timestamp().Logger()

	mqttClient, err = mqtt.New(os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_USERNAME"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC_PREFIX"))
	if err != nil {
		panic(err)
	}

	bmktLogger := logger.With().Str("component", "bmkt").Logger()
	bmktCtx, err = bmkt.New(bmktLogger)
	bmktCtx.IdentifyCallback = identifyCallback
	if err != nil {
		panic(err)
	}
	err = bmktCtx.Open()
	if err != nil {
		panic(err)
	}

	err = mqttClient.Subscribe("fingerprint_in", mqttListener)
	if err != nil {
		panic(err)
	}

	logger.Info().Msg("Entering idle main loop...")

	for {
		time.Sleep(time.Second * 1)
	}
}

type fingerprintCMD struct {
	Command  string `json:"cmd"`
	User     string `json:"user"`
	FingerID int    `json:"finger_id"`
	Progress int    `json:"progress"`
	Error    string `json:"error"`
}

func mqttSend(msg *fingerprintCMD) error {
	dataBytes, err := json.Marshal(msg)
	if err != nil {
		logger.Err(err).Str("component", "bmkt_mqtt_send").Str("op", "json").Send()
		return err
	}
	err = mqttClient.Publish("fingerprint_out", dataBytes)
	if err != nil {
		logger.Err(err).Str("component", "bmkt_mqtt_send").Str("op", "mqtt_publish").Send()
		return err
	}
	return nil
}

func identifyCallback(user string, finger_id int) {
	mqttSend(&fingerprintCMD{
		Command:  "identify",
		User:     user,
		FingerID: finger_id,
		Progress: -1,
		Error:    "",
	})
}

func mqttListener(msg []byte) {
	var data fingerprintCMD
	err := json.Unmarshal(msg, &data)
	if err != nil {
		logger.Err(err).Str("component", "bmkt_mqtt_listener").Str("op", "json").Send()
		return
	}
	data.Error = ""
	data.Progress = -1

	switch data.Command {
	case "enroll":
		err = bmktCtx.Enroll(data.User, data.FingerID)
		if err != nil {
			data.Error = err.Error()
		}
	default:
		data.Error = "Unknown command"
	}

	mqttSend(&data)
}
