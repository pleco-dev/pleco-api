package user

type Service struct {
	UserRepo Repository
}

func NewService(userRepo Repository) *Service {
	return &Service{UserRepo: userRepo}
}

func (s *Service) GetAllUsers(page, limit int, search, role string) ([]User, int64, error) {
	return s.UserRepo.FindAllWithFilter(page, limit, search, role)
}

func (s *Service) DeleteUser(id uint) error {
	return s.UserRepo.Delete(id)
}
