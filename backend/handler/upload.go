package handler

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Vinolia-E/BioTree/backend/util"
	"github.com/google/uuid"
)

/*
UploadHandler handles document uploads and converts them to JSON.

- Accepts multipart form data with a "document" file field (PDF, TXT, DOCX).
- Saves the uploaded file in the "files/" directory with a UUID filename.
- Converts the file to structured JSON via ParseDocumentToJSON().
- Saves the output in the "data/" directory.

Responds:
- 200: {"status": "success"}
- 400/500: {"error": "Error message"}

Used by: POST /upload
*/

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		util.RespondError(w, "Failed to parse form data: "+err.Error())
		return
	}

	// Get file from form data
	file, _, err := r.FormFile("document")
	if err != nil {
		log.Println("Failed to retrieve file from form data:", err)
		util.RespondError(w, "Failed to retrieve file: "+err.Error())
		return
	}
	defer file.Close()

	// Generate unique filename
	filename := uuid.New().String()
	inputPath := filepath.Join("files", filename)
	outputPath := filepath.Join("data", filename+".json")

	// Save uploaded file
	dst, err := os.Create(inputPath)
	if err != nil {
		log.Println("Failed to create destination file:", err)
		util.RespondError(w, "Failed to create destination file: "+err.Error())
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		log.Println("Failed to save file:", err)
		util.RespondError(w, "Failed to save file: "+err.Error())
		return
	}

	// Process file through ParseDocumentToJSON
	if err := util.ParseDocumentToJSON(inputPath, outputPath); err != nil {
		log.Println("Failed to parse document:", err)
		util.RespondError(w, "Failed to parse document")
		return
	}

	util.RespondSuccess(w)
}
