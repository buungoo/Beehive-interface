import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import '../utils/helpers.dart';

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
                    // Add your onPressed code here
                  },
                  child: Icon(CupertinoIcons.add),
                )
              : null,
        )
      : AppBar(
          title: Text(title),
          backgroundColor: bgcolor,
          actions: Action
              ? [IconButton(onPressed: () {}, icon: Icon(Icons.add))]
              : null,
        );
}
