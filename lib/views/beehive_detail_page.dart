import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter/cupertino.dart';
import '../models/beehive.dart';
import '../providers/beehive_data_provider.dart';
import '../utils/helpers.dart';
import '../widgets/shared.dart';

class BeehiveDetailPage extends StatelessWidget {
  final Beehive beehive;

  const BeehiveDetailPage({required this.beehive, super.key});

  @override
  Widget build(BuildContext context) {
    return StreamProvider<int?>(
      // Because the StreamProvider is specified here and not in BeehiveApp
      // class only the BeehiveDetailPage widget can listen to it
      initialData: null, // Nullable initial data
      create: (context) {
        // Setup the Stream which the StreamProvider should listen to
        return BeehiveDataProvider().getTemperatureStream();
      },
      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(context: context, title: beehive.name),
        body: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Beehive ID: ${beehive.id}'),
              const SizedBox(height: 20),
              Consumer<int?>(
                builder: (context, temperature, child) {
                  // Show loading until the stream emits the first data
                  if (temperature == null) {
                    return SharedLoadingIndicator(context: context);
                  }
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
