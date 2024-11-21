# Beehive-interface

## Mockup

![image](https://github.com/user-attachments/assets/3632702c-9ab6-484f-a090-85d735df5c2d)

https://www.figma.com/design/G8Y8VM1ecaz6kD7L0IbYhV/Beehive-app-mockup?node-id=0-1&t=u2Z12erQDL40nr75-1

## Data flowchart

![image](https://github.com/user-attachments/assets/b23924fe-cee5-4670-a239-d57f16a86397)

https://lucid.app/lucidchart/56da49fe-bf38-427f-9650-d7c6f5b29474/edit?viewport_loc=-478%2C12%2C3328%2C1650%2C0_0&invitationId=inv_b8127540-69cb-4d87-9d9b-97718cfe54a9

## App structure

```
lib
├── main.dart
├── models
│   └── README.md
|    This folder contains data objects used in the app. These models define the structure of the app's data.
├── providers
│   └── README.md
|    This folder contains state management logic using Flutter's `provider` package.
|    Providers in this folder manage application state, notifying widgets when data changes.
|    Typical provider classes extend a provider, allowing reactive updates across the app.
├── services
│   └── README.md
|    This folder contains service classes responsible for handling business logic and data fetching.
|    Services may interact with APIs, local storage, or simulate real-time data updates.
├── utils
│   └── README.md
|    This folder contains reusable utility functions that can be used across multiple parts of the app.
├── views
│   └── README.md
|    This folder contains the UI screens (pages) of the app. Each view represents a distinct part of the app's user interface.
└── widgets
    └── README.md
     This folder contains reusable UI components that can be used across multiple views in the app.
```


## Web Api

| **Route**                                   | **HTTP Method** | **Name**                     | **Required Data**                                                                 |
|---------------------------------------------|-----------------|------------------------------|-----------------------------------------------------------------------------------|
| `/register`                                 | POST            | Register                     | - User registration data (e.g., username, password)                              |
| `/login`                                    | POST            | Login                        | - User credentials (e.g., username, password)                                    |
| `/beehive/{beehiveId}/status`               | GET             | Get Beehive Status            | - Beehive ID (integer)                                                           |
| `/beehive/add`                              | POST            | Add Beehive                   | - Beehive data (e.g., name, location, etc.)                                      |
| `/beehive/{beehiveId}/sensor-data/add`      | POST            | Add Sensor Data               | - Beehive ID (integer) <br> - Sensor data (e.g., temperature, humidity, etc.)    |
| `/beehive/list`                             | GET             | Get Beehive List              | - None                                                                            |
| `/beehive/{beehiveId}/sensor-data/{startDate}/{endDate}` | GET | Get Sensor Data by Date Range | - Beehive ID (integer) <br> - Start date and end date (strings in date format)   |
| `/beehive/{beehiveId}/sensor-data/average/{startDate}/{endDate}` | GET | Get Average Sensor Data by Date Range | - Beehive ID (integer) <br> - Start date and end date (strings in date format)   |
| `/beehive/{beehiveId}/sensor-data/latest`    | GET             | Get Latest Sensor Data        | - Beehive ID (integer)                                                           |
| `/beehive/{beehiveId}/{sensorType}/latest`   | GET             | Get Latest Sensor Type Data   | - Beehive ID (integer) <br> - Sensor type (string)                               |
| `/test`                                     | POST            | Test Authentication           | - None                                                                            |
