package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Vinolia-E/BioTree/backend/svgchart"
	"github.com/Vinolia-E/BioTree/backend/util"
)

// ChartRequest represents the request payload for chart generation
type ChartRequest struct {
	DataFile  string `json:"data_file"`
	ChartType string `json:"chart_type"`
	Title     string `json:"title,omitempty"`
	XLabel    string `json:"x_label,omitempty"`
	YLabel    string `json:"y_label,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

// GenerateChartHandler generates SVG charts from processed JSON data
func GenerateChartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ChartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request:", err)
		util.RespondError(w, "Invalid request format")
		return
	}

	// Validate required fields
	if req.DataFile == "" {
		util.RespondError(w, "data_file is required")
		return
	}

	if req.ChartType == "" {
		req.ChartType = "line" // default to line chart
	}

	// Read the JSON data file
	dataPath := filepath.Join("data", req.DataFile)
	jsonData, err := os.ReadFile(dataPath)
	if err != nil {
		log.Println("Failed to read data file:", err)
		util.RespondError(w, "Failed to read data file")
		return
	}

	// Parse JSON data
	var dataPoints []util.DataPoint
	if err := json.Unmarshal(jsonData, &dataPoints); err != nil {
		log.Println("Failed to parse JSON data:", err)
		util.RespondError(w, "Invalid data format")
		return
	}

	// Set up chart options
	var opts []svgchart.Option

	if req.Title != "" {
		opts = append(opts, svgchart.WithTitle(req.Title))
	}

	if req.XLabel != "" {
		opts = append(opts, svgchart.WithXLabel(req.XLabel))
	}

	if req.YLabel != "" {
		opts = append(opts, svgchart.WithYLabel(req.YLabel))
	}

	// Set dimensions (with defaults)
	width := req.Width
	if width == 0 {
		width = 800
	}
	height := req.Height
	if height == 0 {
		height = 400
	}
	opts = append(opts, svgchart.WithDimensions(width, height))

	// Determine chart type
	var chartType svgchart.ChartType
	switch req.ChartType {
	case "bar":
		chartType = svgchart.Bar
	case "line":
		chartType = svgchart.Line
	default:
		chartType = svgchart.Line
	}

	// Generate the chart
	chart, err := svgchart.New(dataPoints, chartType, opts...)
	if err != nil {
		log.Println("Failed to create chart:", err)
		util.RespondError(w, "Failed to create chart")
		return
	}

	svgContent := chart.Generate()

	// Return the SVG content
	response := map[string]interface{}{
		"status": "ok",
		"svg":    svgContent,
		"type":   req.ChartType,
	}

	json.NewEncoder(w).Encode(response)
}

// ListDataFilesHandler returns a list of available processed data files
func ListDataFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Read the data directory
	files, err := os.ReadDir("data")
	if err != nil {
		log.Println("Failed to read data directory:", err)
		util.RespondError(w, "Failed to read data directory")
		return
	}

	var dataFiles []map[string]interface{}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		dataFiles = append(dataFiles, map[string]interface{}{
			"name":     file.Name(),
			"size":     info.Size(),
			"modified": info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	response := map[string]interface{}{
		"status": "ok",
		"files":  dataFiles,
	}

	json.NewEncoder(w).Encode(response)
}

// ProcessAndGenerateHandler combines file processing and chart generation
func ProcessAndGenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Create directories if they don't exist
	if err := os.MkdirAll("files", os.ModePerm); err != nil {
		log.Println("Failed to create files directory:", err)
		util.RespondError(w, "Internal server error")
		return
	}

	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Println("Failed to create data directory:", err)
		util.RespondError(w, "Internal server error")
		return
	}

	// Parse multipart form (10MB max memory)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("Failed to parse form data:", err)
		util.RespondError(w, "Failed to parse form data")
		return
	}

	// Get file from form data
	file, header, err := r.FormFile("document")
	if err != nil {
		log.Println("Failed to retrieve file from form data:", err)
		util.RespondError(w, "Failed to retrieve file")
		return
	}
	defer file.Close()

	// Get chart type from form
	chartType := r.FormValue("chart_type")
	if chartType == "" {
		chartType = "line"
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename)
	inputPath := filepath.Join("files", filename)
	outputPath := filepath.Join("data", filename+".json")

	// Save uploaded file
	if err := saveUploadedFile(file, inputPath); err != nil {
		log.Println("Failed to save file:", err)
		util.RespondError(w, "Failed to save file")
		return
	}

	// Process file through ParseDocumentToJSON
	if err := util.ParseDocumentToJSON(inputPath, outputPath); err != nil {
		log.Println("Failed to parse document:", err)
		util.RespondError(w, "Failed to parse document")
		return
	}

	// Read the processed data
	jsonData, err := os.ReadFile(outputPath)
	if err != nil {
		log.Println("Failed to read processed data:", err)
		util.RespondError(w, "Failed to read processed data")
		return
	}

	var dataPoints []util.DataPoint
	if err := json.Unmarshal(jsonData, &dataPoints); err != nil {
		log.Println("Failed to parse JSON data:", err)
		util.RespondError(w, "Invalid processed data format")
		return
	}

	// Generate chart
	var svgChartType svgchart.ChartType
	switch chartType {
	case "bar":
		svgChartType = svgchart.Bar
	default:
		svgChartType = svgchart.Line
	}

	opts := []svgchart.Option{
		svgchart.WithDimensions(800, 400),
		svgchart.WithTitle("Data Visualization"),
		svgchart.WithXLabel("Categories"),
		svgchart.WithYLabel("Values"),
	}

	chart, err := svgchart.New(dataPoints, svgChartType, opts...)
	if err != nil {
		log.Println("Failed to create chart:", err)
		util.RespondError(w, "Failed to create chart")
		return
	}

	svgContent := chart.Generate()

	// Return success response with SVG
	response := map[string]interface{}{
		"status":     "ok",
		"svg":        svgContent,
		"chart_type": chartType,
		"data_file":  filename + ".json",
	}

	json.NewEncoder(w).Encode(response)
}

// Helper functions
func generateUniqueFilename(originalName string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(originalName)
	base := originalName[:len(originalName)-len(ext)]
	return fmt.Sprintf("%s_%d%s", base, timestamp, ext)
}

func saveUploadedFile(src multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
