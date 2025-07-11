package query

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CircuitProgress struct {
	ID          string     `json:"id"`
	CircuitName string     `json:"circuitName"`
	Status      string     `json:"status"`
	Points      int        `json:"points"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}

type EndUserResponse struct {
	ID              string                 `json:"id"`
	Email           string                 `json:"email"`
	FullName        string                 `json:"fullName"`
	CreatedAt       time.Time              `json:"createdAt"`
	LastLoginAt     time.Time              `json:"lastLoginAt"`
	ExternalId      string                 `json:"externalId"`
	Picture         *string                `json:"picture"`
	CurrentStreak   int                    `json:"currentStreak"`
	LongestStreak   int                    `json:"longestStreak"`
	Metadata        map[string]interface{} `json:"metadata"`
	CircuitProgress []CircuitProgress      `json:"circuitProgress"`
}

type Pagination struct {
	Limit  int
	Offset int
}

type PaginatedEndUserResponse struct {
	Total int                `json:"total"`
	Users []*EndUserResponse `json:"users"`
}

type MetadataItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type endUserRow struct {
	ID              string
	Email           string
	FullName        string
	CreatedAt       time.Time
	LastLoginAt     time.Time
	ExternalId      string
	Picture         *string
	CurrentStreak   int
	LongestStreak   int
	Metadata        json.RawMessage
	CircuitProgress json.RawMessage
}

func ListProjectEndUserQuery(
	db *gorm.DB,
	projectId string,
	pagination Pagination,
	metadataFilters map[string]string,
	circuitId string,
	circuitStep string,
	term string,
) (*PaginatedEndUserResponse, error) {
	baseQuery := db.Table("end_users eu").Where("eu.project_id = ?", projectId)

	if circuitId != "" {
		if circuitStep != "" {
			switch circuitStep {
			case "end":
				baseQuery = baseQuery.Where("EXISTS (SELECT 1 FROM user_circuit_progressions ucp WHERE ucp.end_user_id = eu.id AND ucp.completed_at IS NOT NULL AND ucp.circuit_id = ?)", circuitId)
			case "0":
				baseQuery = baseQuery.Where("NOT EXISTS (SELECT 1 FROM user_circuit_progressions ucp WHERE ucp.end_user_id = eu.id AND ucp.circuit_id = ?)", circuitId)
			default:
				baseQuery = baseQuery.Where("EXISTS (SELECT 1 FROM user_circuit_progressions ucp JOIN user_step_progressions usp ON usp.user_circuit_progression_id = ucp.id WHERE ucp.end_user_id = eu.id AND usp.step_id = ? AND usp.status = 'in_progress')", circuitStep)
			}
		} else {
			baseQuery = baseQuery.Where("EXISTS (SELECT 1 FROM user_circuit_progressions ucp WHERE ucp.end_user_id = eu.id AND ucp.circuit_id = ?)", circuitId)
		}
	}

	if term != "" {
		baseQuery = baseQuery.Where("eu.full_name ILIKE ? OR eu.email ILIKE ?", "%"+term+"%", "%"+term+"%")
	}

	if len(metadataFilters) > 0 {
		i := 0
		for key, value := range metadataFilters {
			alias := fmt.Sprintf("eum%d", i)
			baseQuery = baseQuery.Where(fmt.Sprintf(
				"EXISTS (SELECT 1 FROM end_user_metadatas %s WHERE %s.end_user_id = eu.id AND %s.key_name = ? AND %s.value = ?)",
				alias, alias, alias, alias), key, value)
			i++
		}
	}
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	selectQuery := baseQuery.
		Select(`
			eu.id,
			eu.email,
			eu.full_name,
			eu.created_at,
			eu.last_login_at,
			eu.external_id,
			eu.picture,
			eu.current_streak,
			eu.longest_streak,
			COALESCE(
				json_agg(DISTINCT jsonb_build_object('key', em.key_name, 'value', em.value))
				FILTER (WHERE em.key_name IS NOT NULL), '[]'
			) AS metadata,
			COALESCE(
				json_agg(DISTINCT jsonb_build_object(
					'id', ucp.id,
					'circuitName', c.name,
					'status', CASE WHEN ucp.completed_at IS NULL THEN 'in_progress' ELSE 'completed' END,
					'points', ucp.points,
					'startedAt', to_char(ucp.started_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"'),
					'completedAt', CASE WHEN ucp.completed_at IS NOT NULL THEN to_char(ucp.completed_at AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.MS"Z"') ELSE NULL END
				)) FILTER (WHERE ucp.id IS NOT NULL), '[]'
			) AS circuit_progress
		`).
		Joins("LEFT JOIN end_user_metadatas em ON em.end_user_id = eu.id").
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.end_user_id = eu.id").
		Joins("LEFT JOIN circuits c ON c.id = ucp.circuit_id").
		Group("eu.id").
		Order("eu.created_at DESC").
		Limit(pagination.Limit).
		Offset(pagination.Offset)

	var users []endUserRow
	if err := selectQuery.Scan(&users).Error; err != nil {
		return nil, err
	}

	var result []*EndUserResponse
	for _, row := range users {
		var metaList []MetadataItem
		_ = json.Unmarshal(row.Metadata, &metaList)
		metadata := make(map[string]interface{})
		for _, m := range metaList {
			metadata[m.Key] = m.Value
		}

		var progress []CircuitProgress
		err := json.Unmarshal(row.CircuitProgress, &progress)
		if err != nil {
			fmt.Println("Unmarshal error:", err)
		}

		result = append(result, &EndUserResponse{
			ID:              row.ID,
			Email:           row.Email,
			FullName:        row.FullName,
			CreatedAt:       row.CreatedAt,
			LastLoginAt:     row.LastLoginAt,
			ExternalId:      row.ExternalId,
			Picture:         row.Picture,
			CurrentStreak:   row.CurrentStreak,
			LongestStreak:   row.LongestStreak,
			Metadata:        metadata,
			CircuitProgress: progress,
		})
	}

	return &PaginatedEndUserResponse{
		Total: int(total),
		Users: result,
	}, nil
}
