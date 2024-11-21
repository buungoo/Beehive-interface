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
    double maxValue = values.isNotEmpty ? values.first.value : 0;
    double minValue = values.isNotEmpty ? values.first.value : 0;

    for (var value in values) {
      if (value.value > maxValue) maxValue = value.value;
      if (value.value < minValue) minValue = value.value;
    }

    if (values.isEmpty) {
      return LineChartData(
        lineBarsData: [],
        gridData: FlGridData(show: false),
        titlesData: FlTitlesData(show: false),
        borderData: FlBorderData(show: false),
      );
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
              final formattedDate = DateFormat('y-M-d').format(dateTime);

              return Padding(
                padding: const EdgeInsets.only(
                    top: 10.0), // Add padding to prevent overlap
                child: Transform.rotate(
                  angle: -0.5, // Slight tilt for better readability
                  child: Text(formattedDate, style: TextStyle(fontSize: 10)),
                ),
              );
            },
            reservedSize: 60,
          ),
        ),
      ),
      borderData: FlBorderData(
        show: true,
        border: Border.all(color: Colors.black),
      ),
      maxY: (maxValue + (maxValue / 4)).floorToDouble(),
      minY: (minValue - 10).floorToDouble(),
      lineBarsData: [
        LineChartBarData(
          color: Colors.yellow,
          spots: values
              .map((item) => FlSpot(item.time.millisecondsSinceEpoch.toDouble(),
                  item.value.floorToDouble()))
              .toList(),
          isCurved: false,
          isStepLineChart: true,
          barWidth: 4,
          dotData: FlDotData(show: false),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return FutureProvider<List<SensorValues>?>(
      initialData: [],
      create: (context) async {
        final beehiveId = beehive?.id;
        if (beehiveId == null) return [];
        try {
          final onValue = await BeehiveDataProvider()
              .getBeehiveDataChartStream(beehiveId.toString(), type)
              .first;
          return SensorValues.fromJsonList(onValue);
        } catch (e) {
          print('Error fetching data: $e');
          return [];
        }
      },
      catchError: (context, error) {
        print('Caught error: $error');
        return [];
      },
      child: SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context, title: title, bgcolor: Color(0xFFf4991a)),
        body: _DrawChart(type: type),
      ),
    );
  }

  Widget DrawChart(context) {
    return Consumer<List<SensorValues>?>(
      builder: (context, _data, child) {
        print(_data);
        if (_data == null || _data.length < 1) {
          return Center(child: SharedLoadingIndicator(context: context));
        }
        return Align(
          alignment: Alignment.topCenter,
          child: Container(
            height: MediaQuery.of(context).size.height,
            width: double.infinity,
            padding: const EdgeInsets.all(16.0),
            child: LineChart(buildChartData(_data)),
          ),
        );
      },
    );
  }
}

class _DrawChart extends StatefulWidget {
  final String type;

  const _DrawChart({required this.type});

  @override
  State<_DrawChart> createState() => _DrawChartState();
}

class _DrawChartState extends State<_DrawChart> {
  late double touchedValue;

  final Color? lineColor = Colors.yellow[700];
  final Color? avgLineColor = Colors.yellow[800];
  final Color pointColor = Color(0xFFFFEB3B);

  @override
  void initState() {
    touchedValue = -1;
    super.initState();
  }

  Widget leftTitleWidgets(double value, TitleMeta meta) {
    if (value % 1 != 0) {
      return Container();
    }
    final style = TextStyle(
      color: Colors.green.withOpacity(0.5),
      fontSize: 10,
    );
    String text;
    text = '${value.toInt()}';

    return SideTitleWidget(
      axisSide: meta.axisSide,
      space: 6,
      child: Text(text, style: style, textAlign: TextAlign.center),
    );
  }

  Widget bottomTitleWidgets(double value, TitleMeta meta, List weekday) {
    //const weekday = ['Sat', 'Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri'];

    final isTouched = value == touchedValue;
    final style = TextStyle(
      color: isTouched ? Colors.black : Colors.yellow,
      fontWeight: FontWeight.bold,
    );

    if (value % 1 != 0) {
      return Container();
    }
    return SideTitleWidget(
      space: 4,
      axisSide: meta.axisSide,
      child: Text(
        weekday[value.toInt()],
        style: style,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    //const weekday = ['Sat', 'Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri'];

    return Consumer<List<SensorValues>?>(
      builder: (context, _data, child) {
        if (_data == null || _data.length < 1) {
          return Center(child: SharedLoadingIndicator(context: context));
        }
        var yValues = _data
            .where((item) =>
                item.sensor_id ==
                1) // Filter only items with sensor_id == 1 //TODO: Make this filter correct
            .map((item) => item.value) // Map
            .toList(); // Convert to list

        // generate a weekday list for the x-axis based on yValues date values
        // _data contains DateTime time, convert it to 'Sun', ect..
        var weekday = [];
        _data.toList().forEach((element) {
          if (element.sensor_id != 1) return;
          weekday.add(DateFormat('E').format(element.time));
        });

        var average = yValues.reduce((a, b) => a + b) / yValues.length;
        double maxValue = yValues.isNotEmpty ? yValues[0] : 0;
        double minValue = yValues.isNotEmpty ? yValues[0] : 0;

        for (var value in yValues) {
          if (value > maxValue) maxValue = value;
          if (value < minValue) minValue = value;
        }

        return Column(
          mainAxisSize: MainAxisSize.min,
          children: <Widget>[
            const SizedBox(height: 10),
            Row(
              mainAxisSize: MainAxisSize.min,
              children: <Widget>[
                Text(
                  'Over the last 7 days',
                  style: TextStyle(
                    color: Colors.black,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
            const SizedBox(
              height: 18,
            ),
            AspectRatio(
              aspectRatio: 2,
              child: Padding(
                padding: const EdgeInsets.only(right: 20.0, left: 12),
                child: LineChart(
                  LineChartData(
                    lineTouchData: LineTouchData(
                      getTouchedSpotIndicator:
                          (LineChartBarData barData, List<int> spotIndexes) {
                        return spotIndexes.map((spotIndex) {
                          final spot = barData.spots[spotIndex];
                          if (spot.x == 0 || spot.x == 6) {
                            return null;
                          }
                          return TouchedSpotIndicatorData(
                            FlLine(
                              color: Colors.amber,
                              strokeWidth: 4,
                            ),
                            FlDotData(
                              getDotPainter: (spot, percent, barData, index) {
                                if (index.isEven) {
                                  return FlDotCirclePainter(
                                    radius: 8,
                                    color: Colors.white,
                                    strokeWidth: 5,
                                    strokeColor: Colors.green,
                                  );
                                } else {
                                  return FlDotSquarePainter(
                                    size: 16,
                                    color: Colors.white,
                                    strokeWidth: 5,
                                    strokeColor: Colors.green,
                                  );
                                }
                              },
                            ),
                          );
                        }).toList();
                      },
                      touchTooltipData: LineTouchTooltipData(
                        getTooltipColor: (touchedSpot) => Colors.blueAccent,
                        getTooltipItems: (List<LineBarSpot> touchedBarSpots) {
                          return touchedBarSpots.map((barSpot) {
                            final flSpot = barSpot;
                            if (flSpot.x == 0 || flSpot.x == 6) {
                              return null;
                            }

                            TextAlign textAlign;
                            switch (flSpot.x.toInt()) {
                              case 1:
                                textAlign = TextAlign.left;
                                break;
                              case 5:
                                textAlign = TextAlign.right;
                                break;
                              default:
                                textAlign = TextAlign.center;
                            }

                            return LineTooltipItem(
                              '${weekday[flSpot.x.toInt()]} \n',
                              TextStyle(
                                color: Colors.black,
                                fontWeight: FontWeight.bold,
                              ),
                              children: [
                                TextSpan(
                                  text: flSpot.y.toString(),
                                  style: TextStyle(
                                    color: Colors.black,
                                    fontWeight: FontWeight.w900,
                                  ),
                                ),
                              ],
                              textAlign: textAlign,
                            );
                          }).toList();
                        },
                      ),
                      touchCallback:
                          (FlTouchEvent event, LineTouchResponse? lineTouch) {
                        if (!event.isInterestedForInteractions ||
                            lineTouch == null ||
                            lineTouch.lineBarSpots == null) {
                          setState(() {
                            touchedValue = -1;
                          });
                          return;
                        }
                        final value = lineTouch.lineBarSpots![0].x;

                        if (value == 0 || value == 6) {
                          setState(() {
                            touchedValue = -1;
                          });
                          return;
                        }

                        setState(() {
                          touchedValue = value;
                        });
                      },
                    ),
                    extraLinesData: ExtraLinesData(
                      horizontalLines: [
                        HorizontalLine(
                          y: average, // avg
                          color: avgLineColor,
                          strokeWidth: 3,
                          dashArray: [20, 10],
                        ),
                      ],
                    ),
                    lineBarsData: [
                      LineChartBarData(
                        isStepLineChart: true,
                        spots: yValues.asMap().entries.map((e) {
                          return FlSpot(e.key.toDouble(),
                              e.value.toDouble().floorToDouble());
                        }).toList(),
                        isCurved: false,
                        barWidth: 4,
                        color: lineColor,
                        belowBarData: BarAreaData(
                          show: true,
                          gradient: LinearGradient(
                            colors: [
                              Colors.yellow.withOpacity(0.5),
                              Colors.yellow.withOpacity(0),
                            ],
                            stops: const [0.5, 1.0],
                            begin: Alignment.topCenter,
                            end: Alignment.bottomCenter,
                          ),
                          spotsLine: BarAreaSpotsLine(
                            show: true,
                            flLineStyle: FlLine(
                              color: Colors.yellow,
                              strokeWidth: 2,
                            ),
                            checkToShowSpotLine: (spot) {
                              if (spot.x == 0 || spot.x == 6) {
                                return false;
                              }

                              return true;
                            },
                          ),
                        ),
                        dotData: FlDotData(
                          show: true,
                          getDotPainter: (spot, percent, barData, index) {
                            if (index.isEven) {
                              return FlDotCirclePainter(
                                radius: 6,
                                color: Colors.white,
                                strokeWidth: 3,
                                strokeColor: pointColor,
                              );
                            } else {
                              return FlDotSquarePainter(
                                size: 12,
                                color: Colors.white,
                                strokeWidth: 3,
                                strokeColor: pointColor,
                              );
                            }
                          },
                          checkToShowDot: (spot, barData) {
                            return spot.x != 0 && spot.x != 6;
                          },
                        ),
                      ),
                    ],
                    minY: minValue - 10,
                    maxY: maxValue + 10,
                    borderData: FlBorderData(
                      show: true,
                      border: Border.all(
                        color: Colors.black,
                      ),
                    ),
                    gridData: FlGridData(
                      show: true,
                      drawHorizontalLine: true,
                      drawVerticalLine: true,
                      checkToShowHorizontalLine: (value) => value % 1 == 0,
                      checkToShowVerticalLine: (value) => value % 1 == 0,
                      getDrawingHorizontalLine: (value) {
                        if (value == 0) {
                          return const FlLine(
                            color: Colors.orange,
                            strokeWidth: 2,
                          );
                        } else {
                          return const FlLine(
                            color: Colors.grey,
                            strokeWidth: 0.5,
                          );
                        }
                      },
                      getDrawingVerticalLine: (value) {
                        if (value == 0) {
                          return const FlLine(
                            color: Colors.redAccent,
                            strokeWidth: 10,
                          );
                        } else {
                          return const FlLine(
                            color: Colors.grey,
                            strokeWidth: 0.5,
                          );
                        }
                      },
                    ),
                    titlesData: FlTitlesData(
                      show: true,
                      topTitles: const AxisTitles(
                        sideTitles: SideTitles(showTitles: false),
                      ),
                      rightTitles: const AxisTitles(
                        sideTitles: SideTitles(showTitles: false),
                      ),
                      leftTitles: AxisTitles(
                        sideTitles: SideTitles(
                          showTitles: true,
                          reservedSize: 46,
                          getTitlesWidget: leftTitleWidgets,
                        ),
                      ),
                      bottomTitles: AxisTitles(
                        sideTitles: SideTitles(
                          showTitles: true,
                          reservedSize: 40,
                          getTitlesWidget: (value, meta) =>
                              bottomTitleWidgets(value, meta, weekday),
                        ),
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        );
      },
    );
  }
}
