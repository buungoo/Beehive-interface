import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart'; // GoRouter for navigation
import '../widgets/SharedAppBar.dart';
import 'package:beehive/providers/beehive_user_provider.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  _LoginPageState createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  final TextEditingController _emailController = TextEditingController();
  final TextEditingController _passwordController = TextEditingController();
  String? _errorMessage;

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _login() async {
    final String email = _emailController.text;
    final String password = _passwordController.text;

    try {
      final user = await BeehiveUserProvider().login(email, password);

      if (user == null) {
        throw Exception('Invalid credentials');
      } else {
        context.go('/overview');
      }
      // Navigate to the 'overview' page on successful login
    } catch (e) {
      setState(() {
        _errorMessage = e.toString(); // Capture the error message
      });
    }
  }

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
                key: const Key("usernameField"),
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
                key: const Key("passwordField"),
                decoration: const InputDecoration(
                  border: UnderlineInputBorder(),
                  labelText: 'Enter your password',
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
              onPressed: _login,
              key: const Key("Login"),
              child: const Text('Login'),
            ),
          ],
        ),
      ),
    );
  }
}
