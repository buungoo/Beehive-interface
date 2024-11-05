package db

import (
	"beehive_api/utils"
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handle struct {
	DB *pgx.Conn
}

func InitializeDatabaseConnection() (*pgxpool.Pool, error) {

	// Connect to the database
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config())
	if err != nil {
		utils.LogFatal("Unable to connect to database", err)
	}
	return dbpool, nil

}

func config() *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Get the DATABASE_URL from environment variables
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Println()
		utils.LogFatal("DATABASE_URL environment variable not set", errors.New("empty database url"))
	}

	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		utils.LogFatal("Failed to create a connetionconfig, err: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	// BeforeAcquire is called before a connection is acquired from the pool. It must return true to allow the
	// acquisition or false to indicate that the connection should be destroyed and a different connection should be
	// acquired.
	// dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
	// 	log.Println("Before acquiring the connection pool to the database!!")
	// 	return true
	// }

	// // AfterRelease is called after a connection is released, but before it is returned to the pool. It must return true to
	// // return the connection to the pool or false to destroy the connection.
	// dbConfig.AfterRelease = func(c *pgx.Conn) bool {
	// 	log.Println("After releasing the connection pool to the database!!")
	// 	return true
	// }

	// // BeforeClose is called right before a connection is closed and removed from the pool.
	// dbConfig.BeforeClose = func(c *pgx.Conn) {
	// 	log.Println("Closed the connection pool to the database!!")
	// }

	return dbConfig

}
