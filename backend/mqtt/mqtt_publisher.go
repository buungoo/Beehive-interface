package main

import (
    "fmt"
    "log"
    "math/rand"
    "os"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    // Set MQTT broker URL
    broker := "tcp://broker.hivemq.com:1883" // HiveMQ public broker
    topic := "test/topic"

    opts := mqtt.NewClientOptions()
    opts.AddBroker(broker)
    opts.SetClientID("go_mqtt_publisher")

    // If we set up authentication on the broker
    // opts.SetUsername("your_username")
    // opts.SetPassword("your_password")

    // Create client
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Fatalf("Error connecting to broker: %v", token.Error())
        os.Exit(1)
    }

    // Publish float data to the topic at regular intervals
    for {
        // Generate a random float value
        sensorValue := rand.Float64() * 100
        payload := fmt.Sprintf("%f", sensorValue)
        
        // Publish the value
        token := client.Publish(topic, 0, false, payload)
        token.Wait()

        fmt.Printf("Published %s to topic %s\n", payload, topic)

        time.Sleep(5 * time.Second) // Send every 5 seconds
    }

    // Disconnect when done
    client.Disconnect(250)
}
