package db

import (
	"fmt"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"log"
	"time"
)

type Handle struct {
	DB *pgx.Conn
}

func InitializeDatabaseConnection() (*pgxpool.Pool, error) {
	// // Get the DATABASE_URL from environment variables
	// connStr := os.Getenv("DATABASE_URL")

	// if connStr == "" {
	// 	 fmt.Errorf("DATABASE_URL environment variable not set")
	//	return nil
	// }

	// Connect to the database
	// conn, err := pgx.Connect(context.Background(), connStr)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// dbpool, err := pgxpool.New(context.Background(), connStr)
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	
	return dbpool, nil

}

func config()(*pgxpool.Config) {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Get the DATABASE_URL from environment variables
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		fmt.Errorf("DATABASE_URL environment variable not set")
		return nil
	}

	dbConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Failed to create a connetionconfig, err: ", err)
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
