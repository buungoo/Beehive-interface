
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Set MQTT broker URL
	broker := "broker.hivemq.com:1883" // HiveMQ public broker
	topic := "d0039ebeehive/sensor"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("beehive_publisher")

	// Create client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to broker: %v", token.Error())
		os.Exit(1)
	}

	// Publish data to the topic at regular intervals
	for {
		// Generate the current time
		currentTime := time.Now().Format(time.RFC3339) // ISO 8601 format

		// Construct JSON payload with the current time
		//AwEgAgHs,
		//AgUC
		//BgULwrg=
		payload := fmt.Sprintf(`{
			"applicationID": "1",
			"applicationName": "beehive-sensor-card",
			"data": "AgUC",
			"devEUI": "0080e115000adf82",
			"deviceName": "beehive-sensor-card-dn",
			"fCnt": 200,
			"fPort": 2,
			"rxInfo": [{
				"altitude": 0,
				"latitude": 0,
				"loRaSNR": 5.2,
				"longitude": 0,
				"mac": "24e124fffef0b4f9",
				"name": "24e124fffef0b4f9",
				"rssi": -113,
				"time": "%s"
			}],
			"time": "%s",
			"txInfo": {
				"adr": true,
				"codeRate": "4/5",
				"dataRate": {
					"bandwidth": 125,
					"modulation": "LORA",
					"spreadFactor": 7
				},
				"frequency": 868100000
			}
		}`, currentTime, currentTime)

		// Publish the payload
		token := client.Publish(topic, 0, false, payload)
		token.Wait()

		fmt.Printf("Published %s to topic %s\n", payload, topic)

		// Wait before sending the next message
		time.Sleep(5 * time.Second)
	}

	// Disconnect when done
	client.Disconnect(250)
}

