2	Methods

Bagian ini memaparkan secara terperinci mengenai metodologi penelitian yang digunakan untuk membandingkan performa eksekusi dan produktivitas pengembangan antara backend Golang (Go) dan Dart. Alur metodologi penelitian ini dirancang untuk melakukan pengujian empiris secara objektif dan sistematis, mulai dari persiapan basis data cloud, implementasi kode backend, pengujian unit, hingga analisis statistik data latensi. Secara garis besar, alur kerja penelitian ini digambarkan pada Figure 2. Adapun alur pembahasan dalam penelitian ini terbagi menjadi empat tahap utama: komparasi performa sekuensial (latensi jaringan dan parsing JSON), analisis efisiensi pemrosesan paralel (konkurensi 100 request), evaluasi metrik produktivitas pengembang (Time-to-Test dan Lines of Code), serta sintesis trade-off performa-produktivitas untuk merumuskan rekomendasi backend yang optimal.

[Flowchart Metodologi Penelitian]
Figure 2 Alur Metodologi Penelitian Perbandingan Go dan Dart Backend.

## 2.1 Arsitektur Data dan Skema Basis Data
Penelitian ini menggunakan database berbasis Cloud Backend-as-a-Service (BaaS) Supabase sebagai repositori penyimpanan data telemetri. Untuk merepresentasikan kondisi nyata pada sistem pemantauan transportasi armada, dirancang skema data relasional bersarang 3-level (3-level nested relational schema) yang meliputi entitas armada (fleet), perjalanan (trips), dan data sensor telemetri (trip_telemetry).
1. **Level 1 (Fleet)**: Entitas teratas yang memuat informasi identitas armada kendaraan.
2. **Level 2 (Trips)**: Entitas perantara yang berelasi *one-to-many* dengan `fleet` melalui kunci asing `fleet_id`. Entitas ini menyimpan informasi rute perjalanan.
3. **Level 3 (Trip Telemetry)**: Entitas terbawah yang berelasi *one-to-many* dengan `trips` melalui kunci asing `trip_id`. Entitas ini menyimpan log sensor telemetri dalam kolom bertipe data `JSONB` untuk menampung data berstruktur kompleks dan semi-terstruktur.

Struktur skema relasional bersarang 3-level ini disajikan secara mendetail pada Table 2.

Table 2 Skema Struktur Data Relasional Bersarang 3-Level (Fleet -> Trips -> TripTelemetry)
No	Tabel / Kolom	Tipe Data	Deskripsi / Kunci Relasi
1	**fleet**		-	Tabel Armada (Level 1)
	id	UUID	Primary Key, default: gen_random_uuid()
	name	VARCHAR(100)	Nama armada kendaraan
	license_plate	VARCHAR(20)	Nomor pelat kendaraan
	vehicle_type	VARCHAR(50)	Jenis kendaraan (misal: bus, truk)
	status	VARCHAR(50)	Status operasional armada
	created_at	TIMESTAMP WITH TIME ZONE	Waktu pembuatan data armada
2	**trips**		-	Tabel Rute Perjalanan (Level 2)
	id	UUID	Primary Key, default: gen_random_uuid()
	fleet_id	UUID	Foreign Key ke tabel fleet(id) ON DELETE CASCADE
	origin	VARCHAR(100)	Titik keberangkatan perjalanan
	destination	VARCHAR(100)	Titik tujuan perjalanan
	departure_time	TIMESTAMP WITH TIME ZONE	Waktu keberangkatan armada
	arrival_time	TIMESTAMP WITH TIME ZONE	Waktu kedatangan armada
	passenger_count	INT	Jumlah penumpang kendaraan
	created_at	TIMESTAMP WITH TIME ZONE	Waktu pembuatan data trips
3	**trip_telemetry**		-	Tabel Log Sensor (Level 3)
	id	UUID	Primary Key, default: gen_random_uuid()
	trip_id	UUID	Foreign Key ke tabel trips(id) ON DELETE CASCADE
	sensor_data	JSONB	Data sensor telemetri terstruktur kompleks
	recorded_at	TIMESTAMP WITH TIME ZONE	Waktu sensor direkam

Kolom `sensor_data` pada tabel `trip_telemetry` memuat muatan JSON bertingkat yang mencakup tiga objek parameter sensor:
- **Engine**: Kecepatan putaran mesin dalam `rpm` (integer), suhu mesin dalam `temp` (integer), dan persentase bahan bakar dalam `fuel_level` (integer).
- **Location**: Koordinat posisi armada berupa garis lintang `lat` (double), garis bujur `lng` (double), dan kecepatan kendaraan dalam `speed` (double).
- **Diagnostics**: Riwayat alarm sistem dalam `alerts` (array string) dan status kesehatan sistem dalam `system_health` (string).

## 2.2 Framework Pengembangan Backend dan Mekanisme Deserialisasi
Masing-masing ekosistem backend diimplementasikan untuk melakukan request HTTP GET ke endpoint REST API Supabase guna mengambil data armada secara rekursif menggunakan fitur penanaman relasi (*resource embedding*) bawaan PostgREST: `/fleet?select=*,trips(*,trip_telemetry(*))`.
Perbedaan paradigma deserialisasi data nested JSON antara kedua bahasa pemrograman didefinisikan sebagai berikut:

1. **Golang (Go) Backend**:
   Go memanfaatkan representasi tipe data statis berbasis struct yang bersesuaian dengan format JSON data telemetri. Dengan menggunakan tag JSON pada struct (struct tags), compiler Go dapat memetakan payload secara otomatis. Deserialisasi dilakukan menggunakan library standar `encoding/json` melalui fungsi `json.Unmarshal()`. Proses ini bekerja secara reflektif (*runtime reflection*) untuk mencocokkan kunci JSON dengan field struct yang sesuai.

2. **Dart Backend**:
   Dart menerapkan pemetaan objek model secara manual melalui *factory constructor* (`fromJson`). Bodi respons JSON mentah yang diterima pertama-tama didekode menjadi tipe data dinamis `Map<String, dynamic>` menggunakan fungsi `jsonDecode` dari package `dart:convert`. Selanjutnya, pemetaan data ke objek instans kelas (`Fleet`, `Trip`, `TripTelemetry`, `SensorData`, `Engine`, `Location`, dan `Diagnostics`) dieksekusi secara bertahap atau berantai (*cascading decoding*) dengan dukungan fitur penjamin keamanan nilai null (*sound null safety*).

## 2.3 Metrik Pengukuran dan Prosedur Eksperimen
Evaluasi dalam penelitian ini dibagi menjadi dua kategori utama, yaitu produktivitas pengembangan perangkat lunak (Development Efficiency) dan performa eksekusi data (Execution Performance). Spesifikasi metrik evaluasi beserta alat ukur dan metode pengujian selengkapnya ditunjukkan pada Table 3.

Table 3 Spesifikasi Metrik Evaluasi, Alat Ukur, dan Prosedur Pengujian
No	Kategori Metrik	Parameter Pengukuran	Alat Ukur / Cara Pengukuran
1	**Development Efficiency**	Time-to-Test (TTT)	Stopwatch manual (satuan menit) mencatat waktu penulisan test suite
		Lines of Code (LoC)	VS Code Counter Extension menghitung baris kode pengujian bersih
2	**Execution Performance**	Network Latency ($L_{net}$)	Selisih waktu kirim request HTTP hingga terima respons
		Parsing Latency ($L_{parse}$)	Timer internal berpresisi tinggi untuk deserialisasi JSON ke objek memori
		Concurrency Execution Time	Stopwatch total durasi eksekusi untuk 100 request paralel

Prosedur eksperimen untuk mengumpulkan data performa dilakukan dalam dua skenario pengujian yang diulang sebanyak 100 iterasi:

1. **Skenario Sekuensial (Sequential Experiment)**:
   Backend mengirimkan request HTTP GET ke REST API Supabase secara berurutan sebanyak 100 kali. Pada setiap iterasi, durasi transfer jaringan (Network Latency) dan durasi deserialisasi (Parsing Latency) dicatat secara terpisah menggunakan pencatat waktu internal berpresisi tinggi.
   - Pada Go, pengukuran menggunakan fungsi `time.Now()` dan `time.Since()` dari package standard `time`.
   - Pada Dart, pengukuran menggunakan instans kelas `Stopwatch` dari package standard `dart:core`.

   Formulasi untuk menghitung nilai latensi jaringan ($L_{net}$) disajikan dalam Equal. 1:
   L_net = T_recv_response - T_start_request		(1)
   Keterangan:
   - $T_{start\_request}$ adalah waktu saat backend mulai mengirim request HTTP ke Supabase.
   - $T_{recv\_response}$ adalah waktu saat seluruh data respons dari Supabase selesai diterima sepenuhnya oleh backend.
   
   Sementara itu, formulasi untuk menghitung nilai latensi pemrosesan JSON ($L_{parse}$) disajikan dalam Equal. 2:
   L_parse = T_end_parse - T_start_parse		(2)
   Keterangan:
   - $T_{start\_parse}$ adalah waktu tepat sebelum proses konversi (deserialisasi) data JSON dimulai.
   - $T_{end\_parse}$ adalah waktu ketika data JSON telah selesai diubah sepenuhnya menjadi objek model di memori program.

2. **Skenario Paralel/Konkuren (Concurrency Experiment)**:
   Backend mengirimkan 100 request HTTP GET secara simultan untuk menyimulasikan beban kerja multi-koneksi paralel. Pengukuran total waktu eksekusi konkurensi dihitung sejak request pertama dikirim hingga respons terakhir selesai diproses.
   - Go memanfaatkan model konkurensi multi-threading berbasis Goroutine (green threads) yang dikelola oleh Go runtime scheduler, dengan sinkronisasi eksekusi menggunakan `sync.WaitGroup`.
   - Dart memanfaatkan model asinkronisasi berbasis *single-threaded Event Loop* dengan memanfaatkan fungsi `Future.wait()` untuk menunggu penyelesaian seluruh koneksi I/O asinkron.

## 2.4 Alur Analisis Data
Hasil pengukuran latensi dari 100 iterasi disimpan dalam format JSON (`go_latencies.json` dan `dart_latencies.json`). 

Rata-rata latensi jaringan ($\bar{L}_{net}$) dihitung menggunakan rumus pada Equal. 3:
\bar{L}_net = \frac{1}{N} \sum_{i=1}^{N} L_{net, i}		(3)
Keterangan:
- $N$ adalah jumlah total iterasi ($N = 100$).
- $L_{net, i}$ adalah latensi jaringan pada iterasi ke-$i$.

Sedangkan rata-rata latensi pemrosesan JSON ($\bar{L}_{parse}$) dihitung dengan rumus pada Equal. 4:
\bar{L}_parse = \frac{1}{N} \sum_{i=1}^{N} L_{parse, i}		(4)
Keterangan:
- $L_{parse, i}$ adalah latensi pemrosesan JSON pada iterasi ke-$i$.

Data latensi ini selanjutnya dianalisis untuk melihat tingkat kestabilannya. Untuk mempermudah pembacaan, berkas JSON tersebut divisualisasikan menjadi empat grafik (PNG) beresolusi tinggi menggunakan skrip otomatis (`go-backend/visualization/main.go`) berbasis pustaka `github.com/wcharczuk/go-chart/v2`. Grafik ini terdiri dari grafik batang (perbandingan rata-rata) dan grafik garis (tren kestabilan sepanjang iterasi).

Hasil analisis performa eksekusi dan efisiensi penulisan kode ini akan dibahas lebih mendalam pada bagian Hasil dan Pembahasan (*Results and Discussion*) untuk menentukan backend yang paling optimal.
