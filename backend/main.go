package main

import (
	"beehive_api/api"
	"beehive_api/db"
	"beehive_api/test"
	"beehive_api/utils"
	"log"
	"net/http"
)

func main() {
	// Initialize the logger
	logFile, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logFile.Close()

	// Connect to database
	dbpool, err := db.InitializeDatabaseConnection()
	if err != nil {
		utils.LogError("Error initializing database", err)
	}
	defer dbpool.Close()
	utils.LogInfo("Successfully connected to the database")

	// Initialize App struct with database connection
	//app := &db.Handle{DB: conn}

	// Initialize database tables if not already done
	err = db.InitializeTables(dbpool)
	if err != nil {
		utils.LogError("Database initialization failed: ", err)
	}
	utils.LogInfo("Database tables successfully generated")

	// Inject testdata
	err = test.InjectTestData(dbpool)
	if err != nil {
		utils.LogError("Injection of testdata failed: %v", err)
	}
	utils.LogInfo("Test data successfully injected to the database!")

	// Create a new request multiplexer that takes incoming
	// requests and dispatches them to matching handlers
	mux := http.NewServeMux()

	api.InitRoutes(mux, dbpool)

	utils.LogInfo("Starting http server")

	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		utils.LogError("Http server could not start: ", err)
	}

}
