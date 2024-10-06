package handlers

import (
	"fmt"
	"net/http"
)

func GetSensorData(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, beehive_id string,  sensor string) {
	fmt.Println("Getting sensordata")
}

