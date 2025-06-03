package route

import (
	"net/http"

	"github.com/Vinolia-E/BioTree/backend/handler"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the home page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "frontend/templates/home.html")
	})

	// Serve the upload page
	r.HandleFunc("/upload-page", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/upload-page.html")
	})

	// API endpoints
	r.HandleFunc("/upload", handler.UploadHandler)
	r.HandleFunc("/generate-chart", handler.GenerateChartHandler)
	r.HandleFunc("/api/data-files", handler.ListDataFilesHandler)
	r.HandleFunc("/api/generate-chart", handler.GenerateChartHandler)
	r.HandleFunc("/api/process-and-generate", handler.ProcessAndGenerateHandler)

	// Serve the about page
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/about.html")
	})

	return r
}
