package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Broukt/api-rest-test/internal/database"
	"github.com/Broukt/api-rest-test/internal/handlers"
	"github.com/Broukt/api-rest-test/internal/models"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	database.AutoMigrate(&models.Product{})

	r := gin.Default()
	handlers.RegisterProductRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
