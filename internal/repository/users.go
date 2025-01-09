package repository

import (
	"database/sql"
	"fmt"

	"github.com/lefes/curly-broccoli/internal/domain"
)

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) GetAllUsers() ([]*domain.User, error) {
	query := `
		SELECT id, discord_id, username, role_id, points, respect, daily_messages, last_activity
		FROM users
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.DiscordID,
			&user.Username,
			&user.RoleID,
			&user.Points,
			&user.Respect,
			&user.DailyMessages,
			&user.LastActivity,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersRepo) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (discord_id, username, role_id, points, respect, daily_messages)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, user.DiscordID, user.Username, user.RoleID, user.Points, user.Respect, user.DailyMessages)
	if err != nil {
		return fmt.Errorf("failed to create user with DiscordID %s: %w", user.DiscordID, err)
	}

	return nil
}

func (r *UsersRepo) DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user with ID %s: %w", userID, err)
	}

	return nil
}

func (r *UsersRepo) GetUserByDiscordID(discordID string) (*domain.User, error) {
	query := `
		SELECT id, discord_id, username, role_id, points, respect, daily_messages, last_activity
		FROM users
		WHERE discord_id = ?
	`

	row := r.db.QueryRow(query, discordID)

	user := &domain.User{}

	err := row.Scan(
		&user.ID,
		&user.DiscordID,
		&user.Username,
		&user.RoleID,
		&user.Points,
		&user.Respect,
		&user.DailyMessages,
		&user.LastActivity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by DiscordID %s: %w", discordID, err)
	}

	return user, nil
}

func (r *UsersRepo) UpdateUser(user *domain.User) error {
	query := ` 
		UPDATE users
		SET username = ?, role_id = ?, points = ?, respect = ?, daily_messages = ?, last_activity = ?
		WHERE discord_id = ?`

	_, err := r.db.Exec(query,
		user.Username,
		user.RoleID,
		user.Points,
		user.Respect,
		user.DailyMessages,
		user.LastActivity,
		user.DiscordID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user with DiscordID %s: %w", user.DiscordID, err)
	}

	return nil
}

func (r *UsersRepo) UpdateUserRole(discordID string, roleID int) error {
	query := "UPDATE users SET role_id = ? WHERE discord_id = ?"
	_, err := r.db.Exec(query, roleID, discordID)
	if err != nil {
		return fmt.Errorf("failed to update role for DiscordID %s: %w", discordID, err)
	}
	return nil
}

func (r *UsersRepo) UpdateUserPoints(discordID string, points int) error {
	query := "UPDATE users SET points = points + ? WHERE discord_id = ?"
	_, err := r.db.Exec(query, points, discordID)
	if err != nil {
		return fmt.Errorf("failed to update points for DiscordID %s: %w", discordID, err)
	}
	return nil
}

func (r *UsersRepo) UpdateUserRespect(discordID string, respect int) error {
	query := "UPDATE users SET respect = respect + ? WHERE discord_id = ?"
	_, err := r.db.Exec(query, respect, discordID)
	if err != nil {
		return fmt.Errorf("failed to update respect for DiscordID %s: %w", discordID, err)
	}
	return nil
}

func (r *UsersRepo) UpdateUserDailyMessages(discordID string, dailyMessages int) error {
	query := "UPDATE users SET daily_messages = daily_messages + ? WHERE discord_id = ?"
	_, err := r.db.Exec(query, dailyMessages, discordID)
	if err != nil {
		return fmt.Errorf("failed to update daily messages for DiscordID %s: %w", discordID, err)
	}
	return nil
}
