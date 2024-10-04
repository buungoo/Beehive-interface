// import 'dart.io' show Platform; // To determine platform e.g. Platform.iOS

import 'package:beehive/views/initial_page.dart';
import 'package:beehive/views/signup_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';

import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import 'providers/beehive_list_provider.dart';
import 'views/overview_page.dart';
import 'views/beehive_detail_page.dart';
import 'views/login_page.dart';

import 'utils/helpers.dart';
import 'widgets/shared.dart';

import 'package:beehive/services/BeehiveNotificationService.dart';
import 'package:workmanager/workmanager.dart';

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
        final beehive = context.read<BeehiveListProvider>().findBeehiveById(id);

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
  ],
);

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
              theme: CupertinoThemeData(
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
void callbackDispatcher() {
  Workmanager().executeTask((task, inputData) async {
    print("Background task executed: $task");
    // print task and  current time
    print("Task: $task [${DateTime.now()}]");
    BeeNotification().sendCriticalNotification(
      title: "Beehive Background task",
      body: "Task: $task [${DateTime.now()}]",
    );
    return Future.value(true);
  });
}

const simplePeriodicTask = "com.example.beehive.simplePeriodicTask";

void main() {
  WidgetsFlutterBinding.ensureInitialized();

  Workmanager().initialize(
    callbackDispatcher,
    isInDebugMode: true,
  );

  Workmanager().registerPeriodicTask(
    simplePeriodicTask,
    simplePeriodicTask,
    frequency: Duration(minutes: 5),
  );

  Workmanager().printScheduledTasks();

  print("Init");

  runApp(const BeehiveApp()); // Entry point for the app
}
