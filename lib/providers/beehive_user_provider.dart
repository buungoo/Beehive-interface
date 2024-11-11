import 'package:beehive/models/beehive_user.dart';
import 'package:beehive/config.dart' as config;
import 'package:http/http.dart' as http;
import 'dart:convert';

class BeehiveUserProvider {
  final http.Client client;

  // Accept the client as a constructor parameter (default to `http.Client` if not provided)
  BeehiveUserProvider({http.Client? client})
      : client = client ?? http.Client(); // To allow us to use mock data

  Future<User?> login(String email, String password) async {
    try {
      print(config.BackendServer + '/login');
      //print("Mock login attempt");
      var response =
          await client.post(Uri.parse(config.BackendServer + '/login'),
              headers: <String, String>{
                'Content-Type': 'application/json; charset=UTF-8',
              },
              body: jsonEncode({
                'username': email,
                'password': password,
              }));

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        return await User.fromToken(data['token']);
      } else {
        return null;
      }
    } catch (e) {
      print(e);
      throw Exception(e);
      //print(e);
    }
    return null;
  }

  Future<User> register(String email, String password) async {
    try {
      var response =
          await client.post(Uri.parse(config.BackendServer + '/register'),
              headers: <String, String>{
                'Content-Type': 'application/json; charset=UTF-8',
              },
              body: jsonEncode({
                'username': email,
                'password': password,
              }));

      print(response);

      if (response.statusCode == 200) {
        return User.fromJson(response.body);
      }
    } catch (e) {
      print(e);
    }

    throw Exception('Failed to register');
  }

  Future<void> logout() async {
    await User.removeUser();
  }

  Future<User?> getUser(String token) async {
    try {
      var response = await client.get(Uri.parse(config.BackendServer + '/user'),
          headers: <String, String>{
            'Content-Type': 'application/json; charset=UTF-8',
            'Authorization': 'Bearer $token',
          });

      if (response.statusCode == 200) {
        return User.fromJson(response.body);
      }
    } catch (e) {
      print(e);
    }

    return null;
  }

  Future<User?> getUserFromStorage() async {
    return await User.getUser();
  }
}
