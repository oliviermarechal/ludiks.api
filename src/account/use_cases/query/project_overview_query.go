package query

import (
	"gorm.io/gorm"
)

type StepCompletionRate struct {
	StepId         string `json:"stepId"`
	StepName       string `json:"stepName"`
	CompletionRate int    `json:"completionRate"`
	AverageTime    int    `json:"averageTime"`
}

type CircuitInsights struct {
	CompletionRate        int `json:"completionRate"`
	AverageCompletionTime int `json:"averageCompletionTime"`
}

type CircuitWithInsights struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	CompletionRate        int    `json:"completionRate"`
	AverageCompletionTime int    `json:"averageCompletionTime"`
}

type ProjectKPIs struct {
	Total                 int `json:"total"`
	Active                int `json:"active"`
	Inactive              int `json:"inactive"`
	AverageCompletionRate int `json:"averageCompletionRate"`
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

	// Fetch basic circuit data with insights
	db.Table("circuits c").
		Select(`
		c.id,
		c.name,
		c.type,
		AVG(CASE WHEN ucp.completed_at IS NOT NULL THEN 100 END) as completion_rate,
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
		AVG(CASE WHEN ucp.completed_at IS NOT NULL THEN 100 END) as averageCompletionRate
	`).
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.circuit_id = c.id").
		Where("c.project_id = ?", projectID).
		Scan(&overview.KPIs)

	db.Table("end_users u").
		Select(`
		COUNT(DISTINCT u.id) as total,
		COUNT(DISTINCT CASE WHEN u.last_login_at >= NOW() - INTERVAL '1 week' THEN u.id END) as activeLastWeek,
		COUNT(DISTINCT CASE WHEN ucp.completed_at IS NOT NULL THEN u.id END) as completedAtLeastOne
	`).
		Joins("LEFT JOIN user_circuit_progressions ucp ON ucp.end_user_id = u.id").
		Where("u.project_id = ?", projectID).
		Scan(&overview.Users)

	return overview, nil
}
