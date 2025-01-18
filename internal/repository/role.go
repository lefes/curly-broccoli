package repository

import (
	"database/sql"
	"fmt"

	"github.com/lefes/curly-broccoli/internal/domain"
)

type RoleRepo struct {
	db *sql.DB
}

func NewRoleRepo(db *sql.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) PromoteUser(userID string) error {
	currentRoleID, err := r.GetUserRole(userID)
	if err != nil {
		return fmt.Errorf("failed to get user role for ID %s: %w", userID, err)
	}

	var nextRoleID int
	found := false
	for _, roleID := range domain.RolesIDs {
		if roleID > currentRoleID {
			if !found || roleID < nextRoleID {
				nextRoleID = roleID
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("user with ID %s already has the highest role (ID %d)", userID, currentRoleID)
	}

	query := `
		UPDATE users
		SET role_id = ?
		WHERE discord_id = ?
	`
	_, err = r.db.Exec(query, nextRoleID, userID)
	if err != nil {
		return fmt.Errorf("failed to promote user with ID %s: %w", userID, err)
	}

	return nil
}

func (r *RoleRepo) DemoteUser(userID string) error {
	currentRoleID, err := r.GetUserRole(userID)
	if err != nil {
		return fmt.Errorf("failed to get user role for ID %s: %w", userID, err)
	}

	var prevRoleID int
	found := false
	for _, roleID := range domain.RolesIDs {
		if roleID < currentRoleID {
			if !found || roleID > prevRoleID {
				prevRoleID = roleID
				found = true
			}
		}
	}

	if !found {
		return fmt.Errorf("User has already has the lowest role ")
	}

	query := `
		UPDATE users
		SET role_id = ?
		WHERE discord_id = ?
	`
	_, err = r.db.Exec(query, prevRoleID, userID)
	if err != nil {
		return fmt.Errorf("failed to demote user with ID %s: %w", userID, err)
	}

	return nil
}

func (r *RoleRepo) GetUserRole(userID string) (int, error) {
	query := `
		SELECT role_id
		FROM users
		WHERE discord_id = ?
	`

	var roleID int

	err := r.db.QueryRow(query, userID).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user with ID %s not found", userID)
		}
		return 0, fmt.Errorf("failed to get role for user with ID %s: %w", userID, err)
	}

	return roleID, nil
}
