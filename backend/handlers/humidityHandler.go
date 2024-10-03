package handlers

import (
	"beehive_api/utils"
	"net/http"
)

func HumidityHandler(w http.ResponseWriter, id string, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHumidity()
	case http.MethodPost:
		postHumidity()
	case http.MethodPut:
		updateHumidity()
	case http.MethodDelete:
		deteleHumidity()
	default:
		utils.SendErrorResponse(w, "Not a supported http method", http.StatusBadRequest)
	}

}

func getHumidity() {

}

func postHumidity() {

}

func updateHumidity() {

}

func deteleHumidity() {

}
