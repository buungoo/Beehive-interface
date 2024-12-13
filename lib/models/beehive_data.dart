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
      this.battery = 98.0});

  Map<String, double> toMap() {
    return {
      'temperature': temperature,
      'weight': weight,
      'humidity': humidity,
      'ppm': ppm,
      'battery': 98.0,
    };
  }
}
