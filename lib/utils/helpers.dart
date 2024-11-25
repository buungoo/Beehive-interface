import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';

bool isIOS(BuildContext context) {
  return Theme.of(context).platform == TargetPlatform.iOS;
}

bool isAndroid(BuildContext context) {
  return Theme.of(context).platform == TargetPlatform.android;
}

bool isDarkMode() {
  var brightness =
      SchedulerBinding.instance.platformDispatcher.platformBrightness;
  return brightness == Brightness.dark;
}

String formatHexString(String hex) {
  // Remove leading zeros and split into pairs of two characters.
  var pairs = List.generate(
      hex.length ~/ 2, (i) => hex.substring(i * 2, i * 2 + 2).toUpperCase());
  // Join the pairs with a colon separator.
  return pairs.join(':');
}

List<Color> generateColorsFromString(String input) {
  int hash = input.hashCode;
  List<Color> colors = [];

  for (int i = 0; i < 3; i++) {
    // Extract RGB components and add an offset for variation
    int red = ((hash >> (i * 8) & 0xFF) + (i * 50)) % 256;
    int green = ((hash >> (i * 8 + 8) & 0xFF) + (i * 75)) % 256;
    int blue = ((hash >> (i * 8 + 16) & 0xFF) + (i * 100)) % 256;

    // Rotate the colors for variation
    if (i % 3 == 1) {
      int temp = red;
      red = green;
      green = blue;
      blue = temp;
    } else if (i % 3 == 2) {
      int temp = blue;
      blue = red;
      red = green;
      green = temp;
    }

    // Ensure the colors are vibrant by avoiding greys (where RGB values are close)
    if ((red - green).abs() < 30 &&
        (green - blue).abs() < 30 &&
        (blue - red).abs() < 30) {
      red = (red + 128) % 256; // Push the color out of the grey range
    }

    double opacity = 0.6 + (i * 0.3);

    colors.add(Color.fromRGBO(red, green, blue, opacity));
  }

  return colors;
}

class TimeScaleNotifier extends ChangeNotifier {
  String _timeScale;

  TimeScaleNotifier(this._timeScale);

  String get timeScale => _timeScale;

  void updateTimeScale(String newTimeScale) {
    if (_timeScale != newTimeScale) {
      _timeScale = newTimeScale;
      notifyListeners(); // Notify consumers about the change
    }
  }
}

Duration parseDuration(String timeString) {
  final parts = timeString.split(' '); // Split into parts: ["1", "Day"]
  final int value = int.parse(parts[0]); // Convert "1" to an integer
  final String unit = parts[1].toLowerCase(); // Convert "Day" to lowercase

  switch (unit) {
    case 'day':
      return Duration(days: value);
    case 'week':
      return Duration(days: value * 7);
    case 'month':
      return Duration(days: value * 30); // Approximate 1 month as 30 days
    default:
      throw ArgumentError('Unsupported time unit: $unit');
  }
}
