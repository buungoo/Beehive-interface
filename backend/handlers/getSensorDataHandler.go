package handlers

import (
	"beehive_api/utils"
	"beehive_api/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)


func GetLatestSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer conn.Release()

	beehiveExists, err := verifyBeehiveId(conn.Conn(), beehiveId)
	if err != nil {
		log.Println("Error finding beehive, err: ", err)
		utils.SendErrorResponse(w, "Error finding beehive", http.StatusInternalServerError)
		return
	}

	if !beehiveExists {
		log.Println("Error, beehive doesnt exists")
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
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

func GetLatestOfSensortype(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, sensorType string) {
	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer conn.Release()

	beehiveExists, err := verifyBeehiveId(conn.Conn(), beehiveId)
	if err != nil {
		log.Println("Error finding beehive, err: ", err)
		utils.SendErrorResponse(w, "Error finding beehive", http.StatusInternalServerError)
		return
	}

	if !beehiveExists {
		log.Println("Error, beehive doesnt exists")
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	// Query to find latest temprature for beehive_id
	const sqlQueryFetchTemperature = `SELECT sd.sensor_id, sd.beehive_id, sd.value, sd.time
	FROM sensor_data sd
	JOIN sensors s ON sd.sensor_id = s.id
	WHERE sd.beehive_id = $1
		AND s.type = $2
	ORDER BY sd.time DESC
	LIMIT 1;`


	// Store data in SensorData struct
	var dataResponse models.SensorData

	
	err = conn.QueryRow(context.Background(), sqlQueryFetchTemperature , beehiveId, sensorType).Scan(&dataResponse.SensorID, 
		&dataResponse.BeehiveID, &dataResponse.Value, &dataResponse.Time)
	if err != nil {
		log.Println("Error fetching latest temperature, err: ", err)
		utils.SendErrorResponse(w, "Error fetching temperature", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, dataResponse, http.StatusOK)

}

// Veryfies the provided beehive_id exists in the database
func verifyBeehiveId(conn *pgx.Conn, beehiveId int) (bool, error) {
	const sqlQueryCheckBeehive = `SELECT EXISTS(SELECT 1 FROM beehives WHERE id=$1)`

	// Verify the beehive ID exists
	var exists bool
	err := conn.QueryRow(context.Background(), sqlQueryCheckBeehive, beehiveId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
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
