package invit_membership

type InvitMembershipCommand struct {
	OrganizationID string
	FromUserID     string
	Email          string
	Role           string
}
