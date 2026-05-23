# 📊 Comparative Analysis: Go vs Dart Backend Benchmarks

Penelitian ini membandingkan **efisiensi pengembangan (Development Efficiency)** dan **kecepatan pemrosesan (Processing Speed)** data terstruktur 3-level nested (`Fleet -> Trips -> TripTelemetry`) menggunakan database **Supabase**.

---

## 📁 Struktur Folder (Folderization)

Untuk merapikan proyek dan memisahkan lingkungan pengujian antara Go dan Dart, struktur folder dirancang sebagai berikut:

```text
method/
├── README.md                 # Dokumen petunjuk & tata cara pengujian ini
├── database/
│   └── schema.sql            # Skema DDL tabel database Supabase
├── go-backend/
│   ├── go.mod                # Module Go
│   ├── main.go               # Kode utama Go (parsing & fetch)
│   └── main_test.go          # File benchmark untuk pengujian Go
└── dart-backend/
    ├── pubspec.yaml          # Dependensi Dart
    ├── bin/
    │   └── main.dart         # Kode utama Dart (parsing & fetch)
    └── test/
        └── main_test.dart    # File benchmark untuk pengujian Dart
```

---

## 🛢️ 1. Skema Database (Supabase SQL)
Jalankan DDL SQL berikut di SQL Editor Supabase Anda untuk membuat struktur tabel 3-level nested:

```sql
-- Level 1: Fleet (Armada)
CREATE TABLE fleet (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Level 2: Trips (Perjalanan)
CREATE TABLE trips (
    id SERIAL PRIMARY KEY,
    fleet_id INT REFERENCES fleet(id) ON DELETE CASCADE,
    route_name VARCHAR(100) NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Level 3: Trip Telemetry (Log Sensor)
CREATE TABLE trip_telemetry (
    id SERIAL PRIMARY KEY,
    trip_id INT REFERENCES trips(id) ON DELETE CASCADE,
    latitude NUMERIC(10, 8) NOT NULL,
    longitude NUMERIC(11, 8) NOT NULL,
    speed NUMERIC(5, 2) NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 🚀 2. Cara Menjalankan Pengujian (Testing)

### 🐹 Go Backend

1. Masuk ke folder `go-backend`:
   ```bash
   cd go-backend
   ```
2. Jalankan program biasa untuk tes koneksi:
   ```bash
   go run main.go
   ```
3. Jalankan **Benchmark** (pengujian performa pemrosesan 3-level nested data):
   ```bash
   go test -run=^$ -bench=. -benchtime=10s
   ```

### 🎯 Dart Backend

1. Masuk ke folder `dart-backend`:
   ```bash
   cd dart-backend
   ```
2. Install dependensi:
   ```bash
   dart pub get
   ```
3. Jalankan program biasa:
   ```bash
   dart run bin/main.dart
   ```
4. Jalankan **Benchmark**:
   ```bash
   dart test test/main_test.dart
   ```

---

## 📈 Parameter Pengukuran Penelitian
Berikut adalah metrik yang akan diukur dan dicatat ke dalam draf jurnal Anda setelah pengujian selesai:

| Kategori | Parameter | Deskripsi | Cara Pengukuran |
| :--- | :--- | :--- | :--- |
| **Productivity** | Time-to-Test (TTT) | Waktu yang dihabiskan untuk menulis pengujian | Pengukuran manual dengan Stopwatch (menit) |
| | Lines of Code (LoC) | Baris kode pengujian yang ditulis | VS Code Extension (VS Code Counter) |
| **Performance** | Parsing Latency | Waktu memproses JSON menjadi objek model | Internal timer API (`time.Since` di Go, `Stopwatch` di Dart) |
| | Test Execution Time | Kecepatan eksekusi test runner secara keseluruhan | Waktu penyelesaian command (dalam detik) |
