package main

import (
	"BioTree/svgchart"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Set up MIME types
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".html", "text/html")

	// Create a custom file server that sets the correct headers
	fileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove the X-Content-Type-Options header
		w.Header().Del("X-Content-Type-Options")
		
		// Set the correct content type based on file extension
		path := r.URL.Path
		if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html")
		}
		
		// Strip the /static/ prefix and serve the file
		http.StripPrefix("/static/", fileServer).ServeHTTP(w, r)
	}))

	// API endpoint to generate charts
	http.HandleFunc("/api/generate-chart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			ChartType string             `json:"chartType"`
			Data      []svgchart.DataPoint `json:"data"`
			Options   struct {
				Title   string `json:"title"`
				XLabel  string `json:"xLabel"`
				YLabel  string `json:"yLabel"`
				Width   int    `json:"width"`
				Height  int    `json:"height"`
				ShowGrid bool   `json:"showGrid"`
			} `json:"options"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert data to map for the chart library
		dataMap := make(map[string]float64)
		for _, dp := range request.Data {
			dataMap[dp.Label] = dp.Value
		}

		// Create chart
		chart, err := svgchart.New(
			dataMap,
			svgchart.ChartType(request.ChartType),
			svgchart.WithTitle(request.Options.Title),
			svgchart.WithXLabel(request.Options.XLabel),
			svgchart.WithYLabel(request.Options.YLabel),
			svgchart.WithDimensions(request.Options.Width, request.Options.Height),
			svgchart.WithGrid(request.Options.ShowGrid),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write([]byte(chart.Generate()))
	})

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request for path: %s\n", r.URL.Path)
		
		// Only serve index.html for the root path
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		
		indexPath := "../web/static/index.html"
		fmt.Printf("Attempting to serve: %s\n", indexPath)
		
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			fmt.Printf("File not found: %s\n", indexPath)
			http.Error(w, "Index file not found", http.StatusInternalServerError)
			return
		}
		
		http.ServeFile(w, r, indexPath)
	})

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
