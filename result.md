# 📝 Hasil Pengujian dan Kesimpulan Penelitian (Result & Conclusion - 100 Iterations)

Dokumen ini menyajikan data empiris hasil pengujian perbandingan performa antara **Golang (Go)** dan **Dart** backend dalam mengambil dan memproses data nested 3-level dari database **Supabase** Cloud dengan muatan **100+ baris data telemetri** menggunakan **100 iterasi**.

---

## 📊 1. Data Hasil Pengujian (Empirical Performance)

Pengujian dilakukan dengan menjalankan 100 iterasi request HTTP GET secara berurutan (Sekuensial) dan secara simultan (Paralel):

| Kategori Pengujian | Parameter Pengukuran | Golang Backend | Dart Backend | Selisih Performa / Catatan | Pemenang |
| :--- | :--- | :---: | :---: | :---: | :---: |
| **Sekuensial** | Rata-rata Network Latency | **98.26 ms** | 103.28 ms | Golang ~5% lebih cepat pada jaringan saat pengujian | **Golang** |
| | Rata-rata Parsing Latency | 0.96 ms (964.08 μs) | **0.96 ms** (955.71 μs) | Keduanya sangat cepat (< 1 ms), Dart unggul tipis (~0.8%) | **Dart** |
| **Paralel (Concurrency)** | Total Waktu Eksekusi (100 Req) | 1,680 ms (1.68 s) | **958 ms** (0.96 s) | Dart ~43% lebih cepat pada Event Loop | **Dart** |

### 🔍 Analisis Mendalam Hasil Pengujian:

1.  **Network Latency**:
    *   **Golang** mencatatkan rata-rata network latency **98.26 ms**.
    *   **Dart** mencatatkan rata-rata **103.28 ms**.
    *   *Penjelasan Ilmiah*: Dalam 100 iterasi sekuensial, performa jaringan sangat bergantung pada latensi routing publik HTTP/HTTPS. Kestabilan kedua backend tergolong sangat baik karena mampu menjaga rata-rata latensi di kisaran ~95-105 ms tanpa adanya lonjakan ekstrem (outlier bernilai detik) pada pengujian kali ini.
2.  **Parsing Latency (JSON Decoding)**:
    *   Pada pengujian ini, **Dart** dan **Golang** menunjukkan kecepatan decoding JSON yang hampir setara, masing-masing **0.96 ms (955.71 μs)** dan **0.96 ms (964.08 μs)**. 
    *   Kedua backend terbukti luar biasa efisien karena mampu memetakan 3-level data relasional yang kompleks ke dalam objek memori terstruktur dalam kurun waktu di bawah 1 milidetik.

3.  **Concurrency (Paralel 100 Request)**:
    *   **Dart (Event Loop / Asynchronous)** menyelesaikan 100 request simultan dalam **958 ms**.
    *   **Golang (Goroutines / Multi-threaded)** menyelesaikan 100 request dalam **1,680 ms**.
    *   *Analisis Concurrency*: Dart memanfaatkan *single-threaded Event Loop* dengan taktik *asynchronous I/O* yang sangat efisien dalam menggunakan kembali (*reuse*) koneksi HTTP Keep-Alive yang sudah terjalin. Sebaliknya, Golang dengan 100 Goroutine membuka banyak koneksi TCP secara paralel ke API Gateway Supabase, yang memicu overhead waktu jabat tangan (*handshake*) SSL/TLS secara simultan di sisi jaringan.

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
    **Golang** dan **Dart** menunjukkan performa yang sangat kompetitif dengan kecepatan decoding JSON kompleks di bawah 1 ms. Golang tetap menjadi pilihan tangguh untuk komputasi berat paralel murni berkat manajemen memori bertipe statis yang sangat efisien pada beban multi-threading.
    
2.  **Untuk Kecepatan Rilis & Pengembangan (Developer-Centric)**:  
    **Dart** menawarkan keuntungan produktivitas yang signifikan jika sistem frontend dibangun menggunakan Flutter. Menggunakan Dart di sisi backend memangkas kurva pembelajaran (*learning curve*) pengembang karena hanya perlu menggunakan satu bahasa pemrograman (*single-language stack*) dengan efisiensi asynchronous I/O yang sangat optimal.

Hasil ini menunjukkan bahwa keputusan pemilihan teknologi backend tidak hanya didasarkan pada metrik kecepatan eksekusi semata, melainkan juga mempertimbangkan trade-off antara throughput pemrosesan data, kompleksitas konkurensi jaringan, dan waktu rilis produk (*time-to-market*).
