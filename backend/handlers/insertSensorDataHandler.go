package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsertSensorReading inserts a sensor reading into the `sensor_data` table.
func InsertSensorReading(dbpool *pgxpool.Pool, reading *models.SensorReading) error {
	// Acquire a database connection from the pool
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed to acquire a database connection: %v", err)
	}
	defer conn.Release()

	// Fetch the beehive_id using ParentBeehive (MAC address)
	var beehiveID int
	queryBeehive := `SELECT id FROM beehives WHERE key = $1`
	err = conn.QueryRow(context.Background(), queryBeehive, reading.ParentBeehive.String()).Scan(&beehiveID) // Convert the macaddr back to string for the query
	if err != nil {
		return fmt.Errorf("failed to find beehive with MAC address %s: %v", reading.ParentBeehive.String(), err)
	}

	// Check if the sensor exists in the `sensors` table
	querySensor := `
		SELECT COUNT(*) FROM sensors 
		WHERE id = $1 AND type = $2 AND beehive_id = $3`
	var sensorExists int
	err = conn.QueryRow(context.Background(), querySensor, reading.SensorID, string(reading.SensorType), beehiveID).Scan(&sensorExists)
	if err != nil {
		return fmt.Errorf("error checking sensor existence: %v", err)
	}

	// If the sensor does not exist, add it to the sensors table
	if sensorExists == 0 {
		// Insert the sensor into the `sensors` table
		insertSensorQuery := `
			INSERT INTO sensors (id, type, beehive_id) 
			VALUES ($1, $2, $3)`
		_, err = conn.Exec(context.Background(), insertSensorQuery,
			reading.SensorID, string(reading.SensorType), beehiveID)
		if err != nil {
			return fmt.Errorf("failed to insert sensor: %v", err)
		}
		utils.LogInfo(fmt.Sprintf("Added new sensor: ID=%d, Type=%s, BeehiveID=%d", reading.SensorID, reading.SensorType, beehiveID))
	}

	// Insert the sensor reading into the `sensor_data` table
	insertReadingQuery := `
		INSERT INTO sensor_data (sensor_id, beehive_id, sensor_type, value, time) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err = conn.Exec(context.Background(), insertReadingQuery,
		reading.SensorID, beehiveID, string(reading.SensorType), reading.Value, reading.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert sensor reading: %v", err)
	}

	return nil
}
