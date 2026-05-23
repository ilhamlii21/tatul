class Engine {
  final int rpm;
  final int temp;
  final int fuelLevel;

  Engine({
    required this.rpm,
    required this.temp,
    required this.fuelLevel,
  });

  factory Engine.fromJson(Map<String, dynamic> json) {
    return Engine(
      rpm: json['rpm'] as int,
      temp: json['temp'] as int,
      fuelLevel: json['fuel_level'] as int,
    );
  }
}

class Location {
  final double lat;
  final double lng;
  final double speed;

  Location({
    required this.lat,
    required this.lng,
    required this.speed,
  });

  factory Location.fromJson(Map<String, dynamic> json) {
    return Location(
      lat: (json['lat'] as num).toDouble(),
      lng: (json['lng'] as num).toDouble(),
      speed: (json['speed'] as num).toDouble(),
    );
  }
}

class Diagnostics {
  final List<String> alerts;
  final String systemHealth;

  Diagnostics({
    required this.alerts,
    required this.systemHealth,
  });

  factory Diagnostics.fromJson(Map<String, dynamic> json) {
    final alertsList = json['alerts'] as List? ?? [];
    return Diagnostics(
      alerts: alertsList.map((e) => e.toString()).toList(),
      systemHealth: json['system_health'] as String? ?? 'unknown',
    );
  }
}

class SensorData {
  final Engine engine;
  final Location location;
  final Diagnostics diagnostics;

  SensorData({
    required this.engine,
    required this.location,
    required this.diagnostics,
  });

  factory SensorData.fromJson(Map<String, dynamic> json) {
    return SensorData(
      engine: Engine.fromJson(json['engine'] as Map<String, dynamic>),
      location: Location.fromJson(json['location'] as Map<String, dynamic>),
      diagnostics: Diagnostics.fromJson(json['diagnostics'] as Map<String, dynamic>),
    );
  }
}

// Level 3: Trip Telemetry
class TripTelemetry {
  final String id;
  final String tripId;
  final String recordedAt;
  final SensorData sensorData;

  TripTelemetry({
    required this.id,
    required this.tripId,
    required this.recordedAt,
    required this.sensorData,
  });

  factory TripTelemetry.fromJson(Map<String, dynamic> json) {
    return TripTelemetry(
      id: json['id'] as String,
      tripId: json['trip_id'] as String,
      recordedAt: json['recorded_at'] as String,
      sensorData: SensorData.fromJson(json['sensor_data'] as Map<String, dynamic>),
    );
  }
}

// Level 2: Trips
class Trip {
  final String id;
  final String fleetId;
  final String origin;
  final String destination;
  final String departureTime;
  final String arrivalTime;
  final int passengerCount;
  final String createdAt;
  final List<TripTelemetry> tripTelemetry;

  Trip({
    required this.id,
    required this.fleetId,
    required this.origin,
    required this.destination,
    required this.departureTime,
    required this.arrivalTime,
    required this.passengerCount,
    required this.createdAt,
    required this.tripTelemetry,
  });

  factory Trip.fromJson(Map<String, dynamic> json) {
    final telemetryList = json['trip_telemetry'] as List? ?? [];
    return Trip(
      id: json['id'] as String,
      fleetId: json['fleet_id'] as String,
      origin: json['origin'] as String,
      destination: json['destination'] as String,
      departureTime: json['departure_time'] as String? ?? '',
      arrivalTime: json['arrival_time'] as String? ?? '',
      passengerCount: json['passenger_count'] as int? ?? 0,
      createdAt: json['created_at'] as String,
      tripTelemetry: telemetryList.map((t) => TripTelemetry.fromJson(t as Map<String, dynamic>)).toList(),
    );
  }
}

// Level 1: Fleet
class Fleet {
  final String id;
  final String name;
  final String licensePlate;
  final String vehicleType;
  final String status;
  final String createdAt;
  final List<Trip> trips;

  Fleet({
    required this.id,
    required this.name,
    required this.licensePlate,
    required this.vehicleType,
    required this.status,
    required this.createdAt,
    required this.trips,
  });

  factory Fleet.fromJson(Map<String, dynamic> json) {
    final tripsList = json['trips'] as List? ?? [];
    return Fleet(
      id: json['id'] as String,
      name: json['name'] as String,
      licensePlate: json['license_plate'] as String? ?? '',
      vehicleType: json['vehicle_type'] as String? ?? '',
      status: json['status'] as String,
      createdAt: json['created_at'] as String,
      trips: tripsList.map((t) => Trip.fromJson(t as Map<String, dynamic>)).toList(),
    );
  }
}
