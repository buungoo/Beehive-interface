package main

import (
	"beehive_api/api"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting API server...")

	// Create a new request multiplexer that takes incoming
	// requests and dispatches them to matching handlers
	mux := http.NewServeMux()

	api.InitRoutes(mux)

	// mux.HandleFunc("/beehive", func(w http.ResponseWriter, r *http.Request) {
	// 	randomNum := api.GetRandomNum()
	// 	fmt.Println(w, "Randomnumber: %d", randomNum)
	// })

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Print("Sever error", err)
	}

}
