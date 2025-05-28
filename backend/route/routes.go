package route

import (
	"net/http"

	"github.com/Vinolia-E/BioTree/backend/handler"
)

/*
InitRoutes sets up and returns the HTTP route multiplexer (ServeMux) for the BioTree web application.

Routes:
1. /static/ (GET)
   - Serves static assets (CSS, JS, images) from the frontend/static/ directory.
   - Example: /static/style.css loads frontend/static/style.css.

2. / (GET)
   - Serves the main HTML page (frontend/templates/index.html) when users visit the root URL.

3. /upload (POST)
   - Accepts file uploads via multipart/form-data (PDF, DOCX, TXT).
   - Handles uploaded documents and processes them into structured JSON via UploadHandler.

Returns:
- A configured *http.ServeMux router to be passed to http.ListenAndServe.
*/


func InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("frontend/static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/index.html")
	})

	router.HandleFunc("/upload", handler.UploadHandler)

	return router
}
