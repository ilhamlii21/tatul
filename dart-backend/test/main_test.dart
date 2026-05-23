import 'package:test/test.dart';
import 'package:http/http.dart' as http;
import '../bin/main.dart';

void main() {
  test('Run 50 iterations of data fetch and parse', () async {
    final client = http.Client();
    
    // Load config
    final config = await Config.load();

    const iterations = 50;
    var totalNetworkMs = 0;
    var totalParseUs = 0;

    print('\nMemulai benchmark 50 iterasi untuk Dart...');

    try {
      for (var i = 0; i < iterations; i++) {
        final result = await fetchAndParseData(client, config.supabaseUrl, config.supabaseAnonKey);
        totalNetworkMs += result.networkDuration.inMilliseconds;
        totalParseUs += result.parseDuration.inMicroseconds;
      }

      final avgNetwork = totalNetworkMs / iterations;
      final avgParseUs = totalParseUs / iterations;
      final avgParseMs = avgParseUs / 1000.0;

      print('\n==== HASIL BENCHMARK DART (50 Iterasi) ====');
      print('Rata-rata Network Latency: ${avgNetwork.toStringAsFixed(2)} ms');
      print('Rata-rata Parsing Latency: ${avgParseMs.toStringAsFixed(4)} ms (${avgParseUs.toStringAsFixed(2)} μs)');
      print('===========================================\n');

      expect(iterations, equals(50));
    } catch (e) {
      fail('Benchmark failed: $e');
    } finally {
      client.close();
    }
  });
}
