package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/buungoo/Beehive-interface/api"
	"github.com/buungoo/Beehive-interface/db"
	"github.com/buungoo/Beehive-interface/mqtt"
	"github.com/buungoo/Beehive-interface/test"
	"github.com/buungoo/Beehive-interface/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize the logger
	logFile, err := utils.InitLogger("./logs/logFile.log")
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logFile.Close()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		utils.LogFatal("Failed to load environment variables", err)
	}

	// Connect to the database
	dbpool, err := db.InitializeDatabaseConnection()
	if err != nil {
		utils.LogFatal("Error initializing database connection", err)
	}
	defer dbpool.Close()
	utils.LogInfo("Successfully connected to the database")

	// Initialize database tables
	if err := db.InitializeTables(dbpool); err != nil {
		utils.LogFatal("Error initializing database tables", err)
	}
	utils.LogInfo("Database tables successfully initialized")

	// Inject test data into the database
	if err := test.InjectTestData(dbpool); err != nil {
		utils.LogFatal("Error injecting test data", err)
	}
	utils.LogInfo("Test data successfully injected")

	// Start MQTT subscriber in a separate goroutine
	utils.LogInfo("Starting MQTT subscriber")
	// go func() {
	// 	if err := mqtt.SetupMQTTSubscriber(dbpool); err != nil {
	// 		utils.LogError("MQTT subscriber encountered an error", err)
	// 	}
	// }()
	go func() {
		mqtt.SetupMQTTSubscriber(dbpool)
	}()

	// Initialize HTTP server routes
	mux := http.NewServeMux()
	api.InitRoutes(mux, dbpool)

	// HTTPS server configuration
	serverCert := "certs/server.crt"
	serverKey := "certs/server.key"
	serverAddr := ":8443"

	// Start the HTTPS server in a separate goroutine
	go func() {
		utils.LogInfo("Starting HTTPS server on " + serverAddr)
		if err := http.ListenAndServeTLS(serverAddr, serverCert, serverKey, mux); err != nil {
			utils.LogFatal("Failed to start HTTPS server", err)
		}
	}()

	// Wait for termination signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	utils.LogInfo("Application is running. Press Ctrl+C to shut down.")

	// Block until a termination signal is received
	<-stopChan
	utils.LogInfo("Shutting down application")

	// Perform cleanup (if any)
	utils.LogInfo("Application stopped gracefully")
}
