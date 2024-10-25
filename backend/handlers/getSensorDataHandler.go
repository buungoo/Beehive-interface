package handlers

import (
	"beehive_api/utils"
	"beehive_api/models"
	"context"
	"net/http"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, sensorType string) {
	
	var sensorId int
	var value float64
	var time string
	
	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer conn.Release()

	// First query to get the sensor ID for the beehive
	err = conn.QueryRow(context.Background(), "SELECT id FROM sensors WHERE beehive_id=$1 AND type=$2", beehiveId, sensorType).Scan(&sensorId)
	if err != nil {
		utils.SendErrorResponse(w, "Error fetching sensor ID", http.StatusInternalServerError)
		return
	}

	// Query to fetch the latest sensor data using SQL function
	err = conn.QueryRow(context.Background(), "SELECT value, time FROM fetch_latest_sensor_data_for_beehive($1, $2)", beehiveId, sensorType).Scan(&value, &time)
	if err != nil {
		utils.SendErrorResponse(w, "Error fetching sensor data", http.StatusInternalServerError)
		return
	}

	dataResponse := models.SensorData{
		SensorID:  sensorId,
		BeehiveID: beehiveId,
		Value:     value,
		Time:      time,
	}

	// Send the response as JSON
	utils.SendJSONResponse(w, dataResponse, http.StatusOK)
	
}
