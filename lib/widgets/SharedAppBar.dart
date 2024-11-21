import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import '../utils/helpers.dart';
import 'package:go_router/go_router.dart'; //

PreferredSizeWidget? getNavigationBar(
    {required BuildContext context,
    required String title,
    Color? bgcolor,
    bool Action = false}) {
  return isIOS(context)
      ? CupertinoNavigationBar(
          middle: Text(title),
          backgroundColor: bgcolor,
          trailing: Action
              ? CupertinoButton(
                  padding: EdgeInsets.zero,
                  onPressed: () {
                    context.pushNamed("Camera");
                  },
                  child: Icon(CupertinoIcons.add),
                )
              : null,
        )
      : AppBar(
          title: Text(title),
          backgroundColor: bgcolor,
          actions: Action
              ? [
                  IconButton(
                      onPressed: () {
                        Navigator.of(context).pushNamed('Camera');
                      },
                      icon: Icon(Icons.add))
                ]
              : null,
        );
}
