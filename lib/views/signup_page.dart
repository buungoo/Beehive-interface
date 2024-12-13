import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation
import '../widgets/SharedAppBar.dart';
import 'package:beehive/providers/beehive_user_provider.dart';

class SignupPage extends StatefulWidget {
  const SignupPage({super.key});

  @override
  _SignupPageState createState() => _SignupPageState();
}

class _SignupPageState extends State<SignupPage> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  final TextEditingController _passwordConfirmController =
      TextEditingController();
  String? _errorMessage;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    _passwordConfirmController.dispose();
    super.dispose();
  }

  Future<void> _signup() async {
    final String email = _emailController.text;
    final String password = _passwordController.text;
    final String passwordConfirm = _passwordConfirmController.text;

    try {
      if (password != passwordConfirm || password.isEmpty) {
        throw Exception('Passwords do not match');
      }

      await BeehiveUserProvider().register(email, password);

      context.go('/overview');
    } catch (e) {
      setState(() {
        print(e.toString());
        _errorMessage = e.toString(); // Capture the error message
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: getNavigationBar(context: context, title: 'Sign up'),
      body: Center(
        child: Column(
          mainAxisAlignment:
              MainAxisAlignment.center, // Centers the content vertically
          children: [
            const Text(
              'Signup',
              style: TextStyle(fontSize: 24), // Optional: Add some styling
            ),
            const SizedBox(
                height: 25), // Adds spacing between text and input fields

            // Username Text Field
            SizedBox(
              width: 250,
              child: TextFormField(
                controller: _emailController,
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
                controller: _passwordController,
                decoration: const InputDecoration(
                  border: UnderlineInputBorder(),
                  labelText: 'Enter your password',
                ),
                obscureText: true, // Masks the text input for passwords
                style: const TextStyle(fontSize: 20),
              ),
            ),

            const SizedBox(height: 15), // Adds spacing between input and button

            // Confirm Password Text Field
            SizedBox(
              width: 250,
              child: TextFormField(
                controller: _passwordConfirmController,
                decoration: const InputDecoration(
                  border: UnderlineInputBorder(),
                  labelText: 'Confirm password',
                ),
                obscureText: true, // Masks the text input for passwords
                style: const TextStyle(fontSize: 20),
              ),
            ),

            const SizedBox(height: 15), // Adds spacing between input and button

            // Error Message Display
            if (_errorMessage != null) ...[
              Text(
                _errorMessage!,
                style: const TextStyle(color: Colors.red),
              ),
              const SizedBox(height: 15),
            ],

            // Login Button
            ElevatedButton(
              onPressed: _signup,
              child: const Text('Signup'),
            ),
          ],
        ),
      ),
    );
  }
}
