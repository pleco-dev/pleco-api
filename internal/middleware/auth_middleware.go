package middleware

import (
	"pleco-api/internal/httpx"
	"pleco-api/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpx.Respond(c, 401, "error", "missing token", nil, nil, nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		userIDValue, ok := claims["user_id"].(float64)
		if !ok {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		roleValue, ok := claims["role"].(string)
		if !ok {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		tv, ok := claims["tv"]
		if !ok {
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		var accessTokenVersion uint
		switch value := tv.(type) {
		case float64:
			if value < 0 {
				httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
				c.Abort()
				return
			}
			accessTokenVersion = uint(value)
		case uint:
			accessTokenVersion = value
		case int:
			if value < 0 {
				httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
				c.Abort()
				return
			}
			accessTokenVersion = uint(value)
		default:
			httpx.Respond(c, 401, "error", "invalid token", nil, nil, nil)
			c.Abort()
			return
		}

		userID := uint(userIDValue)
		c.Set("user_id", userID)
		c.Set("role", roleValue)
		c.Set("access_token_version", accessTokenVersion)

		c.Next()
	}
}
