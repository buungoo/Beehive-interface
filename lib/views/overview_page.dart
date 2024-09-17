import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/beehive_list_provider.dart';

class OverviewPage extends StatelessWidget {
  const OverviewPage({super.key});

  @override
  Widget build(BuildContext context) {
    // Through context we acces global providers
    // Start listening to the global BeehiveListProvider
    final beehiveList = context.watch<BeehiveListProvider>().beehives;

    return Scaffold(
      appBar: AppBar(title: const Text('Beehive Overview')),
      body: ListView.builder(
        itemCount: beehiveList.length,
        itemBuilder: (context, index) {
          final beehive = beehiveList[index];
          return ListTile(
            title: Text(beehive.name),
            onTap: () {
              // Navigate to the beehive detail page
              context.go('/beehive/${beehive.id}');
            },
          );
        },
      ),
    );
  }
}
