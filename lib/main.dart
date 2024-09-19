// import 'dart.io' show Platform; // To determine platform e.g. Platform.iOS

import 'package:flutter/material.dart';
import 'package:flutter/cupertino.dart';

import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import 'providers/beehive_list_provider.dart';
import 'views/overview_page.dart';
import 'views/beehive_detail_page.dart';

import 'utils/helpers.dart';

// GoRouter configuration with initial route and named routes
final GoRouter _router = GoRouter(
  initialLocation: '/',
  routes: [
    GoRoute(
      // Name of the route
      name: 'overview',
      // Path we specify for this route
      path: '/',
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
          return Scaffold(
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
    print(isIOS(context));

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

void main() {
  runApp(const BeehiveApp()); // Entry point for the app
}
