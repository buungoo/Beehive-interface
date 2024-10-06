package handlers

import (
	"net/http"
	"fmt"
)

func AddSensorData(w http.ResponseWriter, r *http.Request, conn *pgx.Conn, beehive_id string,  sensor string) {

	fmt.Println("Adding sensordata")
}

