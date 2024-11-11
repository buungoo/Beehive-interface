import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../utils/helpers.dart';

Widget SharedListTile(
    {required BuildContext context,
    required Widget title,
    required bool issue,
    GestureTapCallback? onTap}) {
  return isIOS(context)
      ? Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: CupertinoListTile(
              leading: Icon(
                CupertinoIcons.archivebox,
                color: issue ? Colors.yellow : Colors.red,
                size: 32,
              ),
              title: title,
              onTap: onTap))
      : Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: ListTile(
              leading: Icon(Icons.inventory_2), title: title, onTap: onTap),
        );
}
