package services

import (
	"sort"

	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type RoleService struct {
	repo   repository.Roles
	logger *logging.Logger
}

func NewRoleService(repo repository.Roles, l *logging.Logger) *RoleService {
	return &RoleService{repo: repo, logger: l}
}

func (s *RoleService) GetUserRole(userID string) (*domain.Role, error) {
	roleID, err := s.repo.GetUserRole(userID)
	if err != nil {
		return nil, err
	}

	for _, role := range domain.Roles {
		if role.ID == roleID {
			return &role, nil
		}
	}

	return nil, nil
}

func (s *RoleService) WillGetPromotion(userID string, respectToAdd int) (bool, *domain.Role) {
	currentRespect, err := s.repo.GetUserRespect(userID)
	if err != nil {
		s.logger.Errorf("Error getting user respect from database: %s", err)
		return false, nil
	}

	newRespect := currentRespect + respectToAdd
	if len(domain.Roles) == 0 {
		s.logger.Errorf("domain.Roles is empty! Cannot determine promotion.")
		return false, nil
	}
	sort.SliceStable(domain.Roles, func(i, j int) bool {
		return domain.Roles[i].RespectRequired < domain.Roles[j].RespectRequired
	})

	var currentRole *domain.Role
	var nextRole *domain.Role

	for _, role := range domain.Roles {
		if currentRespect >= role.RespectRequired {
			currentRole = &role
		}

		if newRespect >= role.RespectRequired && currentRespect < role.RespectRequired {
			nextRole = &role
			break
		}
	}

	if nextRole.ID != 0 {
		s.logger.Infof("User %s will be promoted to %s!", userID, nextRole.Name)
		return true, nextRole
	}

	return false, currentRole
}
