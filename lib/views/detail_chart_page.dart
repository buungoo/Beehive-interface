import 'package:beehive/widgets/shared_dropdown.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:beehive/models/beehive.dart';
import 'package:provider/provider.dart';
import 'package:flutter/material.dart';
import 'package:beehive/widgets/shared.dart';
import 'package:beehive/models/SensorValues.dart';
import 'package:intl/intl.dart';
import 'package:beehive/providers/beehive_data_provider.dart';

class BeeChartDataProvider with ChangeNotifier {
  List<SensorValues> _sensorValues = [];
  String _selectedTimescale = '1 Day';
  bool _isLoading = false;
  String sensor;
  String beehiveID;

  // constructor with default values
  BeeChartDataProvider({this.sensor = "temperature", this.beehiveID = "1"});

  List<SensorValues> get sensorValues => _sensorValues;

  String get selectedTimescale => _selectedTimescale;

  bool get isLoading => _isLoading;

  void setIsLoading(bool value) {
    _isLoading = value;
    notifyListeners();
  }

  void setSensor(String value) {
    sensor = value;
  }

  void fetchData() async {
    setIsLoading(true);

    //TODO: Make sure this works with the API, atm we do not have any live data I can test with
    var fetchedData = await BeehiveDataProvider().fetchBeehiveDataChart(
        beehiveId: beehiveID, sensor: sensor, timescale: _selectedTimescale);
    try {
      setValues(fetchedData);
    } catch (e) {
      //print("Error fetching data");
      setValues([]);
    } finally {
      setIsLoading(false);
    }
  }

  void setValues(List<SensorValues> values) {
    _sensorValues = values;
    notifyListeners();
  }

  void setTimescale(String range) {
    _selectedTimescale = range;
    fetchData();
  }
}

class BeeChartPage extends StatelessWidget {
  final String? id;
  final Beehive beehive;
  final String title;
  final String type; // temperature, weight, humidity, ppm

  const BeeChartPage(
      {this.id,
      required this.beehive,
      super.key,
      required this.title,
      required this.type});

  @override
  Widget build(BuildContext context) {
    final dataProvider =
        BeeChartDataProvider(beehiveID: beehive.id, sensor: type);

    return ChangeNotifierProvider<BeeChartDataProvider>.value(
        value: dataProvider,
        child: Consumer<BeeChartDataProvider>(builder:
            (BuildContext context, BeeChartDataProvider Data, Widget? child) {
          return beeContent(context, Data);
        }));
  }

  Widget beeContent(BuildContext context, BeeChartDataProvider dataProvider) {
    return SharedScaffold(
        context: context,
        appBar: getNavigationBar(
            context: context, title: title, bgcolor: const Color(0xFFf4991a)),
        body: Center(
            child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
              _Settings(dataProvider: dataProvider),
              _ChartView(dataProvider: dataProvider)
            ])));
  }
}

class _Settings extends StatelessWidget {
  final BeeChartDataProvider dataProvider;

  const _Settings({required this.dataProvider});

  @override
  Widget build(BuildContext context) {
    void onTimeRangeChange(String value) {
      context.read<BeeChartDataProvider>().setTimescale(value);
    }

    final List<String> timeVariables = <String>[
      '1 Day',
      '1 Week',
      '1 Month',
    ];

    return Row(
      mainAxisSize: MainAxisSize.min,
      children: <Widget>[
        const Text(
          'Over the last ',
          style: TextStyle(
            color: Colors.black,
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
        SharedDropdownMenu(
          itemList: timeVariables,
          onItemChanged: onTimeRangeChange,
        ),
      ],
    );
  }
}

class _ChartView extends StatefulWidget {
  final BeeChartDataProvider dataProvider;

  const _ChartView({required this.dataProvider});

  @override
  State<StatefulWidget> createState() => _ChartViewState();
}

class _ChartViewState extends State<_ChartView> {
  late double touchedValue;

  final Color? lineColor = Colors.yellow[700];
  final Color? avgLineColor = Colors.yellow[800];
  final Color pointColor = const Color(0xFFFFEB3B);

  @override
  void initState() {
    super.initState();
    touchedValue = -1;
    widget.dataProvider.fetchData();
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
    return Consumer<BeeChartDataProvider>(
      builder:
          (BuildContext context, BeeChartDataProvider value, Widget? child) {
        if (!value.sensorValues.isNotEmpty && !value.isLoading) {
          return const Center(
            child: Text(
              "No data available",
              style: TextStyle(fontSize: 16, color: Colors.grey),
            ),
          );
        }

        if (value.isLoading) {
          return const Center(
            child: CircularProgressIndicator(),
          );
        }

        var yValues = value.sensorValues
            .where((item) =>
                item.sensor_id ==
                1) // Filter only items with sensor_id == 1 //TODO: Make this filter correct
            .map((item) => item.value) // Map
            .toList(); // Convert to list

        // generate a weekday list for the x-axis based on yValues date values
        // _data contains DateTime time, convert it to 'Sun', ect..
        var weekday = [];
        if (value.selectedTimescale == '1 Day') {
          value.sensorValues.toList().forEach((element) {
            if (element.sensor_id != 1) return;
            weekday.add(DateFormat('H').format(element.time));
          });
        } else {
          value.sensorValues.toList().forEach((element) {
            if (element.sensor_id != 1) return;
            weekday.add(DateFormat('E').format(element.time));
          });
        }

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
                            const FlLine(
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
                                    strokeColor: Colors.yellow,
                                  );
                                } else {
                                  return FlDotSquarePainter(
                                    size: 16,
                                    color: Colors.white,
                                    strokeWidth: 5,
                                    strokeColor: Colors.yellow,
                                  );
                                }
                              },
                            ),
                          );
                        }).toList();
                      },
                      touchTooltipData: LineTouchTooltipData(
                        getTooltipColor: (touchedSpot) => Colors.orangeAccent,
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
                              '${weekday[flSpot.x.toInt()]}\n',
                              const TextStyle(
                                color: Colors.black,
                                fontWeight: FontWeight.bold,
                              ),
                              children: [
                                TextSpan(
                                  text: flSpot.y.toString(),
                                  style: const TextStyle(
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
                            flLineStyle: const FlLine(
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
