package test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"beehive_api/models"
    "time"
)

// Inject harmful data into the database for testing

func InjectTestData(dbpool *pgxpool.Pool) error {

    // Insert test data into the users table
    var user1ID, user2ID int
    err := dbpool.QueryRow(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", "user1", "pass1").Scan(&user1ID)
    if err != nil {
        return fmt.Errorf("failed to insert user1: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", "user2", "pass2").Scan(&user2ID)
    if err != nil {
        return fmt.Errorf("failed to insert user2: %v", err)
    }

    // Insert test data into the beehives table
    var beehive1ID, beehive2ID int
    err = dbpool.QueryRow(context.Background(), "INSERT INTO beehives (name, user_id) VALUES ($1, $2) RETURNING id", "Beehive A", user1ID).Scan(&beehive1ID)
    if err != nil {
        return fmt.Errorf("failed to insert beehive1: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO beehives (name, user_id) VALUES ($1, $2) RETURNING id", "Beehive B", user2ID).Scan(&beehive2ID)
    if err != nil {
        return fmt.Errorf("failed to insert beehive2: %v", err)
    }

    // Insert test data into the sensors table
    var sensor1ID, sensor2ID, sensor3ID, sensor4ID, sensor5ID, sensor6ID, sensor7ID, sensor8ID int
    err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "temperature", beehive1ID).Scan(&sensor1ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor1: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "humidity", beehive1ID).Scan(&sensor2ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor2: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "weight", beehive1ID).Scan(&sensor3ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor3: %v", err)
    }
	err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "oxygen", beehive1ID).Scan(&sensor4ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor4: %v", err)
    }
	err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "temperature", beehive2ID).Scan(&sensor5ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor5: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "humidity", beehive2ID).Scan(&sensor6ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor6: %v", err)
    }
    err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "weight", beehive2ID).Scan(&sensor7ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor7: %v", err)
    }
	err = dbpool.QueryRow(context.Background(), "INSERT INTO sensors (type, beehive_id) VALUES ($1, $2) RETURNING id", "oxygen", beehive2ID).Scan(&sensor8ID)
    if err != nil {
        return fmt.Errorf("failed to insert sensor8: %v", err)
    }

    // Initialize the string slice with date-time values
    testTimes := []string{
        "2024-10-01 10:45:00", "2024-10-02 11:50:00", "2024-10-03 12:55:00", 
        "2024-10-04 12:45:00", "2024-10-04 10:30:00", "2024-10-05 11:55:00", 
        "2024-10-05 12:45:00", "2024-10-05 12:59:00",
    }

    // Slice to hold parsed time.Time values
    testTimesTimeformat := make([]time.Time, len(testTimes))

    
    layout := "2006-01-02 15:04:05"

    // Parse each time string into a time.Time value
    for i := 0; i < len(testTimes); i++ {
        parsedTime, err := time.Parse(layout, testTimes[i])
        if err != nil {
            fmt.Println("Error parsing time:", err)
            return err
        }
        testTimesTimeformat[i] = parsedTime
    }
        
	
	

    // Insert test data into the sensor_data table
    testData := []models.SensorData{
        {BeehiveID: beehive1ID, SensorID: sensor1ID, Value: 23.4, Time: testTimesTimeformat[0]},
        {BeehiveID: beehive1ID, SensorID: sensor2ID, Value: 55.8, Time: testTimesTimeformat[1]},
        {BeehiveID: beehive1ID, SensorID: sensor3ID, Value: 12.3, Time: testTimesTimeformat[2]},
		{BeehiveID: beehive1ID, SensorID: sensor4ID, Value: 23.4, Time: testTimesTimeformat[3]},
        {BeehiveID: beehive2ID, SensorID: sensor5ID, Value: 23.8, Time: testTimesTimeformat[4]},
        {BeehiveID: beehive2ID, SensorID: sensor6ID, Value: 54.2, Time: testTimesTimeformat[5]},
		{BeehiveID: beehive2ID, SensorID: sensor7ID, Value: 12.4, Time: testTimesTimeformat[6]},
        {BeehiveID: beehive2ID, SensorID: sensor8ID, Value: 23.8, Time: testTimesTimeformat[7]},
    }

    for _, data := range testData {
        insertSQL := `
            INSERT INTO sensor_data (sensor_id, beehive_id, value, time)
            VALUES ($1, $2, $3, $4)
        `
        _, err := dbpool.Exec(context.Background(), insertSQL, data.SensorID, data.BeehiveID, data.Value, data.Time)
        if err != nil {
            return fmt.Errorf("failed to insert test sensor data: %v", err)
        }
    }

    return nil
}

