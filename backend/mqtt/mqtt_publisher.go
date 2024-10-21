package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	broker := "tcp://broker.hivemq.com:1883" // HiveMQ public broker
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)

	topic := "test/topic"

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker:", token.Error()))
	}

	for i := 1; i <= 10; i++ {
		message := fmt.Sprintf("Publishing message %d", i)
		token := client.Publish(topic, 0, false, message)
		token.Wait()

		fmt.Println("Published:", message)
		time.Sleep(1 * time.Second)
	}

	client.Disconnect(250)
}
