import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:beehive/utils/helpers.dart';
import 'package:flutter/material.dart';

class BeeNotification {
  final FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin =
      FlutterLocalNotificationsPlugin();

  int id = 0;

  init(BuildContext context) async {
    if (isIOS(context)) {
      this.askPermIOS();
    }
  }

  askPermIOS() async {
    await flutterLocalNotificationsPlugin
        .resolvePlatformSpecificImplementation<
            IOSFlutterLocalNotificationsPlugin>()
        ?.requestPermissions(
          alert: true,
          badge: true,
          sound: true,
          critical: true,
        );
  }

  askPermAndroid() async {
    final bool granted = await flutterLocalNotificationsPlugin
            .resolvePlatformSpecificImplementation<
                AndroidFlutterLocalNotificationsPlugin>()
            ?.areNotificationsEnabled() ??
        false;

    if (granted) return;

    final AndroidFlutterLocalNotificationsPlugin? androidImplementation =
        flutterLocalNotificationsPlugin.resolvePlatformSpecificImplementation<
            AndroidFlutterLocalNotificationsPlugin>();

    final bool? grantedNotificationPermission =
        await androidImplementation?.requestNotificationsPermission();
  }

  Future<void> sendCriticalNotification(
      {required String title, required String body}) async {
    id++;

    print("Sending Notification...");
    const AndroidNotificationDetails androidNotificationDetails =
        AndroidNotificationDetails(
      'beehive',
      'BeeHive',
      channelDescription: 'Beehive notification',
      importance: Importance.max,
      priority: Priority.high,
      ticker: 'ticker',
      icon: 'ic_launcher',
    );

    const DarwinNotificationDetails iOSNotificationDetails =
        DarwinNotificationDetails(
            presentBanner: true,
            presentList: true,
            presentAlert: true,
            presentSound: true,
            interruptionLevel: InterruptionLevel.critical,
            threadIdentifier: 'beeHive');

    const NotificationDetails notificationDetails = NotificationDetails(
        android: androidNotificationDetails, iOS: iOSNotificationDetails);

    await flutterLocalNotificationsPlugin.show(
        id, title, body, notificationDetails);
  }
}
