import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:dart_backend/models/fleet.dart';

class Config {
  final String supabaseUrl;
  final String supabaseAnonKey;

  Config({required this.supabaseUrl, required this.supabaseAnonKey});

  static Future<Config> load() async {
    // Try ../config.json first (when running from bin/ or test/), then local config.json
    var file = File('../config.json');
    if (!await file.exists()) {
      file = File('config.json');
      if (!await file.exists()) {
        throw Exception('config.json not found!');
      }
    }
    final content = await file.readAsString();
    final json = jsonDecode(content);
    return Config(
      supabaseUrl: json['SUPABASE_URL'] as String,
      supabaseAnonKey: json['SUPABASE_ANON_KEY'] as String,
    );
  }
}

class FetchResult {
  final List<Fleet> fleets;
  final Duration networkDuration;
  final Duration parseDuration;

  FetchResult({
    required this.fleets,
    required this.networkDuration,
    required this.parseDuration,
  });
}

Future<FetchResult> fetchAndParseData(http.Client client, String baseUrl, String anonKey) async {
  final url = Uri.parse('${baseUrl}fleet?select=*,trips(*,trip_telemetry(*))');
  
  final networkStopwatch = Stopwatch()..start();
  final response = await client.get(
    url,
    headers: {
      'apikey': anonKey,
      'Authorization': 'Bearer $anonKey',
    },
  );
  networkStopwatch.stop();

  if (response.statusCode != 200) {
    throw Exception('Failed to load data: ${response.statusCode}');
  }

  final parseStopwatch = Stopwatch()..start();
  final List<dynamic> decodedJson = jsonDecode(response.body);
  final fleets = decodedJson.map((f) => Fleet.fromJson(f as Map<String, dynamic>)).toList();
  parseStopwatch.stop();

  return FetchResult(
    fleets: fleets,
    networkDuration: networkStopwatch.elapsed,
    parseDuration: parseStopwatch.elapsed,
  );
}

void main() async {
  final client = http.Client();
  try {
    final config = await Config.load();
    final result = await fetchAndParseData(client, config.supabaseUrl, config.supabaseAnonKey);
    print('Koneksi sukses!');
    print('Jumlah Fleet Terbaca: ${result.fleets.length}');
    print('Network Latency: ${result.networkDuration.inMilliseconds} ms');
    print('Parsing Latency: ${result.parseDuration.inMicroseconds} μs (${result.parseDuration.inMilliseconds} ms)');

    if (result.fleets.isNotEmpty) {
      final f = result.fleets.first;
      print('\nContoh Data Level 1 (Fleet): ${f.name} (Status: ${f.status})');
      if (f.trips.isNotEmpty) {
        final t = f.trips.first;
        print('  └─ Level 2 (Trip): Route ${t.origin} -> ${t.destination}');
        if (t.tripTelemetry.isNotEmpty) {
          final tel = t.tripTelemetry.first;
          print('      └─ Level 3 (Telemetry): Lat ${tel.sensorData.location.lat}, Lng ${tel.sensorData.location.lng}, Speed ${tel.sensorData.location.speed} km/h');
        }
      }
    }
  } catch (e) {
    print('Error: $e');
  } finally {
    client.close();
  }
}
