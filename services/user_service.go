package services

import (
	"go-auth-app/models"
	"go-auth-app/repositories"
)

type UserService struct {
	UserRepo repositories.UserRepository
}

func (s *UserService) GetAllUsers(page, limit int, search, role string) ([]models.User, int64, error) {
	return s.UserRepo.FindAllWithFilter(page, limit, search, role)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.UserRepo.Delete(id)
}

func (s *UserService) GetProfile(userID uint) (*models.User, error) {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// contoh future logic
	// - audit log
	// - enrich data
	// - cache

	return user, nil
}
