package auth

import (
	"gorm.io/gorm"
)

// AuthRepository defines repository-level methods for auth module (currently empty, extend as needed)
type AuthRepository interface {
	// add method signatures if needed
}

type authRepo struct {
	db *gorm.DB
}

// NewRepository returns a new implementation of AuthRepository
func NewRepository(db *gorm.DB) AuthRepository {
	return &authRepo{db: db}
}
