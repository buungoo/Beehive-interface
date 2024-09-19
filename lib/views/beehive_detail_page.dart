import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter/cupertino.dart';
import '../models/beehive.dart';
import '../models/beehive_data.dart';
import '../providers/beehive_data_provider.dart';
import '../widgets/shared.dart';
import '../widgets/BeeDetailView/detailgrid.dart';

class BeehiveDetailPage extends StatelessWidget {
  final Beehive beehive;

  const BeehiveDetailPage({required this.beehive, super.key});

  @override
  Widget build(BuildContext context) {
    return StreamProvider<BeehiveData?>(
      // Because the StreamProvider is specified here and not in BeehiveApp
      // class only the BeehiveDetailPage widget can listen to it
      initialData: null, // Nullable initial data
      create: (context) {
        // Setup the Stream which the StreamProvider should listen to
        return BeehiveDataProvider().getBeehiveDataStream();
      },
      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(context: context, title: beehive.name),
        body: DetailGrid(),
      ),
    );
  }
}
