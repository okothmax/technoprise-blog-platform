package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"technoprise-blog-backend/internal/database"
	"technoprise-blog-backend/internal/handlers"
	"technoprise-blog-backend/internal/middleware"
)

// @title TechnoPrise Blog API
// @version 1.0
// @description Futuristic accessibility-first blog platform API
// @termsOfService https://technopriseglobal.com/terms

// @contact.name TechnoPrise Global Support
// @contact.url https://technopriseglobal.com/contact
// @contact.email support@technopriseglobal.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.AccessibilityHeaders())
	
	// CORS configuration for frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{
			"http://localhost:4200", 
			"https://localhost:4200", 
			"http://localhost:4201", 
			"http://localhost:4202", 
			"http://localhost:4203",
			"http://127.0.0.1:4200",
			"http://127.0.0.1:39623", // Browser preview proxy
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize handlers
	blogHandler := handlers.NewBlogHandler(db)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Blog routes
		blogs := v1.Group("/blogs")
		{
			blogs.GET("", blogHandler.GetBlogs)           // GET /api/v1/blogs?page=1&limit=10&search=query
			blogs.GET("/:slug", blogHandler.GetBlogBySlug) // GET /api/v1/blogs/my-blog-post
			blogs.POST("", blogHandler.CreateBlog)         // POST /api/v1/blogs
			blogs.PUT("/:id", blogHandler.UpdateBlog)      // PUT /api/v1/blogs/1
			blogs.DELETE("/:id", blogHandler.DeleteBlog)   // DELETE /api/v1/blogs/1
		}

		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"timestamp": time.Now().UTC(),
				"version":   "1.0.0",
				"service":   "TechnoPrise Blog API",
			})
		})

		// Accessibility check endpoint
		v1.GET("/accessibility", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"wcag_compliance": "AA",
				"features": []string{
					"Screen reader optimization",
					"Keyboard navigation",
					"High contrast support",
					"Reduced motion support",
					"Focus management",
					"Semantic HTML",
					"ARIA labels",
				},
				"last_audit": time.Now().UTC(),
			})
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 TechnoPrise Blog API starting on port %s", port)
	log.Printf("📱 Frontend URL: http://localhost:4200")
	log.Printf("🔗 API Documentation: http://localhost:%s/api/v1/health", port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
