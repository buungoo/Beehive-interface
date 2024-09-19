//CupertinoPageScaffold

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../utils/helpers.dart';

Widget SharedScaffold(
    {required BuildContext context, dynamic appBar, required Widget body}) {
  return isIOS(context)
      ? CupertinoPageScaffold(navigationBar: appBar, child: body)
      : Scaffold(appBar: appBar, body: body);
}
