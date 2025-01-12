package routes

import (
	"apps90-hms/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()
	//router.Use(middlewares.APIResponseMiddleware())
	router.Use(middlewares.APIErrorMiddleware())

	// Register routes
	AuthRoutes(router)
	EntityRoutes(router)

	return router
}
