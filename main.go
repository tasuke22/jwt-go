package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tasuke/go-auth/controllers"
	"github.com/tasuke/go-auth/middlewares"
	"github.com/tasuke/go-auth/models"
)

func main() {
	models.ConnectDataBase()

	r := gin.Default()

	public := r.Group("/api")

	public.POST("/register", controllers.Register)
	public.POST("/login", controllers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)

	r.Run(":8080")
}
