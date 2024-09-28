package handlers

import (
	"beehive_api/utils"
	"net/http"
)

func OxygenHandler(w http.ResponseWriter, id string, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getOxygen()
	case http.MethodPost:
		postOxygen()
	case http.MethodPut:
		updateOxygen()
	case http.MethodDelete:
		deteleOxygen()
	default:
		utils.SendErrorResponse(w, "Not a supported http method", http.StatusBadRequest)
	}

}

func getOxygen() {

}

func postOxygen() {

}

func updateOxygen() {

}

func deteleOxygen() {

}
