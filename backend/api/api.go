package api

import (
	"beehive_api/handlers"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"strings"
)

func InitRoutes(mux *http.ServeMux) {
	// Register routes
	mux.HandleFunc("/beehive/", beehiveHandler)

}

// Direct to correct handler based on http request
func beehiveHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if pathParts[1] == "login" {
		// Here we authenticate the user with token or something
		// Need a return here
	}
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

	beehive_id := pathParts[1]
	sensor := pathParts[2]

	switch r.Method {
	case http.MethodGet:
		handlers.GetSensorData(w, beehive_id, r, sensor)

	case http.MethodPost:
		handlers.AddSensorData(w, beehive_id, r, sensor)

	case http.MethodPut:
		handlers.UpdateSensorData(w, beehive_id, r, sensor)

	case http.MethodDelete:
		handlers.DeleteSensorData(w, beehive_id, r, sensor)

	default:
		utils.SendErrorResponse(w, "HTTP method not found", http.StatusBadRequest)

	}

}
