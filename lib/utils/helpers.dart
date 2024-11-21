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
