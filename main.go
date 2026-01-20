package main

import (
	"leaderboard-backend/handlers"
	"leaderboard-backend/services"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the leaderboard service
	leaderboardService := services.NewLeaderboardService()

	// Seed initial data (10,000 users)
	log.Println("Seeding 10,000 users...")
	count := leaderboardService.SeedUsers(10000, 100, 5000)
	log.Printf("Seeded %d users successfully\n", count)

	// Initialize handlers
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)

	// Set up Gin router
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		// Leaderboard endpoints
		api.GET("/leaderboard", leaderboardHandler.GetLeaderboard)
		api.GET("/search", leaderboardHandler.SearchUsers)
		api.GET("/user/:id/rank", leaderboardHandler.GetUserRank)
		api.POST("/user/:id/rating", leaderboardHandler.UpdateRating)
		
		// Utility endpoints
		api.POST("/seed", leaderboardHandler.SeedUsers)
		api.GET("/stats", leaderboardHandler.GetStats)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
