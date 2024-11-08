package models

import (
	"time"
	//"errors"
	//"fmt"
)

// type Datetime time.Time

// func (d *Datetime) UnmarshalJSON(b []byte) error {
// 	if len(b) < 2 || b[0] != '"' || b[len(b)-1] != '"' {
// 		return errors.New("not a json string")
// 	}

// 	// 1. Strip the double quotes from the JSON string.
// 	b = b[1:len(b)-1]

// 	layout := "2006-01-02T15:04:05"

// 	// 2. Parse the result using our desired format.
// 	t, err := time.Parse(layout, string(b))
// 	if err != nil {
// 		return fmt.Errorf("failed to parse time: %w", err)
// 	}

// 	// finally, assign t to *d
// 	*d = Datetime(t)

// 	return nil
// }

type Beehives struct {
	Id     int    `json: "id"`
	Name   string `json: "name`
	UserID int    `json: "user_id`
}

type SensorData struct {
	SensorID   int       `json:"sensor_id"`
	BeehiveID  int       `json:"beehive_id"`
	SensorType string    `json:"sensor_type"`
	Value      float64   `json:"value"`
	Time       time.Time `json:"time"`
}

type SensorType string

const (
	SensorTypeTemperature SensorType = "temperature"
	SensorTypeHumidity    SensorType = "humidity"
	SensorTypeOxygen      SensorType = "oxygen"
	SensorTypeWeight      SensorType = "weight"
	SensorTypeMicrophone  SensorType = "microphone"
)

var validSensorTypes = map[SensorType]bool{
	SensorTypeTemperature: true,
	SensorTypeHumidity:    true,
	SensorTypeOxygen:      true,
	SensorTypeWeight:      true,
	SensorTypeMicrophone:  true,
}

func (st SensorType) IsValid() bool {
	return validSensorTypes[st]
}
