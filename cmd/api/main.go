package main

import (
	"dairanotes/internal/auth"
	"dairanotes/internal/controller"
	"dairanotes/internal/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	db, err := database.ConnectDatabase()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}

	authController := controller.NewAuthController(db)
	noteController := controller.NewNotesController(db)
	userController := controller.NewUserController(db)

	r.POST("/login", authController.Login)

	noteGroup := r.Group("/notes")
	noteGroup.Use(auth.JwtMiddleware())
	noteGroup.GET("/", noteController.Index)
	noteGroup.POST("/", noteController.Store)
	noteGroup.GET("/:id", noteController.Show)
	noteGroup.PATCH("/:id", noteController.Update)
	noteGroup.DELETE("/:id", noteController.Destroy)

	userGroup := r.Group("/users")
	userGroup.Use(auth.JwtMiddleware())
	userGroup.POST("/", userController.Store)
	userGroup.PATCH("/:id", userController.Update)
	userGroup.DELETE("/:id", userController.Destroy)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server on port 8080
	r.Run(":8080")
}
