package handlers

import (
	"beehive_api/authentication"
	"beehive_api/models"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"github.com/jackc/pgx/v5"
	"encoding/json"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	
	//w.Header().Set("Content-Type", "application/json")
	
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("The user request value %v", user)
	
	if user.Username == "Emil" && user.Password == "123456" {
		tokenString, err := authentication.CreateToken(user.Username)
		if err != nil {

			utils.SendErrorResponse(w, "No username found", http.StatusInternalServerError)
			//fmt.Errorf("No username found")
		}
		// Return the token in the response
		response := map[string]string{
			"message": "User validated",
			"token":   tokenString,
		}
		utils.SendJSONResponse(w, response, http.StatusOK) 
		//w.WriteHeader(http.StatusOK)
		//fmt.Fprint(w, tokenString)
		return
	} else {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized)
		//w.WriteHeader(http.StatusUnauthorized)
		//fmt.Fprint(w, "Invalid credentials")
	}
}
//utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
