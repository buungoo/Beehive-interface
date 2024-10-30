import 'package:beehive/models/beehive_user.dart';
import 'package:beehive/config.dart' as config;
import 'package:http/http.dart';
import 'dart:convert';

class BeehiveUserProvider {
  Future<User?> login(String email, String password) async {
    print(email);

    try {
      print(config.BackendServer + '/login');
      var response = await post(Uri.parse(config.BackendServer + '/login'),
          headers: <String, String>{
            'Content-Type': 'application/json; charset=UTF-8',
          },
          body: jsonEncode({
            'Username': email,
            'Password': password,
          }));

      final data = json.decode(response.body);

      if (response.statusCode == 200) {
        print(data);
        return await User.fromToken(data['token']);
      }
    } catch (e) {
      throw Exception(e);
      //print(e);
    }
    return null;
  }

  Future<User> register(String email, String password) async {
    try {
      var response = await post(Uri.parse(config.BackendServer + '/register'),
          headers: <String, String>{
            'Content-Type': 'application/json; charset=UTF-8',
          },
          body: jsonEncode({
            'email': email,
            'password': password,
          }));

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
      var response = await get(Uri.parse(config.BackendServer + '/user'),
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
