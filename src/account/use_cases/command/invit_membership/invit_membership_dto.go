package invit_membership

import "github.com/go-playground/validator/v10"

type InvitMembershipDTO struct {
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

func (g *InvitMembershipDTO) Validate() error {
	validator := validator.New()
	return validator.Struct(g)
}
