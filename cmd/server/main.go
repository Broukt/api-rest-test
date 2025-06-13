package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Broukt/api-rest-test/internal/database"
	"github.com/Broukt/api-rest-test/internal/handlers"
	"github.com/Broukt/api-rest-test/internal/models"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	database.AutoMigrate(&models.Product{})

	mux := http.NewServeMux()
	handlers.RegisterProductRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
