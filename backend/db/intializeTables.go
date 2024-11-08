package db

import (
	"beehive_api/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeTables(dbpool *pgxpool.Pool) error {
	// Create the tables if they don't already exist
	createTablesSQL := `
	CREATE EXTENSION IF NOT EXISTS timescaledb;

	CREATE TABLE IF NOT EXISTS "users" (
		"id" SERIAL PRIMARY KEY,
		"username" VARCHAR(255) UNIQUE NOT NULL CHECK (username <> ''),
		"password" BYTEA NOT NULL
	);

	CREATE TABLE IF NOT EXISTS "beehives" (
		"id" INTEGER PRIMARY KEY,
		"name" VARCHAR NOT NULL,
		"user_id" INTEGER REFERENCES "users" ("id")
	);

	CREATE TABLE IF NOT EXISTS "sensors" (
		"id" INTEGER UNIQUE,
		"type" VARCHAR NOT NULL,
		"beehive_id" INTEGER REFERENCES "beehives" ("id"),
		PRIMARY KEY ("id", "type")
	);

	CREATE TABLE IF NOT EXISTS "beehive_status" (
		"beehive_id" INTEGER REFERENCES "beehives" ("id"),
		"sensor_id" INTEGER REFERENCES "sensors" ("id"),
		"solved" BOOLEAN,
		"read" BOOLEAN,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS "sensor_data" (
		"sensor_id" INTEGER REFERENCES "sensors" ("id"),
		"beehive_id" INTEGER REFERENCES "beehives" ("id"),
		"sensor_type" VARCHAR NOT NULL,
		"value" FLOAT NOT NULL,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("sensor_id","beehive_id", "time"),
		FOREIGN KEY ("sensor_id", "sensor_type") REFERENCES "sensors" ("id", "type")
	);
	`

	// Execute the SQL to create the tables
	_, err := dbpool.Exec(context.Background(), createTablesSQL)
	if err != nil {
		utils.LogFatal("Failed to create tables: ", err)
	}

	// Check if hypertable exists
	hypertableCheckSQL := `
	SELECT EXISTS (SELECT 1 FROM timescaledb_information.hypertables WHERE hypertable_name = 'sensor_data');
	`

	var hypertableExists bool
	err = dbpool.QueryRow(context.Background(), hypertableCheckSQL).Scan(&hypertableExists)
	if err != nil {
		utils.LogFatal("Failed to check for hypertable: ", err)
	}

	// Create the hypertable if it doesn't already exist
	if !hypertableExists {
		createHypertable := `
		SELECT create_hypertable('sensor_data', 'time');
		`

		_, err = dbpool.Exec(context.Background(), createHypertable)
		if err != nil {
			utils.LogFatal("Failed to create hypertable: ", err)
		}

	}

	return nil
}
