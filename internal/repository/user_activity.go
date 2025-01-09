package repository

import (
	"fmt"
	"time"

	"github.com/lefes/curly-broccoli/internal/domain"
)

type UserActivityRepo struct {
	activities *domain.UserActivities
}

func NewUserActivityRepo(a *domain.UserActivities) *UserActivityRepo {
	return &UserActivityRepo{activities: a}
}

func (r *UserActivityRepo) AddOrUpdateUser(userID string) *domain.UserActivity {
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

func (r *UserActivityRepo) MarkLimitReached(userID string) {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	r.activities.LimitReachedIDs[userID] = true
}

func (r *UserActivityRepo) Reset() []*domain.UserActivity {
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

func (r *UserActivityRepo) GetMaxMessages() int {
	r.activities.Mu.Lock()
	defer r.activities.Mu.Unlock()

	return r.activities.MaxMessages
}
