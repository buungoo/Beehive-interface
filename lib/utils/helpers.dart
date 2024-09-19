import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';

bool isIOS(BuildContext context) {
  return Theme.of(context).platform == TargetPlatform.iOS;
}

bool isDarkMode() {
  var brightness =
      SchedulerBinding.instance.platformDispatcher.platformBrightness;
  return brightness == Brightness.dark;
}
