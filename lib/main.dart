// import 'dart.io' show Platform; // To determine platform e.g. Platform.iOS

import 'package:beehive/views/initial_page.dart';
import 'package:beehive/views/signup_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';
import 'dart:io';
import 'package:flutter/services.dart';

import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import 'providers/beehive_list_provider.dart';
import 'views/overview_page.dart';
import 'views/beehive_detail_page.dart';
import 'views/login_page.dart';
import 'package:beehive/views/detail_chart_page.dart';
import 'package:beehive/views/camera.dart';

import 'package:shared_preferences/shared_preferences.dart';

import 'utils/helpers.dart';
import 'widgets/shared.dart';

import 'package:beehive/models/beehive.dart';
import 'package:beehive/services/BeehiveNotificationService.dart';
import 'package:workmanager/workmanager.dart';
import 'dart:convert';
import 'dart:typed_data';

// GoRouter configuration with initial route and named routes
final GoRouter _router = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
        name: 'initial-page',
        path: '/',
        builder: (context, state) => const InitialPage()),
    GoRoute(
        name: 'signup-page',
        path: '/signup_page',
        builder: (context, state) => const SignupPage()),
    GoRoute(
        name: 'login page',
        path: '/login_page',
        builder: (context, state) => const LoginPage()),
    GoRoute(
      // Name of the route
      name: 'overview',
      // Path we specify for this route
      path: '/overview',
      // Widget that we bind to the path
      builder: (context, state) => const OverviewPage(),
    ),
    GoRoute(
      name: 'beehive-detail',
      path: '/beehive/:id',
      builder: (context, state) {
        // Retrieve the path parameter 'id'
        final String id = state.pathParameters['id']!;
        // Fetch the beehive from the provider
        Beehive? beehive =
            context.read<BeehiveListProvider>().findBeehiveById(id);

        // If beehive is not found, display an error message
        if (beehive == null) {
          return SharedScaffold(
            context: context,
            appBar: AppBar(
              title: const Text('Error'),
              leading: IconButton(
                icon: const Icon(Icons.arrow_back),
                // Navigate to overview page using go to get a clean stack
                onPressed: () => context.go('/'),
              ),
            ),
            body: const Center(child: Text('Beehive not found!')),
          );
        }

        // If beehive is found, return the beehive detail page
        return BeehiveDetailPage(beehive: beehive);
      },
    ),
    GoRoute(
        name: "Camera",
        path: "/camera",
        builder: (context, state) {
          return const Camera();
        }),
    GoRoute(
      name: "testing",
      path: '/beehive/test/:id/:type',
      builder: (context, state) {
        final String id = state.pathParameters['id']!;
        final String type = state.pathParameters['type']!;

        final beehive = context.read<BeehiveListProvider>().findBeehiveById(id);

        if (beehive == null) {
          return SharedScaffold(
            context: context,
            appBar: AppBar(
              title: const Text('Error'),
              leading: IconButton(
                icon: const Icon(Icons.arrow_back),
                onPressed: () => context.go('/'),
              ),
            ),
            body: const Center(child: Text('Beehive not found!')),
          );
        }

        return BeeChartPage(
            beehive: beehive,
            title: type[0].toUpperCase() + type.substring(1),
            type: type);
      },
    ),
  ],
);

class DevHttpOverrides extends HttpOverrides {
  final String pemString;

  DevHttpOverrides(this.pemString);

  @override
  HttpClient createHttpClient(final SecurityContext? context) {
    return super.createHttpClient(context)
      ..badCertificateCallback = (X509Certificate cert, String host, int port) {
        return pemString.compareTo(cert.pem) == 1;
      };
  }
}

// Main app class that sets up providers and routing
class BeehiveApp extends StatelessWidget {
  const BeehiveApp({super.key});

  @override
  Widget build(BuildContext context) {
    // Ask for perm
    BeeNotification().init(context);

    return MultiProvider(
      // The providers that are specified here are globally available
      // This means every widget in the app can listen to them
      providers: [
        // See overview_page for use of this global provider
        ChangeNotifierProvider(create: (context) => BeehiveListProvider()),
      ],
      child: isIOS(context)
          ? CupertinoApp.router(
              localizationsDelegates: const <LocalizationsDelegate>[
                DefaultMaterialLocalizations.delegate,
                DefaultWidgetsLocalizations.delegate,
                DefaultCupertinoLocalizations.delegate,
              ],
              routerConfig: _router, // Pass the GoRouter configuration
              title: 'Beehive App', // App title
              theme: const CupertinoThemeData(
                primaryColor: CupertinoColors.systemYellow,
              ),
            )
          : MaterialApp.router(
              routerConfig: _router, // Pass the GoRouter configuration
              title: 'Beehive App', // App title
              theme: ThemeData(
                primarySwatch: Colors.yellow,
              ),
            ),
    );
  }
}

@pragma('vm:entry-point')
void callbackDispatcher() async {
  print("callbackDispatcher was called");
  Workmanager().executeTask((task, inputData) {
    try {
      return BeeNotification().checkIssues();
    } catch (e) {
      print(e);
      BeeNotification().sendCriticalNotification(
          title: "Unable to contact RockPI",
          body: "Unable to fetch latest status from RockPI");
      return Future.value(false);
    }

    //return Future.value(true);
  });
}

const simplePeriodicTask = "com.example.beehive.simplePeriodicTask";

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  ByteData data = await PlatformAssetBundle().load('assets/ca/sigma.pem');
  String pemString = utf8.decode(data.buffer.asUint8List());
  //print(pemString);

  HttpOverrides.global = DevHttpOverrides(pemString);
// Replace with: https://stackoverflow.com/a/69481863*/

  Workmanager().initialize(
    callbackDispatcher,
    isInDebugMode: false,
  );
  Workmanager().registerPeriodicTask(
    simplePeriodicTask,
    simplePeriodicTask,
    initialDelay: const Duration(seconds: 30),
    //frequency: config.bgWorkerFetchRate,
    constraints: Constraints(
      networkType: NetworkType.connected,
    ),
  );

  // one off task
  /*Workmanager().registerOneOffTask("com.example.beehive.rescheduledTask",
      "com.example.beehive.rescheduledTask",
      initialDelay: const Duration(seconds: 10),
      constraints: Constraints(
        networkType: NetworkType.connected,
      ));*/

  //Workmanager().printScheduledTasks();

  runApp(const BeehiveApp()); // Entry point for the app
}
