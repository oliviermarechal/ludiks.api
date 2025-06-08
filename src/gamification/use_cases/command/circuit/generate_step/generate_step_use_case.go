package generate_step

import (
	"math"
	"strconv"

	"github.com/google/uuid"

	"ludiks/src/gamification/domain/models"
	domain_repositories "ludiks/src/gamification/domain/repositories"
)

type GenerateStepUseCase struct {
	circuitRepository domain_repositories.CircuitRepository
}

func NewGenerateStepUseCase(circuitRepository domain_repositories.CircuitRepository) *GenerateStepUseCase {
	return &GenerateStepUseCase{
		circuitRepository: circuitRepository,
	}
}

func getDefaultExponent(curve CurveType, exponent *float64) float64 {
	if exponent != nil {
		return *exponent
	}

	switch curve {
	case Exponential:
		return 2.0
	case Logarithmic:
		return 2.0
	default:
		return 1.0
	}
}

func calculateCumulativePoints(curve CurveType, totalPoints int, numberOfSteps int, exponent *float64) []int {
	exp := getDefaultExponent(curve, exponent)
	points := make([]int, numberOfSteps)

	switch curve {
	case Linear:
		pointsPerStep := totalPoints / numberOfSteps
		for i := 0; i < numberOfSteps; i++ {
			points[i] = pointsPerStep * (i + 1)
		}
	case Exponential:
		for i := 0; i < numberOfSteps; i++ {
			ratio := math.Pow(float64(i+1)/float64(numberOfSteps), exp)
			points[i] = int(math.Round(float64(totalPoints) * ratio))
		}
	case Logarithmic:
		for i := 0; i < numberOfSteps; i++ {
			ratio := math.Pow(float64(i+1)/float64(numberOfSteps), 1/exp)
			points[i] = int(math.Round(float64(totalPoints) * ratio))
		}
	default:
		pointsPerStep := totalPoints / numberOfSteps
		for i := 0; i < numberOfSteps; i++ {
			points[i] = pointsPerStep * (i + 1)
		}
	}

	return points
}

func (u *GenerateStepUseCase) Execute(command *GenerateStepCommand) (*models.Circuit, error) {
	circuit, err := u.circuitRepository.Find(command.CircuitID)
	if err != nil {
		return nil, err
	}

	cumulativePoints := calculateCumulativePoints(command.Curve, command.MaxValue, command.NumberOfSteps, command.Exponent)

	steps := make([]*models.Step, command.NumberOfSteps)
	for i := 0; i < command.NumberOfSteps; i++ {
		description := "Description of step " + strconv.Itoa(i+1)
		steps[i] = models.CreateStep(
			uuid.New().String(),
			"Step "+strconv.Itoa(i+1),
			&description,
			cumulativePoints[i],
			circuit.ID,
			i+1,
			command.EventName,
		)
	}

	steps, err = u.circuitRepository.CreateMultipleSteps(steps)
	if err != nil {
		return nil, err
	}

	return u.circuitRepository.Find(circuit.ID)
}
