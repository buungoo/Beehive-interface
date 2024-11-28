import 'package:provider/provider.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:beehive/widgets/shared.dart';

class Statusbox extends StatefulWidget {
  final List<Map<String, dynamic>> data;
  const Statusbox({super.key, required this.data});

  @override
  State<Statusbox> createState() => _Statusbox();
}

class _Statusbox extends State<Statusbox> {
  double _sheetPosition = 0.15;
  final double _dragSensitivity = 700;

  @override
  Widget build(BuildContext context) {
    return Container(
        height: 700,
        child: DraggableScrollableSheet(
          initialChildSize: _sheetPosition,
          minChildSize: 0.15,
          builder: (BuildContext context, ScrollController scrollController) {
            return ColoredBox(
              color: Colors.white,
              child: Column(
                children: <Widget>[
                  Grabber(
                    onVerticalDragUpdate: (DragUpdateDetails details) {
                      setState(() {
                        _sheetPosition -= details.delta.dy / _dragSensitivity;
                        if (_sheetPosition < 0.15) {
                          _sheetPosition = 0.15;
                        }
                        if (_sheetPosition > 1.0) {
                          _sheetPosition = 1.0;
                        }
                      });
                    },
                    onVerticalDragEnd: (DragEndDetails details) {
                      if (details.primaryVelocity! < -_dragSensitivity) {
                        setState(() {
                          _sheetPosition = 1.0;
                        });
                      }
                    },
                  ),
                  Flexible(
                      child: Column(
                        children: widget.data.map((issue) {
                          return Status(
                            name: issue['SensorType'],
                            value: issue['Description'],
                            dateTime: DateTime.parse(issue['TimeOfError']),
                            description: issue['Description'],
                          );
                        }).toList(),
                  )),
                ],
              ),
            );
          },
        ));
  }
}

// Widget for Status, should contain a title aka name, value that's weird and a
// date/time and also small description of what it means

class Status extends StatelessWidget {
  final String name;
  final String value;
  final DateTime dateTime;
  final String description;

  const Status(
      {super.key,
      required this.name,
      required this.value,
      required this.dateTime,
      required this.description});

  @override
  Widget build(BuildContext context) {
    final formatter = DateFormat('yyyy-MM-dd');
    String formattedDate = formatter.format(dateTime);

    return Padding(
        padding: EdgeInsets.all(8.0),
        child: Container(
          width: double.infinity,
          padding: EdgeInsets.all(8.0),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(4.0),
            color: Colors.grey[800],
          ),
          child: Column(children: [
            Text(name),
            Text(value, style: TextStyle(color: Colors.red)),
            Text(formattedDate),
            Text(description), // This should be a weird value like "123456789"
          ]),
        ));
  }
}

class Grabber extends StatelessWidget {
  const Grabber({
    super.key,
    required this.onVerticalDragUpdate,
    required this.onVerticalDragEnd,
  });

  final ValueChanged<DragUpdateDetails> onVerticalDragUpdate;
  final ValueChanged<DragEndDetails> onVerticalDragEnd;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onVerticalDragUpdate: onVerticalDragUpdate,
      onVerticalDragEnd: onVerticalDragEnd,
      child: Container(
        width: double.infinity,
        color: Colors.orange,
        child: Align(
          alignment: Alignment.topCenter,
          child: Container(
            margin: const EdgeInsets.symmetric(vertical: 8.0),
            width: 32.0,
            height: 4.0,
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(8.0),
            ),
          ),
        ),
      ),
    );
  }
}
