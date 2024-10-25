package handlers

import (
	"beehive_api/utils"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
}