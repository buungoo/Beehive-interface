package handlers

import (
	"beehive_api/utils"
	"net/http"
)

// Mock data for testing
var mockTemperatureData = map[string]interface{}{
	"id":          "1",
	"temperature": 13.37,
}

type TemperatureData struct {
	ID          string  `json:"id"`
	Temperature float32 `json:"temperature"`
}

func TemperatureHandler(w http.ResponseWriter, id string, r *http.Request) {

	// Mock data for testing
	// tempData := TemperatureData{
	// 	ID:          id,
	// 	Temperature: 13.37,
	// }

	switch r.Method {
	case http.MethodGet:
		getTemperature(w, mockTemperatureData)
	case http.MethodPost:
		postTemperature()
	case http.MethodPut:
		updateTemperature()
	case http.MethodDelete:
		deteleTemperature()
	default:
		utils.SendErrorResponse(w, "Not a supported http method", http.StatusBadRequest)
	}

	// // Convert to JSON format for Responewriter
	// if err := json.NewEncoder(w).Encode(tempData); err != nil {
	// 	http.Error(w, "Could not encode response to JSON", http.StatusInternalServerError)
	// 	return
	// }

}

func getTemperature(w http.ResponseWriter, data interface{}) {
	// This is simply for testing, will be replaced with real logic
	utils.SendJSONResponse(w, data, 200)

}

func postTemperature() {

}

func updateTemperature() {

}

func deteleTemperature() {

}
