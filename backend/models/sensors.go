package models

import "time"

type SensorData struct {
	SensorID		int 		`json:"sensor_id"`
	BeehiveID 		int 		`json:"beehive_id"`
	Value			float64 	`json:"value"`
	Time			time.Time 	`json:"time"`
}

type SensorType string

const (
	SensorTypeTemperature SensorType = "temperature"
	SensorTypeHumidity SensorType = "humidity"
	SensorTypeOxygen SensorType = "oxygen"
	SensorTypeWeight SensorType = "weight"
)

var validSensorTypes = map[SensorType]bool{
	SensorTypeTemperature: true,
	SensorTypeHumidity: true,
	SensorTypeOxygen: true,
	SensorTypeWeight: true,
}

func(st SensorType) IsValid() bool {
	return validSensorTypes[st]
}