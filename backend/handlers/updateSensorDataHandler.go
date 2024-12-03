package handlers

import (
	"net/http"

	"github.com/buungoo/Beehive-interface/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateSensorData(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool, beehive_id int, sensor string) {
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)

}
