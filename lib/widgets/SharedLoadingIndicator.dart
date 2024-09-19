//CupertinoPageScaffold

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import '../utils/helpers.dart';

Widget SharedLoadingIndicator({required BuildContext context}) {
  return isIOS(context)
      ? const CupertinoActivityIndicator()
      : const CircularProgressIndicator();
}
