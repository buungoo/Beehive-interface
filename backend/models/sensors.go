package models

import "time"

type SensorData struct {
	SensorID		int 		`json:"sensor_id"`
	BeehiveID 		int 		`json:"beehive_id"`
	Value			float64 	`json:"value"`
	Time			time.Time 	`json:"time"`
}