package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("test/topic", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func main() {
	// Set MQTT broker URL
	broker := "tcp://broker.hivemq.com:1883" // HiveMQ public broker
	topic := "test/topic"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go_mqtt_subscriber")

	// If we set up authentication on the broker
	// opts.SetUsername("your_username")
	// opts.SetPassword("your_password")

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// Create client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
		os.Exit(1)
	}

	// Subscribe to the topic
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Error subscribing to topic: %v", token.Error())
		os.Exit(1)
	}

	// Keep the subscriber running
	for {
		time.Sleep(1 * time.Second)
	}

	// Disconnect when done
	// client.Disconnect(250)
}
