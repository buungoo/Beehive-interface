package test

import (
	//"beehive_api/models"

	"beehive_api/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// Inject harmful data into the database for testing

func InjectTestData(dbPool *pgxpool.Pool) error {

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	// Hash password with bcrypt
	password1, err := bcrypt.GenerateFromPassword([]byte("pass1"), bcrypt.DefaultCost)
	if err != nil {
		utils.LogFatal("Error hashing password:", err)
	}

	password2, err := bcrypt.GenerateFromPassword([]byte("pass2"), bcrypt.DefaultCost)
	if err != nil {
		utils.LogFatal("Error hashing password:", err)
	}

	// Insert test data into the users table
	var user1ID, user2ID int
	err = conn.QueryRow(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", "user1", password1).Scan(&user1ID)
	if err != nil {
		return fmt.Errorf("failed to insert user1: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", "user2", password2).Scan(&user2ID)
	if err != nil {
		return fmt.Errorf("failed to insert user2: %v", err)
	}

	// Insert test data into the beehives table
	var beehive1ID = 1
	var beehive2ID = 2
	_, err = conn.Exec(context.Background(), "INSERT INTO beehives (id, name, user_id) VALUES ($1, $2, $3)", beehive1ID, "Beehive A", user1ID)
	if err != nil {
		return fmt.Errorf("failed to insert beehive1: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO beehives (id, name, user_id) VALUES ($1, $2, $3)", beehive2ID, "Beehive B", user2ID)
	if err != nil {
		return fmt.Errorf("failed to insert beehive2: %v", err)
	}

	// Insert test data into the sensors table

	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 1, "temperature", beehive1ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor1: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 2, "humidity", beehive1ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor2: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 3, "weight", beehive1ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor3: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 4, "oxygen", beehive1ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor4: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 5, "temperature", beehive2ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor5: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 6, "humidity", beehive2ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor6: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 7, "weight", beehive2ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor7: %v", err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO sensors (id, type, beehive_id) VALUES ($1, $2, $3) ", 8, "oxygen", beehive2ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor8: %v", err)
	}

	// // Read JSON file
	jsonData, err := os.Open("/test/test_data.json")
	if err != nil {
		utils.LogFatal("Error opening json file, err: ", err)
	}
	defer jsonData.Close()

	// Read the content of the file into a byte slice
	data, err := io.ReadAll(jsonData)
	if err != nil {
		utils.LogFatal("Error reading json file, err: %v", err)
	}

	// Unmarshal JSON data
	//var readings []models.SensorData
	var readings []struct {
		SensorID   int     `json:"sensor_id"`
		BeehiveID  int     `json:"beehive_id"`
		SensorType string  `json:"sensor_type"`
		Value      float64 `json:"value"`
		Time       string  `json:"time"`
	}
	err = json.Unmarshal(data, &readings)
	if err != nil {
		utils.LogFatal("Error unmarshaling json data: %v", err)
	}

	// Parse each time string into a time.Time value and insert into the database
	layout := "2006-01-02 15:04:05"
	for _, reading := range readings {
		parsedTime, err := time.Parse(layout, reading.Time)
		if err != nil {
			utils.LogError("Error parsing time %s: %v", err)
			continue
		}

		_, err = conn.Exec(context.Background(),
			"INSERT INTO sensor_data (sensor_id, beehive_id, sensor_type, value, time) VALUES ($1, $2, $3, $4, $5)",
			reading.SensorID, reading.BeehiveID, reading.SensorType, reading.Value, parsedTime)
		if err != nil {
			utils.LogFatal("Failed to insert sensor data for sensor_id %d: %v", err)
		}
	}

	return nil
}
