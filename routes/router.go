package routes

import (
	"apps90-hms/middlewares"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()
	//router.Use(middlewares.APIResponseMiddleware())
	router.Use(middlewares.APIErrorMiddleware())

	// Serve static images from the "assets/images" folder
	router.Static("/images", "./assets/images")

	// CORS Middleware
	FrontendUrl := os.Getenv("FRONTEND_URL")
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{FrontendUrl}, // Allow frontend domain
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // If using cookies or auth headers
		MaxAge:           12 * time.Hour,
	}))

	// Register routes
	AuthRoutes(router)
	EntityRoutes(router)
	PatientRoutes(router)

	return router
}
