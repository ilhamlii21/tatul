package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

type LatencyLog struct {
	Network []float64 `json:"network"`
	Parsing []float64 `json:"parsing"`
}

func BenchmarkGetSupabaseAPI(b *testing.B) {
	cfg, err := LoadConfig()
	if err != nil {
		b.Fatalf("Gagal memuat config.json: %v", err)
	}

	// Optimasi HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _, err := FetchAndParseData(client, cfg.SupabaseURL, cfg.SupabaseAnonKey)
		if err != nil {
			b.Fatalf("Request gagal pada perulangan ke-%d: %v", i, err)
		}
	}
}

func TestAverageLatency100Iterations(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Gagal memuat config.json: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	const iterations = 100
	var totalNetwork time.Duration
	var totalParse time.Duration

	logData := LatencyLog{
		Network: make([]float64, iterations),
		Parsing: make([]float64, iterations),
	}

	fmt.Printf("\nMemulai benchmark %d iterasi untuk Go...\n", iterations)

	for i := 0; i < iterations; i++ {
		_, netTime, parseTime, err := FetchAndParseData(client, cfg.SupabaseURL, cfg.SupabaseAnonKey)
		if err != nil {
			t.Fatalf("Error pada iterasi ke-%d: %v", i+1, err)
		}
		totalNetwork += netTime
		totalParse += parseTime

		// Catat waktu dalam milidetik
		logData.Network[i] = float64(netTime.Nanoseconds()) / 1e6
		logData.Parsing[i] = float64(parseTime.Nanoseconds()) / 1e6
	}

	// Simpan ke file JSON
	logBytes, err := json.MarshalIndent(logData, "", "  ")
	if err == nil {
		_ = os.WriteFile("../go_latencies.json", logBytes, 0644)
	}

	avgNetwork := float64(totalNetwork.Milliseconds()) / float64(iterations)
	avgParseMicro := float64(totalParse.Microseconds()) / float64(iterations)
	avgParseMs := avgParseMicro / 1000.0

	fmt.Printf("\n==== HASIL BENCHMARK GO (%d Iterasi Sekuensial) ====\n", iterations)
	fmt.Printf("Rata-rata Network Latency: %.2f ms\n", avgNetwork)
	fmt.Printf("Rata-rata Parsing Latency: %.4f ms (%.2f μs)\n", avgParseMs, avgParseMicro)
	fmt.Printf("=====================================================\n\n")
}

func TestConcurrentLatency100Iterations(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("Gagal memuat config.json: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	const iterations = 100
	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("Memulai benchmark %d iterasi PARALEL (Goroutine) untuk Go...\n", iterations)

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(iter int) {
			defer wg.Done()
			_, _, _, err := FetchAndParseData(client, cfg.SupabaseURL, cfg.SupabaseAnonKey)
			if err != nil {
				t.Errorf("Error pada iterasi paralel ke-%d: %v", iter+1, err)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("\n==== HASIL BENCHMARK GO (%d Iterasi Paralel) ====\n", iterations)
	fmt.Printf("Total Waktu Eksekusi Paralel: %d ms (%v)\n", elapsed.Milliseconds(), elapsed)
	fmt.Printf("===================================================\n\n")
}
