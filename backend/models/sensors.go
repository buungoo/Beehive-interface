package models

type SensorData struct {
	SensorID		string 	`json:"sensor_id"`
	BeehiveID		int 	`json:"beehive_id"`
	Value			float64 `json:"value"`
	Time			string 	`json:"time"`
}