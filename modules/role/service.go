package role

type Service struct {
	RoleRepo Repository
}

func NewService(roleRepo Repository) *Service {
	return &Service{RoleRepo: roleRepo}
}

func (s *Service) FindByID(id uint) (*Role, error) {
	return s.RoleRepo.FindByID(id)
}
