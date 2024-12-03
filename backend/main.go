package main

import (
	"beehive_api/api"
	"beehive_api/db"
	"beehive_api/mqtt"

	// "beehive_api/test"
	"beehive_api/utils"
	"log"
	"net/http"

	"github.com/buungoo/Beehive-interface/api"
	"github.com/buungoo/Beehive-interface/db"
	"github.com/buungoo/Beehive-interface/utils"

	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the logger
	logFile, err := utils.InitLogger("./logs/logFile.log")
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logFile.Close()

	// Initialize environment variables
	err = godotenv.Load()
	if err != nil {
		utils.LogFatal("failed to get environment variables", err)
	}

	// Connect to database
	dbpool, err := db.InitializeDatabaseConnection()
	if err != nil {
		utils.LogError("Error initializing database", err)
		return
	}
	defer dbpool.Close()
	utils.LogInfo("Successfully connected to the database")

	// Initialize database tables
	if err := db.InitializeTables(dbpool); err != nil {
		utils.LogError("Database initialization failed: ", err)
		return
	}
	utils.LogInfo("Database tables successfully generated")

	// Inject test data
	// if err := test.InjectTestData(dbpool); err != nil {
	// 	utils.LogError("Injection of test data failed: ", err)
	// 	return
	// }
	// utils.LogInfo("Test data successfully injected to the database!")

	// HTTP server setup
	mux := http.NewServeMux()
	api.InitRoutes(mux, dbpool)
	go func() {
		utils.LogInfo("Starting HTTP server on port 8080")
		if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
			utils.LogError("HTTP server failed: ", err)
			os.Exit(1) // Exit on HTTP server failure
		}
	}()

	utils.LogInfo("Starting MQTT subscriber")
	// MQTT Subscriber setup
	go mqtt.SetupMQTTSubscriber(dbpool)

	// Wait for termination signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	utils.LogInfo("Shutting down application")
}
