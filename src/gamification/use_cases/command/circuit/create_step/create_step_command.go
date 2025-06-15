package create_step

type CreateStepCommand struct {
	CircuitID           string
	Name                string
	Description         *string
	CompletionThreshold int
	EventName           string
	StepNumber          *int
}
