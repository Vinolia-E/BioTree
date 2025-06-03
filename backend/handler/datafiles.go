package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Vinolia-E/BioTree/backend/util"
)

// FileInfo represents information about a data file
type FileInfo struct {
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
	Units    []string  `json:"units"`
}

// ListDataFilesHandler returns a list of all data files
func ListDataFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers if needed
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Get all files in the data directory
	files, err := os.ReadDir("data")
	if err != nil {
		log.Println("Failed to read data directory:", err)
		util.RespondError(w, "Failed to read data directory")
		return
	}

	var fileInfos []FileInfo

	for _, file := range files {
		// Skip directories and non-JSON files
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Get file info
		filePath := filepath.Join("data", file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		// Get units from file
		units, err := util.GetUnitsFromFile(filePath)
		if err != nil {
			log.Printf("Failed to get units from file %s: %v", file.Name(), err)
			units = []string{} // Empty array if units can't be read
		}

		// Add file info to the list
		fileInfos = append(fileInfos, FileInfo{
			Name:     file.Name(),
			Size:     fileInfo.Size(),
			Modified: fileInfo.ModTime(),
			Units:    units,
		})
	}

	// Return the list of files
	responseWithCompression(w, r, map[string]interface{}{
		"status": "ok",
		"files":  fileInfos,
	})
}
