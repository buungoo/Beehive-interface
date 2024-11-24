package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect a beehive to a user by using the sensor-cards mac address
func AddBeehive(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Retrieve the macaddress from http-body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	macAddrStruct := struct {
		Addr string `json:"macaddress"`
	}{}
	err := decoder.Decode(&macAddrStruct)
	if err != nil {
		utils.LogError("Could not decode macaddress", err)
		utils.SendErrorResponse(w, "Could not decode macaddress", http.StatusInternalServerError)
		return
	}
	utils.LogInfo("Macaddress is: " + macAddrStruct.Addr)

	// Acquire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool!!", errors.New("error while acquiring a connection from the pool"))
	}
	defer conn.Release()

	// Fetch userid
	userId, err := utils.GetUserId(conn.Conn(), username)
	if err != nil {
		utils.LogError("Error fetching user id, err: ", err)
		utils.SendErrorResponse(w, "Error fetching user id", http.StatusInternalServerError)
		return
	}

	// Verify the mac address is correct
	beehiveExists, err := utils.VerifyBeehive(conn.Conn(), macAddrStruct.Addr)
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

	const sqlQueryFetchBeehiveId = `SELECT id FROM beehives WHERE key = $1::macaddr8`

	var beehiveId int
	// Fetch beehiveId for beehive
	err = conn.QueryRow(context.Background(), sqlQueryFetchBeehiveId, macAddrStruct.Addr).Scan(&beehiveId)
	if err != nil {
		utils.LogError("This is the error", err)
		utils.LogError("Error, beehive doesnt exists", errors.New("beehive doesn't exist"))
		utils.SendErrorResponse(w, "Beehive doesn't exist", http.StatusNotFound)
		return
	}

	const sqlQueryAddBeehive = `INSERT INTO user_beehive (user_id, beehive_id) VALUES ($1, $2)`

	_, err = conn.Exec(context.Background(), sqlQueryAddBeehive, userId, beehiveId)
	if err != nil {
		utils.LogError("Error adding beehive to user, error: ", err)
		utils.SendErrorResponse(w, "Error adding beehive to user", http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, "Beehive added to user", http.StatusOK)

}

// Returns the beehive_status table which shows issues
func GetBeehiveStatus(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
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

	const sqlQueryFetchBeehiveStatus = `SELECT * FROM beehive_status 
										WHERE beehive_id=$1
										ORDER BY time_of_error DESC
										LIMIT 1; `

	// Hold the data
	var data models.BeehiveStatus

	// Fetch all data
	err = conn.QueryRow(context.Background(), sqlQueryFetchBeehiveStatus, beehiveId).Scan(&data.IssueId, &data.SensorId,
		&data.BeehiveId, &data.SensorType, &data.Description, &data.Solved, &data.Read, &data.TimeOfError, &data.TimeRead)
	if err != nil {
		utils.LogError("error reading beehivestatus: ", err)
		utils.SendErrorResponse(w, "error reading beehivestatus", http.StatusInternalServerError)
		return
	}

	if !data.Read {
		err = updateBeehiveStatusOnRead(dbPool, data)
		if err != nil {
			utils.LogError("error updating beehive_status: ", err)
		}
	}

	utils.SendJSONResponse(w, data, http.StatusOK)
	//utils.SendErrorResponse(w, "Under development", http.StatusNotFound)

}

// Updates the beehive_status table when a value outside of the limits has been receive from the sensors
func UpdateBeehiveStatusOnAdd(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int, statusMessage string, data models.SensorData) {

	// Acquire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	const sqlQueryUpdateBeehiveStatus = `INSERT INTO beehive_status (sensor_id, beehive_id, sensor_type, description, solved, read, time_of_error) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	// Insert username and password
	_, err = conn.Exec(context.Background(), sqlQueryUpdateBeehiveStatus, data.SensorID, data.BeehiveID, data.SensorType, statusMessage, false, false, data.Time)
	if err != nil {
		utils.LogError("Error updating status of beehive: ", err)
		utils.SendErrorResponse(w, "Error updating status of beehive", http.StatusBadRequest)
		return
	}
}

// Updates the read time after the issue has been read
func updateBeehiveStatusOnRead(dbPool *pgxpool.Pool, data models.BeehiveStatus) error {

	// Acquire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	const sqlQueryUpdateBeehiveStatus = `UPDATE beehive_status 
									SET read = $1, time_read = $2 
									WHERE issue_id = $3 AND beehive_id = $4 AND sensor_id = $5`

	// Update beehive_status
	_, err = conn.Exec(context.Background(), sqlQueryUpdateBeehiveStatus, true, time.Now(), data.IssueId, data.BeehiveId, data.SensorId)
	if err != nil {
		return err
	}

	return nil
}

// Returns a list of the beehives connected to the user
func GetBeehiveList(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Acquire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool!!", errors.New("error while acquiring a connection from the pool"))
	}
	defer conn.Release()

	// Fetch userid
	userId, err := utils.GetUserId(conn.Conn(), username)
	if err != nil {
		utils.LogError("Error fetching user id, err: ", err)
		utils.SendErrorResponse(w, "Error fetching user id", http.StatusInternalServerError)
	}

	// Fetch all beehives connected to the user
	const sqlQueryFetchAllBeehives = `SELECT b.id, b.name 
						FROM beehives b 
						JOIN user_beehive ub ON ub.beehive_id = b.id 
						WHERE ub.user_id=$1`

	rows, err := conn.Query(context.Background(), sqlQueryFetchAllBeehives, userId)
	if err != nil {
		utils.LogError("Error fetching all beehives, err: ", err)
		utils.SendErrorResponse(w, "Error fetching all beehives", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	beehives, err := iterateBeehives(rows)
	if err != nil {
		utils.LogError("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, beehives, http.StatusOK)

}

func iterateBeehives(rows pgx.Rows) ([]models.Beehives, error) {
	// Slice to hold the data from returned rows
	var dataResponse []models.Beehives

	for rows.Next() {
		var beehive models.Beehives
		if err := rows.Scan(&beehive.Id, &beehive.Name); err != nil {
			return dataResponse, err
		}
		dataResponse = append(dataResponse, beehive)
	}
	if err := rows.Err(); err != nil {
		return dataResponse, err
	}

	return dataResponse, nil
}
