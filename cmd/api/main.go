package main

import (
	"go-api-starterkit/internal/appsetup"
)

func main() {
	if err := appsetup.RunAPI(appsetup.RegisterDocsFromDisk); err != nil {
		panic(err)
	}
}
