import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:flutter/cupertino.dart';
import 'package:beehive/models/beehive.dart';
import 'package:beehive/models/beehive_data.dart';
import '../providers/beehive_data_provider.dart';
import '../widgets/shared.dart';
import '../widgets/BeeDetailView/detailgrid.dart';
import 'package:beehive/widgets/BeeDetailView/statusbox.dart';

class BeehiveDetailPage extends StatelessWidget {
  final Beehive beehive;

  const BeehiveDetailPage({required this.beehive, super.key});

  @override
  Widget build(BuildContext context) {
    return StreamProvider<BeehiveData>(
      // Because the StreamProvider is specified here and not in BeehiveApp
      // class only the BeehiveDetailPage widget can listen to it
      initialData:
          new BeehiveData(temperature: 0, weight: 0, humidity: 0, ppm: 0),
      create: (context) {
        return BeehiveDataProvider().getBeehiveDataStream(beehive.id);
      },
      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context,
            title: beehive.name,
            bgcolor: const Color(0xFFf4991a)),
        body: Stack(
          children: [
            _buildDetailGrid(beehive.id),
            _buildStatusBox(beehive.id),
          ],
        ),
      ),
    );
  }

  Widget _buildDetailGrid(String id) {
    return Positioned.fill(
      child: DetailGrid(id: id),
    );
  }

  Widget _buildStatusBox(String id) {
    return FutureBuilder<List<dynamic>>(
      future: BeehiveDataProvider().fetchBeehiveIssueStatusesList(id),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return Container(
            child: Text("Loading"),
          ); // or a loading indicator
        } else if (snapshot.hasError) {
          return Container(
            child: Text("Error"),
          ); // or an error message
        } else if (snapshot.hasData && snapshot.data!.isNotEmpty) {
          return Positioned(
            bottom: 0,
            left: 0,
            right: 0,
            child: Statusbox(data: snapshot.data!),
          );
        } else {
          return Container(
            child: Text("No data"),
          );
        }
      },
    );
  }
}
