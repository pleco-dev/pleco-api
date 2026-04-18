package permission

type Checker interface {
	HasPermission(roleName, permission string) (bool, error)
}

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) HasPermission(roleName, permission string) (bool, error) {
	if roleName == "superadmin" {
		return true, nil
	}

	return s.Repo.HasRolePermission(roleName, permission)
}
