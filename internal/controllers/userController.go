package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pattern/internal/models"
	"pattern/internal/repositories"
	"pattern/internal/utils"
	"strconv"
	"time"
)

type AuthController struct {
	userRepository repositories.UserRepository
}

func NewAuthController(userRepository repositories.UserRepository) *AuthController {
	return &AuthController{
		userRepository: userRepository,
	}
}

func (c *AuthController) RegisterUserHandler(ctx *gin.Context) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	hashedPass, _ := utils.HashPassword(user.Password)
	var newUser models.User
	newUser.Username = user.Username
	newUser.Email = user.Email
	newUser.Password = hashedPass
	newUser.Balance = user.Balance
	newUser.Currency = user.Currency
	if err := c.userRepository.CreateUser(&newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	signedToken, _ := utils.CreateToken(string(newUser.ID), newUser.Email)
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/app/v1",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(ctx.Writer, &cookie)
	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "data": signedToken})
}

func (c *AuthController) LoginHandler(ctx *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := c.userRepository.GetUserByEmail(loginRequest.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	if !utils.CheckPasswordHash(loginRequest.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	token, err := utils.CreateToken(strconv.Itoa(int(user.ID)), user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
