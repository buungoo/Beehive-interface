import '../models/beehive.dart';
import 'package:http/http.dart';
import 'package:beehive/config.dart' as config;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';

class BeehiveApi {
  BeehiveApi();

  /// Return an array of all beehives that are available for the user
  Future<List<Beehive>> GetHives() async {
    // need user token or smt

    try {
      final prefs = await SharedPreferences.getInstance();
      final token = prefs.getString('token');

      //print("${config.BackendServer}/beehive/list");

      var response = await get(
          Uri.parse("${config.BackendServer}/beehive/list"),
          headers: <String, String>{
            'Content-Type': 'application/json; charset=UTF-8',
            'Authorization': 'Bearer $token',
          });

      return (jsonDecode(response.body) as List)
          .map((e) => Beehive.fromJson(e))
          .toList();
    } catch (e) {
      print(e);
      return [];
    }
  }

  Future<bool> verifyUser() async {
    final perf = await SharedPreferences.getInstance();
    final token = perf.getString('token');

    //print(Uri.parse("${config.BackendServer}/test"));
    //print(token);

    if (token == null) {
      return false;
    }

    var response = await post(Uri.parse("${config.BackendServer}/test"),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Authorization': 'Bearer $token',
        });

    return response.statusCode == 200;
  }

/**
 * Return a single beehive by its id
 */
/*GetHive(int id) {
    return}*/
}
