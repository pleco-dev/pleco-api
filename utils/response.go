package utils

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(200, gin.H{
		"status": "success",
		"data":   data,
		"meta":   meta,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(400, gin.H{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errors,
	})
}
