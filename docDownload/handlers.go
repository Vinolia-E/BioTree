package main

import "net/http"

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	

	file, header, err := r.FormFile("document")

	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !IsSupportedFile(header.Filename) {
		http.Error(w, "Unsupported file format", http.StatusBadRequest)
		return
	}
}
