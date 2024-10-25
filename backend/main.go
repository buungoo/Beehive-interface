package main

import (
	"beehive_api/api"
	"beehive_api/db"
	"beehive_api/test"
	"fmt"
	"net/http"
)

func main() {
	// Connect to database
	dbpool, err := db.InitializeDatabaseConnection()
	if err != nil {
		fmt.Printf("Error initializing database", err)
	}
	defer dbpool.Close()
	fmt.Println("Successfully connected to the database")

	// Initialize App struct with database connection
	//app := &db.Handle{DB: conn}

	// Initialize database tables if not already done
	err = db.InitializeTables(dbpool)
	if err != nil {
		fmt.Printf("Database initialization failed: %v", err)
	}
	fmt.Println("Tables successfully generated")

	// Inject testdata
	err = test.InjectTestData(dbpool)
	if err != nil {
		fmt.Printf("Injection of testdata failed: %v", err)
	}

	fmt.Println("Test data successfully injected!")
	
	// Create a new request multiplexer that takes incoming
	// requests and dispatches them to matching handlers
	mux := http.NewServeMux()

	api.InitRoutes(mux, dbpool)
	
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		fmt.Print("Sever error", err)
	}
	fmt.Println("Http API server started")
	

}
