import 'dart:async';
import 'package:beehive/models/beehive_data.dart';
import 'dart:math';
import 'package:beehive/config.dart' as config;

import '../models/beehive_data.dart';

class BeehiveDataProvider {
  // Simulates a stream of nullable temperature data
  Stream<BeehiveData> getBeehiveDataStream() async* {
    var random = Random();
    int temp = 35; // Initial temperature in °C
    int weight = 50; // Initial hive weight in kg
    int humidity = 60; // Initial humidity percentage
    int ppm = 416; // Initial CO2 concentration in ppm

    bool init = false;

    while (true) {
      if (init) {
        await Future.delayed(config.refreshRate);
      } else {
        init = true;
      }

      print("Fetch Data");
      // Simulate small fluctuations in temperature (±0.5°C)
      temp += random.nextInt(3) - 1;

      // Simulate weight change, gradual increase or decrease (±1 kg)
      weight += random.nextInt(3) - 1;

      // Simulate humidity change (±2%)
      humidity += random.nextInt(5) - 2;

      // Simulate CO2 concentration change (±5 ppm)
      ppm += random.nextInt(2) - 1;

      // Ensure realistic bounds
      //temp = temp.clamp(30, 40); // Keep within realistic hive temperature range
      //weight = weight.clamp(45, 60); // Hive weight fluctuation range
      //humidity = humidity.clamp(50, 80); // Humidity percentage range
      //ppm = ppm.clamp(400, 450); // CO2 concentration range

      yield BeehiveData(
          temperature: temp, weight: weight, humidity: humidity, ppm: ppm);
    }
  }

  Stream<String> getBeehiveSensorData(String type) async* {

    yield "Hello";

  }
}
