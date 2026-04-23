package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed docs/openapi.yaml docs/swagger-ui.html
var docsFS embed.FS

func init() {
	registerDocsRoutes = func(router *gin.Engine) {
		router.GET("/docs", func(c *gin.Context) {
			page, err := fs.ReadFile(docsFS, "docs/swagger-ui.html")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load Swagger UI"})
				return
			}

			c.Data(http.StatusOK, "text/html; charset=utf-8", page)
		})

		router.GET("/docs/openapi.yaml", func(c *gin.Context) {
			spec, err := fs.ReadFile(docsFS, "docs/openapi.yaml")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load OpenAPI spec"})
				return
			}

			c.Header("Content-Disposition", "inline; filename=openapi.yaml")
			c.Data(http.StatusOK, "text/plain; charset=utf-8", spec)
		})
	}
}
