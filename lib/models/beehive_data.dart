class BeehiveData {
  final double temperature;
  final double weight;
  final double humidity;
  final double ppm;
  final double battery;

  BeehiveData(
      {required this.temperature,
      required this.weight,
      required this.humidity,
      required this.ppm,
      required this.battery});

  Map<String, double> toMap() {
    return {
      'temperature': temperature,
      'weight': weight,
      'humidity': humidity,
      'ppm': ppm,
      'battery': battery,
    };
  }
}
