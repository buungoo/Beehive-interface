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
		"id" SERIAL,
		"username" VARCHAR(255) UNIQUE NOT NULL CHECK (username <> ''),
		"password" BYTEA NOT NULL,
		PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "beehives" (
		"id" SERIAL,
		"name" VARCHAR NOT NULL,
		"key" MACADDR8 UNIQUE NOT NULL,
		PRIMARY KEY ("id")
	);

	CREATE TABLE IF NOT EXISTS "user_beehive" (
		"user_id" INTEGER NOT NULL,
		"beehive_id" INTEGER NOT NULL,
		PRIMARY KEY ("user_id", "beehive_id"),
		FOREIGN KEY ("user_id") REFERENCES "users" ("id"),
		FOREIGN KEY ("beehive_id") REFERENCES "beehives" ("id")
	);

	CREATE TABLE IF NOT EXISTS "sensors" (
		"id" INTEGER NOT NULL,
		"type" VARCHAR NOT NULL,
		"beehive_id" INTEGER NOT NULL,
		PRIMARY KEY ("id", "type", "beehive_id"),
		FOREIGN KEY ("beehive_id") REFERENCES "beehives" ("id")
	);

	CREATE TABLE IF NOT EXISTS "beehive_status" (
		"issue_id" SERIAL,
		"sensor_id" INTEGER NOT NULL,
		"beehive_id" INTEGER NOT NULL,
		"sensor_type" VARCHAR NOT NULL,
		"description" VARCHAR,
		"solved" BOOLEAN,
		"read" BOOLEAN,
		"time_of_error" TIMESTAMPTZ,
		"time_read" TIMESTAMPTZ,
		PRIMARY KEY ("issue_id", "beehive_id", "sensor_id"),
		FOREIGN KEY ("beehive_id") REFERENCES "beehives" ("id"), 
		FOREIGN KEY ("sensor_id", "sensor_type", "beehive_id") REFERENCES "sensors" ("id", "type", "beehive_id")
	);

	CREATE TABLE IF NOT EXISTS "sensor_data" (
		"sensor_id" INTEGER NOT NULL,
		"beehive_id" INTEGER NOT NULL,
		"sensor_type" VARCHAR NOT NULL,
		"value" FLOAT NOT NULL,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("sensor_id","beehive_id", "time"),
		FOREIGN KEY ("sensor_id", "sensor_type", "beehive_id") REFERENCES "sensors" ("id", "type", "beehive_id")
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
