package handlers

import (
	"beehive_api/utils"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateSensorData(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool, beehive_id int,  sensor string) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)

}

