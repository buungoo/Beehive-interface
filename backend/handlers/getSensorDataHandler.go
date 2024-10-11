package handlers

import (
	"beehive_api/utils"
	"net/http"
	"github.com/jackc/pgx/v5"
)

func GetSensorData(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, beehive_id string,  sensor string) {
	
	utils.SendErrorResponse(w, "Under development", http.StatusNotFound)
}

