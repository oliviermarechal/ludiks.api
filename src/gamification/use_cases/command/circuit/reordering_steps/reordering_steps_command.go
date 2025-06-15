package reordering_steps

type ReorderingStepsCommand struct {
	CircuitID string
	Steps     []struct {
		StepId     string
		StepNumber int
	}
}
