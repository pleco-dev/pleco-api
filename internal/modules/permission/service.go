package permission

import (
	"context"
	"fmt"
	"time"

	"pleco-api/internal/cache"
)

type Checker interface {
	HasPermission(roleName, permission string) (bool, error)
}

type Service struct {
	Repo  Repository
	Cache cache.Store
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) HasPermission(roleName, permission string) (bool, error) {
	if roleName == "superadmin" {
		return true, nil
	}

	if s.Cache != nil {
		var allowed bool
		key := fmt.Sprintf("role:permission:%s:%s", roleName, permission)
		if ok, err := s.Cache.GetJSON(context.Background(), key, &allowed); err == nil && ok {
			return allowed, nil
		}
		allowed, err := s.Repo.HasRolePermission(roleName, permission)
		if err != nil {
			return false, err
		}
		_ = s.Cache.SetJSON(context.Background(), key, allowed, 10*time.Minute)
		return allowed, nil
	}

	return s.Repo.HasRolePermission(roleName, permission)
}

func (s *Service) ListAll() ([]Permission, error) {
	return s.Repo.ListAllPermissions()
}

func (s *Service) ListRolePermissionsByName(roleName string) ([]string, error) {
	if s.Cache != nil {
		var permissions []string
		key := fmt.Sprintf("role:permissions:%s", roleName)
		if ok, err := s.Cache.GetJSON(context.Background(), key, &permissions); err == nil && ok {
			return permissions, nil
		}
		found, err := s.Repo.ListRolePermissionsByName(roleName)
		if err != nil {
			return nil, err
		}
		_ = s.Cache.SetJSON(context.Background(), key, found, 10*time.Minute)
		return found, nil
	}

	return s.Repo.ListRolePermissionsByName(roleName)
}
