package update_step

type UpdateStepCommand struct {
	StepId              string
	Name                string
	Description         string
	CompletionThreshold int
	EventName           string
}
