package controllers

import (
	"go-auth-app/dto"
	"go-auth-app/models"
	"go-auth-app/repositories"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserRepo repositories.UserRepository
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func (a *AuthController) Register(c *gin.Context) {
	var input dto.RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	err := a.UserRepo.Create(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered"})
}

func (a *AuthController) Login(c *gin.Context) {
	var input dto.LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user, err := a.UserRepo.FindByEmail(input.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// generate token (tetap sama seperti sebelumnya)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	accessTokenString, _ := accessToken.SignedString(jwtKey)
	refreshTokenString, _ := refreshToken.SignedString(jwtKey)

	// simpan refresh token via repo
	user.RefreshToken = refreshTokenString
	a.UserRepo.Update(user)

	c.JSON(200, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")

	user, err := a.UserRepo.FindByID(userID.(uint))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user.RefreshToken = ""
	a.UserRepo.Update(user)

	c.JSON(200, gin.H{"message": "Logged out"})
}

func (a *AuthController) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// ambil user dari repository
	user, err := a.UserRepo.FindByID(userID.(uint))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// response (hindari kirim password!)
	c.JSON(200, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

func (a *AuthController) Dashboard(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Admin dashboard"})
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
}
