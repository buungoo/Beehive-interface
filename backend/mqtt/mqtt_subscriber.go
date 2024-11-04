package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Initialize the log file
var logFile *os.File

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// Get the current timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Log the received message to the log file with timestamp
	if _, err := logFile.WriteString(fmt.Sprintf("%s - Received message: %s from topic: %s\n", timestamp, msg.Payload(), msg.Topic())); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func main() {
	// Set MQTT broker URL
	// broker := "tcp://localhost:1883" // or any broker you want
	// broker := "broker.emqx.io:1883"//"tcp://broker.hivemq.com:1883"   // HiveMQ public broker
	broker := "tcp://broker.hivemq.com:1883"   // HiveMQ public broker
	topic := "d0039ebeehive/sensor"

	// Open log file
	var err error
	logFile, err = os.OpenFile("./subscriberlog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("go_mqtt_subscriber")

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

	// Channel for graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Keep the subscriber running
	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()

	// Wait for shutdown signal
	<-stopChan

	// Cleanup before exiting
	client.Disconnect(250)
	fmt.Println("Subscriber disconnected and exiting")
}
