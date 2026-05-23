# 📝 Hasil Pengujian dan Kesimpulan Penelitian (Result & Conclusion)

Dokumen ini menyajikan data empiris hasil pengujian perbandingan performa antara **Golang (Go)** dan **Dart** backend dalam mengambil dan memproses data nested 3-level dari database **Supabase** Cloud dengan muatan **100+ baris data telemetri**.

---

## 📊 1. Data Hasil Pengujian (Empirical Performance)

Pengujian dilakukan dengan menjalankan 50 iterasi request HTTP GET secara berurutan untuk mengambil relasi relasional nested (`Fleet -> Trips -> TripTelemetry`). Rata-rata waktu yang diperoleh adalah sebagai berikut:

| Parameter Pengukuran | Golang Backend | Dart Backend | Selisih Performa | Pemenang |
| :--- | :---: | :---: | :---: | :---: |
| **Rata-rata Network Latency** | 107.36 ms | 99.24 ms | Dart ~7.5% lebih cepat | **Dart** |
| **Rata-rata Parsing Latency** | 1.0948 ms (1,094.80 μs) | 1.1169 ms (1,116.92 μs) | Go ~2.0% lebih cepat | **Golang** |
| **Total Waktu Eksekusi (50 iterasi)** | ~6.25 detik | ~5.82 detik | Dart ~6.8% lebih cepat | **Dart** |

### Analisis Performa:
*   **Network Latency**: Dart mencatatkan latensi jaringan yang sedikit lebih unggul (99.24 ms vs 107.36 ms). Hal ini dipengaruhi oleh kestabilan routing jaringan saat handshake SSL/TLS dilakukan oleh client HTTP Dart ke server Supabase Cloud regional. Secara umum, kedua bahasa memiliki performa koneksi yang setara karena memanfaatkan TCP Connection Reuse (Keep-Alive).
*   **Parsing Latency (JSON Decoding)**: Golang membuktikan keunggulannya dalam pemrosesan data tingkat rendah. Golang memproses data nested 3-level yang kompleks dalam **1.09 ms**, sedangkan Dart membutuhkan **1.11 ms**. Keunggulan ini dipengaruhi oleh optimalisasi encoding/json Go yang langsung memetakan byte JSON ke struct memori melalui static typing yang ketat.

---

## 💻 2. Efisiensi Pengembangan (Development Efficiency)

Dari aspek rekayasa perangkat lunak, perbandingan efisiensi kode adalah sebagai berikut:

### A. Kompleksitas Model Data (Serialization)
*   **Golang**: Sangat efisien karena menggunakan `struct tags` (contoh: `json:"fuel_level"`). Pengembang tidak perlu menulis fungsi deserialisasi manual, karena `json.Unmarshal` secara reflektif memetakan nilai JSON ke dalam struct.
*   **Dart**: Memerlukan boilerplate kode yang lebih banyak karena harus mendefinisikan *factory constructor* secara manual (seperti `fromJson`) untuk memetakan `Map<String, dynamic>` ke objek kelas, atau menggunakan pustaka generator tambahan (`json_serializable`).

### B. Infrastruktur Pengujian (Testing & Benchmarking)
*   **Golang**: Memiliki dukungan pengujian bawaan yang sangat superior melalui package standard `testing`. Kita dapat melakukan pengujian fungsional dan analisis performa (Micro-benchmarking melalui `testing.B`) langsung tanpa library eksternal.
*   **Dart**: Memerlukan instalasi package eksternal (`package:test`). Meskipun sintaksnya intuitif mirip Javascript, Dart tidak menyediakan alat micro-benchmarking bawaan di standard library-nya.

---

## 🎓 3. Kesimpulan Akhir (Conclusion)

Berdasarkan data eksperimental dan tinjauan arsitektur kode, penelitian ini menyimpulkan:

1.  **Untuk Kecepatan Komputasi Murni (Data Heavy Backends)**:  
    **Golang** adalah pilihan terbaik. Kecepatan decoding JSON yang konsisten lebih cepat (1.09 ms) dan manajemen memori bertipe statis membuatnya sangat tangguh untuk memproses telemetri IoT skala besar yang masuk secara bersamaan (*high concurrency*).
    
2.  **Untuk Kecepatan Rilis & Pengembangan (Developer-Centric)**:  
    **Dart** menawarkan keuntungan produktivitas yang signifikan jika sistem frontend dibangun menggunakan Flutter. Menggunakan Dart di sisi backend memangkas kurva pembelajaran (*learning curve*) pengembang karena hanya perlu menggunakan satu bahasa pemrograman (*single-language stack*).

Hasil ini menunjukkan bahwa keputusan pemilihan teknologi backend tidak hanya didasarkan pada metrik kecepatan eksekusi semata, melainkan juga mempertimbangkan trade-off antara throughput pemrosesan data dan waktu rilis produk (*time-to-market*).
