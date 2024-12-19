package mqtt

import (
	"github.com/buungoo/Beehive-interface/handlers"
	"github.com/buungoo/Beehive-interface/models"
	"github.com/buungoo/Beehive-interface/utils"
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
		fmt.Println("Error parsing sensor message", err)
		return
	}

	for _, reading := range readings {
		utils.LogInfo(fmt.Sprintf("Sensor Reading: %+v", reading))
		fmt.Printf("Sensor Reading: %+v\n", reading)
		// Insert the reading into the database
		err := handlers.InsertSensorReading(dbpool, reading)
		if err != nil {
			utils.LogError("Failed to insert sensor reading: ", err)
			fmt.Println("Failed to insert sensor reading: ", err)
		} else {
			utils.LogInfo(fmt.Sprintf("Successfully inserted reading into the database: %+v", reading))
			fmt.Println("Successfully inserted reading into the database:", reading)
		}
	}
}

func parseSensorMessage(message SensorMessage) ([]*models.SensorReading, error) {
	// Decode the base64 data into decimal
	decodedData, err := base64.StdEncoding.DecodeString(message.Data)
	if err != nil {
		utils.LogError("Error decoding base64 data", err)
		fmt.Println("Error decoding base64 data", err)
		return nil, fmt.Errorf("error decoding base64 data: %v", err)
	}

	// Format the time we received as it couldnt be parsed correctly otherwise
	if dotIndex := strings.Index(message.Time, "."); dotIndex != -1 {
		message.Time = message.Time[:dotIndex+4] + "Z" // We only need microseconds so truncate and add "Z" to indicate UTC
	}

	// Parse the time string into the correct type
	timeStamp, err := time.Parse(time.RFC3339, message.Time)
	if err != nil {
		utils.LogError("Error parsing time", err)
		fmt.Println("Error parsing time", err)
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

		sensorId := decodedData[i+1]
		rawValue := decodedData[i+2]

		builder := models.NewSensorReadingBuilder(sensorType, timeStamp).
			SetSensorID(sensorId).
			SetDevEUI(message.DevEUI).
			SetValue(rawValue)

		// Group each individual reading into a list
		readings = append(readings, builder.Build())

		utils.LogInfo(fmt.Sprintf("Processed reading: SensorType=%v, SensorID=%v, RawValue=%v", sensorType, sensorId, rawValue))
		fmt.Printf("Processed reading: SensorType=%v, SensorID=%v, RawValue=%v\n", sensorType, sensorId, rawValue) // print and log
	}

	return readings, nil
}

type SensorMessage struct {
	ApplicationName string `json:"applicationName"`
	Data            string `json:"data"`
	Time            string `json:"time"`
	DevEUI          string `json:"devEUI"`
}

// We have to use a closure function since the client should be able to use the dbpool
func createMessagePubHandler(dbpool *pgxpool.Pool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		var sensorMessage SensorMessage
		if err := json.Unmarshal(msg.Payload(), &sensorMessage); err != nil {
			utils.LogError("Error parsing JSON", err)
			fmt.Println("Error parsing JSON", err)
			return
		}

		utils.LogInfo(fmt.Sprintf("Received message: %+v", sensorMessage))
		fmt.Printf("Received message: %+v\n", sensorMessage) // print and log
		handleSensorMessage(sensorMessage, dbpool)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	utils.LogInfo("Connected to MQTT broker")
	fmt.Println("Connected to MQTT broker") // print and log
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	utils.LogError("Connection lost", err)
	fmt.Printf("Connection lost: %v\n", err)

	// Retry to connect if connection was lost
	for {
		utils.LogInfo("Attempting to reconnect...")
		fmt.Println("Attempting to reconnect...")
		if token := client.Connect(); token.Wait() && token.Error() == nil {
			utils.LogInfo("Reconnected successfully")
			fmt.Println("Reconnected successfully")
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func retryConnect(client mqtt.Client) error {
	for {
		// Attempt to connect to the broker
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			utils.LogError("Failed to connect to MQTT broker", token.Error())
			fmt.Println("Failed to connect to MQTT broker", token.Error())
			// Wait before retrying
			utils.LogInfo("Retrying to connect to MQTT broker...")
			fmt.Println("Retrying to connect to MQTT broker...")
			time.Sleep(5 * time.Second) // Sleep for 5 seconds before retrying
			continue                    // Retry connection
		}

		// If connection was successful, break out of the loop
		break
	}
	return nil
}

func SetupMQTTSubscriber(dbpool *pgxpool.Pool) {
	broker := "broker.hivemq.com:1883"
	topic := "d0039ebeehive/sensor"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("fkjsdhflafhlhgds") //beehive_subscriber")
	opts.SetDefaultPublishHandler(createMessagePubHandler(dbpool))

	opts.OnConnect = func(client mqtt.Client) {
		utils.LogInfo("Connected to MQTT broker")
		fmt.Println("Connected to MQTT broker") // print and log
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		utils.LogError("Connection to MQTT broker lost", err)
		fmt.Println("Connection to MQTT broker lost", err) // print and log
	}

	client := mqtt.NewClient(opts)

	// Attempt to connect to the broker
	for {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			utils.LogError("Failed to connect to MQTT broker", token.Error())
			fmt.Println("Failed to connect to MQTT broker", token.Error())
			// Wait before retrying
			utils.LogInfo("Retrying to connect to MQTT broker...")
			fmt.Println("Retrying to connect to MQTT broker...")
			time.Sleep(5 * time.Second)
			continue
		}

		// If connection was successful, break out of the loop
		break
	}

	// Subscribe to the topic
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		utils.LogError("Failed to subscribe to MQTT topic", token.Error())
		fmt.Println("Failed to subscribe to MQTT topic", token.Error())
		return
	}

	utils.LogInfo("Subscribed to MQTT topic")
	fmt.Println("Subscribed to MQTT topic")

	// Use the main application's stop channel to block
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	client.Disconnect(250)
	utils.LogInfo("MQTT subscriber disconnected")
	fmt.Println("MQTT subscriber disconnected")
}
