import 'package:flutter/material.dart';
import 'package:beehive/models/beehive.dart';
import 'package:beehive/services/BeehiveApiService.dart';

class BeehiveListProvider extends ChangeNotifier {
  // A hardcoded list of beehives that the provider manages
  final List<Beehive> _beehives = [
    Beehive(id: '1', name: 'Beehive 1'),
    Beehive(id: '2', name: 'Beehive 2'),
  ];

  //List<Beehive> _beehives = BeehiveApi().GetHives();

  // Getter that returns the list of beehives
  Future<List<Beehive>> get beehives => BeehiveApi().GetHives();

  // Finds a beehive by its ID. Throws an exception if the beehive isn't found.
  Beehive? findBeehiveById(String id) {
    return _beehives.firstWhere(
      (beehive) => beehive.id == id, // Check if the beehive ID matches
      orElse: () =>
          throw Exception('Beehive not found'), // Throws error if no match
    );
  }

  void addBeehive(Beehive beehive) {
    _beehives.add(beehive);
    notifyListeners();
  }
}
