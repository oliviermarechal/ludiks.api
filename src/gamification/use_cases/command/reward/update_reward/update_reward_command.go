package update_reward

type UpdateRewardCommand struct {
	RewardID                  string
	Name                      string
	Description               *string
	StepID                    *string
	UnlockOnCircuitCompletion bool
}
