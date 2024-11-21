import 'dart:async';
import 'package:beehive/models/beehive_data.dart';
import 'dart:convert';
import 'package:http/http.dart';
import 'package:beehive/config.dart' as config;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:intl/intl.dart';

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

        print(response.body);
        final data = json.decode(response.body);

        double temp = data[0]['value']; // Temperature in °C
        double weight = data[1]['value']; // Weight in grams
        double humidity = data[2]['value']; // Humidity in %
        double ppm = data[3]['value']; // Particles Per Million (PPM)

        yield BeehiveData(
            temperature: temp, weight: weight, humidity: humidity, ppm: ppm);
      } catch (e) {
        print(e);
      }
    }
  }

  Stream<String> getBeehiveDataChartStream(
      String beehiveid, String sensor) async* {
    while (true) {
      try {
        final prefs = await SharedPreferences.getInstance();
        final token = prefs.getString('token');

        var date1 = DateTime.now().subtract(Duration(days: 10));
        var date2 = DateTime.now();

        // parse it to string in 2006-01-02 format
        final formatter = DateFormat('yyyy-MM-dd');
        String formattedDate1 = formatter.format(date1);
        String formattedDate2 = formatter.format(date2);

        print("formattedDate1: $formattedDate1");
        print("formattedDate2: $formattedDate2");

        //TODO: Make sure to only extract the "type" from the request
        final uri = Uri.parse(
            '${config.BackendServer}/beehive/$beehiveid/sensor-data/$formattedDate1/$formattedDate2');

        print(uri);

        var response = await get(uri, headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        });
        yield response.body.toString();
      } catch (e) {
        print(e);
      }
      await Future.delayed(config.refreshRate);
    }
  }

  Stream<String> getBeehiveSensorData(String type) async* {
    yield "Hello";
  }
}
