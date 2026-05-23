import 'dart:io';
import 'dart:convert';
import 'package:test/test.dart';
import 'package:http/http.dart' as http;
import '../bin/main.dart';

void main() {
  test('Run 100 iterations of data fetch and parse (Sequential)', () async {
    final client = http.Client();
    
    // Load config
    final config = await Config.load();

    const iterations = 100;
    var totalNetworkMs = 0;
    var totalParseUs = 0;

    final networkLog = List<double>.filled(iterations, 0.0);
    final parsingLog = List<double>.filled(iterations, 0.0);

    print('\nMemulai benchmark $iterations iterasi sekuensial untuk Dart...');

    try {
      for (var i = 0; i < iterations; i++) {
        final result = await fetchAndParseData(client, config.supabaseUrl, config.supabaseAnonKey);
        totalNetworkMs += result.networkDuration.inMilliseconds;
        totalParseUs += result.parseDuration.inMicroseconds;

        // Catat latensi milidetik
        networkLog[i] = result.networkDuration.inMicroseconds / 1000.0;
        parsingLog[i] = result.parseDuration.inMicroseconds / 1000.0;
      }

      // Simpan ke file JSON
      final logData = {
        'network': networkLog,
        'parsing': parsingLog,
      };
      final file = File('../dart_latencies.json');
      await file.writeAsString(jsonEncode(logData));

      final avgNetwork = totalNetworkMs / iterations;
      final avgParseUs = totalParseUs / iterations;
      final avgParseMs = avgParseUs / 1000.0;

      print('\n==== HASIL BENCHMARK DART ($iterations Iterasi Sekuensial) ====');
      print('Rata-rata Network Latency: ${avgNetwork.toStringAsFixed(2)} ms');
      print('Rata-rata Parsing Latency: ${avgParseMs.toStringAsFixed(4)} ms (${avgParseUs.toStringAsFixed(2)} μs)');
      print('=======================================================\n');

      expect(iterations, equals(100));
    } catch (e) {
      fail('Benchmark failed: $e');
    } finally {
      client.close();
    }
  });

  test('Run 100 concurrent requests (Parallel/Asynchronous)', () async {
    final client = http.Client();
    final config = await Config.load();
    const iterations = 100;

    print('\nMemulai benchmark $iterations iterasi PARALEL (Event Loop) untuk Dart...');
    final stopwatch = Stopwatch()..start();

    try {
      await Future.wait(
        List.generate(iterations, (i) => fetchAndParseData(client, config.supabaseUrl, config.supabaseAnonKey))
      );

      stopwatch.stop();
      print('\n==== HASIL BENCHMARK DART ($iterations Iterasi Paralel) ====');
      print('Total Waktu Eksekusi Paralel: ${stopwatch.elapsedMilliseconds} ms (${stopwatch.elapsed})');
      print('=====================================================\n');

      expect(iterations, equals(100));
    } catch (e) {
      fail('Concurrent benchmark failed: $e');
    } finally {
      client.close();
    }
  });
}
