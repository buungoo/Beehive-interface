package api

import (
	"beehive_api/handlers"
	"net/http"
	"strings"
)

func InitRoutes(mux *http.ServeMux) {
	// Register routes
	mux.HandleFunc("/beehive/", beehiveHandler)

}

// Direct to correct beehive-handler based on id
func beehiveHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) > 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
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
		http.Error(w, "Unknown sensor type", http.StatusBadRequest)
	}
}
