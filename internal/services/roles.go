package services

import (
	"github.com/lefes/curly-broccoli/internal/repository"
)

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) PromoteUser(userID string) error {
	return s.repo.PromoteUser(userID)
}

func (s *RoleService) DemoteUser(userID string) error {
	return s.repo.DemoteUser(userID)
}

func (s *RoleService) GetUserRole(userID string) (int, error) {
	return s.repo.GetUserRole(userID)
}
