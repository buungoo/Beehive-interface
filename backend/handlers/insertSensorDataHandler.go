package handlers

import (
	"github.com/buungoo/Beehive-interface/models"
	"github.com/buungoo/Beehive-interface/utils"

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
	err = conn.QueryRow(context.Background(), queryBeehive, reading.BeehiveID.String()).Scan(&beehiveID) // Convert the macaddr back to string for the query
	if err != nil {
		return fmt.Errorf("failed to find beehive with MAC address %s: %v", reading.BeehiveID.String(), err)
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

	// Verify the sensor reading
	isValid, verificationMessage := reading.VerifyInputData()

	if !isValid {
		// If verification fails, log the error in the beehive_status table
		insertErrorQuery := `
			INSERT INTO beehive_status (sensor_id, beehive_id, sensor_type, description, solved, read, time_of_error) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = conn.Exec(context.Background(), insertErrorQuery,
			reading.SensorID, beehiveID, string(reading.SensorType), verificationMessage, false, false, reading.Time)
		if err != nil {
			utils.LogWarn(fmt.Sprintf("Failed to log error in beehive_status table: %v", err))
		} else {
			utils.LogInfo(fmt.Sprintf("Logged error in beehive_status: %s", verificationMessage))
		}
	} else {
		// If the reading is valid, check for active issues in beehive_status
		queryActiveIssue := `
			SELECT COUNT(*) FROM beehive_status 
			WHERE sensor_id = $1 AND beehive_id = $2 AND solved = false`
		var activeIssueCount int
		err = conn.QueryRow(context.Background(), queryActiveIssue, reading.SensorID, beehiveID).Scan(&activeIssueCount)
		if err != nil {
			utils.LogError("Error checking for active issues in beehive_status: ", err)
			return fmt.Errorf("error checking active issues: %v", err)
		}

		// Mark the issue as solved if there's an active issue
		if activeIssueCount > 0 {
			updateIssueQuery := `
				UPDATE beehive_status 
				SET solved = true, time_of_resolution = $1 
				WHERE sensor_id = $2 AND beehive_id = $3 AND solved = false`
			_, err = conn.Exec(context.Background(), updateIssueQuery, reading.Time, reading.SensorID, beehiveID)
			if err != nil {
				utils.LogError("Failed to mark active issue as solved: ", err)
			} else {
				utils.LogInfo(fmt.Sprintf("Marked active issue as solved for SensorID=%d, BeehiveID=%d", reading.SensorID, beehiveID))
			}
		}
	}

	// Insert the sensor reading into the `sensor_data` table
	insertReadingQuery := `
		INSERT INTO sensor_data (sensor_id, beehive_id, sensor_type, value, time) 
		VALUES ($1, $2, $3, $4, $5)`
	_, err = conn.Exec(context.Background(), insertReadingQuery,
		reading.SensorID, beehiveID, string(reading.SensorType), reading.Value, reading.Time)
	if err != nil {
		return fmt.Errorf("failed to insert sensor reading: %v", err)
	}

	return nil
}

