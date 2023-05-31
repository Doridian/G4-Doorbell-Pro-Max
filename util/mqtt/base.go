package mqtt

import (
	"fmt"
	"os"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type MQTTContext struct {
	client      paho.Client
	topicPrefix string
}

type MQTTHandler func(msg []byte)

func New(host string, username string, password string, topicPrefix string) (*MQTTContext, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	opts := paho.NewClientOptions()
	opts.SetWriteTimeout(time.Second * 10)
	opts.SetClientID(fmt.Sprintf("g4adv-%s", hostname))

	opts.AddBroker(fmt.Sprintf("tcp://%s", host))
	opts.SetUsername(username)
	opts.SetPassword(password)

	client := paho.NewClient(opts)

	t := client.Connect()
	t.Wait()
	err = t.Error()
	if err != nil {
		return nil, err
	}

	return &MQTTContext{
		client:      client,
		topicPrefix: topicPrefix,
	}, nil
}

func (ctx *MQTTContext) Close() {
	client := ctx.client
	ctx.client = nil
	if client == nil {
		return
	}
	client.Disconnect(1000)
}

func (ctx *MQTTContext) makeTopic(topic string) string {
	return ctx.topicPrefix + topic
}

func (ctx *MQTTContext) Publish(topic string, msg []byte) error {
	t := ctx.client.Publish(ctx.makeTopic(topic), 0, false, msg)
	t.Wait()
	return t.Error()
}

func (ctx *MQTTContext) Subscribe(topic string, callback MQTTHandler) error {
	pahoHandler := func(clt paho.Client, msg paho.Message) {
		callback(msg.Payload())
		msg.Ack()
	}
	t := ctx.client.Subscribe(ctx.makeTopic(topic), 0, pahoHandler)
	t.Wait()
	return t.Error()
}
