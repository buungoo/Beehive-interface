import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:beehive/providers/beehive_user_provider.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:beehive/models/beehive_user.dart';
import 'package:shared_preferences/shared_preferences.dart';

// Step 3: Generate the mock class for http.Client
@GenerateMocks([http.Client])
import 'beehive_user_provider_test.mocks.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  setUp(() async {
    // Set up a mock for shared preferences before running any tests
    SharedPreferences.setMockInitialValues({});
  });

  group('BeehiveUserProvider login', () {
    late BeehiveUserProvider beehiveUserProvider;
    late MockClient mockClient;

    setUp(() {
      mockClient = MockClient();
      beehiveUserProvider = BeehiveUserProvider(client: mockClient);
    });

    test('Valid login works', () async {
      // Arrange
      const email = 'aValidUser';
      const password = 'aValidPass';
      final responseBody =
          jsonEncode({'message': 'User Validated', 'token': 'a mock token'});

      when(mockClient.post(
        Uri.parse('https://rockpi.bungos.duckdns.org/login'),
        headers: anyNamed('headers'),
        body: anyNamed('body'),
      )).thenAnswer((_) async => http.Response(responseBody, 200));

      // Act
      final result = await beehiveUserProvider.login(email, password);
      print(result);

      expect(result, isNotNull);
      //expect(result?.message, 'User validated');
    });

    test('should return null when login fails', () async {
      // Arrange
      const email = 'wrong_user';
      const password = 'wrong_pass';
      when(mockClient.post(
        Uri.parse('https://rockpi.bungos.duckdns.org/login'),
        headers: anyNamed('headers'),
        body: anyNamed('body'),
      )).thenAnswer((_) async => http.Response('Unauthorized', 401));

      // Act
      final result = await beehiveUserProvider.login(email, password);

      // Assert
      expect(result, isNull);
    });
  });
}
