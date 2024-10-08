package handlers

import (
	"beehive_api/handlers"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"strings"
	"github.com/jackc/pgx/v5"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	func LoginHandler(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
	   
		var u User
		json.NewDecoder(r.Body).Decode(&u)
		fmt.Printf("The user request value %v", u)
		
		if u.Username == "Chek" && u.Password == "123456" {
		  tokenString, err := CreateToken(u.Username)
		  if err != nil {
			utils.SendErrorResponse(w, "No username found", http.StatusInternalServerError)
			//fmt.Errorf("No username found")
		   }
		  utils.SendJSONResponse(w, "User validated", http.StatusOK) 
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
}