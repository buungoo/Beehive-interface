import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import '../utils/helpers.dart';

PreferredSizeWidget? getNavigationBar(
    {required BuildContext context, required String title, Color? bgcolor}) {
  return isIOS(context)
      ? CupertinoNavigationBar(
          middle: Text(title),
          backgroundColor: bgcolor,
        )
      : AppBar(
          title: Text(title),
          backgroundColor: bgcolor,
        );
}
