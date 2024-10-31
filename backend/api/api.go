package api

import (
	"beehive_api/authentication"
	"beehive_api/handlers"
	"beehive_api/models"
	"beehive_api/utils"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"log"
)
// Register routes and send to correct handler
func InitRoutes(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	
	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, dbPool)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, dbPool)
	})

	mux.HandleFunc("GET /beehive/list", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetBeehiveList(w, r, dbPool)
	}))

	mux.HandleFunc("GET /beehive/{beehiveId}/sensor-data/{startDate}/{endDate}", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}

		date1 := r.PathValue("startDate")
		date2 := r.PathValue("endDate")
	
		handlers.GetDataByDate(w,r, dbPool, beehiveId, date1, date2)

	}))

	mux.HandleFunc("GET /beehive/{beehiveId}/sensor-data/latest", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.GetLatestSensorData(w, r, dbPool, beehiveId)
	}))

	mux.HandleFunc("GET /beehive/{beehiveId}/{sensorType}/latest", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		// Validate the sensortype
		sensorType := models.SensorType(r.PathValue("sensorType"))
		if !sensorType.IsValid() {
			log.Println("Requested invalid sensortype")
			utils.SendErrorResponse(w, "Invalid sensortype", http.StatusBadRequest)
			return
		}
		sensorTypeString := string(sensorType)
		handlers.GetLatestOfSensortype(w, r, dbPool, beehiveId, sensorTypeString)
	}))

	mux.HandleFunc("POST /test", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		testAuthentication(w, r)
	}))


}


func testAuthentication(w http.ResponseWriter, r *http.Request,) {
	utils.SendJSONResponse(w, "Token is Valid", http.StatusOK)
	return
}

