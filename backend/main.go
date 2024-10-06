package main

import (
	"beehive_api/api"
	"beehive_api/db"
	"beehive_api/test"
	"context"
	"fmt"
	"net/http"
)

func main() {
	// Connect to database
	conn, err := db.InitializeDatabaseConnection()
	if err != nil {
		fmt.Printf("Error initializing database", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Successfully connected to the database")

	// Initialize App struct with database connection
	//app := &db.Handle{DB: conn}

	// Initialize database tables if not already done
	err = db.InitializeTables(conn)
	if err != nil {
		fmt.Printf("Database initialization failed: %v", err)
	}
	fmt.Println("Tables successfully generated")

	// Inject testdata
	err = test.InjectTestData(conn)
	if err != nil {
		fmt.Printf("Injection of testdata failed: %v", err)
	}

	fmt.Println("Test data successfully injected!")
	
	// Create a new request multiplexer that takes incoming
	// requests and dispatches them to matching handlers
	mux := http.NewServeMux()

	api.InitRoutes(mux, conn)
	
	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		fmt.Print("Sever error", err)
	}
	fmt.Println("Http API server started")
	

}
