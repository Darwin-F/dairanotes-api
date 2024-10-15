package controller

import (
	"dairanotes/internal/auth"
	"dairanotes/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	ba entities.UserMethodsInterface
}

func NewAuthController(db *sqlx.DB) *AuthController {
	methods := entities.NewUserMethods(db)
	return &AuthController{ba: methods}
}

func (a *AuthController) Login(c *gin.Context) {
	var json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	password, err := a.ba.GetPasswordByUserName(c.Request.Context(), json.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(json.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(json.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}
