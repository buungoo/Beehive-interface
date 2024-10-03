package handlers

import (
	"beehive_api/utils"
	"net/http"
)

func WeightHandler(w http.ResponseWriter, id string, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getWeight()
	case http.MethodPost:
		postWeight()
	case http.MethodPut:
		updateWeight()
	case http.MethodDelete:
		deteleWeight()
	default:
		utils.SendErrorResponse(w, "Not a supported http method", http.StatusBadRequest)
	}

}

func getWeight() {

}

func postWeight() {

}

func updateWeight() {

}

func deteleWeight() {

}
