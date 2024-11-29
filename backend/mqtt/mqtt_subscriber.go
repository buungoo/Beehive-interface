package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/buungoo/Beehive-interface/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	SensorType    Sensor
	SensorID      uint8
	Value         interface{}      // Allows for the different types we need, i.e. uint8, int8, bool
	Timestamp     time.Time        // To store the timestamp of the reading
	ParentBeehive net.HardwareAddr // MAC address of the parent Beehive
}

// Builder pattern for SensorReading
type SensorReadingBuilder struct {
	sensorType    Sensor
	sensorID      uint8
	value         interface{}
	timestamp     time.Time
	parentBeehive net.HardwareAddr
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
			utils.LogWarn("Received invalid value type for Temperature. Expected int8.")
		}
	case Microphone:
		if v, ok := value.(uint8); ok {
			b.value = v == 1 // Microphone can be either 0 or 1 and we should assign a boolean value.
		} else {
			utils.LogWarn("Invalid value type for Temperature. Expected int8.")
		}
	default:
		if v, ok := value.(uint8); ok {
			b.value = v
		} else {
			utils.LogWarn("Received invalid value type for default sensors. Expected uint8.")
		}
	}
	return b
}

func (b *SensorReadingBuilder) SetDevEUI(parentBeehive string) *SensorReadingBuilder {
	// Ensure the string is the correct length for a MAC address
	if len(parentBeehive) != 16 {
		utils.LogWarn("Invalid DevEUI length. Expected 16 characters.")
		return b
	}

	// Insert colons to format as a MAC address
	macFormatted := strings.ToLower(parentBeehive[:2] + ":" +
		parentBeehive[2:4] + ":" +
		parentBeehive[4:6] + ":" +
		parentBeehive[6:8] + ":" +
		parentBeehive[8:10] + ":" +
		parentBeehive[10:12])

	// Parse the formatted MAC address
	mac, err := net.ParseMAC(macFormatted)
	if err != nil {
		utils.LogWarn("Failed to parse DevEUI as MAC address.")
	} else {
		b.parentBeehive = mac
	}

	return b
}

func (b *SensorReadingBuilder) Build() *SensorReading {
	return &SensorReading{
		SensorType:    b.sensorType,
		SensorID:      b.sensorID,
		Value:         b.value,
		Timestamp:     b.timestamp,
		ParentBeehive: b.parentBeehive,
	}
}

func parseSensorMessage(message SensorMessage) ([]*SensorReading, error) {
	decodedData, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 data: %v", err)
	}

	if dotIndex := strings.Index(message.Time, "."); dotIndex != -1 {
		message.Time = message.Time[:dotIndex] + "Z" // Add "Z" to indicate UTC
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

		builder := NewSensorReadingBuilder(sensorType, timeStamp /* net.HardwareAddr(message.DevEUI) */).SetSensorID(sensorId).SetDevEUI(message.DevEUI)
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

func handleSensorMessage(message SensorMessage) {
	readings, err := parseSensorMessage(message)
	if err != nil {
		utils.LogError("Error parsing sensor message: %v", err)
		return
	}

	for _, reading := range readings {
		utils.LogInfo(fmt.Sprintf("Sensor Reading: %+v", reading))
		// Parse the message into sensor objects

		// Insert the parsed reading into the database
		// if err := insertSensorReading(reading); err != nil {
		// 	utils.LogError("Failed to insert sensor reading: %v", err)
		// }
	}
}

type SensorMessage struct {
	ApplicationName string `json:"applicationName"`
	Data            string `json:"data"`
	Time            string `json:"time"`
	DevEUI          string `json:"devEUI"`
}

// Initialize the log file
// var logFile *os.File

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// Get the current timestamp
	// not needed if we use the new logging function
	// timestamp := time.Now().Format("2006-01-02 15:04:05")

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

	// utils.LogInfo(fmt.Sprintf("Received message - applicationName: %s, data: %s, time: %s",
	// 	sensorMessage.ApplicationName, sensorMessage.Data, sensorMessage.Time))
	fmt.Println("Received message:", sensorMessage)
	utils.LogInfo(fmt.Sprintf("Received message: %+v", sensorMessage))

	handleSensorMessage(sensorMessage) //.Data)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	utils.LogError("Connection lost", err)
	fmt.Printf("Connection lost: %v\n", err)

	// Retry to connect if connection was lost
	for {
		fmt.Println("Attempting to reconnect...")
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			fmt.Println("Reconnected successfully")
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // Add a task to the WaitGroup

	// HiveMQ public broker
	broker := "broker.hivemq.com:1883"
	topic := "d0039ebeehive/sensor"

	// Initialize the logger
	logFile, err := utils.InitLogger("./logs/subscriberlog")
	if err != nil {
		utils.LogFatal("Failed to initialize logger", err)
	}
	// Ensure the log file is closed when the program terminates
	defer logFile.Close()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("beehive_subscriber")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// Create client
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		utils.LogError("Error connecting to broker", token.Error())
	}

	// Subscribe to the topic
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		utils.LogError("Error subscribing to topic", token.Error())
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Keep the subscriber running
	go func() {
		// Wait for a termination signal
		<-stopChan
		// Mark the task as done
		wg.Done()
	}()

	// Wait for all tasks in the WaitGroup to complete
	wg.Wait()
	close(stopChan)

	// Disconnect from broker
	client.Disconnect(250)
	utils.LogInfo("Subscriber disconnected and exiting")
}
