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

// Direct to correct beehive-handler based on id
func beehiveHandler(w http.ResponseWriter, r *http.Request) {
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

	id := pathParts[1]
	sensor := pathParts[2]

	switch sensor {
	case "humidity":
		handlers.HumidityHandler(w, id, r)

	case "temperature":
		handlers.TemperatureHandler(w, id, r)

	case "oxygen":
		handlers.OxygenHandler(w, id, r)

	case "weight":
		handlers.WeightHandler(w, id, r)

	case "microphone":
		handlers.SoundHandler(w, id, r)
	default:
		utils.SendErrorResponse(w, "Sensor not found", http.StatusNotFound)

		//http.Error(w, "Unknown sensor type", http.StatusBadRequest)
	}

}
