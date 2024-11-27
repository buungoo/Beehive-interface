// Package models contains models and methods that is used to handle sensordata to and from the database.
package models

import (
	"fmt"
	"time"
	//"errors"
	//"fmt"
)

// BeehiveStatus is a struct that is used for the beehive_status table.
type BeehiveStatus struct {
	IssueId     int        `json: "issue_id`
	SensorId    int        `json: "sensor_id"`
	BeehiveId   int        `json: "beehive_id"`
	SensorType  string     `json: "sensor_type"`
	Description string     `json: "description"`
	Solved      bool       `json: "solved"`
	Read        bool       `json: "read"`
	TimeOfError *time.Time `json: "time_of_error, omitempty"`
	TimeRead    *time.Time `json: "time_read, omitempty"`
}

// Beehives is never used.
type Beehives struct {
	Id   int    `json: "id"`
	Name string `json: "name`
}

// Season is used to implement seasons and map data.
type Season struct {
	Name      string
	LowTemp   float32
	HighTemp  float32
	LowHumid  float32
	HighHumid float32
}

// Limits for temperature
const (
	WinterLowTemp  float32 = -40.0
	WinterHighTemp float32 = 40.0
	SpringLowTemp  float32 = -30.0
	SpringHighTemp float32 = 40.0
	SummerLowTemp  float32 = 0.0
	SummerHighTemp float32 = 40.0
	FallLowTemp    float32 = -30.0
	FallHighTemp   float32 = 30.0
)

// Limits for humidity
const (
	WinterLowHumidity  float32 = 5.0
	WinterHighHumidity float32 = 50.0
	SpringLowHumidity  float32 = 5.0
	SpringHighHumidity float32 = 60.0
	SummerLowHumidity  float32 = 10.0
	SummerHighHumidity float32 = 70.0
	FallLowHumidity    float32 = 5.0
	FallHighHumidity   float32 = 60.0
)

// Limits for oxygen
const (
	LowOxygen  float32 = 18.0
	HighOxygen float32 = 25.0
)

// Limits for weight
const (
	LowWeight  float32 = 0.0
	HighWeight float32 = 10.0
)

// Limits for microphone
const LowMicNoise float32 = 0.0

// Pointers to each season
var (
	winter = &Season{Name: "winter", LowTemp: WinterLowTemp, HighTemp: WinterHighTemp,
		LowHumid: WinterLowHumidity, HighHumid: WinterHighHumidity}
	summer = &Season{Name: "summer", LowTemp: SummerLowTemp, HighTemp: SummerHighTemp,
		LowHumid: SummerLowHumidity, HighHumid: SummerHighHumidity}
	fall = &Season{Name: "fall", LowTemp: FallLowTemp, HighTemp: FallHighTemp,
		LowHumid: FallLowHumidity, HighHumid: FallHighHumidity}
	spring = &Season{Name: "spring", LowTemp: SpringLowTemp, HighTemp: SpringHighTemp,
		LowHumid: SpringLowHumidity, HighHumid: SpringHighHumidity}
)

// seasons maps each month to a season with a month as key a Season pointer as value
var seasons = map[time.Month]*Season{
	time.January: winter, time.February: winter,
	time.March: spring, time.April: spring, time.May: spring,
	time.June: summer, time.July: summer, time.August: summer,
	time.September: fall, time.October: fall, time.November: fall,
	time.December: winter,
}

// SensorData is u
type SensorData struct {
	SensorID   int       `json:"sensor_id"`
	BeehiveID  int       `json:"beehive_id"`
	SensorType string    `json:"sensor_type"`
	Value      float32   `json:"value"`
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

// IsValid returns true if the sensortype exists.
func (st SensorType) IsValid() bool {
	return validSensorTypes[st]
}

// String convert SensorType to string.
func (st SensorType) String() string {
	return string(st)
}

// VerifyInputData verifies sensorvalue and returns true if everything looks good, else return false with message
func (sd SensorData) VerifyInputData() (bool, string) {
	switch sd.SensorType {
	case SensorTypeTemperature.String():
		return sd.verifyTemperature()
	case SensorTypeHumidity.String():
		return sd.verifyHumidity()
	case SensorTypeOxygen.String():
		return sd.verifyOxygen()
	case SensorTypeWeight.String():
		return sd.verifyWeight()
	case SensorTypeMicrophone.String():
		return sd.verifyMicrophone()
	default:
		return false, "semthing went wrong while verifying error"
	}
}

// Verify temperature sensorvalues
func (sd SensorData) verifyTemperature() (bool, string) {
	month := sd.Time.Month()
	season, exists := seasons[month]
	if !exists {
		return false, "Error finding month"
	}
	if sd.Value < season.LowTemp {
		return false, "temperature is below " + fmt.Sprintf("%f", season.LowTemp) + " Celsius"

	} else if sd.Value > season.HighTemp {
		return false, "temperature is above " + fmt.Sprintf("%f", season.HighTemp) + " Celsius"
	} else {
		return true, "something went wrong while checking temperature"
	}
}

// Verify humidity sensorvalues
func (sd SensorData) verifyHumidity() (bool, string) {
	month := sd.Time.Month()
	season, exists := seasons[month]
	if !exists {
		return false, "Error finding month"
	}
	if sd.Value < season.LowHumid {
		return false, "humidity is below " + fmt.Sprintf("%f", season.LowHumid) + " Celsius"

	} else if sd.Value > season.HighHumid {
		return false, "temperature is above " + fmt.Sprintf("%f", season.HighHumid) + " Celsius"
	} else {
		return true, "humiditylevels are within limits"
	}

}

// Verify oxygen sensorvalues
func (sd SensorData) verifyOxygen() (bool, string) {
	if sd.Value < LowOxygen {
		return false, "humiditylevel is below " + fmt.Sprintf("%f", LowOxygen) + "%"
	} else if sd.Value > HighOxygen {
		return false, "humiditylevel is above " + fmt.Sprintf("%f", HighOxygen) + "%"
	} else {
		return true, "oxygenlevels are within limits"
	}
}

// Verify weight sensorvalues
func (sd SensorData) verifyWeight() (bool, string) {
	if sd.Value < LowWeight {
		return false, "weight is below " + fmt.Sprintf("%f", LowWeight) + "kg"
	} else if sd.Value > HighWeight {
		return false, "weight is above " + fmt.Sprintf("%f", LowOxygen) + "kg"
	}
	return true, "weight is within limits"
}

// Verify microphone sensorvalues
func (sd SensorData) verifyMicrophone() (bool, string) {
	if sd.Value > LowMicNoise {
		return false, "microphone is detecting noise"
	}
	return true, "microphone is not detecting any noise"
}
