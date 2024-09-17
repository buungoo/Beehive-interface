import 'package:flutter/material.dart';
import '../models/beehive.dart';

class BeehiveListProvider extends ChangeNotifier {
  final List<Beehive> _beehives = [
    Beehive(id: '1', name: 'Beehive 1'),
    Beehive(id: '2', name: 'Beehive 2'),
  ];

  List<Beehive> get beehives => _beehives;

  // Simulates finding a beehive by ID
  Beehive? findBeehiveById(String id) {
    return _beehives.firstWhere(
      (beehive) => beehive.id == id,
      orElse: () => throw Exception('Beehive not found'),
    );
  }
}
