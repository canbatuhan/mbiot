package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	BROKER = "broker.hivemq.com"
	PORT   = 1883
	TOPIC  = "mbiot/soil_moisture"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func SetMQTTOptions(mqtt_options *mqtt.ClientOptions) {
	mqtt_options.AddBroker(fmt.Sprintf("tcp://%s:%d", BROKER, PORT))
	mqtt_options.SetClientID("go_mqtt_client")
	mqtt_options.SetDefaultPublishHandler(messagePubHandler)
	mqtt_options.OnConnect = connectHandler
	mqtt_options.OnConnectionLost = connectLostHandler
}

func subscribe(client mqtt.Client) {
	token := client.Subscribe(TOPIC, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n", TOPIC)
}

func main() {
	mqttOptions := mqtt.NewClientOptions()
	SetMQTTOptions(mqttOptions)

	client := mqtt.NewClient(mqttOptions)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	subscribe(client)
	for client.IsConnected() {
		// infinite loop
	}
}
