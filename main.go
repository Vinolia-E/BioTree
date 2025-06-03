package main

import (
	"log"
	"net/http"

	"github.com/Vinolia-E/BioTree/backend/route"
)

func main() {
	router := route.InitRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("Starting server on http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
