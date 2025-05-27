package route

import (
	"net/http"

	"github.com/Vinolia-E/BioTree/backend/handler"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/templates/index.html")
	})

	r.HandleFunc("/upload", handler.UploadHandler)

	return r
}
