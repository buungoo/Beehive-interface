package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/buungoo/Beehive-interface/models"
	"github.com/buungoo/Beehive-interface/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetLatestSensorData(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
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
		utils.LogError("Error, beehive doesnt exists", errors.New("beehive doesn't exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	const sqlQueryFetchLatestData = `SELECT DISTINCT ON (sensor_id) sensor_id, beehive_id, sensor_type, value, time
		FROM sensor_data
		WHERE beehive_id = $1
		ORDER BY sensor_id, time DESC;
		`
	// Fetch all data
	rows, err := conn.Query(context.Background(), sqlQueryFetchLatestData, beehiveId)
	if err != nil {
		utils.LogError("Error fetching data", err)
		utils.SendErrorResponse(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	data, err := iterateData(rows)
	if err != nil {
		utils.LogError("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, data, http.StatusOK)

}

func GetLatestOfSensortype(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, sensorType string) {
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
	}

	beehiveExists, err := utils.VerifyBeehiveId(conn.Conn(), beehiveId, userId)
	if err != nil {
		utils.LogError("Error finding beehive, err: ", err)
		utils.SendErrorResponse(w, "Error finding beehive", http.StatusInternalServerError)
		return
	}

	if !beehiveExists {
		utils.LogError("Error, beehive doesnt exists", errors.New("beehive doesn't exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	// Query to find latest value of a sensortype for beehive_id
	const sqlQueryFetchLatestSensorValueByType = `SELECT sensor_id, beehive_id, sensor_type, value, time
	FROM sensor_data 
	WHERE beehive_id = $1 AND sensor_type = $2
	ORDER BY time DESC
	LIMIT 1;`

	// Store data in SensorData struct
	var dataResponse models.SensorData

	err = conn.QueryRow(context.Background(), sqlQueryFetchLatestSensorValueByType, beehiveId, sensorType).Scan(&dataResponse.SensorID,
		&dataResponse.BeehiveID, &dataResponse.SensorType, &dataResponse.Value, &dataResponse.Time)
	if err != nil {
		utils.LogError("Error fetching latest sensorvalue, err: ", err)
		utils.SendErrorResponse(w, "Error fetching sensorvalue", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, dataResponse, http.StatusOK)

}

func GetDataByDate(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, date1 string, date2 string) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Verify and parse the input dates
	parsedDate1, parsedDate2, err := verifyDates(date1, date2)
	if err != nil {
		utils.LogError("Error parsing the dates: ", err)
		utils.SendErrorResponse(w, "Wrong format of the dates, or wrong order", http.StatusBadRequest)
	}

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
		utils.LogError("Error, beehive doesnt exists: ", errors.New("beehive doesnt exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	const sqlQueryFetchDataBetweenDates = `SELECT sensor_id, beehive_id, sensor_type, value, time
	FROM sensor_data 
	WHERE beehive_id=$1 AND time BETWEEN $2 AND $3
	ORDER BY time;
	`

	// Fetch all data
	rows, err := conn.Query(context.Background(), sqlQueryFetchDataBetweenDates, beehiveId, parsedDate1, parsedDate2)
	if err != nil {
		utils.LogError("Error fetching data", err)
		utils.SendErrorResponse(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	data, err := iterateData(rows)
	if err != nil {
		utils.LogError("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, data, http.StatusOK)

}

func GetAverageDataByDate(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, date1 string, date2 string) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Verify and parse the input dates
	parsedDate1, parsedDate2, err := verifyDates(date1, date2)
	if err != nil {
		utils.LogError("Error parsing the dates: ", err)
		utils.SendErrorResponse(w, "Wrong format of the dates, or wrong order", http.StatusBadRequest)
	}

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
		utils.LogError("Error, beehive doesnt exists: ", errors.New("beehive doesnt exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	const sqlQueryFetchAverageDataBetweenDates = `SELECT sensor_id,
		DATE(time) AS day,
		AVG(value) AS daily_average
		FROM sensor_data
		WHERE beehive_id = $1 
		AND time BETWEEN $2 AND $3  
		GROUP BY sensor_id, day
		ORDER BY sensor_id, day;
		`

	// Fetch all data
	rows, err := conn.Query(context.Background(), sqlQueryFetchAverageDataBetweenDates, beehiveId, parsedDate1, parsedDate2)
	if err != nil {
		utils.LogError("Error fetching data", err)
		utils.SendErrorResponse(w, "Error fetching data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	data, err := iterateData(rows)
	if err != nil {
		utils.LogError("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, data, http.StatusOK)

}

// Helper functions

func iterateData(rows pgx.Rows) ([]models.SensorData, error) {
	// Slice to hold the data from returned rows
	var dataResponse []models.SensorData

	for rows.Next() {
		var data models.SensorData
		if err := rows.Scan(&data.SensorID, &data.BeehiveID, &data.SensorType, &data.Value, &data.Time); err != nil {
			return dataResponse, err
		}
		dataResponse = append(dataResponse, data)
	}
	if err := rows.Err(); err != nil {
		return dataResponse, err
	}

	return dataResponse, nil
}

func verifyDates(date1 string, date2 string) (time.Time, time.Time, error) {

	layout := "2006-01-02"

	// Verify the date is in correct format
	parsedDate1, err := time.Parse(layout, date1)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error wrong format of date1: %w", err)
	}

	parsedDate2, err := time.Parse(layout, date2)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("error wrong format of date2: %w", err)
	}

	// Verify the dates are in correct order
	if parsedDate2.Before(parsedDate1) {

		return time.Time{}, time.Time{}, errors.New("date 2 is before date 1")
	}
	return parsedDate1, parsedDate2, nil

}
