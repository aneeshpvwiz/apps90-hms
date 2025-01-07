package routes

import (
	"apps90-hms/controllers"

	"github.com/gin-gonic/gin"
)

func EntityRoutes(r *gin.Engine) {
	entity := r.Group("/entity")
	{
		entity.POST("/", controllers.CreateEntity)
		entity.POST("/user", controllers.CreateUserEntity)
	}
}
