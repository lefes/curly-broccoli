package repository

import (
	"database/sql"

	"github.com/lefes/curly-broccoli/internal/logging"
)

type RolesRepo struct {
	db     *sql.DB
	logger *logging.Logger
}

func NewRolesRepo(db *sql.DB, l *logging.Logger) *RolesRepo {
	return &RolesRepo{
		db:     db,
		logger: l,
	}
}

func (r *RolesRepo) GetUserRespect(discrodID string) (int, error) {
	query := `
	SELECT respect
	FROM users
	WHERE discord_id = $1
	`

	var respect int
	err := r.db.QueryRow(query, discrodID).Scan(&respect)
	if err != nil {
		return 0, err
	}

	return respect, nil
}

func (r *RolesRepo) GetUserRole(discordID string) (int, error) {
	query := `
	SELECT role_id
	FROM users
	WHERE discord_id = $1
	`

	var roleID int
	err := r.db.QueryRow(query, discordID).Scan(&roleID)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

func (r *RolesRepo) AddUserRespect(discordID string, respect int) error {
	query := `
	UPDATE users
	SET respect = respect + $1
	WHERE discord_id = $2
	`

	_, err := r.db.Exec(query, respect, discordID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepo) AddDayUserRespect(discordID string, respect int) error {
	query := `
	UPDATE users
	SET respect_today = respect_today + $1
	WHERE discord_id = $2
	`

	_, err := r.db.Exec(query, respect, discordID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepo) RemoveUserRespect(discordID string, respect int) error {
	query := `
	UPDATE users
	SET respect = CASE
	    WHEN respect - $1 < 0 THEN 0
	    ELSE respect - $1
	END
	WHERE discord_id = $2
	`

	_, err := r.db.Exec(query, respect, discordID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RolesRepo) RemoveDayUserRespect(discordID string, respect int) error {
	query := `
	UPDATE users
	SET respect = CASE
	    WHEN respect_today - $1 < 0 THEN 0
	    ELSE respect_today - $1
	END
	WHERE discord_id = $2
	`
	_, err := r.db.Exec(query, respect, discordID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RolesRepo) UpdateUserRole(discordID string, roleID int) error {
	query := `
	UPDATE users
	SET role_id = $1
	WHERE discord_id = $2
	`

	_, err := r.db.Exec(query, roleID, discordID)
	if err != nil {
		return err
	}

	return nil
}
