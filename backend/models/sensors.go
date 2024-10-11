package models

type SensorData struct {
	BeehiveID 		int 	`json: "beehive_id"`
	SensorID		int 	`json:"sensor_id"`
	Value			float64 `json:"value"`
	Time			string 	`json:"time"`
}