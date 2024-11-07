package main

import (
	"encoding/base64"
	// "encoding/binary"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	// "strings"
	"syscall"
	"time"
)

type SensorMessage struct {
	ApplicationName string `json:"applicationName"`
	Data            string `json:"data"`
	Time            string `json:"time"`
}

// Initialize the log file
var logFile *os.File

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// Get the current timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Parse the received message payload
	var sensorMessage SensorMessage
	if err := json.Unmarshal(msg.Payload(), &sensorMessage); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return
	}

	// Parse the message as it is received like this:
	// {
	// "applicationID":"1",
	// "applicationName":"beehive-sensor-card",
	// "data":"YWFhYWE=",
	// "devEUI":"0080e115000adf82",
	// "deviceName":"beehive-sensor-card-dn",
	// "fCnt":319,
	// "fPort":2,
	// "rxInfo":[{"altitude":0,"latitude":0,"loRaSNR":7.5,"longitude":0,"mac":"24e124fffef0b4f9","name":"24e124fffef0b4f9","rssi":-109,"time":"2024-11-05T14:01:49.217376Z"}],
	// "time":"2024-11-05T14:01:49.217376Z",
	// "txInfo":{"adr":true,"codeRate":"4/5","dataRate":{"bandwidth":125,"modulation":"LORA","spreadFactor":7},"frequency":868300000}
	// }
	logMessage := fmt.Sprintf("%s - applicationName: %s, data: %s, time: %s\n",
		timestamp, sensorMessage.ApplicationName, sensorMessage.Data, sensorMessage.Time)

	// Log the received message to the log file with timestamp
	// if _, err := logFile.WriteString(fmt.Sprintf("%s - Received message: %s from topic: %s\n", timestamp, msg.Payload(), msg.Topic())); err != nil {
	// 	log.Printf("Error writing to log file: %v", err)
	// }

	// Log the parsed message to the log file
	if _, err := logFile.WriteString(logMessage); err != nil {
		log.Printf("Error writing to log file: %v", err)
	}

	// fmt.Println(logMessage)
	// hx := hex.EncodeToString([]byte(sensorMessage.Data))
	// fmt.Println(sensorMessage.Data + " ==> " + hx)

	h, err := base64.StdEncoding.DecodeString(sensorMessage.Data)
	if err != nil {
		// handle error
	}

	// p, err := base64.StdEncoding.DecodeString(sensorMessage.Data)
	// if err != nil {
	// 	// handle error
	// }
	// h := hex.EncodeToString(p)
	fmt.Println(h) // prints 415256494e

	// // Split the hex string into slices of 2 characters each and convert to integers
	// var intValues []int
	// for i := 0; i < len(h); i += 2 {
	// 	hexPair := h[i : i+2]
	// 	intValue, err := strconv.ParseInt(hexPair, 16, 32)
	// 	if err != nil {
	// 		// handle error
	// 		fmt.Println("Error parsing hex pair:", err)
	// 		return
	// 	}
	// 	intValues = append(intValues, int(intValue))
	// }

	// Print the resulting integers
	// fmt.Println("Integer values:", intValues)
	// // myslice :=
	// byteData := binary.BigEndian.Uint16([]byte(h))
	// fmt.Println(byteData)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func main() {
	// Set MQTT broker URL
	// broker := "tcp://localhost:1883" // For testing hosting the broker on a local machine
	// broker := "broker.emqx.io:1883" // emqx public broker
	broker := "broker.hivemq.com:1883" // HiveMQ public broker
	topic := "d0039ebeehive/sensor"

	// Open log file
	var err error
	logFile, err = os.OpenFile("./logs/subscriberlog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("local_subscriber2")
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
