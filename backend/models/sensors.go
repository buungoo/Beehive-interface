package models

import (
	"fmt"
	"time"
	//"errors"
	//"fmt"
)

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

type Beehives struct {
	Id   int    `json: "id"`
	Name string `json: "name`
}

type Season struct {
	Name      string
	LowTemp   float64
	HighTemp  float64
	LowHumid  float64
	HighHumid float64
}

// Limits for temperature
const (
	WinterLowTemp  float64 = -40.0
	WinterHighTemp float64 = 40.0
	SpringLowTemp  float64 = -30.0
	SpringHighTemp float64 = 40.0
	SummerLowTemp  float64 = 0.0
	SummerHighTemp float64 = 40.0
	FallLowTemp    float64 = -30.0
	FallHighTemp   float64 = 30.0
)

// Limits for humidity
const (
	WinterLowHumidity  float64 = 5.0
	WinterHighHumidity float64 = 50.0
	SpringLowHumidity  float64 = 5.0
	SpringHighHumidity float64 = 60.0
	SummerLowHumidity  float64 = 10.0
	SummerHighHumidity float64 = 70.0
	FallLowHumidity    float64 = 5.0
	FallHighHumidity   float64 = 60.0
)

// Limits for oxygen
const (
	LowOxygen  float64 = 18.0
	HighOxygen float64 = 25.0
)

// Limits for weight
const (
	LowWeight  float64 = 0.0
	HighWeight float64 = 10.0
)

// Limits for microphone
const LowMicNoise float64 = 0.0

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

var seasons = map[time.Month]*Season{
	time.January: winter, time.February: winter,
	time.March: spring, time.April: spring, time.May: spring,
	time.June: summer, time.July: summer, time.August: summer,
	time.September: fall, time.October: fall, time.November: fall,
	time.December: winter,
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

func (st SensorType) String() string {
	return string(st)
}

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

func (sd SensorData) verifyOxygen() (bool, string) {
	if sd.Value < LowOxygen {
		return false, "humiditylevel is below " + fmt.Sprintf("%f", LowOxygen) + "%"
	} else if sd.Value > HighOxygen {
		return false, "humiditylevel is above " + fmt.Sprintf("%f", HighOxygen) + "%"
	} else {
		return true, "oxygenlevels are within limits"
	}
}

func (sd SensorData) verifyWeight() (bool, string) {
	if sd.Value < LowWeight {
		return false, "weight is below " + fmt.Sprintf("%f", LowWeight) + "kg"
	} else if sd.Value > HighWeight {
		return false, "weight is above " + fmt.Sprintf("%f", LowOxygen) + "kg"
	}
	return true, "weight is within limits"
}

func (sd SensorData) verifyMicrophone() (bool, string) {
	if sd.Value > LowMicNoise {
		return false, "microphone is detecting noise"
	}
	return true, "microphone is not detecting any noise"
}
