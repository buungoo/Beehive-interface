import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation
import 'package:beehive/services/BeehiveApiService.dart';
import 'package:beehive/providers/beehive_data_provider.dart';
import 'package:beehive/services/BeehiveNotificationService.dart';

class InitialPage extends StatelessWidget {
  const InitialPage({super.key});

  void checkAuth(BuildContext context) async {
    print("Checking Auth");
    final authenticated = await BeehiveApi().verifyUser();
    print("Authenticated: $authenticated");
    if (authenticated) {
      await BeeNotification().checkIssues();
      // Navigate to the 'overview' page on successful login
      context.go('/overview');
    }
  }

  @override
  Widget build(BuildContext context) {
    checkAuth(context);

    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center, // Center Y axis
          children: [
            Text('Hello'),
            SizedBox(height: 25),
            ElevatedButton(
              key: Key("Login"),
              onPressed: () {
                context.push('/login_page');
              },
              child: Text('Login'),
            ),
            SizedBox(
              height: 15,
            ),
            ElevatedButton(
              key: Key("Signin"),
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
