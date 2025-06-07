package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Vinolia-E/BioTree/backend/route"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := route.InitRoutes()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
