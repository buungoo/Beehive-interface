package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func InitializeTables(conn *pgx.Conn) error {
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
		"type" varchar NOT NULL,
		"beehive_id" integer REFERENCES "beehives" ("id")
	);

	CREATE TABLE IF NOT EXISTS "sensor_data" (
		"sensor_id" integer REFERENCES "sensors" ("id"),
		"value" float NOT NULL,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("sensor_id", "time")
	);
	`

	// Execute the SQL to create the tables
	_, err := conn.Exec(context.Background(), createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	// Check if hypertable exists
	hypertableCheckSQL := `
	SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'sensor_data');
	`

	var hypertableExists bool
	err = conn.QueryRow(context.Background(), hypertableCheckSQL).Scan(&hypertableExists)
	if err != nil {
		return fmt.Errorf("failed to check for hypertable: %v", err)
	}

	// Create the hypertable if it doesn't already exist
	if !hypertableExists {
		createHypertable := `
		SELECT create_hypertable('sensor_data', 'time');
		`

		_, err = conn.Exec(context.Background(), createHypertable)
		if err != nil {
			return fmt.Errorf("failed to create hypertable: %v", err)
		}

		// Add hash partitioning for sensor_id
		addHashPartitioning := `
		SELECT add_dimension('sensor_data', 'sensor_id', 4);
		`

		_, err = conn.Exec(context.Background(), addHashPartitioning)
		if err != nil {
			return fmt.Errorf("failed to add hash partitioning: %v", err)
		}
	}

	return nil
}
