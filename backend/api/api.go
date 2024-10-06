package api

import (
	"beehive_api/handlers"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"strings"
)

func InitRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	// Register routes
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerHandler(w, r, conn)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, conn)
	})

	mux.HandleFunc("/beehive/", func(w http.ResponseWriter, r *http.Request) {
		beehiveHandler(w, r, conn)
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
}
func loginHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
}

// Direct to correct handler based on http request
func beehiveHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(pathParts) != 3 {
		if len(pathParts) > 3 {
			fmt.Println("Len > 3")
			utils.SendErrorResponse(w, "URL is to long", http.StatusBadRequest)
			//http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		} else {
			fmt.Println("Len < 3")
			utils.SendErrorResponse(w, "URL is to short", http.StatusBadRequest)
			return
		}
	}

	beehiveId := pathParts[1]
	sensorType := pathParts[2]

	switch r.Method {
	case http.MethodGet:
		handlers.GetSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodPost:
		handlers.AddSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodPut:
		handlers.UpdateSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodDelete:
		handlers.DeleteSensorData(w, r, conn, beehiveId, sensorType)

	default:
		utils.SendErrorResponse(w, "HTTP method not found", http.StatusBadRequest)

	}

}

