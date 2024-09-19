import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../utils/helpers.dart';

Widget SharedListTile(
    {required BuildContext context,
    required Widget title,
    GestureTapCallback? onTap}) {
  return isIOS(context)
      ? CupertinoListTile(title: title, onTap: onTap)
      : ListTile(title: title, onTap: onTap);
}
