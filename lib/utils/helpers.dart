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
