package main

import (
	"log"
	"os"

	"property-listing/cache"
	"property-listing/db"
	"property-listing/handlers"
	"property-listing/middleware"

	// "property-listing/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.InitDB()

	cache.InitRedis()

	// Import CSV data if needed

	// if os.Getenv("IMPORT_CSV") == "true" {
	// 	err := utils.ImportPropertiesFromCSV("C:/Users/ASHISH TIWARI/Downloads/db424fd9fb74_1748258398689.csv")
	// 	if err != nil {
	// 		log.Printf("Error importing CSV: %v", err)
	// 	}
	// }

	r := gin.Default()

	// Public routes
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Property routes
		authorized.POST("/properties", handlers.CreateProperty)
		authorized.GET("/properties", handlers.GetProperties)
		authorized.GET("/properties/:id", handlers.GetProperty)
		authorized.PUT("/properties/:id", handlers.UpdateProperty)
		authorized.DELETE("/properties/:id", handlers.DeleteProperty)

		// Favorites routes
		authorized.POST("/favorites/:id", handlers.AddFavorite)
		authorized.DELETE("/favorites/:id", handlers.RemoveFavorite)
		authorized.GET("/favorites", handlers.GetFavorites)

		// Recommendations routes
		authorized.POST("/properties/:id/recommend", handlers.RecommendProperty)
		authorized.GET("/recommendations", handlers.GetRecommendations)
		authorized.PUT("/recommendations/:id/read", handlers.MarkRecommendationAsRead)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
