package route

import (
	"net/http"

	"github.com/Vinolia-E/BioTree/backend/handler"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/index.html")
	})

	r.HandleFunc("/upload", handler.UploadHandler)

	// SVG chart generation endpoints
	r.HandleFunc("/api/generate-chart", handler.GenerateChartHandler)
	r.HandleFunc("/api/data-files", handler.ListDataFilesHandler)
	r.HandleFunc("/api/process-and-generate", handler.ProcessAndGenerateHandler)

	return r
}
