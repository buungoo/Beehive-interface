package main

import (
	"beehive_api/api"
	"beehive_api/db"
	"beehive_api/mqtt"
	"beehive_api/test"
	"beehive_api/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Initialize the logger
	logFile, err := utils.InitLogger("./logs/logFile.log")
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

	// Set up MQTT subscriber
	go func() {
		mqtt.SetupMQTTSubscriber(dbpool)
	}()

	// Wait for termination signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-stopChan
	utils.LogInfo("Shutting down application")

}
