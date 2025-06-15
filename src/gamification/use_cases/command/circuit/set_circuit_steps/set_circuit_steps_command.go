package set_circuit_steps

type SetCircuitStepsCommand struct {
	CircuitID string
	Steps     []struct {
		Name                string
		Description         *string
		CompletionThreshold int
		EventName           string
		StepNumber          *int
	}
}
