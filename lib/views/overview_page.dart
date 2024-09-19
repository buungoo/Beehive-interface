import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/beehive_list_provider.dart';
import '../widgets/SharedListView.dart';
import '../widgets/SharedAppBar.dart';

class OverviewPage extends StatelessWidget {
  const OverviewPage({super.key});

  @override
  Widget build(BuildContext context) {
    // Through context we can access global providers
    // Start listening to the global BeehiveListProvider
    final beehiveList = context.watch<BeehiveListProvider>().beehives;

    return Scaffold(
      appBar: getNavigationBar(context: context, title: 'Beehive Overview'),
      body: ListView.builder(
        itemCount: beehiveList.length,
        itemBuilder: (context, index) {
          final beehive = beehiveList[index];
          return SharedListTile(
            context: context,
            title: Text(beehive.name),
            onTap: () {
              // Navigate to the beehive detail page using GoRouter pathing
              context.pushNamed(
                'beehive-detail', // The name of the route
                pathParameters: {
                  'id': beehive.id
                }, // Use 'pathParameters' to pass the id
              );
            },
          );
        },
      ),
    );
  }
}
