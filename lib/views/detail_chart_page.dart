import 'package:fl_chart/fl_chart.dart';
import '../models/beehive.dart';
import 'package:provider/provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:beehive/widgets/shared.dart';

class BeeChartPage extends StatelessWidget {
  final String? id;
  final Beehive? beehive;

  const BeeChartPage({this.id, this.beehive, super.key});

  LineChartData get data => LineChartData(
        gridData: FlGridData(show: true),
        titlesData: FlTitlesData(
          leftTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: true, reservedSize: 40),
          ),
          bottomTitles: AxisTitles(
            sideTitles: SideTitles(showTitles: true),
          ),
        ),
        borderData: FlBorderData(
          show: true,
          border: Border.all(color: Colors.black),
        ),
        maxY: 6,
        lineBarsData: [
          LineChartBarData(
            color: Colors.yellow,
            spots: [
              FlSpot(0, 1),
              FlSpot(1, 3),
              FlSpot(2, 2),
              FlSpot(3, 5),
              FlSpot(4, 4),
              FlSpot(5, 3),
              FlSpot(6, 4),
            ],
            isCurved: true,
            barWidth: 4,
            dotData: FlDotData(show: false),
          ),
        ],
      );

  @override
  Widget build(BuildContext context) {
    return StreamProvider<String?>(
      // Because the StreamProvider is specified here and not in BeehiveApp
      // class only the BeehiveDetailPage widget can listen to it
      initialData: "", // Nullable initial data
      create: (context) {
        return Stream.value("hi");
        // Setup the Stream which the StreamProvider should listen to
        //return BeehiveDataProvider().getBeehiveDataStream();
      },
      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context, title: "temperature", bgcolor: Color(0xFFf4991a)),
        body: Align(
          alignment: Alignment.topCenter,
          child: Container(
            height: 300,
            padding: const EdgeInsets.all(16.0),
            child: LineChart(data),
          ),
        ),
      ),
    );
  }
}
