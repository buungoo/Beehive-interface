package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	// Fetch userid
	userId, err := utils.GetUserId(conn.Conn(), username)
	if err != nil {
		utils.LogError("Error fetching user id, err: ", err)
		utils.SendErrorResponse(w, "Error fetching user id", http.StatusInternalServerError)
		return
	}

	// Verify the beehive exists and that the user has access to said beehive
	beehiveExists, err := utils.VerifyBeehiveId(conn.Conn(), beehiveId, userId)
	if err != nil {
		utils.LogError("Error finding beehive, err: ", err)
		utils.SendErrorResponse(w, "Error finding beehive", http.StatusInternalServerError)
		return
	}

	if !beehiveExists {
		utils.LogError("Error, beehive doesnt exists", errors.New("beehives don't exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	// To store data if input is multiple sensors
	var inputArray []models.SensorData

	// Put input in byte slice
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		utils.LogError("Error reading req body, err: ", err)
		utils.SendErrorResponse(w, "Error decoding payload", http.StatusBadRequest)
		return
	}

	const sqlQueryInsertNewData = `INSERT INTO sensor_data (sensor_id, beehive_id, sensor_type, value, time) VALUES($1, $2, $3, $4, $5)`

	// Check if its a single or multiple input and handle each case accordingly
	if len(reqBody) > 0 && reqBody[0] == '[' {
		if err := json.Unmarshal(reqBody, &inputArray); err != nil {
			utils.LogError("Error decoding payload, err: ", err)
			utils.SendErrorResponse(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		for _, data := range inputArray {
			//err := dataValidation(data)
			IsValid, message := data.VerifyInputData()
			if !IsValid {
				UpdateBeehiveStatusOnAdd(w, r, dbPool, beehiveId, message, data)
				utils.LogInfo(message)

			}
			_, err = conn.Exec(context.Background(), sqlQueryInsertNewData, data.SensorID, data.BeehiveID, data.SensorType, data.Value, data.Time)
			if err != nil {
				utils.LogError("Error inserting data, err: ", err)
				utils.SendErrorResponse(w, "Error inserting data", http.StatusInternalServerError)
				return
			}
		}

		utils.SendJSONResponse(w, "Data succesfully added", http.StatusOK)
		return

	} else {
		var inputObject models.SensorData
		if err := json.Unmarshal(reqBody, &inputObject); err != nil {
			utils.LogError("Error decoding payload, err: ", err)
			utils.SendErrorResponse(w, "Invalid payload", http.StatusBadRequest)
			return
		}
		//err := dataValidation(inputObject)
		IsValid, message := inputObject.VerifyInputData()
		if !IsValid {
			UpdateBeehiveStatusOnAdd(w, r, dbPool, beehiveId, message, inputObject)
			utils.LogInfo(message)

		}
		_, err = conn.Exec(context.Background(), sqlQueryInsertNewData, inputObject.SensorID, inputObject.BeehiveID, inputObject.SensorType, inputObject.Value, inputObject.Time)
		if err != nil {
			utils.LogError("Error inserting data, err: ", err)
			utils.SendErrorResponse(w, "Error inserting data", http.StatusInternalServerError)
			return
		}

		utils.SendJSONResponse(w, "Data successfully added", http.StatusOK)
		return
	}

}
