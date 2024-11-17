package main

import (
	"encoding/base64"
	"strings"

	// "encoding/binary"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	// "strings"
	"syscall"
	"time"
)

// Sensor enum to represent sensor types
type Sensor uint8

const (
	LoadCell    Sensor = 1
	Temperature        = 2
	Humidity           = 3
	Microphone         = 4
	Oxygen             = 5
)

// SensorReading struct to hold a parsed sensor reading
type SensorReading struct {
	SensorType Sensor
	SensorID   uint8
	Value      interface{} // Allows for different types, like int8 or uint8
	Timestamp  time.Time   // To store the timestamp of the reading
}

// Builder pattern for SensorReading
type SensorReadingBuilder struct {
	sensorType Sensor
	sensorID   uint8
	value      interface{}
	timestamp  time.Time
}

func NewSensorReadingBuilder(sensorType Sensor, timestamp time.Time) *SensorReadingBuilder {
	return &SensorReadingBuilder{sensorType: sensorType, timestamp: timestamp}
}
func (b *SensorReadingBuilder) SetSensorID(id uint8) *SensorReadingBuilder {
	b.sensorID = id
	return b
}

func (b *SensorReadingBuilder) SetValue(value interface{}) *SensorReadingBuilder {
	switch b.sensorType {
	case Temperature:
		if v, ok := value.(int8); ok {
			b.value = v
		} else {
			log.Println("Received invalid value type for Temperature. Expected int8.")
		}
	case Microphone:
		if v, ok := value.(uint8); ok {
			b.value = v == 1 // Microphone can be either 0 or 1. We should assign a boolean value.
		} else {
			log.Println("Received invalid value type for Microphone. Expected uint8.")
		}
	default:
		if v, ok := value.(uint8); ok {
			b.value = v
		} else {
			log.Println("Received invalid value type for default sensors. Expected uint8.")
		}
	}
	return b
}

func (b *SensorReadingBuilder) Build() *SensorReading {
	return &SensorReading{
		SensorType: b.sensorType,
		SensorID:   b.sensorID,
		Value:      b.value,
		Timestamp:  b.timestamp,
	}
}

func parseSensorMessage(message SensorMessage) ([]*SensorReading, error) {
	decodedData, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 data: %v", err)
	}

	if dotIndex := strings.Index(message.Time, "."); dotIndex != -1 {
		message.Time = message.Time[:dotIndex] + "Z" // Add "Z" to indicate UTC
		// fmt.Println(message.Time)
	}

	timeStamp, err := time.Parse("2006-01-02T15:04:05Z", message.Time)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}

	var readings []*SensorReading
	for i := 0; i < len(decodedData); i += 3 {
		sensorType := Sensor(decodedData[i])
		sensorId := decodedData[i+1]
		rawValue := decodedData[i+2]

		builder := NewSensorReadingBuilder(sensorType, timeStamp).SetSensorID(sensorId)
		switch sensorType {
		case Temperature:
			builder.SetValue(int8(rawValue)) // Temperature uses int8
		default:
			builder.SetValue(rawValue) // Other sensors use uint8
		}

		readings = append(readings, builder.Build())
	}

	return readings, nil
}

// Mock function to simulate message handling
func handleSensorMessage(message SensorMessage) {
	readings, err := parseSensorMessage(message)
	if err != nil {
		log.Printf("Error parsing sensor message: %v", err)
		return
	}

	for _, reading := range readings {
		fmt.Printf("Sensor Reading: %+v\n", reading)
		// Parse the message into sensor objects
	}
}

type SensorMessage struct {
	ApplicationName string `json:"applicationName"`
	Data            string `json:"data"`
	Time            string `json:"time"`
}

// type Sensor uint8
//
// const (
// 	LoadCell    Sensor = 1
// 	Temperature        = 2
// 	Humidity           = 3
// 	Microphone         = 4
// )

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

	handleSensorMessage(sensorMessage) //.Data)

	// h, err := base64.StdEncoding.DecodeString(sensorMessage.Data)
	// if err != nil {
	// 	// handle error
	// }
	//
	// fmt.Println(h)

	// sensorId byte := 0
	// sensorType := 0
	// sensorValue := 0
	// sensorId, sensorType, sensorValue byte
	// sensorId, sensorType, sensorValue := byte(0), byte(0), byte(0)
	//
	// for index, element := range h {
	// 	fmt.Println("Index:", index, "Element:", element)
	//
	// 	switch {
	// 	case index%3 == 1:
	// 		sensorType = element
	// 		fmt.Println("sensorType assigned:", sensorType)
	// 	case index%3 == 2:
	// 		sensorId = element
	// 		fmt.Println("sensorId assigned:", sensorId)
	// 	case index%3 == 0:
	// 		sensorValue = element
	// 		fmt.Println("sensorValue assigned:", sensorValue)
	// 	}
	// }
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
	opts.SetClientID("local_subscriber")
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
