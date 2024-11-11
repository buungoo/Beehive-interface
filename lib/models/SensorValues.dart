import 'dart:convert';

class SensorValues {
  final int sensor_id;
  final int beehive_id;
  final double value;
  final DateTime time;

  SensorValues(
      {required this.sensor_id,
      required this.beehive_id,
      required this.value,
      required this.time});

  factory SensorValues.fromJson(Map<String, dynamic> json) {
    return SensorValues(
        sensor_id: json['sensor_id'],
        beehive_id: json['beehive_id'],
        value: json['value'],
        time: DateTime.parse(json['time']));
  }

  // turn json array to list of SensorValues
  static List<SensorValues> fromJsonList(String json) {
    final parsed = jsonDecode(json);
    final values = parsed.map((item) => SensorValues.fromJson(item)).toList();
    // Cast to SendorValues to infer type
    return values.cast<SensorValues>();
  }
}
