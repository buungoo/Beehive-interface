package db

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v5"
	"os"
)

type Handle struct {
	DB *pgx.Conn
}

func InitializeDatabaseConnection() (*pgx.Conn, error) {
	// Get the DATABASE_URL from environment variables
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	
	return conn, nil

}
