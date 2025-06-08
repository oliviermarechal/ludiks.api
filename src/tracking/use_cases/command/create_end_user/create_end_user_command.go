package create_end_user

type CreateEndUserCommand struct {
	ProjectID string
	ID        string
	FullName  string
	Email     *string
	Picture   *string
}
