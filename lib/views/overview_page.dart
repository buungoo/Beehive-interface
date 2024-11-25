import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import '../providers/beehive_list_provider.dart';
import '../widgets/shared.dart';
import 'package:beehive/models/beehive.dart';
import 'package:beehive/services/BeehiveNotificationService.dart';
import 'dart:io';
import 'package:beehive/widgets/shared_button.dart';

const simplePeriodicTask = "com.example.beehive.simplePeriodicTask";

class OverviewPage extends StatefulWidget {
  const OverviewPage({super.key});

  @override
  _OverviewPageState createState() => _OverviewPageState();
}

class _OverviewPageState extends State<OverviewPage> {
  late Future<List<Beehive>> beehiveList;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    beehiveList = context.watch<BeehiveListProvider>().beehives;
  }

  @override
  Widget build(BuildContext context) {
    // Through context we can access global providers
    // Start listening to the global BeehiveListProvider
    Future<List<Beehive>> beehiveList =
        context.watch<BeehiveListProvider>().beehives;

    return SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context,
            title: 'Beehive Overview',
            bgcolor: Color(0xFFf4991a),
            ActionButtn: SharedButton(
              //context: context,
              onPressed: () async {
                // Navigate to the camera page using GoRouter pathing
                final result = await context.pushNamed('Camera');
                print('Camera result: $result');
                if (result == true) {
                  setState(() {
                    beehiveList = context.read<BeehiveListProvider>().beehives;
                  });
                }
              },
            )),
        body: RefreshIndicator(
          onRefresh: () async {
            setState(() {
              beehiveList = context.read<BeehiveListProvider>().beehives;
            });
          },
          child: FutureBuilder<List<Beehive>>(
            future: beehiveList,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return Center(child: SharedLoadingIndicator(context: context));
              } else if (snapshot.hasError) {
                return Center(child: Text('Error: ${snapshot.error}'));
              } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
                return Center(child: Text('No beehives found'));
              } else {
                // If the future completed successfully, build the ListView
                final beehives = snapshot.data!;
                return ListView.builder(
                  itemCount: beehives.length,
                  itemBuilder: (context, index) {
                    final beehive = beehives[index];
                    return SharedListTile(
                      context: context,
                      title: Text(beehive.name),
                      issue: false,
                      onTap: () {
                        // Navigate to the beehive detail page using GoRouter pathing
                        context.pushNamed(
                          'beehive-detail', // The name of the route
                          pathParameters: {
                            'id': beehive.id.toString(),
                            // Ensure id is a string if needed
                          }, // Use 'pathParameters' to pass the id
                        );
                      },
                    );
                  },
                );
              }
            },
          ),
        ));
  }
}
