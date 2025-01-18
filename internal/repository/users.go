package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lefes/curly-broccoli/internal/domain"
)

type UsersRepo struct {
	db         *sql.DB
	activities *domain.UserActivities
}

func NewUsersRepo(db *sql.DB, a *domain.UserActivities) *UsersRepo {
	return &UsersRepo{
		db:         db,
		activities: a,
	}
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

/* func (r *UsersRepo) UpdateUser(user *domain.User) error { */
/* query := `  */
/* UPDATE users */
/* SET username = ?, role_id = ?, points = ?, respect = ?, daily_messages = ?, last_activity = ? */
/* WHERE discord_id = ?` */

/* _, err := r.db.Exec(query, */
/* user.Username, */
/* user.RoleID, */
/* user.Points, */
/* user.Respect, */
/* user.DailyMessages, */
/* user.LastActivity, */
/* user.DiscordID, */
/* ) */
/* if err != nil { */
/* return fmt.Errorf("failed to update user with DiscordID %s: %w", user.DiscordID, err) */
/* } */

/* return nil */
/* } */

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
	query := "UPDATE users SET daily_messages = ? WHERE discord_id = ?"
	_, err := r.db.Exec(query, dailyMessages, discordID)
	if err != nil {
		return fmt.Errorf("failed to update daily messages for DiscordID %s: %w", discordID, err)
	}
	return nil
}

func (r *UsersRepo) AddOrUpdateUserActivity(userID string) *domain.UserActivity {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	if r.activities.LimitReachedIDs[userID] {
		return nil
	}

	activity, exists := r.activities.Activities[userID]
	if !exists {
		activity = &domain.UserActivity{
			UserID:          userID,
			LastMessageTime: time.Now(),
			NextMessageTime: time.Now(),
			MessageCount:    0,
		}
		r.activities.Activities[userID] = activity
	} else {
		fmt.Println("Updating activity for user", userID)
		activity.LastMessageTime = time.Now()
	}

	return activity
}

func (r *UsersRepo) Reset() []*domain.UserActivity {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	var result []*domain.UserActivity
	for _, activity := range r.activities.Activities {
		result = append(result, activity)
	}

	r.activities.Activities = make(map[string]*domain.UserActivity)
	r.activities.LimitReachedIDs = make(map[string]bool)

	return result
}

func (r *UsersRepo) GetMaxMessages() int {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	return r.activities.MaxMessages
}

func (r *UsersRepo) IsLimitReached(userID string) bool {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	return r.activities.LimitReachedIDs[userID]
}

func (r *UsersRepo) MarkLimitReached(userID string) {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	r.activities.LimitReachedIDs[userID] = true
}
