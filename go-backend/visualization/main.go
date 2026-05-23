package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wcharczuk/go-chart/v2"
)

type LatencyData struct {
	Network []float64 `json:"network"`
	Parsing []float64 `json:"parsing"`
}

func readLatencies(filePath string) (LatencyData, error) {
	var data LatencyData
	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&data)
	return data, err
}

func main() {
	// Ambil folder root method/
	rootPath, err := filepath.Abs("../../")
	if err != nil {
		rootPath = "."
	}

	// ==========================================
	// 1. GENERATE BAR CHARTS (Seperti Sebelumnya)
	// ==========================================

	// Bar Chart Network Latency
	networkChart := chart.BarChart{
		Title: "Rata-rata Network Latency (100 Iterasi)",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 20,
				Left:   20,
				Right:  20,
			},
		},
		Height:   400,
		Width:    600,
		BarWidth: 80,
		Bars: []chart.Value{
			{Value: 94.14, Label: "Go (94.14 ms)"},
			{Value: 103.30, Label: "Dart (103.30 ms)"},
		},
	}

	networkFile := filepath.Join(rootPath, "network_latency.png")
	f1, err := os.Create(networkFile)
	if err == nil {
		_ = networkChart.Render(chart.PNG, f1)
		f1.Close()
		fmt.Printf("Bar Chart Network Latency disimpan: %s\n", networkFile)
	}

	// Bar Chart Parsing Latency
	parsingChart := chart.BarChart{
		Title: "Rata-rata Parsing Latency (100 Iterasi)",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 20,
				Left:   20,
				Right:  20,
			},
		},
		Height:   400,
		Width:    600,
		BarWidth: 80,
		Bars: []chart.Value{
			{Value: 1.8410, Label: "Go (1.84 ms)"},
			{Value: 1.1938, Label: "Dart (1.19 ms)"},
		},
	}

	parsingFile := filepath.Join(rootPath, "parsing_latency.png")
	f2, err := os.Create(parsingFile)
	if err == nil {
		_ = parsingChart.Render(chart.PNG, f2)
		f2.Close()
		fmt.Printf("Bar Chart Parsing Latency disimpan: %s\n", parsingFile)
	}

	// ==========================================
	// 2. GENERATE LINE CHARTS (Berdasarkan Log)
	// ==========================================

	goLat, errGo := readLatencies(filepath.Join(rootPath, "go_latencies.json"))
	dartLat, errDart := readLatencies(filepath.Join(rootPath, "dart_latencies.json"))

	if errGo != nil || errDart != nil {
		fmt.Printf("Warning: File log latensi tidak lengkap (Go err: %v, Dart err: %v). Lewati pembuatan Line Chart.\n", errGo, errDart)
		fmt.Println("Silakan jalankan 'go test' di go-backend dan 'dart test' di dart-backend terlebih dahulu.")
		return
	}

	// Buat nilai sumbu X (1 sampai 100)
	iterations := len(goLat.Network)
	if len(dartLat.Network) < iterations {
		iterations = len(dartLat.Network)
	}
	xValues := make([]float64, iterations)
	for i := 0; i < iterations; i++ {
		xValues[i] = float64(i + 1)
	}

	// Line Chart 1: Network Latency Over Time
	networkLineChart := chart.Chart{
		Title: "Stabilitas Network Latency (100 Iterasi)",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 20,
				Left:   30,
				Right:  20,
			},
		},
		XAxis: chart.XAxis{
			Name:      "Iterasi Ke-",
			NameStyle: chart.Style{TextRotationDegrees: 0},
		},
		YAxis: chart.YAxis{
			Name: "Network Latency (ms)",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "Golang",
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
					StrokeWidth: 2,
				},
				XValues: xValues,
				YValues: goLat.Network[:iterations],
			},
			chart.ContinuousSeries{
				Name: "Dart",
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					StrokeWidth: 2,
				},
				XValues: xValues,
				YValues: dartLat.Network[:iterations],
			},
		},
	}
	networkLineChart.Elements = []chart.Renderable{
		chart.Legend(&networkLineChart),
	}

	netLineFile := filepath.Join(rootPath, "network_latency_line.png")
	f3, err := os.Create(netLineFile)
	if err == nil {
		_ = networkLineChart.Render(chart.PNG, f3)
		f3.Close()
		fmt.Printf("Line Chart Network Latency disimpan: %s\n", netLineFile)
	}

	// Line Chart 2: Parsing Latency Over Time
	parsingLineChart := chart.Chart{
		Title: "Stabilitas Parsing Latency (100 Iterasi)",
		Background: chart.Style{
			Padding: chart.Box{
				Top:    40,
				Bottom: 20,
				Left:   30,
				Right:  20,
			},
		},
		XAxis: chart.XAxis{
			Name: "Iterasi Ke-",
		},
		YAxis: chart.YAxis{
			Name: "Parsing Latency (ms)",
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name: "Golang",
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
					StrokeWidth: 2,
				},
				XValues: xValues,
				YValues: goLat.Parsing[:iterations],
			},
			chart.ContinuousSeries{
				Name: "Dart",
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					StrokeWidth: 2,
				},
				XValues: xValues,
				YValues: dartLat.Parsing[:iterations],
			},
		},
	}
	parsingLineChart.Elements = []chart.Renderable{
		chart.Legend(&parsingLineChart),
	}

	parseLineFile := filepath.Join(rootPath, "parsing_latency_line.png")
	f4, err := os.Create(parseLineFile)
	if err == nil {
		_ = parsingLineChart.Render(chart.PNG, f4)
		f4.Close()
		fmt.Printf("Line Chart Parsing Latency disimpan: %s\n", parseLineFile)
	}
}
