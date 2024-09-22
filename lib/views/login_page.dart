import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation
import '../widgets/SharedAppBar.dart';

class LoginPage extends StatelessWidget {
  const LoginPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: getNavigationBar(context: context, title: 'Login'),
      body: Center(
        child: Column(
          mainAxisAlignment:
              MainAxisAlignment.center, // Centers the content vertically
          children: [
            const Text(
              'Login',
              style: TextStyle(fontSize: 24), // Optional: Add some styling
            ),
            const SizedBox(
                height: 25), // Adds spacing between text and input fields

            // Username Text Field
            SizedBox(
              width: 250,
              child: TextFormField(
                decoration: const InputDecoration(
                  border: UnderlineInputBorder(),
                  labelText: 'Enter your username',
                ),
                style: const TextStyle(fontSize: 20),
              ),
            ),

            const SizedBox(height: 15), // Adds spacing between input fields

            // Password Text Field
            SizedBox(
              width: 250,
              child: TextFormField(
                decoration: const InputDecoration(
                  border: UnderlineInputBorder(),
                  labelText: 'Enter your password',
                ),
                obscureText: true, // Masks the text input for passwords
                style: const TextStyle(fontSize: 20),
              ),
            ),

            const SizedBox(height: 15), // Adds spacing between input and button

            // Login Button
            ElevatedButton(
              onPressed: () {
                // Navigate to the 'overview' page
                // Here we need to talk to API and do som checks
                context.go('/overview');
              },
              child: const Text('Login'),
            ),
          ],
        ),
      ),
    );
  }
}
