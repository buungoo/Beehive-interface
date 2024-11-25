package mqtt

import (
	"beehive_api/handlers"
	"beehive_api/models"
	"beehive_api/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "log"
	"os"
	"os/signal"
	"strings"
	// "sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jackc/pgx/v5/pgxpool"
)

func handleSensorMessage(message SensorMessage, dbpool *pgxpool.Pool) {
	readings, err := parseSensorMessage(message)
	if err != nil {
		utils.LogError("Error parsing sensor message", err)
		return
	}

	for _, reading := range readings {
		utils.LogInfo(fmt.Sprintf("Sensor Reading: %+v", reading))
		// Insert the reading into the database
		err := handlers.InsertSensorReading(dbpool, reading)
		if err != nil {
			utils.LogError("Failed to insert sensor reading: ", err)
		}
	}
}

func parseSensorMessage(message SensorMessage) ([]*models.SensorReading, error) {
	// Decode the base64 data into decimal
	decodedData, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 data: %v", err)
	}

	// Format the time we received as it couldnt be parsed correctly otherwise
	if dotIndex := strings.Index(message.Time, "."); dotIndex != -1 {
		message.Time = message.Time[:dotIndex] + "Z" // Add "Z" to indicate UTC
	}

	timeStamp, err := time.Parse("2006-01-02T15:04:05Z", message.Time)
	if err != nil {
		return nil, fmt.Errorf("error parsing time: %v", err)
	}

	// Each sensor contains 6 bytes, 2 for sensor type, 2 for id, 2 for value
	var readings []*models.SensorReading
	for i := 0; i < len(decodedData); i += 3 {
		var sensorType models.Sensor
		switch decodedData[i] {
		case 1:
			sensorType = models.LoadCell
		case 2:
			sensorType = models.Temperature
		case 3:
			sensorType = models.Humidity
		case 4:
			sensorType = models.Microphone
		case 5:
			sensorType = models.Oxygen
		default:
			sensorType = "Unknown"
		}

		// (decodedData[i])
		sensorId := decodedData[i+1]
		rawValue := decodedData[i+2]

		builder := models.NewSensorReadingBuilder(sensorType, timeStamp).
			SetSensorID(sensorId).
			SetDevEUI(message.DevEUI).
			SetValue(rawValue)

		// Group each individual reading into a list
		readings = append(readings, builder.Build())
	}

	return readings, nil
}

type SensorMessage struct {
	ApplicationName string `json:"applicationName"`
	Data            string `json:"data"`
	Time            string `json:"time"`
	DevEUI          string `json:"devEUI"`
}

// Initialize the log file
// var logFile *os.File

// var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	// Get the current timestamp
// 	// not needed if we use the new logging function
// 	// timestamp := time.Now().Format("2006-01-02 15:04:05")
//
// 	// Parse the received message payload
// 	var sensorMessage SensorMessage
// 	if err := json.Unmarshal(msg.Payload(), &sensorMessage); err != nil {
// 		log.Printf("Error parsing JSON: %v", err)
// 		return
// 	}
//
// 	// Parse the message as it is received like this:
// 	// {
// 	// "applicationID":"1",
// 	// "applicationName":"beehive-sensor-card",
// 	// "data":"YWFhYWE=",
// 	// "devEUI":"0080e115000adf82",
// 	// "deviceName":"beehive-sensor-card-dn",
// 	// "fCnt":319,
// 	// "fPort":2,
// 	// "rxInfo":[{"altitude":0,"latitude":0,"loRaSNR":7.5,"longitude":0,"mac":"24e124fffef0b4f9","name":"24e124fffef0b4f9","rssi":-109,"time":"2024-11-05T14:01:49.217376Z"}],
// 	// "time":"2024-11-05T14:01:49.217376Z",
// 	// "txInfo":{"adr":true,"codeRate":"4/5","dataRate":{"bandwidth":125,"modulation":"LORA","spreadFactor":7},"frequency":868300000}
// 	// }
//
// 	// utils.LogInfo(fmt.Sprintf("Received message - applicationName: %s, data: %s, time: %s",
// 	// 	sensorMessage.ApplicationName, sensorMessage.Data, sensorMessage.Time))
// 	fmt.Println("Received message:", sensorMessage)
// 	utils.LogInfo(fmt.Sprintf("Received message: %+v", sensorMessage))
//
// 	handleSensorMessage(sensorMessage) //.Data)
// }

// We have to use a closure function since the client should be able to use the dbpool
func createMessagePubHandler(dbpool *pgxpool.Pool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var sensorMessage SensorMessage
		if err := json.Unmarshal(msg.Payload(), &sensorMessage); err != nil {
			utils.LogError("Error parsing JSON", err)
			return
		}

		utils.LogInfo(fmt.Sprintf("Received message: %+v", sensorMessage))
		handleSensorMessage(sensorMessage, dbpool)
	}
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

// func SetupMQTTSubscriber(dbpool *pgxpool.Pool) {
// 	broker := "broker.hivemq.com:1883"
// 	topic := "d0039ebeehive/sensor"
//
// 	opts := mqtt.NewClientOptions()
// 	opts.AddBroker(broker)
// 	opts.SetClientID("beehive_subscriber")
// 	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
// 		// Parse the received message payload
// 		var sensorMessage SensorMessage
// 		if err := json.Unmarshal(msg.Payload(), &sensorMessage); err != nil {
// 			utils.LogError("Error parsing JSON", err)
// 			return
// 		}
//
// 		handleSensorMessage(sensorMessage, dbpool)
// 	})
// 	opts.OnConnect = func(client mqtt.Client) {
// 		utils.LogInfo("Connected to MQTT broker")
// 	}
// 	opts.OnConnectionLost = func(client mqtt.Client, err error) {
// 		utils.LogError("Connection lost", err)
// 	}
//
// 	client := mqtt.NewClient(opts)
// 	if token := client.Connect(); token.Wait() && token.Error() != nil {
// 		utils.LogError("Error connecting to broker", token.Error())
// 		return
// 	}
//
// 	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
// 		utils.LogError("Error subscribing to topic", token.Error())
// 		return
// 	}
//
// 	utils.LogInfo("Subscribed to MQTT topic")
//
// 	// Block this goroutine until termination
// 	stopChan := make(chan os.Signal, 1)
// 	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
// 	<-stopChan
//
// 	client.Disconnect(250)
// 	utils.LogInfo("MQTT subscriber disconnected")
// }

func SetupMQTTSubscriber(dbpool *pgxpool.Pool) {
	broker := "broker.hivemq.com:1883"
	topic := "d0039ebeehive/sensor"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("beehive_subscriber")
	opts.SetDefaultPublishHandler(createMessagePubHandler(dbpool))

	opts.OnConnect = func(client mqtt.Client) {
		utils.LogInfo("Connected to MQTT broker")
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		utils.LogError("Connection lost", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		utils.LogError("Error connecting to broker", token.Error())
		return
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		utils.LogError("Error subscribing to topic", token.Error())
		return
	}

	utils.LogInfo("Subscribed to MQTT topic")

	// Block this goroutine until termination
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	client.Disconnect(250)
	utils.LogInfo("MQTT subscriber disconnected")
}

