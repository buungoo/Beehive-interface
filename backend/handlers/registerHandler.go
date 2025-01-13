package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/buungoo/Beehive-interface/models"
	"github.com/buungoo/Beehive-interface/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	// Decode input to a user struct
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash password with bcrypt
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogFatal("Error hashing password:", err)
	}

	// Prepared statements
	const sqlQueryCheckUsername = `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	const sqlQueryInsertNewUser = `INSERT INTO users (username, password) VALUES($1, $2)`

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	// Check if the username already exists
	var exists bool
	err = conn.QueryRow(context.Background(), sqlQueryCheckUsername, user.Username).Scan(&exists)
	if err != nil {
		utils.LogError("Error checking username:", err)
		utils.SendErrorResponse(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	if exists {
		utils.LogError("Username already exists, err:", err)
		utils.SendErrorResponse(w, "User already exists", http.StatusConflict)
		return
	}

	// Insert username and password
	_, err = conn.Exec(context.Background(), sqlQueryInsertNewUser, user.Username, hashedPW)
	if err != nil {
		utils.LogError("Error registering user, error: ", err)
		utils.SendErrorResponse(w, "Error registering user", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"message": "User registered",
		"Code":    "200",
	}
	utils.SendJSONResponse(w, response, http.StatusOK)

}
