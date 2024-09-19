import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import '../utils/helpers.dart';

PreferredSizeWidget? getNavigationBar(
    {required BuildContext context, required String title}) {
  return isIOS(context)
      ? CupertinoNavigationBar(middle: Text(title))
      : AppBar(title: Text(title));
}
