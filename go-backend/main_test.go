package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

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

func TestAverageLatency50Iterations(t *testing.T) {
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

	const iterations = 50
	var totalNetwork time.Duration
	var totalParse time.Duration

	fmt.Printf("\nMemulai benchmark 50 iterasi untuk Go...\n")

	for i := 0; i < iterations; i++ {
		_, netTime, parseTime, err := FetchAndParseData(client, cfg.SupabaseURL, cfg.SupabaseAnonKey)
		if err != nil {
			t.Fatalf("Error pada iterasi ke-%d: %v", i+1, err)
		}
		totalNetwork += netTime
		totalParse += parseTime
	}

	avgNetwork := float64(totalNetwork.Milliseconds()) / float64(iterations)
	avgParseMicro := float64(totalParse.Microseconds()) / float64(iterations)
	avgParseMs := avgParseMicro / 1000.0

	fmt.Printf("\n==== HASIL BENCHMARK GO (50 Iterasi) ====\n")
	fmt.Printf("Rata-rata Network Latency: %.2f ms\n", avgNetwork)
	fmt.Printf("Rata-rata Parsing Latency: %.4f ms (%.2f μs)\n", avgParseMs, avgParseMicro)
	fmt.Printf("=========================================\n\n")
}
