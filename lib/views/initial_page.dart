import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation
import 'package:beehive/services/BeehiveApiService.dart';

class InitialPage extends StatelessWidget {
  const InitialPage({super.key});

  void checkAuth(BuildContext context) async {
    //print("Checking Auth");
    final authenticated = await BeehiveApi().verifyUser();
    //print("Authenticated: $authenticated");
    if (authenticated) {
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
            const Text('Hello'),
            const SizedBox(height: 25),
            ElevatedButton(
              key: const Key("Login"),
              onPressed: () {
                context.push('/login_page');
              },
              child: const Text('Login'),
            ),
            const SizedBox(
              height: 15,
            ),
            ElevatedButton(
              key: const Key("Signin"),
              onPressed: () {
                context.push('/signup_page');
              },
              child: const Text('Sign Up'),
            ),
          ],
        ),
      ),
    );
  }
}
