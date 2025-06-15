package log_end_user

type LogEndUserCommand struct {
	ProjectID string
	ID        string
	FullName  string
	Email     *string
	Picture   *string
	Metadata  map[string]interface{}
}
