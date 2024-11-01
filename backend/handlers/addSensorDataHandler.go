package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Error while acquiring connection from the database pool!!")
	}
	defer conn.Release()

	// Fetch userid
	userId, err := utils.GetUserId(conn.Conn(), username)
	if err != nil {
		log.Println("Error fetching user id, err: ", err)
		utils.SendErrorResponse(w, "Error fetching user id", http.StatusInternalServerError)
		return
	}

	// Verify the beehive exists and that the user has access to said beehive
	beehiveExists, err := utils.VerifyBeehiveId(conn.Conn(), beehiveId, userId)
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

	// To store data if input is multiple sensors
	var inputArray []models.SensorData

	// Put input in byte slice
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading req body, err: ", err, " at time: ", time.Now().Format("2006-01-02 15:04:05"))
		utils.SendErrorResponse(w, "Error decoding payload", http.StatusBadRequest)
		return
	}

	const sqlQueryInsertNewData = `INSERT INTO sensor_data (sensor_id, beehive_id, value, time) VALUES($1, $2, $3, $4)`

	// Check if its a single or multiple input and handle each case accordingly
	if len(reqBody) > 0 && reqBody[0] == '[' {
		if err := json.Unmarshal(reqBody, &inputArray); err != nil {
			log.Println("Error decoding payload, err: ", err, " at time: ", time.Now().Format("2006-01-02 15:04:05"))
			utils.SendErrorResponse(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		for _, data := range inputArray {
			_, err = conn.Exec(context.Background(), sqlQueryInsertNewData, data.SensorID, data.BeehiveID, data.Value, data.Time)
			if err != nil {
				log.Println("Error inserting data, err: ", err, " at time: ", time.Now().Format("2006-01-02 15:04:05"))
				utils.SendErrorResponse(w, "Error inserting data", http.StatusInternalServerError)
				return
			}
		}

		utils.SendJSONResponse(w, "Data succesfully added", http.StatusOK)
		return

	} else {
		var inputObject models.SensorData
		if err := json.Unmarshal(reqBody, &inputObject); err != nil {
			log.Println("Error decoding payload, err: ", err, " at time: ", time.Now().Format("2006-01-02 15:04:05"))
			utils.SendErrorResponse(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		_, err = conn.Exec(context.Background(), sqlQueryInsertNewData, inputObject.SensorID, inputObject.BeehiveID, inputObject.Value, inputObject.Time)
		if err != nil {
			log.Println("Error inserting data, err: ", err, " at time: ", time.Now().Format("2006-01-02 15:04:05"))
			utils.SendErrorResponse(w, "Error inserting data", http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponse(w, "Data successfully added", http.StatusOK)
		return
	}

}

// func dataValidation() error {

// 	return nil
// }
