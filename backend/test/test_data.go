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
	var beehive1ID, beehive2ID int
	err = conn.QueryRow(context.Background(), "INSERT INTO beehives (name, user_id) VALUES ($1, $2) RETURNING id", "Beehive A", user1ID).Scan(&beehive1ID)
	if err != nil {
		return fmt.Errorf("failed to insert beehive1: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO beehives (name, user_id) VALUES ($1, $2) RETURNING id", "Beehive B", user2ID).Scan(&beehive2ID)
	if err != nil {
		return fmt.Errorf("failed to insert beehive2: %v", err)
	}

	// Insert test data into the sensors table
	var sensor1ID, sensor2ID, sensor3ID, sensor4ID, sensor5ID, sensor6ID, sensor7ID, sensor8ID int
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "temperature", beehive1ID).Scan(&sensor1ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor1: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "humidity", beehive1ID).Scan(&sensor2ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor2: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "weight", beehive1ID).Scan(&sensor3ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor3: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "oxygen", beehive1ID).Scan(&sensor4ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor4: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "temperature", beehive2ID).Scan(&sensor5ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor5: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "humidity", beehive2ID).Scan(&sensor6ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor6: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "weight", beehive2ID).Scan(&sensor7ID)
	if err != nil {
		return fmt.Errorf("failed to insert sensor7: %v", err)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "oxygen", beehive2ID).Scan(&sensor8ID)
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
	var readings []struct {
		SensorID  int     `json:"sensor_id"`
		BeehiveID int     `json:"beehive_id"`
		Value     float64 `json:"value"`
		Time      string  `json:"time"`
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
			"INSERT INTO sensor_data (sensor_id, beehive_id, value, time) VALUES ($1, $2, $3, $4)",
			reading.SensorID, reading.BeehiveID, reading.Value, parsedTime)
		if err != nil {
			utils.LogFatal("Failed to insert sensor data for sensor_id %d: %v", err)
		}
	}

	return nil
}
