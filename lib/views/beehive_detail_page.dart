import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/beehive.dart';

class BeehiveDetailPage extends StatelessWidget {
  final Beehive beehive;

  const BeehiveDetailPage({required this.beehive, super.key});

  @override
  Widget build(BuildContext context) {
    return StreamProvider<int>(
      initialData: 0, // Simulated initial data
      create: (context) {
        // Simulate a stream of data (e.g., temperature updates)
        return Stream.periodic(
            const Duration(seconds: 1), (count) => 20 + count);
      },
      child: Scaffold(
        appBar: AppBar(title: Text(beehive.name)),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Beehive ID: ${beehive.id}'),
              const SizedBox(height: 20),
              Consumer<int>(
                builder: (context, temperature, child) {
                  return Text('Simulated Temperature: $temperature°C');
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
