package middleware

import (
	"fmt"
	"go-auth-app/config"
	"go-auth-app/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing Authorization header"})
			return
		}

		// ✅ harus "Bearer xxx"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		if tokenString == "" || tokenString == "null" || tokenString == "undefined" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token is missing or invalid"})
			return
		}

		// ✅ parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ✅ VALIDASI signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return config.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		// ✅ ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid claims"})
			return
		}

		// ✅ safe parsing
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid user_id"})
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid role"})
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Set("role", role)

		c.Next()
	}
}

func RefreshToken(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Bind request
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Parse token
	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	// 🔍 Ambil user dari DB
	var user models.User
	config.DB.Where("id = ?", claims["user_id"]).First(&user)

	if user.ID == 0 {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// ✅ VALIDASI refresh token harus sama dengan DB
	if user.RefreshToken != body.RefreshToken {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// 🔐 Generate access token baru
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenString, _ := newAccessToken.SignedString(config.JWTSecret)

	// (Optional tapi recommended 🔥)
	// 🔄 Rotate refresh token
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, _ := newRefreshToken.SignedString(config.JWTSecret)

	// Simpan refresh token baru ke DB
	user.RefreshToken = refreshTokenString
	config.DB.Save(&user)

	// Response
	c.JSON(200, gin.H{
		"access_token":  accessTokenString,
		"refresh_token": refreshTokenString,
	})
}
