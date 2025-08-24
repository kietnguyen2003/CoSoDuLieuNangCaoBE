package main

// package main

// import (
// 	"log"

// 	"clinic-management/internal/config"
// 	"clinic-management/internal/database"
// 	"clinic-management/internal/routes"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	cfg := config.Load()

// 	db, err := database.Connect(cfg.DatabaseURL)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}
// 	defer db.Close()

// 	router := gin.Default()
// 	routes.SetupRoutes(router, db)

// 	log.Printf("Server starting on port %s", cfg.Port)
// 	if err := router.Run(":" + cfg.Port); err != nil {
// 		log.Fatal("Failed to start server:", err)
// 	}
// }
