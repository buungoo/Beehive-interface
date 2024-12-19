// Package models contains models and methods that is used to handle sensordata to and from the database.
package models

import (
	"github.com/buungoo/Beehive-interface/utils"

	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Sensor string

const (
	LoadCell    Sensor = "loadcell"
	Temperature Sensor = "temperature"
	Humidity    Sensor = "humidity"
	Microphone  Sensor = "microphone"
	Oxygen      Sensor = "oxygen"
	Battery     Sensor = "battery"
)

func (s Sensor) IsValid() bool {
	switch s {
	case LoadCell, Temperature, Humidity, Microphone, Oxygen, Battery:
		return true
	default:
		return false
	}
}

// SensorReading struct to hold a parsed sensor reading
type SensorReading struct {
	SensorType Sensor           `json:"sensor_type"`
	SensorID   int              `json:"sensor_id"`  // TODO: Dirty solution, we receive uint8 id but databse is taking int
	Value      interface{}      `json:"value"`      // Allows for the different types we need, i.e., uint8, int8, bool
	Time       time.Time        `json:"time"`       // To store the timestamp of the reading
	BeehiveID  net.HardwareAddr `json:"beehive_id"` // MAC address of the parent Beehive, I could have used a string but this makes sure it is parsed as a valid macaddr
	// SensorID   int       `json:"sensor_id"`
	// BeehiveID  int       `json:"beehive_id"`
	// SensorType string    `json:"sensor_type"`
	// Value      float32   `json:"value"`
	// Time       time.Time `json:"time"`
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
	return &SensorReadingBuilder{
		sensorType: sensorType,
		timestamp:  timestamp,
	}
}

func (b *SensorReadingBuilder) SetSensorID(id uint8) *SensorReadingBuilder {
	b.sensorID = id
	return b
}

func loadValuesCommaSeparated(filename string) ([]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		for _, part := range parts {
			v, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid number in file: %s", err)
			}
			values = append(values, v)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func interpolateYToX(xs, ys []float64, yTarget float64) (float64, error) {
	// yMin := ys[0]
	yMax := ys[len(ys)-1]

	if yTarget == yMax {
		return xs[len(xs)-1], nil
	}

	i := findInterval(ys, yTarget)
	if i == -1 {
		return 0, fmt.Errorf("target out of range")
	}

	y0 := ys[i]
	y1 := ys[i+1]
	x0 := xs[i]
	x1 := xs[i+1]

	t := (yTarget - y0) / (y1 - y0)
	return x0 + t*(x1-x0), nil
}

func findInterval(values []float64, target float64) int {
	for i := 0; i < len(values)-1; i++ {
		if target >= values[i] && target <= values[i+1] {
			return i
		}
	}
	return -1
}

func (b *SensorReadingBuilder) SetValue(value interface{}) *SensorReadingBuilder {
	switch b.sensorType {
	case Temperature:
		switch v := value.(type) {
		case int8:
			b.value = v
		default:
			fmt.Println("Invalid value type for Temperature. Expected int8.")
		}
	case Battery:
		switch v := value.(type) {
		case uint16:
			// Load LUT values from files
			xs, err := loadValuesCommaSeparated("models/lut_x.txt")
			if err != nil {
				fmt.Printf("Error loading LUT X values: %s\n", err)
				return b
			}
			ys, err := loadValuesCommaSeparated("models/lut_y.txt")
			if err != nil {
				fmt.Printf("Error loading LUT Y values: %s\n", err)
				return b
			}

			if len(xs) != len(ys) {
				fmt.Println("X and Y arrays must have the same length.")
				return b
			}

			// Interpolate value
			interpolatedValue, err := interpolateYToX(xs, ys, float64(v))
			if err != nil {
				fmt.Printf("Error interpolating value: %s\n", err)
				return b
			}

			b.value = interpolatedValue // Set the resulting float64 value
		default:
			fmt.Println("Invalid value type for Battery. Expected uint16.")
		}
	default:
		switch v := value.(type) {
		case uint8:
			b.value = v
		default:
			fmt.Println("Invalid value type for sensor. Expected uint8.")
		}
	}
	return b
}

func (b *SensorReadingBuilder) SetDevEUI(parentBeehive string) *SensorReadingBuilder {
	// Make sure the length of the macaddr is valid
	if len(parentBeehive) != 16 {
		utils.LogWarn("Invalid DevEUI length. Expected 16 characters.")
		return b
	}

	// Format the macaddr as it is received without colons
	macFormatted := strings.ToLower(parentBeehive[:2] + ":" +
		parentBeehive[2:4] + ":" +
		parentBeehive[4:6] + ":" +
		parentBeehive[6:8] + ":" +
		parentBeehive[8:10] + ":" +
		parentBeehive[10:12] + ":" +
		parentBeehive[12:14] + ":" +
		parentBeehive[14:16])

	// Parse the string into macaddr8 type
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
		SensorType: b.sensorType,
		SensorID:   int(b.sensorID), // TODO: Dirty solution, we shouldn't have to parse it
		Value:      b.value,
		Time:       b.timestamp,
		BeehiveID:  b.parentBeehive,
	}
}

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

// Season is used to implement seasons and map data.
type Season struct {
	Name      string
	LowTemp   int8
	HighTemp  int8
	LowHumid  uint8
	HighHumid uint8
}

// Limits for temperature (int8)
const (
	WinterLowTemp  int8 = -40
	WinterHighTemp int8 = 40
	SpringLowTemp  int8 = -30
	SpringHighTemp int8 = 40
	SummerLowTemp  int8 = 0
	SummerHighTemp int8 = 40
	FallLowTemp    int8 = -30
	FallHighTemp   int8 = 30
)

// Limits for humidity (uint8)
const (
	WinterLowHumidity  uint8 = 5
	WinterHighHumidity uint8 = 50
	SpringLowHumidity  uint8 = 5
	SpringHighHumidity uint8 = 60
	SummerLowHumidity  uint8 = 10
	SummerHighHumidity uint8 = 70
	FallLowHumidity    uint8 = 5
	FallHighHumidity   uint8 = 60
)

// Limits for oxygen (uint8)
const (
	LowOxygen  uint8 = 18
	HighOxygen uint8 = 25
)

// Limits for weight (uint8)
const (
	LowWeight  uint8 = 0
	HighWeight uint8 = 100
)

// Limits for microphone (bool)
const (
	// LowMicNoise bool = false
	MicLow  uint8 = 0
	MicHigh uint8 = 1
)

const (
	LowBattery  uint8 = 0
	HighBattery uint8 = 100
)

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

// VerifyInputData verifies sensorvalue and returns true if everything looks good, else return false with message
func (reading SensorReading) VerifyInputData() (bool, string) {
	switch reading.SensorType {
	case Temperature:
		return reading.verifyTemperature()
	case Humidity:
		return reading.verifyHumidity()
	case Oxygen:
		return reading.verifyOxygen()
	case LoadCell:
		return reading.verifyWeight()
	case Microphone:
		return reading.verifyMicrophone()
	case Battery:
		return reading.VerifyBattery()
	default:
		return false, "semthing went wrong while verifying error"
	}
}

// Verify temperature sensorvalues
func (reading SensorReading) verifyTemperature() (bool, string) {
	month := reading.Time.Month()
	season, exists := seasons[month]
	if !exists {
		return false, "Error finding month"
	}

	// Type assertion for int8
	temp, ok := reading.Value.(int8)
	if !ok {
		return false, "Invalid value type for temperature. Expected int8."
	}

	if temp < season.LowTemp {
		return false, "temperature is below " + fmt.Sprintf("%f", season.LowTemp) + " Celsius"
	} else if temp > season.HighTemp {
		return false, "temperature is above " + fmt.Sprintf("%f", season.HighTemp) + " Celsius"
	}
	return true, "temperature is within limits"
}

// Verify humidity sensorvalues
func (reading SensorReading) verifyHumidity() (bool, string) {
	month := reading.Time.Month()
	season, exists := seasons[month]
	if !exists {
		return false, "Error finding month"
	}

	// Type assertion for uint8
	humid, ok := reading.Value.(uint8)
	if !ok {
		return false, "Invalid value type for humidity. Expected uint8."
	}

	if humid < season.LowHumid {
		return false, "humidity is below " + fmt.Sprintf("%f", season.LowHumid) + "%"
	} else if humid > season.HighHumid {
		return false, "humidity is above " + fmt.Sprintf("%f", season.HighHumid) + "%"
	}
	return true, "humidity levels are within limits"
}

// Verify oxygen sensorvalues
func (reading SensorReading) verifyOxygen() (bool, string) {
	// Type assertion for uint8
	oxygen, ok := reading.Value.(uint8)
	if !ok {
		return false, "Invalid value type for oxygen. Expected uint8."
	}

	if oxygen < LowOxygen {
		return false, "oxygen level is below " + fmt.Sprintf("%f", LowOxygen) + "%"
	} else if oxygen > HighOxygen {
		return false, "oxygen level is above " + fmt.Sprintf("%f", HighOxygen) + "%"
	}
	return true, "oxygen levels are within limits"
}

// Verify weight sensorvalues
func (reading SensorReading) verifyWeight() (bool, string) {
	// Type assertion for uint8
	weight, ok := reading.Value.(uint8)
	if !ok {
		return false, "Invalid value type for weight. Expected uint8."
	}

	if weight < LowWeight {
		return false, "weight is below " + fmt.Sprintf("%f", LowWeight) + "kg"
	} else if weight > HighWeight {
		return false, "weight is above " + fmt.Sprintf("%f", HighWeight) + "kg"
	}
	return true, "Weight is within limit"
}

// Verify microphone sensorvalues
// func (reading SensorReading) verifyMicrophone() (bool, string) {
// 	// Type assertion for bool
// 	mic, ok := reading.Value.(bool)
// 	if !ok {
// 		return false, "Invalid value type for microphone. Expected bool."
// 	}
//
// 	if mic {
// 		return false, "The bees are angy"
// 	}
// 	return true, "The bees are happi"
// }

// Verify microphone sensor values
func (reading SensorReading) verifyMicrophone() (bool, string) {
	// Type assertion for int
	mic, ok := reading.Value.(uint8)
	if !ok {
		return false, "Invalid value type for microphone. Expected int."
	}

	switch mic {
	case 0:
		return true, "The bees are happi"
	case 1:
		return false, "The bees are angy"
	default:
		return false, "Invalid microphone value. Expected 0  or 1."
	}
}

// Verify battery level
func (reading SensorReading) VerifyBattery() (bool, string) {
	// Type assertion for uint8
	battery, ok := reading.Value.(uint8)
	if !ok {
		return false, "Invalid value type for battery. Expected uint8."
	}

	if battery < LowBattery || battery > HighBattery {
		return false, fmt.Sprintf("battery level is out of range: %d%%. Expected between 0-100%%", battery)
	}
	return true, fmt.Sprintf("battery level is within range: %d%%", battery)
}
