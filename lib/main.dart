// import 'dart.io' show Platform;

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import 'providers/beehive_list_provider.dart';
import 'views/overview_page.dart';
import 'views/beehive_detail_page.dart';

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
        // Query parameters for widgets that need it
        final String id = state.pathParameters['id']!;
        final beehive = context.read<BeehiveListProvider>().findBeehiveById(id);

        // If beehive is not found, display an error message
        if (beehive == null) {
          return Scaffold(
            appBar: AppBar(
              title: const Text('Error'),
              leading: IconButton(
                icon: const Icon(Icons.arrow_back),
                onPressed: () => context.go('/'), // Navigate to overview page
              ),
            ),
            body: const Center(child: Text('Beehive not found!')),
          );
        }

        // Wrap the detail page with a Scaffold and an AppBar
        return Scaffold(
          appBar: AppBar(
            title: const Text('Beehive Details'),
            leading: IconButton(
              icon: const Icon(Icons.arrow_back),
              onPressed: () {
                context.go('/'); // Navigate back to the overview page
              },
            ),
          ),
          body: BeehiveDetailPage(beehive: beehive),
        );
      },
    ),
  ],
);

// Main app class that sets up providers and routing
class BeehiveApp extends StatelessWidget {
  const BeehiveApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      // The providers that are specified here are globally available
      // This means every widget in the app can listen to them
      providers: [
        ChangeNotifierProvider(create: (context) => BeehiveListProvider()),
      ],
      child: MaterialApp.router(
        routerConfig: _router, // Pass the GoRouter configuration
        title: 'Beehive App', // App title
      ),
    );
  }
}

void main() {
  runApp(const BeehiveApp()); // Entry point for the app
}
