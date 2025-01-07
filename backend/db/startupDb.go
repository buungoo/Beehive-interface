package db

import (
	"context"

	"github.com/buungoo/Beehive-interface/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StartupDb adds needed data to the database at startup.
func StartupDb(dbPool *pgxpool.Pool) error {
	var macAddr string = "0080e115000adf82"
	var beehiveName string = "Beehive A"

	// Acquire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
		return err
	}
	defer conn.Release()

	const sqlQuery = `
	INSERT INTO beehives (name, key) 
	VALUES ($1, $2)
	ON CONFLICT (key) DO NOTHING;
`

	_, err = conn.Exec(context.Background(), sqlQuery, beehiveName, macAddr)
	if err != nil {
		utils.LogError("failed to insert beehive A", err)
		return err
	}

	return nil
}
