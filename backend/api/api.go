package api

import (
	"beehive_api/authentication"
	"beehive_api/handlers"
	"beehive_api/models"
	"beehive_api/utils"
	"net"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Middleware for Docker bridge network authentication
func DockerBridgeAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract IP from "ip:port"
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			utils.SendErrorResponse(w, "Forbidden: Invalid Remote Address", http.StatusForbidden)
			return
		}

		// Define the Docker bridge network CIDR
		bridgeCIDR := "172.18.0.0/16"
		_, bridgeNet, err := net.ParseCIDR(bridgeCIDR)
		if err != nil {
			utils.SendErrorResponse(w, "Server Error: Invalid Network Configuration", http.StatusInternalServerError)
			return
		}

		// Check if the request IP is within the bridge network
		parsedIP := net.ParseIP(ip)
		if !bridgeNet.Contains(parsedIP) {
			utils.SendErrorResponse(w, "Forbidden: Access Denied", http.StatusForbidden)
			return
		}

		// Call the next handler if the IP is valid
		next(w, r)
	}
}

// Register routes and send to correct handler
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

	mux.HandleFunc("POST /beehive/{beehiveId}/add", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
		if err != nil {
			utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
			return
		}
		handlers.AddBeehive(w, r, dbPool, beehiveId)
	}))

	// mux.HandleFunc("POST /beehive/{beehiveId}/sensor-data/add", authentication.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
	// 	beehiveId, err := strconv.Atoi(r.PathValue("beehiveId"))
	// 	if err != nil {
	// 		utils.SendErrorResponse(w, "Invalid Beehive id", http.StatusBadRequest)
	// 		return
	// 	}
	// 	handlers.AddSensorData(w, r, dbPool, beehiveId)
	// }))
	mux.HandleFunc("POST /beehive/{beehiveId}/sensor-data/add", DockerBridgeAuth(func(w http.ResponseWriter, r *http.Request) {
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
