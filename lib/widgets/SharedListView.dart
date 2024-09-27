import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../utils/helpers.dart';

Widget SharedListTile(
    {required BuildContext context,
    required Widget title,
    GestureTapCallback? onTap}) {
  return isIOS(context)
      ? Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: CupertinoListTile(
              leading: const Icon(
                CupertinoIcons.archivebox,
                color: Colors.red,
                size: 32,
              ),
              title: title,
              onTap: onTap))
      : Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: ListTile(
              leading: Icon(Icons.inventory_2), title: title, onTap: onTap),
        );

  return isIOS(context)
      ? CupertinoListTile(title: title, onTap: onTap)
      : ListTile(title: title, onTap: onTap);
}
