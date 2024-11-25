import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import '../utils/helpers.dart';
import 'package:go_router/go_router.dart'; //

PreferredSizeWidget? getNavigationBar(
    {required BuildContext context,
    required String title,
    Color? bgcolor,
    Widget? ActionButtn}) {
  return isIOS(context)
      ? CupertinoNavigationBar(
          middle: Text(title),
          backgroundColor: bgcolor,
          trailing: ActionButtn,
          //trailing: ActionButtn == null ? [ActionButtn] : null,
        )
      : AppBar(
          title: Text(title),
          backgroundColor: bgcolor,
          actions: ActionButtn != null ? [ActionButtn] : null,
          //actions: ActionButtn != null ? [ActionButtn] : null,
        );
}
