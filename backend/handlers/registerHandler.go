package handlers

import (
	"beehive_api/models"
	"beehive_api/utils"
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	// Decode input to a user structure
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Instantiate hashfunction
	hash := sha256.New()

	// Coerse string to bytes and add to hash
	hash.Write([]byte(user.Password))

	// Complete the hash and get the hashed password as bytes 
	hashedPw := hash.Sum(nil)

	// Prepared statements 
	const sqlQueryCheckUsername = `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	const sqlQueryInsertNewUSer = `INSERT INTO users (username, password) VALUES($1, $2)`

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err!=nil {
	 	log.Fatal("Error while acquiring connection from the database pool!")
	} 
	defer conn.Release()

	// Check if the username already exists
    var exists bool
    err = dbPool.QueryRow(context.Background(), sqlQueryCheckUsername , user.Username).Scan(&exists)
    if err != nil {
        log.Println("Error checking username:", err)
        utils.SendErrorResponse(w, "Error registering user", http.StatusInternalServerError)
        return
    }

	if exists {
		log.Println("Username already exists, err:", err)
		utils.SendErrorResponse(w, "User already exists", http.StatusConflict)
		return
	}
	
	// Insert username and password
	_, err = conn.Exec(context.Background(), sqlQueryInsertNewUSer, user.Username, hashedPw)
	if err != nil {
		log.Println("Error registering user, error: ", err)
		utils.SendErrorResponse(w, "Error registering user", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "User registered",
		"Code":   "200",
	}
	utils.SendJSONResponse(w, response, http.StatusOK) 
	
}