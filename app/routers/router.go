package routers

import (
	"api/app/controllers"
	"api/app/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		// Login route
		v1.POST("/login", controllers.Login)

		// User route
		v1.POST("/users", controllers.CreateUser)

		// Tasks routes
		v1.POST("/tasks", middleware.AuthUser(), controllers.CreateTask)
		v1.GET("/tasks", middleware.AuthUser(), controllers.GetAllTasks)
		v1.GET("/tasks/:id", middleware.AuthUser(), controllers.GetTask)
		v1.PUT("/tasks/:id", middleware.AuthUser(), controllers.UpdateTask)
		v1.DELETE("/tasks/:id", middleware.AuthUser(), controllers.DeleteTasks)
		v1.GET("/user_tasks", middleware.AuthUser(), controllers.GetTasksByUser)
	}
}
