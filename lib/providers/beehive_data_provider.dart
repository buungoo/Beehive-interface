import 'dart:async';

class BeehiveDataProvider {
  // Simulates a stream of nullable temperature data
  Stream<int?> getTemperatureStream() {
    // This will periodically fetch data and notify listeners
    return Stream.periodic(
      const Duration(seconds: 1),
      (count) => 20 + count,
    );
  }
}
