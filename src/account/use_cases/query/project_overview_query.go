package query

import (
	"gorm.io/gorm"
)

type StepCompletionRate struct {
	StepId         string  `json:"stepId"`
	StepName       string  `json:"stepName"`
	CompletionRate float64 `json:"completionRate"`
	AverageTime    int     `json:"averageTime"`
}

type CircuitInsights struct {
	CompletionRate        float64 `json:"completionRate"`
	AverageCompletionTime int     `json:"averageCompletionTime"`
}

type CircuitWithInsights struct {
	Id                    string  `json:"id"`
	Name                  string  `json:"name"`
	Type                  string  `json:"type"`
	Active                bool    `json:"active"`
	CompletionRate        float64 `json:"completionRate"`
	AverageCompletionTime int     `json:"averageCompletionTime"`
}

type ProjectKPIs struct {
	Total                 int     `json:"total"`
	Active                int     `json:"active"`
	Inactive              int     `json:"inactive"`
	AverageCompletionRate float64 `json:"averageCompletionRate"`
}

type ProjectUsers struct {
	Total               int `json:"total"`
	ActiveLastWeek      int `json:"activeLastWeek"`
	CompletedAtLeastOne int `json:"completedAtLeastOne"`
}

type ProjectOverview struct {
	KPIs     ProjectKPIs           `json:"KPIs"`
	Users    ProjectUsers          `json:"users"`
	Circuits []CircuitWithInsights `json:"circuits"`
}

func ProjectOverviewQuery(db *gorm.DB, projectID string) (ProjectOverview, error) {
	var overview ProjectOverview
	overview.Circuits = make([]CircuitWithInsights, 0)

	db.Table("circuits c").
		Select(`
		c.id,
		c.name,
		c.type,
		c.active,
		AVG(
			CASE 
				WHEN ucp.id IS NOT NULL THEN 
					(CAST((
						SELECT COUNT(*) 
						FROM user_step_progressions usp 
						WHERE usp.user_circuit_progression_id = ucp.id 
						AND usp.completed_at IS NOT NULL
					) AS FLOAT) / 
					CAST((
						SELECT COUNT(*) 
						FROM steps s 
						WHERE s.circuit_id = c.id
					) AS FLOAT)) * 100
				ELSE 0 
			END
		) as completion_rate,
		EXTRACT(EPOCH FROM AVG(CASE WHEN ucp.completed_at IS NOT NULL THEN ucp.completed_at - ucp.started_at END)) as average_completion_time
		`).
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.circuit_id = c.id").
		Group("c.id").
		Where("c.project_id = ?", projectID).
		Order("c.created_at DESC").
		Limit(10).
		Scan(&overview.Circuits)

	db.Table("circuits c").
		Select(`
		COUNT(DISTINCT c.id) as total,
		COUNT(DISTINCT CASE WHEN c.active = true THEN c.id END) as active,
		COUNT(DISTINCT CASE WHEN c.active = false THEN c.id END) as inactive,
		AVG(
			CASE 
				WHEN ucp.id IS NOT NULL THEN 
					(CAST((
						SELECT COUNT(*) 
						FROM user_step_progressions usp 
						WHERE usp.user_circuit_progression_id = ucp.id 
						AND usp.completed_at IS NOT NULL
					) AS FLOAT) / 
					CAST((
						SELECT COUNT(*) 
						FROM steps s 
						WHERE s.circuit_id = c.id
					) AS FLOAT)) * 100
				ELSE 0 
			END
		) as averageCompletionRate
	`).
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.circuit_id = c.id").
		Where("c.project_id = ?", projectID).
		Scan(&overview.KPIs)

	var total int64
	db.Table("end_users u").
		Select("COUNT(DISTINCT u.id) as total").
		Where("u.project_id = ?", projectID).
		Scan(&total)
	overview.Users.Total = int(total)

	var activeLastWeek int64
	db.Table("end_users u").
		Select("COUNT(DISTINCT u.id) as activeLastWeek").
		Where("u.project_id = ? AND u.last_login_at >= NOW() - INTERVAL '1 week'", projectID).
		Scan(&activeLastWeek)
	overview.Users.ActiveLastWeek = int(activeLastWeek)

	var completedAtLeastOne int64
	db.Table("end_users u").
		Select("COUNT(DISTINCT u.id) as completedAtLeastOne").
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.end_user_id = u.id").
		Where("u.project_id = ? AND ucp.completed_at IS NOT NULL", projectID).
		Scan(&completedAtLeastOne)
	overview.Users.CompletedAtLeastOne = int(completedAtLeastOne)

	return overview, nil
}
