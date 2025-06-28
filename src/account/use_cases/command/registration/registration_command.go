package registration

type RegistrationCommand struct {
	Email    string
	Password string
	InviteID *string
}
