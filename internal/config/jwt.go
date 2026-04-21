package config

import (
	"log"
)

func mustSecret(key string) []byte {
	secret := GetEnv(key, "")
	if secret == "" {
		log.Fatalf("%s is not set", key)
	}
	return []byte(secret)
}

func MustJWTSecret() []byte {
	return mustSecret("JWT_SECRET")
}
