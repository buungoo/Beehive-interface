import 'dart:async';
import 'package:beehive/models/beehive_data.dart';
import 'dart:convert';
import 'package:http/http.dart';
import 'package:beehive/config.dart' as config;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';
import 'package:beehive/models/SensorValues.dart';

import 'package:beehive/utils/helpers.dart';

class BeehiveDataProvider {
  // Simulates a stream of nullable temperature data
  Stream<BeehiveData> getBeehiveDataStream(String beehiveid) async* {
    bool init = false;

    while (true) {
      if (init) {
        await Future.delayed(config.refreshRate);
      } else {
        init = true;
      }

      try {
        final prefs = await SharedPreferences.getInstance();
        final token = prefs.getString('token');

        final uri = Uri.parse(
            '${config.BackendServer}/beehive/$beehiveid/sensor-data/latest');

        var response = await get(uri, headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        });

        final data = json.decode(response.body);

        double temp = data[0]['value']; // Temperature in °C
        double weight = data[1]['value']; // Weight in grams
        double humidity = data[2]['value']; // Humidity in %
        double ppm = data[3]['value']; // Particles Per Million (PPM)

        yield BeehiveData(
            temperature: temp, weight: weight, humidity: humidity, ppm: ppm);
      } catch (e) {
        //print(e);
      }
    }
  }

  Future<List<SensorValues>> fetchBeehiveDataChart(
      {required String beehiveId,
      required String sensor,
      String timescale = '1 Week'}) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('token');

      Duration timeRange = parseDuration(timescale);

      var date1 = DateTime.now().subtract(timeRange);
      var date2 = DateTime.now();

      // Format dates to "yyyy-MM-dd"
      final formatter = DateFormat('yyyy-MM-dd');
      String formattedDate1 = formatter.format(date1);
      String formattedDate2 = formatter.format(date2);

      final uri = Uri.parse(
        '${config.BackendServer}/beehive/$beehiveId/sensor-data/$formattedDate1/$formattedDate2',
      );

      print(uri);

      var response = await get(
        uri,
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        },
      );

      // using timescale, if 1 day. take all data for today.
      // else, take the average of the data for each day in the timescale

      final values = SensorValues.fromJsonList(response.body);

      print("INDAYS: ${timeRange.inDays}");
      if (timeRange.inDays == 1) {
        return values;
      } else {
        Map<String, List<SensorValues>> groupedData = {};
        for (var value in values) {
          String date = formatter.format(value.time);
          if (!groupedData.containsKey(date)) {
            groupedData[date] = [];
          }
          groupedData[date]!.add(value);
        }

        List<SensorValues> averagedValues = [];
        groupedData.forEach((date, values) {
          double avgValue = values.map((v) => v.value).reduce((a, b) => a + b) /
              values.length;

          averagedValues.add(SensorValues(
            sensor_id: values.first.sensor_id,
            beehive_id: values.first.beehive_id,
            value: avgValue,
            time: DateTime.parse(date),
          ));
        });

        return averagedValues;
      }

      // Return the response body as a string
      //return response.body.toString();
    } catch (e) {
      // Handle error gracefully
      print("Error fetching beehive data: $e");
      //return 'Error fetching data';
      return [];
    }
  }

  Future<dynamic> fetchBeehiveIssueStatuses(String beehiveId, bool readOnly) async {
    var path = '/beehive/$beehiveId/status';
    if (readOnly) {
      path += '/list';
    }

    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('token');

      final uri = Uri.parse('${config.BackendServer}/$path');

      var response = await get(uri, headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      });

      //final data = json.decode(response.body);

      //return data;
      if(readOnly) {
        return {
          "IssueId": 1,
          "SensorId": 5,
          "BeehiveId": 2,
          "SensorType": "temperature",
          "Description": "temperature is above 40.000000 Celsius",
          "Solved": false,
          "Read": false,
          "TimeOfError": "2025-12-01T05:15:10Z",
          "TimeRead": null
        };
      } else {
        return [
        {
          "IssueId": 1,
    "SensorId": 5,
    "beehive_id": 2,
    "SensorType": "temperature",
    "Description": "temperature is above 40.000000 Celsius",
    "Solved": false,
    "Read": true,
    "TimeOfError": "2025-12-01T05:15:10Z",
    "TimeRead": "2024-11-28T10:52:35.778848Z"
  },
    {
    "IssueId": 3,
    "SensorId": 6,
    "beehive_id": 2,
    "SensorType": "humidity",
    "Description": "temperature is above 60.000000 Celsius",
    "Solved": false,
    "Read": false,
    "TimeOfError": "2025-10-11T05:15:12Z",
    "TimeRead": null
    }];
      }


    } catch (e) {
      print("Error fetching beehive issue statuses: $e");
      return null;
    }



  }


  Stream<String> getBeehiveSensorData(String type) async* {
    yield "Hello";
  }
}
