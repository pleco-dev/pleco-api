package appsetup

import (
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterDocsFromDisk(router *gin.Engine) {
	router.GET("/docs", func(c *gin.Context) {
		page, err := os.ReadFile("docs/swagger-ui.html")
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to load Swagger UI"})
			return
		}

		c.Data(200, "text/html; charset=utf-8", page)
	})

	router.GET("/docs/openapi.yaml", func(c *gin.Context) {
		spec, err := os.ReadFile("docs/openapi.yaml")
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to load OpenAPI spec"})
			return
		}

		c.Data(200, "application/yaml; charset=utf-8", spec)
	})
}
