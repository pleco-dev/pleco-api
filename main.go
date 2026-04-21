package main

import (
	"go-api-starterkit/internal/appsetup"

	"github.com/gin-gonic/gin"
)

var registerDocsRoutes = func(_ *gin.Engine) {}

func main() {
	if err := appsetup.RunAPI(registerDocsRoutes); err != nil {
		panic(err)
	}
}
