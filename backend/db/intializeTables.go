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
    "id" serial PRIMARY KEY,
    "username" VARCHAR(255) UNIQUE NOT NULL CHECK (username <> ''),
    "password" BYTEA NOT NULL
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
		"beehive_id" integer REFERENCES "beehives" ("id"),
		"value" float NOT NULL,
		"time" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		PRIMARY KEY ("sensor_id", "time")
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
