package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()

	// Register routes
	AuthRoutes(r)
	EntityRoutes(r)

	return r
}
