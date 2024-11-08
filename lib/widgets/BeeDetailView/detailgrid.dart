import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';
import 'package:beehive/models/beehive_data.dart';
import 'package:beehive/widgets/shared.dart';
import 'package:beehive/utils/helpers.dart';
import 'dart:ui';
import 'package:go_router/go_router.dart'; // GoRouter for navigation

class DetailGrid extends StatelessWidget {
  final String id;
  DetailGrid({required this.id, super.key});

  @override
  Widget build(BuildContext context) {
    // Use Consumer to listen for changes from the StreamProvider
    return Consumer<BeehiveData?>(
      builder: (context, beehiveData, child) {
        if (beehiveData == null) {
          return Center(child: SharedLoadingIndicator(context: context));
        }
        return Padding(
          padding: const EdgeInsets.all(8.0),
          child: GridView.builder(
            gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 2,
              crossAxisSpacing: 10.0,
              mainAxisSpacing: 10.0,
              childAspectRatio: 1,
            ),
            itemCount: 5, // Example grid size
            itemBuilder: (context, index) {
              switch (index) {
                case 0:
                  return GestureDetector(
                    child: FrostedGlassBox(
                      title: 'Temperature',
                      value: "${beehiveData.temperature.toString()}°C",
                      colors: [
                        Colors.green.withOpacity(0.2),
                        Colors.orange.withOpacity(0.3),
                        Colors.red.withOpacity(0.2),
                      ],
                    ),
                    onTap: () {
                      context.pushNamed(
                        'testing',
                        pathParameters: {
                          'id': id,
                          'type': "temperature",

                          // Ensure id is a string if needed
                        }, // Use 'pathParameters' to pass the id
                      );
                    },
                  );
                case 1:
                  return GestureDetector(
                    child: FrostedGlassBox(
                      title: 'Weight',
                      value: "${beehiveData.weight.toString()} kg",
                      colors: [
                        Colors.deepPurple.withOpacity(0.2),
                        Colors.blueAccent.withOpacity(0.3),
                        Colors.cyanAccent.withOpacity(0.2),
                      ],
                    ),
                    onTap: () {
                      context.pushNamed(
                        'testing',
                        pathParameters: {
                          'id': id,
                          'type': "weight",

                          // Ensure id is a string if needed
                        }, // Use 'pathParameters' to pass the id
                      );
                    },
                  );
                case 2:
                  return GestureDetector(
                    child: FrostedGlassBox(
                      title: 'Humidity',
                      value: "${beehiveData.humidity.toString()} %",
                      colors: [
                        Colors.blue.withOpacity(0.2),
                        Colors.lightBlue.withOpacity(0.3),
                        Colors.lightBlueAccent.withOpacity(0.2),
                      ],
                    ),
                    onTap: () {
                      context.pushNamed(
                        'testing',
                        pathParameters: {
                          'id': id,
                          'type': "humidity",

                          // Ensure id is a string if needed
                        }, // Use 'pathParameters' to pass the id
                      );
                    },
                  );
                case 3:
                  return GestureDetector(
                    child: FrostedGlassBox(
                      title: 'Oxygen',
                      value: "${beehiveData.ppm.toString()} ppm",
                      colors: [
                        Colors.grey.withOpacity(0.2),
                        Colors.grey.withOpacity(0.3),
                        Colors.grey.withOpacity(0.2),
                      ],
                    ),
                    onTap: () {
                      context.pushNamed(
                        'testing',
                        pathParameters: {
                          'id': id,
                          'type': "oxygen",

                          // Ensure id is a string if needed
                        }, // Use 'pathParameters' to pass the id
                      );
                    },
                  );
                case 4:
                  return GestureDetector(
                    child: FrostedGlassBox(
                      title: 'Battery',
                      value: "98 %",
                      colors: [
                        Colors.green.withOpacity(0.2),
                        Colors.green.withOpacity(0.3),
                        Colors.green.withOpacity(0.2),
                      ],
                    ),
                    onTap: () {
                      context.pushNamed(
                        'testing',
                        pathParameters: {
                          'id': id,
                          'type': "battery",

                          // Ensure id is a string if needed
                        }, // Use 'pathParameters' to pass the id
                      );
                    },
                  );
                default:
                  return FrostedGlassBox(title: 'null', value: 'null');
              }
            },
          ),
        );
      },
    );
  }
}

class _DataBox extends StatelessWidget {
  final String title;
  final String value;

  _DataBox({required this.title, required this.value});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: isDarkMode() ? Colors.yellow[800] : Colors.yellow[100],
        borderRadius: BorderRadius.circular(10.0),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              title,
              style: TextStyle(
                fontSize: 16.0,
                fontWeight: FontWeight.bold,
              ),
            ),
            SizedBox(height: 10.0),
            Text(
              value,
              style: TextStyle(
                fontSize: 20.0,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class FrostedGlassBox extends StatelessWidget {
  String title;
  String value;
  final List<Color> colors;

  FrostedGlassBox(
      {required this.title,
      required this.value,
      this.colors = const [Colors.white30]});

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 300,
      height: 200,
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: colors,
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(15),
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(15),
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 10.0, sigmaY: 10.0),
          child: Container(
              color: Colors.black.withOpacity(0.2),
              alignment: Alignment.center,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      fontSize: 16.0,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  SizedBox(height: 10.0),
                  Text(
                    value,
                    style: TextStyle(
                      fontSize: 20.0,
                    ),
                  ),
                ],
              )),
        ),
      ),
    );
  }
}
