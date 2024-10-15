package controller

import (
	"dairanotes/internal/business"
	"dairanotes/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type UserController struct {
	bu business.UserBusinessInterface
}

func NewUserController(db *sqlx.DB) *UserController {
	methods := entities.NewUserMethods(db)
	bu := business.NewUserBusiness(methods)
	return &UserController{bu: bu}
}

func (uc *UserController) Store(c *gin.Context) {
	var newUser entities.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request",
		})
		return
	}

	err := uc.bu.Store(c.Request.Context(), newUser)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "User created",
	})

	return
}

func (uc *UserController) Update(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request",
		})
		return
	}

	var user entities.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request",
		})
		return
	}

	err = uc.bu.Update(c.Request.Context(), userID, user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User updated",
	})
}

func (uc *UserController) Destroy(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Bad request",
		})
		return
	}

	err = uc.bu.Destroy(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted",
	})
}
