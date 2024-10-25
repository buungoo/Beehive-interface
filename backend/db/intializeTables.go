package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeTables(dbpool *pgxpool.Pool) error {
	// Create the tables if they don't already exist
	createTablesSQL := `
	CREATE EXTENSION IF NOT EXISTS timescaledb;

	CREATE TABLE IF NOT EXISTS "users" (
		"id" serial PRIMARY KEY,
		"username" varchar NOT NULL,
		"password" varchar NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "beehives" (
		"id" serial PRIMARY KEY,
		"name" varchar NOT NULL,
		"user_id" integer REFERENCES "users" ("id")
	);

	CREATE TABLE IF NOT EXISTS "sensors" (
		"id" serial PRIMARY KEY,
		"type" varchar NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "sensor_data" (
		"sensor_id" integer REFERENCES "sensors" ("id"),
		"beehive_id" integer REFERENCES "beehives" ("id"),
		"value" float NOT NULL,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("sensor_id", "time")
	);

	-- SQL Function to Fetch Latest Sensor Data for a Beehive
	CREATE OR REPLACE FUNCTION fetch_latest_sensor_data_for_beehive(
		p_beehive_id INT,
		p_sensor_type VARCHAR
	)
	RETURNS TABLE(sensor_value FLOAT, measurement_time TIMESTAMPTZ) AS $$
	BEGIN
		RETURN QUERY
		SELECT sd.value, sd.time
		FROM sensor_data sd
		JOIN sensors s ON sd.sensor_id = s.id
		WHERE sd.beehive_id = p_beehive_id
		AND s.type = p_sensor_type
		ORDER BY sd.time DESC
		LIMIT 1;
	END;
	$$ LANGUAGE plpgsql;
	`

	// Execute the SQL to create the tables
	_, err := dbpool.Exec(context.Background(), createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Check if hypertable exists
	hypertableCheckSQL := `
	SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'sensor_data');
	`

	var hypertableExists bool
	err = dbpool.QueryRow(context.Background(), hypertableCheckSQL).Scan(&hypertableExists)
	if err != nil {
		return fmt.Errorf("failed to check for hypertable: %v", err)
	}

	// Create the hypertable if it doesn't already exist
	if !hypertableExists {
		createHypertable := `
		SELECT create_hypertable('sensor_data', 'time');
		`

		_, err = dbpool.Exec(context.Background(), createHypertable)
		if err != nil {
			return fmt.Errorf("failed to create hypertable: %v", err)
		}

		// Add hash partitioning for sensor_id
		addHashPartitioning := `
		SELECT add_dimension('sensor_data', 'sensor_id', 4);
		`

		_, err = dbpool.Exec(context.Background(), addHashPartitioning)
		if err != nil {
			return fmt.Errorf("failed to add hash partitioning: %v", err)
		}
	}

	return nil
}
