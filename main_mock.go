package main

import (
	"log"

	"clinic-management/internal/config"
	"clinic-management/internal/routes"
	"clinic-management/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Initialize mock data service instead of database connection
	mockService, err := services.NewMockDataService()
	if err != nil {
		log.Fatal("Failed to initialize mock data service:", err)
	}

	router := gin.Default()
	routes.SetupMockRoutes(router, mockService)

	log.Printf("Server starting on port %s (using mock data)", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}