package handlers

import (
	"net/http"
	"fmt"
)

func AddSensorData(w http.ResponseWriter, beehive_id string, r *http.Request, sensor string) {

	fmt.Println("Adding sensordata")
}

