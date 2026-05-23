# 📊 Comparative Analysis: Go vs Dart Backend Benchmarks

Penelitian ini membandingkan **efisiensi pengembangan (Development Efficiency)** dan **kecepatan pemrosesan (Processing Speed)** data terstruktur 3-level nested (`Fleet -> Trips -> TripTelemetry`) menggunakan database **Supabase**.

---

## 📁 Struktur Folder (Folderization)

Struktur proyek saat ini menerapkan pemisahan modular arsitektur bersih (MVC/Clean Architecture style) dan isolasi kredensial:

```text
method/
├── README.md                 # Dokumen petunjuk & tata cara pengujian ini
├── .gitignore                # File Git Ignore untuk mengamankan kredensial
├── config.json               # Konfigurasi kredensial Supabase (DIIABAIKAN OLEH GIT)
├── database/
│   └── schema.sql            # Skema DDL tabel database Supabase (UUID & JSONB)
├── go-backend/
│   ├── go.mod                # Module Go
│   ├── models/
│   │   └── fleet.go          # Model data Go terpisah (Level 1-3)
│   ├── main.go               # Controller utama Go (fetch & print data)
│   ├── main_test.go          # Test suite 50 iterasi & benchmark Go
│   └── visualization/
│       └── main.go           # Script visualisasi penghasil grafik PNG
└── dart-backend/
    ├── pubspec.yaml          # Dependensi Dart
    ├── lib/
    │   └── models/
    │       └── fleet.dart    # Model data Dart terpisah (Level 1-3)
    ├── bin/
    │   └── main.dart         # Controller utama Dart (fetch & print data)
    └── test/
        └── main_test.dart    # Test suite 50 iterasi Dart
```

---

## 🔑 1. Setup Kredensial API

Sebelum menjalankan pengujian, buat sebuah file bernama **`config.json`** di root direktori proyek ini (sejajar dengan `README.md`):

```json
{
  "SUPABASE_URL": "https://<url-supabase-anda>.supabase.co/rest/v1/",
  "SUPABASE_ANON_KEY": "<anon-key-supabase-anda>"
}
```

*Catatan: File ini secara otomatis diabaikan oleh `.gitignore` agar aman dari commit publik.*

---

## 🛢️ 2. Skema Database (Supabase SQL)

Jalankan DDL SQL berikut di SQL Editor Supabase Anda untuk membuat struktur tabel 3-level nested dengan UUID dan data telemetry berbasis JSONB:

```sql
-- Level 1: Fleet (Armada)
CREATE TABLE fleet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    license_plate VARCHAR(20),
    vehicle_type VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Level 2: Trips (Perjalanan)
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID REFERENCES fleet(id) ON DELETE CASCADE,
    origin VARCHAR(100) NOT NULL,
    destination VARCHAR(100) NOT NULL,
    departure_time TIMESTAMP WITH TIME ZONE,
    arrival_time TIMESTAMP WITH TIME ZONE,
    passenger_count INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Level 3: Trip Telemetry (Log Sensor - menggunakan JSONB)
CREATE TABLE trip_telemetry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trip_id UUID REFERENCES trips(id) ON DELETE CASCADE,
    sensor_data JSONB NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## 🚀 3. Cara Menjalankan Pengujian (Testing)

### 🐹 Go Backend

1. Masuk ke folder `go-backend`:
   ```bash
   cd go-backend
   ```
2. Jalankan program biasa untuk tes koneksi dan fetch data tunggal:
   ```bash
   go run main.go
   ```
3. Jalankan **Latency Test (50 Iterasi)** untuk merekam rata-rata waktu parsing dan network:
   ```bash
   go test -v -run=TestAverageLatency50Iterations
   ```
4. Jalankan **Micro-Benchmark** (pengujian performa throughput Go):
   ```bash
   go test -run=^$ -bench=. -benchtime=10s
   ```

### 🎯 Dart Backend

1. Masuk ke folder `dart-backend`:
   ```bash
   cd dart-backend
   ```
2. Install dependensi package:
   ```bash
   dart pub get
   ```
3. Jalankan program biasa untuk tes koneksi:
   ```bash
   dart run bin/main.dart
   ```
4. Jalankan **Latency Test (50 Iterasi)** untuk merekam rata-rata waktu:
   ```bash
   dart test test/main_test.dart
   ```

---

## 📈 4. Visualisasi Grafik

Untuk men-generate file grafik batang (`network_latency.png` dan `parsing_latency.png`), jalankan kode visualisasi berikut:

1. Masuk ke direktori visualization Go:
   ```bash
   cd go-backend/visualization
   ```
2. Jalankan script:
   ```bash
   go run main.go
   ```
Grafik perbandingan akan langsung disimpan di root folder proyek ini.

---

## 📊 Parameter Pengukuran Penelitian
Berikut adalah metrik yang diukur dalam penelitian jurnal Anda:

| Kategori | Parameter | Deskripsi | Cara Pengukuran |
| :--- | :--- | :--- | :--- |
| **Productivity** | Time-to-Test (TTT) | Waktu yang dihabiskan untuk menulis pengujian | Pengukuran manual dengan Stopwatch (menit) |
| | Lines of Code (LoC) | Baris kode pengujian yang ditulis | VS Code Extension (VS Code Counter) |
| **Performance** | Parsing Latency | Waktu memproses JSON menjadi objek model | Internal timer API (`time.Since` di Go, `Stopwatch` di Dart) |
| | Network Latency | Waktu yang dihabiskan selama transfer HTTP data | Selisih waktu mulai kirim s/d penerimaan respons body |
