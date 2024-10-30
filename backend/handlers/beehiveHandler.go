package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"log"
)

func GetBeehiveStatus(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
	return
}

func GetBeehiveList(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	// Retrieve the username from the request context
	username := r.Context().Value("username").(string)

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 log.Fatal("Error while acquiring connection from the database pool!!")
	} 
	defer conn.Release()
	
	const sqlQueryFetchUserID = `SELECT id FROM users WHERE username=$1`

	var userID int
	// Fetch userid for user
	err = conn.QueryRow(context.Background(), sqlQueryFetchUserID, username).Scan(&userID)
	if err != nil {
		log.Println("Error fetching user id, err: ", err)
		utils.SendErrorResponse(w, "Error fetching user id", http.StatusInternalServerError)
		return
	}


	const sqlQueryFetchAllBeehives = `SELECT id, name, user_id FROM beehives WHERE user_id=$1`

	rows, err := conn.Query(context.Background(), sqlQueryFetchAllBeehives, userID)
	if err != nil {
		log.Println("Error fetching all beehives, err: ", err)
		utils.SendErrorResponse(w, "Error fetching all beehives", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Put all data into struct before returning to client
	beehives, err := iterateBeehives(rows)
	if err != nil {
		log.Println("Error iterating data, err: ", err)
		utils.SendErrorResponse(w, "Error iterating data", http.StatusInternalServerError)
		return
	}

	// Return the data
	utils.SendJSONResponse(w, beehives, http.StatusOK)
	return
	
}

func iterateBeehives(rows pgx.Rows) ([]models.Beehives, error){
	// Slice to hold the data from returned rows
	var dataResponse []models.Beehives

	for rows.Next() {
		var beehive models.Beehives
		if err := rows.Scan(&beehive.Id, &beehive.Name, &beehive.UserID); err !=nil {
			return dataResponse, err
		}
		dataResponse = append(dataResponse, beehive)
	}
	if err := rows.Err(); err != nil {
		return dataResponse, err
	}

	return dataResponse, nil
}
