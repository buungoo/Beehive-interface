package api

import (
	"fmt"
	"beehive_api/handlers"
	"net/http"
	"strings"
	"beehive_api/utils"
)

func InitRoutes(mux *http.ServeMux) {
	// Register routes
	mux.HandleFunc("/beehive/", beehiveHandler)

}

// Direct to correct beehive-handler based on id
func beehiveHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 4 {
		fmt.Println("Len > 3")
		utils.SendErrorResponse(w, "URL is to long", http.StatusBadRequest)
		//http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}
	if len(pathParts) < 3 {
		fmt.Println("Len < 3")
		utils.SendErrorResponse(w, "URL is not correct", http.StatusBadRequest)
		return
	}
	
	id := pathParts[2]
	sensor := pathParts[3]

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
		utils.SendErrorResponse(w, "Unknown sensor type", http.StatusBadRequest)

		//http.Error(w, "Unknown sensor type", http.StatusBadRequest)
	}
	
}
