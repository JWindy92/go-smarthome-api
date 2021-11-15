package mqtt_utils

import (
	"fmt"

	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var broker = "10.0.0.228"
var port = 9002

var Zap = zap.NewLogger("mqtt_utils.go")

type MqttHandler struct {
	Client *mqtt.Client
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s", msg.Payload())
	Zap.Logger.Infow(
		"Received message",
		"topic", msg.Topic(),
		"message", string(msg.Payload()),
	)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	Zap.Logger.Infow(
		"Connected to mqtt broker",
		"host", broker,
		"port", port,
	)
	client.Subscribe("test", 1, messageHandler)
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	Zap.Logger.Warnw(
		"Connection lost to mqtt broker",
		"error", err,
	)
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	Zap.Logger.Infow(
		"Received message",
		"topic", msg.Topic(),
		"message", string(msg.Payload()),
	)
}

func MqttInit() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ws://%s:%d", broker, port))
	opts.SetClientID("smarthome_mqtt_client")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
