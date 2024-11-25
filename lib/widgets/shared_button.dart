import 'package:beehive/utils/helpers.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class SharedButton extends StatelessWidget {
  final void Function()? onPressed;
  //final BuildContext context;

  SharedButton({this.onPressed});

  @override
  Widget build(BuildContext context) {
    return isIOS(context) ? _buildCupertinoButton() : _buildMaterialButton();
  }

  Widget _buildCupertinoButton() {
    return CupertinoButton(
      padding: EdgeInsets.zero,
      onPressed: onPressed,
      child: Icon(CupertinoIcons.add),
    );
  }

  Widget _buildMaterialButton() {
    return IconButton(onPressed: onPressed, icon: Icon(Icons.add));
  }
}
