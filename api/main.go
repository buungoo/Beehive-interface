package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("hello world")
	// Create a new request multiplexer
	// Take incoming request and dispatch them to matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/recipe", &RecipesHandler{})
	mux.Handle("/recipes/", &RecipesHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)
}

type RecipesHandler struct{}

func (h *RecipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my recipe page"))
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
