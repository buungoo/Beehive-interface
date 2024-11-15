package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect a beehive to a user by using the sensor-cards mac address
func AddBeehive(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool, beehiveId int) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Retrieve the macaddress from http-body
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	macAddr := struct {
		Addr string `json:"macaddress"`
	}{}
	err := decoder.Decode(&macAddr)
	if err != nil {
		utils.LogError("Could not decode macaddress", err)
		utils.SendErrorResponse(w, "Could not decode macaddress", http.StatusInternalServerError)
		return
	}
	utils.LogInfo("Macaddress is: " + macAddr.Addr)
	macAddress := macAddr.Addr

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
	beehiveExists, err := utils.VerifyBeehive(conn.Conn(), beehiveId, macAddress)
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

	const sqlQueryAddBeehive = `INSERT INTO user_beehive (user_id, beehive_id) VALUES ($1, $2)`

	_, err = conn.Exec(context.Background(), sqlQueryAddBeehive, userId, beehiveId)
	if err != nil {
		utils.LogError("Error adding beehive to user, error: ", err)
		utils.SendErrorResponse(w, "Error adding beehive to user", http.StatusBadRequest)
		return
	}

	utils.SendJSONResponse(w, "Beehive added to user", http.StatusOK)

}

func GetBeehiveStatus(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)

}

func UpdateBeehiveStatus(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

}

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
