package handlers

import (
	"beehive_api/utils"
	"beehive_api/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"time"
)

func GetSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, sensorType string) {
	
	var sensorId int
	var value float64
	var time time.Time
	
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
		BeehiveID: beehiveId,
		SensorID:  sensorId,
		Value:     value,
		Time:      time,
	}

	// Send the response as JSON
	utils.SendJSONResponse(w, dataResponse, http.StatusOK)
	
}

func GetLatestSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int){
	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer conn.Release()

	const sqlQueryCheckBeehive = `SELECT EXISTS(SELECT 1 FROM beehives WHERE id=$1)`

	//Verify the beehive ID exists
	var exists bool
	err = conn.QueryRow(context.Background(), sqlQueryCheckBeehive , beehiveId).Scan(&exists)
	if err != nil {
		log.Println("Error checking beehive, err: ", err)
		utils.SendErrorResponse(w, "Error finding beehive", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Println("Beehive does not exist, err:", err)
		utils.SendErrorResponse(w, "Beehive does not exists", http.StatusNotFound)
		return
	}

	const sqlQueryGetLatestData = `SELECT DISTINCT ON (sensor_id) sensor_id, beehive_id, value, time
		FROM sensor_data
		WHERE beehive_id = $1
		ORDER BY sensor_id, time DESC;
		`
	// Fetch all data
	rows, err := conn.Query(context.Background(), sqlQueryGetLatestData, beehiveId)
	if err != nil {
		log.Println("Error fetching data", err)
		utils.SendErrorResponse(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	data, err := iterateData(rows)
	if err != nil {
		log.Println("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, data, http.StatusOK)
	return

}

func iterateData(rows pgx.Rows) ([]models.SensorData, error){
	// Slice to hold the data from returned rows
	var dataResponse []models.SensorData

	for rows.Next() {
		var data models.SensorData
		if err := rows.Scan(&data.SensorID, &data.BeehiveID, &data.Value, &data.Time); err !=nil {
			return dataResponse, err
		}
		dataResponse = append(dataResponse, data)
	}
	if err := rows.Err(); err != nil {
		return dataResponse, err
	}

	return dataResponse, nil
}
