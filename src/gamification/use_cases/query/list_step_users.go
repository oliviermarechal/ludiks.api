package query

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type StepUser struct {
	ID            string          `json:"id"`
	FullName      string          `json:"fullName"`
	Email         string          `json:"email"`
	Picture       string          `json:"picture"`
	CreatedAt     time.Time       `json:"createdAt"`
	LastLoginAt   time.Time       `json:"lastLoginAt"`
	CurrentStreak int             `json:"currentStreak"`
	LongestStreak int             `json:"longestStreak"`
	Metadata      json.RawMessage `json:"metadata"`
	StartedAt     time.Time       `json:"startedAt"`
}

type PaginatedStepUserResponse struct {
	Data   []StepUser `json:"data"`
	Total  int64      `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}

func ListStepUsers(db *gorm.DB, StepId string, limit int, offset int, filters map[string]string) (PaginatedStepUserResponse, error) {
	var users []StepUser
	var total int64

	baseQuery := db.Table("user_step_progressions usp").
		Joins("LEFT JOIN user_circuit_progressions ucp ON usp.user_circuit_progression_id = ucp.id").
		Joins("LEFT JOIN end_users eu ON ucp.end_user_id = eu.id").
		Where("usp.step_id = ?", StepId).
		Where("usp.status = ?", "in_progress")

	if len(filters) > 0 {
		i := 0
		for key, value := range filters {
			alias := fmt.Sprintf("eum%d", i)
			baseQuery = baseQuery.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM end_user_metadatas %s WHERE %s.end_user_id = eu.id AND %s.key_name = ? AND %s.value = ?)",
				alias, alias, alias, alias), key, value)
			i++
		}
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return PaginatedStepUserResponse{}, err
	}

	err := baseQuery.
		Select(`
			eu.id,
			eu.full_name as "FullName",
			eu.email,
			eu.picture,
			eu.created_at as "CreatedAt",
			eu.last_login_at as "LastLoginAt",
			eu.current_streak as "CurrentStreak",
			eu.longest_streak as "LongestStreak",
			COALESCE(
				(
					SELECT json_object_agg(eum.key_name, eum.value)
					FROM end_user_metadatas eum
					WHERE eum.end_user_id = eu.id
				),
				'{}'::json
			) as "Metadata",
			usp.started_at as "StartedAt"
		`).
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return PaginatedStepUserResponse{}, err
	}

	return PaginatedStepUserResponse{
		Data:   users,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
