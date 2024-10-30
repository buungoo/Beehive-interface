import 'package:fl_chart/fl_chart.dart';
import '../models/beehive.dart';
import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';

class BeeChartPage extends StatelessWidget {
  final String? id;
  final Beehive? beehive;

  const BeeChartPage({this.id, this.beehive, super.key});

  LineTouchData get lineTouchData1 => LineTouchData(
        handleBuiltInTouches: true,
        touchTooltipData: LineTouchTooltipData(
          getTooltipColor: (touchedSpot) => Colors.blueGrey.withOpacity(0.8),
        ),
      );

  FlGridData get gridData => const FlGridData(show: false);

  LineChartData get sampleData1 => LineChartData(
        lineTouchData: lineTouchData1,
        gridData: gridData,
        minX: 0,
        maxX: 14,
        maxY: 4,
        minY: 0,
      );

  @override
  Widget build(BuildContext context) {
    return LineChart(
      sampleData1,
      duration: const Duration(milliseconds: 250),
    );
  }
}
