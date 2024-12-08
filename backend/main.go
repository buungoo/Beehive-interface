package main

import (
	"log"
	"net/http"

	"github.com/buungoo/Beehive-interface/api"
	"github.com/buungoo/Beehive-interface/db"
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

	// Initialize environment variables
	err = godotenv.Load()
	if err != nil {
		utils.LogFatal("failed to get environment variables", err)
	}

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
        
	server_krt := "certs/server.crt"
	server_key := "certs/server.key"

	if err := http.ListenAndServeTLS(":8443",server_krt,server_key, mux); err != nil {
		utils.LogError("Http server could not start: ", err)
	}

}
