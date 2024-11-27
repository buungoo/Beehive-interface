// Package api contains the all the endpoints and handles routing http-request to the correct handler.
//
// This package is the entrypoint for all Api calls. It sends requests to the correct handler.
// In som occasions it performes some checks before passing it along to a handler.
package api

import (
	"beehive_api/authentication"
	"beehive_api/handlers"
	"beehive_api/models"
	"beehive_api/utils"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitRouters registers routes and sends http-requests to the correct handler
func InitRoutes(mux *http.ServeMux, dbPool *pgxpool.Pool) {

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, dbPool)
	})

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, dbPool)
	})

	mux.HandleFunc("GET /beehive/{beehiveId}/status", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.GetBeehiveStatus(w, r, dbPool, beehiveId)
	}))

	mux.HandleFunc("GET /beehive/{beehiveId}/status/list", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.GetBeehiveStatusList(w, r, dbPool, beehiveId)
	}))

	mux.HandleFunc("POST /beehive/add", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		handlers.AddBeehiveToUser(w, r, dbPool)
	}))

	mux.HandleFunc("POST /beehive/{beehiveId}/remove", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.RemoveBeehiveFromUser(w, r, dbPool, beehiveId)
	}))

	mux.HandleFunc("POST /beehive/{beehiveId}/sensor-data/add", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.AddSensorData(w, r, dbPool, beehiveId)
	}))

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

		handlers.GetDataByDate(w, r, dbPool, beehiveId, date1, date2)

	}))

	mux.HandleFunc("GET /beehive/{beehiveId}/sensor-data/average/{startDate}/{endDate}", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}

		date1 := r.PathValue("startDate")
		date2 := r.PathValue("endDate")

		handlers.GetAverageDataByDate(w, r, dbPool, beehiveId, date1, date2)

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

func testAuthentication(w http.ResponseWriter, r *http.Request) {
	utils.SendJSONResponse(w, "Token is Valid", http.StatusOK)
	return
}
