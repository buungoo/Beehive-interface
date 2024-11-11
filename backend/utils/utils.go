package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"`
}

// This is used for every successfull request
func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		SendErrorResponse(w, "Could not encode JSON response", http.StatusInternalServerError)
	}
}

// This is used for every error response
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Genereate error message in JSON format
	errorResponse := ErrorResponse{
		Error: message,
		Code:  statusCode,
	}

	json.NewEncoder(w).Encode(errorResponse)

}

// Returns the userid
func GetUserId(conn *pgx.Conn, username string) (int, error) {
	const sqlQueryFetchUserID = `SELECT id FROM users WHERE username=$1`

	var userID int
	// Fetch userid for user
	err := conn.QueryRow(context.Background(), sqlQueryFetchUserID, username).Scan(&userID)
	if err != nil {
		return userID, err
	}
	return userID, nil

}

// Veryfies the provided beehive_id exists in the database
func VerifyBeehiveId(conn *pgx.Conn, beehiveId int, userId int) (bool, error) {
	const sqlQueryCheckBeehive = `SELECT EXISTS(SELECT 1 FROM user_beehive WHERE beehive_id=$1 AND user_id=$2)`

	// Verify the beehive ID exists
	var exists bool
	err := conn.QueryRow(context.Background(), sqlQueryCheckBeehive, beehiveId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
