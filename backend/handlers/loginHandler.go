package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/buungoo/Beehive-interface/authentication"
	"github.com/buungoo/Beehive-interface/models"
	"github.com/buungoo/Beehive-interface/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	// Decode input to a user struct
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepared statement
	const sqlQueryGetPassword = `SELECT password FROM users WHERE username=$1`

	// Acuire connection from the connection pool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		utils.LogFatal("Error while acquiring connection from the database pool: ", err)
	}
	defer conn.Release()

	var password []byte
	// Fetch password from db if user exists
	err = conn.QueryRow(context.Background(), sqlQueryGetPassword, user.Username).Scan(&password)
	if err != nil {
		if err == pgx.ErrNoRows {
			utils.SendErrorResponse(w, "Username does not exists", http.StatusBadRequest)
			return
		}
		utils.LogError("Error fetching password", err)
		utils.SendErrorResponse(w, "Error retreiving user", http.StatusInternalServerError)
		return
	}

	// Compare the provided password with the stored hash
	// If the passwords match, create and return a JWT to the user
	err = bcrypt.CompareHashAndPassword(password, []byte(user.Password))
	if err != nil {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusBadRequest)
		return
	} else {
		tokenString, err := authentication.CreateToken(user.Username)
		if err != nil {
			utils.SendErrorResponse(w, "Error creating token", http.StatusInternalServerError)
		}
		// Return the token in the response
		response := map[string]string{
			"message": "User validated",
			"token":   tokenString,
		}
		utils.SendJSONResponse(w, response, http.StatusOK)
		return
	}
}
