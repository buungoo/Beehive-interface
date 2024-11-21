package main

import (
	"fmt"
	"log"
	// "math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Set MQTT broker URL
	// broker := "tcp://localhost:1883" // If we are running our own local broker
	// broker := "broker.emqx.io:1883" // emqx public broker
	broker := "broker.hivemq.com:1883" // HiveMQ public broker
	topic := "d0039ebeehive/sensor"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("beehive_publisher")

	// If we set up authentication on the broker
	// opts.SetUsername("your_username")
	// opts.SetPassword("your_password")

	// Create client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
		os.Exit(1)
	}

	// JSON payload
	sensorValue := `{"applicationID":"1","applicationName":"beehive-sensor-card","data":"AwEgAgHs","devEUI":"0080e115000adf82","deviceName":"beehive-sensor-card-dn","fCnt":200,"fPort":2,"rxInfo":[{"altitude":0,"latitude":0,"loRaSNR":5.2,"longitude":0,"mac":"24e124fffef0b4f9","name":"24e124fffef0b4f9","rssi":-113,"time":"2024-11-07T13:39:29.776959Z"}],"time":"2024-11-07T13:39:29.776959Z","txInfo":{"adr":true,"codeRate":"4/5","dataRate":{"bandwidth":125,"modulation":"LORA","spreadFactor":7},"frequency":868100000}}`

	// Publish float data to the topic at regular intervals
	for {
		// // Generate a random float value
		// sensorValue := rand.Float64() * 100
		// payload := fmt.Sprintf("%f", sensorValue)
		//
		// // Publish the value
		// token := client.Publish(topic, 0, false, payload)
		// token.Wait()
		//
		// fmt.Printf("Published %s to topic %s\n", payload, topic)
		//
		// time.Sleep(5 * time.Second) // Send every 5 seconds

		// Publish the value
		token := client.Publish(topic, 0, false, sensorValue)
		token.Wait()

		fmt.Printf("Published %s to topic %s\n", sensorValue, topic)

		time.Sleep(5 * time.Second) // Send every 5 seconds
	}

	// Disconnect when done
	client.Disconnect(250)
}
