import 'dart:async';
import 'package:beehive/models/beehive_data.dart';
import 'dart:convert';
import 'package:http/http.dart';
import 'package:beehive/config.dart' as config;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';
import 'package:beehive/models/SensorValues.dart';

import 'package:beehive/utils/helpers.dart';

class TestSensorData {
  final String sensorId;
  final String beehiveId;
  final String sensorType;
  final double value;
  final DateTime time;

  TestSensorData({
    required this.sensorId,
    required this.beehiveId,
    required this.sensorType,
    required this.value,
    required this.time,
  });

  factory TestSensorData.fromJson(Map<String, dynamic> json) {
    return TestSensorData(
      sensorId: json['sensor_id'],
      beehiveId: json['beehive_id'],
      sensorType: json['sensor_type'],
      value: num.tryParse(json['value'])?.toDouble() ?? 0.0,
      time: DateTime.parse(json['time']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'sensor_id': sensorId,
      'beehive_id': beehiveId,
      'sensor_type': sensorType,
      'value': value,
      'time': time.toIso8601String(),
    };
  }
}

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

        final List<dynamic> data = json.decode(response.body);
        //final test = data.map((item) => item as Map<String, dynamic>).toList();

        int tempIndex = data
            .indexWhere((element) => element['sensor_type'] == 'temperature');
        int weightIndex =
            data.indexWhere((element) => element['sensor_type'] == 'weight');
        int humidityIndex =
            data.indexWhere((element) => element['sensor_type'] == 'humidity');
        int ppmIndex =
            data.indexWhere((element) => element['sensor_type'] == 'oxygen');

        int batteryIndex =
            data.indexWhere((element) => element['sensor_type'] == 'battery');

        // set value to  0 if index not found
        double temp = tempIndex != -1
            ? data[tempIndex]['value'].toDouble()
            : 0.0; // Temperature in °C
        double weight = weightIndex != -1
            ? data[weightIndex]['value'].toDouble()
            : 0.0; // Weight in grams or smt
        double humidity = humidityIndex != -1
            ? data[humidityIndex]['value'].toDouble()
            : 0.0; // Humidity in %
        double ppm = ppmIndex != -1
            ? data[ppmIndex]['value'].toDouble()
            : 0.0; // Particles Per Million (PPM)
        double battery = batteryIndex != -1
            ? data[batteryIndex]['value'].toDouble()
            : 0.0; // Battery in %

        yield BeehiveData(
            temperature: temp,
            weight: weight,
            humidity: humidity,
            ppm: ppm,
            battery: battery);
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

  Future<Map<String, dynamic>> fetchBeehiveIssueStatus(String beehiveId) async {
    var path = 'beehive/$beehiveId/status';

    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('token');

      final uri = Uri.parse('${config.BackendServer}/$path');

      var response = await get(uri, headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      });

      return json.decode(response.body);
    } catch (e) {
      print("Error fetching beehive issue status: $e");
      return {};
    }
  }

  Future<List<dynamic>> fetchBeehiveIssueStatusesList(String beehiveId) async {
    var path = '/beehive/$beehiveId/status/list';

    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('token');

      final uri = Uri.parse('${config.BackendServer}/$path');

      var response = await get(uri, headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
        'Authorization': 'Bearer $token',
      });

      final data = json.decode(response.body);
      return data;
    } catch (e) {
      print("Error fetching beehive issue statuses: $e");
      return [];
    }
  }

  Stream<String> getBeehiveSensorData(String type) async* {
    yield "Hello";
  }
}
