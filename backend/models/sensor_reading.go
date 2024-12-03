package models

import (
	"beehive_api/utils"
	"net"
	"strings"
	"time"
)

// Sensor enum to represent sensor types
type Sensor string


const (
	LoadCell	Sensor = "loadcell"
	Temperature Sensor = "temperature"
	Humidity    Sensor = "humidity"
	Microphone  Sensor = "microphone"
	Oxygen      Sensor = "oxygen"
	Battery		Sensor = "battery"
)

// const (
// 	LoadCell    Sensor = 1
// 	Temperature        = 2
// 	Humidity           = 3
// 	Microphone         = 4
// 	Oxygen             = 5
// )

// SensorReading struct to hold a parsed sensor reading
type SensorReading struct {
	SensorType    Sensor
	SensorID      uint8
	Value         interface{}      // Allows for the different types we need, i.e., uint8, int8, bool
	Timestamp     time.Time        // To store the timestamp of the reading
	ParentBeehive net.HardwareAddr // MAC address of the parent Beehive, I could have used a string but this makes sure it is parsed as a valid macaddr
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
	return &SensorReadingBuilder{sensorType: sensorType, timestamp: timestamp}
}

func (b *SensorReadingBuilder) SetSensorID(id uint8) *SensorReadingBuilder {
	b.sensorID = id
	return b
}

func (b *SensorReadingBuilder) SetValue(value interface{}) *SensorReadingBuilder {
	switch b.sensorType {
	case Temperature:
		if v, ok := value.(int8); ok {
			b.value = v
		} else {
			utils.LogWarn("Invalid value type for Temperature. Expected int8.")
		}
	case Microphone:
		if v, ok := value.(uint8); ok {
			b.value = v == 1 // Microphone can be either 0 or 1 as we want to map it to boolean.
		} else {
			utils.LogWarn("Invalid value type for Microphone. Expected uint8.")
		}
	default:
		if v, ok := value.(uint8); ok {
			b.value = v
		} else {
			utils.LogWarn("Invalid value type for sensor. Expected uint8.")
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
		SensorType:    b.sensorType,
		SensorID:      b.sensorID,
		Value:         b.value,
		Timestamp:     b.timestamp,
		ParentBeehive: b.parentBeehive,
	}
}
