package handler

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Vinolia-E/BioTree/backend/svgchart"
	"github.com/Vinolia-E/BioTree/backend/util"
	"golang.org/x/time/rate"
)

const (
	maxFileSize   = 10 << 20 // 10MB
	maxWidth      = 4096
	maxHeight     = 4096
	defaultWidth  = 800
	defaultHeight = 400
	cacheExpiry   = 1 * time.Hour
)

var (
	validChartTypes = map[string]bool{
		"line": true,
		"bar":  true,
		"pie":  true,
	}

	// Cache for generated SVGs
	svgCache = &cache{
		items: make(map[string]cacheItem),
	}

	// Rate limiter: 100 requests per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute/100), 1)
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

type cache struct {
	sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	svg       string
	createdAt time.Time
}

// validate checks if the chart request is valid
func (r *ChartRequest) validate() error {
	if r.DataFile == "" {
		return errors.New("data_file is required")
	}

	if r.ChartType != "" && !validChartTypes[r.ChartType] {
		return fmt.Errorf("invalid chart type: %s", r.ChartType)
	}

	if r.Width < 0 || r.Width > maxWidth {
		return fmt.Errorf("width must be between 0 and %d", maxWidth)
	}

	if r.Height < 0 || r.Height > maxHeight {
		return fmt.Errorf("height must be between 0 and %d", maxHeight)
	}

	return nil
}

// GenerateChartHandler generates SVG charts from processed JSON data
func GenerateChartHandler(w http.ResponseWriter, r *http.Request) {
	// Apply rate limiting
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Set CORS headers if needed
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}
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

	// Validate request
	if err := req.validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		util.RespondError(w, err.Error())
		return
	}

	// Check cache first
	cacheKey := fmt.Sprintf("%s-%s-%d-%d", req.DataFile, req.ChartType, req.Width, req.Height)
	svgCache.RLock()
	if item, ok := svgCache.items[cacheKey]; ok {
		if time.Since(item.createdAt) < cacheExpiry {
			svgCache.RUnlock()
			responseWithCompression(w, r, map[string]interface{}{
				"status": "ok",
				"svg":    item.svg,
				"type":   req.ChartType,
				"cached": true,
			})
			return
		}
	}
	svgCache.RUnlock()

	// Read the JSON data file
	dataPath := filepath.Clean(filepath.Join("data", req.DataFile))
	if !strings.HasPrefix(dataPath, filepath.Clean("data/")) {
		log.Println("Path traversal attempt detected")
		util.RespondError(w, "Invalid file path")
		return
	}

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
	case "pie":
		chartType = svgchart.Pie
	default:
		log.Printf("Invalid chart type requested: %s", req.ChartType)
		util.RespondError(w, fmt.Sprintf("Invalid chart type: %s. Valid types are: line, bar, pie", req.ChartType))
		return
	}

	// Generate the chart
	chart, err := svgchart.New(dataPoints, chartType, opts...)
	if err != nil {
		log.Println("Failed to create chart:", err)
		util.RespondError(w, "Failed to create chart")
		return
	}

	svgContent := chart.Generate()

	// Cache the result
	svgCache.Lock()
	svgCache.items[cacheKey] = cacheItem{
		svg:       svgContent,
		createdAt: time.Now(),
	}
	svgCache.Unlock()

	// Return the SVG content with compression
	responseWithCompression(w, r, map[string]interface{}{
		"status": "ok",
		"svg":    svgContent,
		"type":   req.ChartType,
		"cached": false,
	})
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

// responseWithCompression writes a JSON response with optional gzip compression
func responseWithCompression(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	// Check if client accepts gzip
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		json.NewEncoder(gz).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
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

	// Parse multipart form with size limit
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
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
