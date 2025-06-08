package create_reward

type CreateRewardCommand struct {
	CircuitID                 string
	Name                      string
	Description               *string
	StepID                    *string
	UnlockOnCircuitCompletion bool
}
