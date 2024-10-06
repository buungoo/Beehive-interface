package handlers

import (
	"fmt"
	"net/http"
)

func GetSensorData(w http.ResponseWriter, beehive_id string, r *http.Request, sensor string) {
	fmt.Println("Getting sensordata")
}

