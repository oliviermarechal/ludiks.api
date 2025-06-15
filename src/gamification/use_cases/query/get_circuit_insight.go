package query

import (
	"gorm.io/gorm"
)

type CircuitInsight struct {
	ID                string               `json:"id"`
	ActiveUsers       *int                 `json:"activeUsers"`
	Name              string               `json:"name"`
	Type              string               `json:"type"`
	CompletionRate    *float64             `json:"completionRate"`
	AvgCompletionTime *int                 `json:"avgCompletionTime"`
	Steps             []CircuitStepInsight `json:"steps" gorm:"-"`
}

type CircuitStepInsight struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	CompletionThreshold int      `json:"completionThreshold"`
	UsersCompleted      int      `json:"usersCompleted"`
	UsersOnThisStep     int      `json:"usersOnThisStep"`
	CompletionRate      *float64 `json:"completionRate"`
	AvgTime             *int     `json:"avgTime"`
}

func GetCircuitInsight(db *gorm.DB, circuitID string) (*CircuitInsight, error) {
	var circuitInsight CircuitInsight

	err := db.Table("circuits c").
		Select(`
			c.id,
			c.name,
			c.type,
			COUNT(DISTINCT ucp.end_user_id) AS "ActiveUsers",
			CASE 
				WHEN COUNT(DISTINCT ucp.end_user_id) = 0 THEN 0
				ELSE ROUND(
					100.0 * COUNT(DISTINCT CASE WHEN ucp.completed_at IS NOT NULL THEN ucp.end_user_id END) / COUNT(DISTINCT ucp.end_user_id),
					2
				)
			END AS "CompletionRate",
			ROUND(AVG(EXTRACT(EPOCH FROM (ucp.completed_at - ucp.started_at))) FILTER (WHERE ucp.completed_at IS NOT NULL)) AS "AvgCompletionTime"
		`).
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.circuit_id = c.id").
		Where("c.id = ?", circuitID).
		Group("c.id, c.name, c.type").
		Scan(&circuitInsight).Error

	if err != nil {
		return &CircuitInsight{}, err
	}

	orderBy := "s.completion_threshold ASC"
	if circuitInsight.Type == "objective" {
		orderBy = "s.step_number ASC"
	}

	err = db.Table("steps s").
		Select(`
			s.id,
			s.name,
			s.completion_threshold,
			COUNT(DISTINCT CASE WHEN usp.completed_at IS NOT NULL THEN ucp.end_user_id END) AS "UsersCompleted",
			COUNT(DISTINCT CASE WHEN usp.completed_at IS NULL THEN ucp.end_user_id END) AS "UsersOnThisStep",
			CASE 
				WHEN total_users.total IS NULL OR total_users.total = 0 THEN 0
				ELSE ROUND(
					100.0 * COUNT(DISTINCT CASE WHEN usp.completed_at IS NOT NULL THEN ucp.end_user_id END) / total_users.total,
					2
				)
			END AS "CompletionRate",
			ROUND(AVG(EXTRACT(EPOCH FROM (usp.completed_at - usp.started_at))) FILTER (WHERE usp.completed_at IS NOT NULL)) AS "AvgTime"
		`).
		Joins("LEFT JOIN user_step_progressions usp ON usp.step_id = s.id").
		Joins("LEFT JOIN user_circuit_progressions ucp ON usp.user_circuit_progression_id = ucp.id").
		Joins(`LEFT JOIN (
			SELECT circuit_id, COUNT(DISTINCT end_user_id) AS total
			FROM user_circuit_progressions
			WHERE circuit_id = ?
			GROUP BY circuit_id
		) AS total_users ON total_users.circuit_id = s.circuit_id`, circuitID).
		Where("s.circuit_id = ?", circuitID).
		Group("s.id, s.name, s.completion_threshold, s.step_number, total_users.total").
		Order(orderBy).
		Scan(&circuitInsight.Steps).Error

	return &circuitInsight, err
}
