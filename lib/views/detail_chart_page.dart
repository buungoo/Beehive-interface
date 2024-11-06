import 'package:fl_chart/fl_chart.dart';
import '../models/beehive.dart';
import 'package:provider/provider.dart';
import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:beehive/widgets/shared.dart';
import 'package:beehive/providers/beehive_data_provider.dart';
import 'package:beehive/models/SensorValues.dart';
import 'package:intl/intl.dart';

class BeeChartPage extends StatelessWidget {
  final String? id;
  final Beehive? beehive;
  final String title;
  final String type; // temperature, weight, humidity, ppm

  const BeeChartPage(
      {this.id,
      this.beehive,
      super.key,
      required this.title,
      required this.type});

  LineChartData buildChartData(List<SensorValues> values) {
    // find max value in values
    double maxValue = 0;
    for (var value in values) {
      if (value.value > maxValue) {
        maxValue = value.value;
      }
    }

    return LineChartData(
      gridData: FlGridData(show: true),
      titlesData: FlTitlesData(
        rightTitles: AxisTitles(
          sideTitles: SideTitles(showTitles: false, reservedSize: 40),
        ),
        topTitles: AxisTitles(
          sideTitles: SideTitles(showTitles: false, reservedSize: 40),
        ),
        leftTitles: AxisTitles(
            sideTitles: SideTitles(
          showTitles: true,
          getTitlesWidget: (value, meta) {
            return Text(value.toString(), style: TextStyle(fontSize: 10));
          },
          reservedSize: 30,
        )),
        bottomTitles: AxisTitles(
          sideTitles: SideTitles(
            showTitles: true,
            getTitlesWidget: (value, meta) {
              final DateTime dateTime =
                  DateTime.fromMillisecondsSinceEpoch(value.toInt());
              final formattedDate = DateFormat('y-M-d')
                  .format(dateTime); // Format the date as needed
              return Text(formattedDate, style: TextStyle(fontSize: 10));
            },
            reservedSize: 40,
          ),
        ),
      ),
      borderData: FlBorderData(
        show: true,
        border: Border.all(color: Colors.black),
      ),
      maxY: maxValue,
      lineBarsData: [
        LineChartBarData(
          color: Colors.yellow,
          spots: values
              .map((item) => FlSpot(
                  item.time.millisecondsSinceEpoch.toDouble(), item.value))
              .toList(),
          isCurved: true,
          barWidth: 4,
          dotData: FlDotData(show: false),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return StreamProvider<List<SensorValues>?>(
      // Because the StreamProvider is specified here and not in BeehiveApp
      // class only the BeehiveDetailPage widget can listen to it
      initialData: [],
      create: (context) async* {
        final beehiveId = beehive?.id;
        if (beehiveId == null) {
          yield [];
        } else {
          try {
            final onValue = await BeehiveDataProvider()
                .getBeehiveDataChartStream(beehiveId.toString(), type)
                .first;
            yield SensorValues.fromJsonList(onValue);
          } catch (e) {
            // Handle the error appropriately
            print('Error fetching data: $e');
            throw e; // Re-throw the error if needed
          }
        }
      },

      catchError: (context, error) {
        print('Caught error: $error');
        return null;
      },

      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context, title: title, bgcolor: Color(0xFFf4991a)),
        body: DrawChart(context),
      ),
    );
  }

  Widget DrawChart(context) {
    return Consumer<List<SensorValues>?>(
      builder: (context, _data, child) {
        print(_data);
        if (_data == null) {
          return Center(child: SharedLoadingIndicator(context: context));
        }
        return Align(
          alignment: Alignment.topCenter,
          child: Container(
            height: 300,
            width: double.infinity,
            padding: const EdgeInsets.all(16.0),
            child: LineChart(buildChartData(_data)),
          ),
        );
      },
    );
  }
}
