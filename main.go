package main

import (
	"log"
	"os"
	"spotify-clone/database"
	"spotify-clone/handlers"
	"spotify-clone/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize databases
	if err := database.InitMySQL(); err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := database.InitMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := database.InitNeo4j(); err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}

	defer database.Close()

	// Setup Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"mysql":   database.MySQL != nil,
			"mongodb": database.MongoDB != nil,
			"neo4j":   database.Neo4j != nil,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Tracks routes (public read access)
		tracks := v1.Group("/tracks")
		{
			tracks.GET("", handlers.GetTracks)
			tracks.GET("/:id", handlers.GetTrackByID)
			tracks.GET("/:id/similar", handlers.GetSimilarTracks)
			tracks.POST("/add", handlers.AddTrackWithValidation) // Uses stored procedure with validation
		}

		// Artists routes (public read access)
		artists := v1.Group("/artists")
		{
			artists.GET("", handlers.GetArtists)
			artists.GET("/:id", handlers.GetArtistByID)
			artists.GET("/:id/stats", handlers.GetArtistStats) // Uses stored procedure
		}

		// Albums routes (public read access)
		albums := v1.Group("/albums")
		{
			albums.GET("", handlers.GetAlbums)
			albums.GET("/:id/stats", handlers.GetAlbumStats)       // Uses trigger-maintained data
			albums.GET("/:id/duration", handlers.GetAlbumDuration) // Uses SQL function
		}

		// Search (public access)
		v1.GET("/search", handlers.Search)

		// Trending and genre recommendations (public access)
		recommendations := v1.Group("/recommendations")
		{
			recommendations.GET("/trending", handlers.GetTrendingTracks)
			recommendations.GET("/genre/:genre", handlers.GetGenreRecommendations)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			profile := protected.Group("/profile")
			{
				profile.GET("", handlers.GetProfile)
				profile.PUT("/preferences", handlers.UpdatePreferences)
			}

			// User playlists
			playlists := protected.Group("/playlists")
			{
				playlists.POST("", handlers.CreatePlaylist)
				playlists.GET("", handlers.GetUserPlaylists)
				playlists.GET("/:id", handlers.GetPlaylistByID)
				playlists.PUT("/:id", handlers.UpdatePlaylist)
				playlists.DELETE("/:id", handlers.DeletePlaylist)
				playlists.POST("/:id/tracks", handlers.AddTrackToPlaylist)
				playlists.DELETE("/:id/tracks/:trackId", handlers.RemoveTrackFromPlaylist)
			}

			// Recording plays
			protected.POST("/tracks/:id/play", handlers.RecordPlay)

			// Personalized recommendations
			protected.GET("/recommendations", handlers.GetRecommendations)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
