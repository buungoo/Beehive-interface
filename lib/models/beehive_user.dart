import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:beehive/config.dart' as config;

class User {
  final String? uid;
  final String? username;
  final String? token;

  User({this.uid, this.username, this.token});

  saveUser() async {
    final perf = await SharedPreferences.getInstance();
    perf.setString('user', jsonEncode(this));
  }

  static fromJson(String token) async {
    // convert string to json
    //final perf = await SharedPreferences.getInstance();
    //perf.setString('user', json);
    //var decoded = jsonDecode(json);

    return User(uid: '2', username: 'Emil', token: token);
  }

  static fromToken(String token) async {
    final perf = await SharedPreferences.getInstance();
    await perf.setString('token', token);
    return User(uid: "0", username: '', token: token);
  }

  static getUser() async {
    // check if user is already stored in shared preferences
    final perf = await SharedPreferences.getInstance();
    var token = perf.getString('token');
    print(token);
    if (token == null) {
      return null;
    }

    return fromJson(token);
  }

  static removeUser() async {
    final perf = await SharedPreferences.getInstance();
    perf.remove('user');
  }
}
