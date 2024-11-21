import 'package:beehive/utils/helpers.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('formatHexString', () {
    test('should format hex string with correct formatting', () {
      expect(formatHexString('0080e115000adf82'), '00:80:E1:15:00:0A:DF:82');
    });

    test('should handle string with all zeros', () {
      expect(formatHexString('0000000000000000'), '00:00:00:00:00:00:00:00');
    });

    test('should handle string with mixed digits', () {
      expect(formatHexString('1234567890abcdef'), '12:34:56:78:90:AB:CD:EF');
    });

    test('should handle odd-length hex string', () {
      expect(formatHexString('1'), '');
    });

    test('should handle empty string', () {
      expect(formatHexString(''), '');
    });
  });
}
