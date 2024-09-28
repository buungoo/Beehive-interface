package handlers

import (
	"beehive_api/utils"
	"net/http"
)

func SoundHandler(w http.ResponseWriter, id string, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getSound()
	case http.MethodPost:
		postSound()
	case http.MethodPut:
		updateSound()
	case http.MethodDelete:
		deteleSound()
	default:
		utils.SendErrorResponse(w, "Not a supported http method", http.StatusBadRequest)
	}

}

func getSound() {

}

func postSound() {

}

func updateSound() {

}

func deteleSound() {

}
