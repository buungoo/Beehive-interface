import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter/cupertino.dart';
import '../models/beehive.dart';
import '../providers/beehive_data_provider.dart';
import '../utils/helpers.dart';
import '../widgets/SharedAppBar.dart';

class LoginPage extends StatelessWidget{
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context){
    return Scaffold(
      appBar: getNavigationBar(context: context, title: 'Login'),
      body: Text("Welcome")
    );

  }


}