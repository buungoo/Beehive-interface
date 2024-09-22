import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation


class InitialPage extends StatelessWidget{
  const InitialPage({super.key});

  void printText(){
    print("Yeet");
  }

  @override
  Widget build(BuildContext context){
    return Scaffold(
      body: Center(
          child: Column(
          mainAxisAlignment: MainAxisAlignment.center,  // Center Y axis
          children: [
            Text('Hello'),
            SizedBox(height: 25),
            ElevatedButton(
              onPressed: () {
                context.push('/login_page');
              },
              child: Text('Login'),
            ),
            SizedBox(height: 15,),
            ElevatedButton(
              onPressed: () {
                context.push('/signup_page');
              },
              child: Text('Sign Up'),
            ),
          ],
        ),
      ),
    );


  }


}