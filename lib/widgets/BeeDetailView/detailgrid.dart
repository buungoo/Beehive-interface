import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'package:provider/provider.dart';
import 'package:beehive/models/beehive_data.dart';
import 'package:beehive/widgets/shared.dart';
import 'package:beehive/utils/helpers.dart';

class DetailGrid extends StatelessWidget {
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
            itemCount: 3, // Example grid size
            itemBuilder: (context, index) {
              switch (index) {
                case 0:
                  return DataBox(
                    title: "Temperature",
                    value: "${beehiveData.temperature}°C",
                  );
                case 1:
                  return DataBox(
                    title: "Weight",
                    value: "${beehiveData.weight}kg",
                  );
                case 2:
                  return DataBox(
                    title: "Humidity",
                    value: "${beehiveData.humidity}%",
                  );
                default:
                  return DataBox(title: "null", value: "null");
              }
            },
          ),
        );
      },
    );
  }
}

class DataBox extends StatelessWidget {
  final String title;
  final String value;

  DataBox({required this.title, required this.value});

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
