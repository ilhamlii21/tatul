package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wcharczuk/go-chart/v2"
)

func main() {
	// Ambil folder root method/
	rootPath, err := filepath.Abs("../../")
	if err != nil {
		rootPath = "."
	}

	// 1. CHART NETWORK LATENCY
	networkChart := chart.BarChart{
		Title: "Rata-rata Network Latency (50 Iterasi)",
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
			{Value: 104.84, Label: "Go (104.84 ms)"},
			{Value: 100.28, Label: "Dart (100.28 ms)"},
		},
	}

	networkFile := filepath.Join(rootPath, "network_latency.png")
	f1, err := os.Create(networkFile)
	if err != nil {
		fmt.Printf("Error create network_latency.png: %v\n", err)
		return
	}
	defer f1.Close()

	err = networkChart.Render(chart.PNG, f1)
	if err != nil {
		fmt.Printf("Error render network chart: %v\n", err)
		return
	}
	fmt.Printf("Visualisasi Network Latency sukses disimpan di: %s\n", networkFile)

	// 2. CHART PARSING LATENCY
	parsingChart := chart.BarChart{
		Title: "Rata-rata Parsing Latency (50 Iterasi)",
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
			{Value: 1.1419, Label: "Go (1.14 ms)"},
			{Value: 1.3616, Label: "Dart (1.36 ms)"},
		},
	}

	parsingFile := filepath.Join(rootPath, "parsing_latency.png")
	f2, err := os.Create(parsingFile)
	if err != nil {
		fmt.Printf("Error create parsing_latency.png: %v\n", err)
		return
	}
	defer f2.Close()

	err = parsingChart.Render(chart.PNG, f2)
	if err != nil {
		fmt.Printf("Error render parsing chart: %v\n", err)
		return
	}
	fmt.Printf("Visualisasi Parsing Latency sukses disimpan di: %s\n", parsingFile)
}
