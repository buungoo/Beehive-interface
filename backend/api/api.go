package api

import (
	"beehive_api/authentication"
	"beehive_api/handlers"
	"beehive_api/utils"
	"fmt"
	"net/http"
	"strings"
	"github.com/jackc/pgx/v5"
)

func InitRoutes(mux *http.ServeMux, conn *pgx.Conn) {
	// Register routes
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, conn)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, conn)
	})

	mux.HandleFunc("/beehive/", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveHandler(w, r, conn)
	}))

	mux.HandleFunc("/test", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		testAuthentication(w, r)
	}))


}



// Direct to correct handler based on http request
func beehiveHandler(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
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

	beehiveId := pathParts[1]
	sensorType := pathParts[2]

	switch r.Method {
	case http.MethodGet:
		handlers.GetSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodPost:
		handlers.AddSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodPut:
		handlers.UpdateSensorData(w, r, conn, beehiveId, sensorType)

	case http.MethodDelete:
		handlers.DeleteSensorData(w, r, conn, beehiveId, sensorType)

	default:
		utils.SendErrorResponse(w, "HTTP method not found", http.StatusBadRequest)

	}

}

func testAuthentication(w http.ResponseWriter, r *http.Request,) {
	utils.SendJSONResponse(w, "Token is Valid", http.StatusOK)
	return
}

