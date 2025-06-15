package log_end_user

import (
	"fmt"
	"ludiks/src/tracking/domain/models"
	domain_repositories "ludiks/src/tracking/domain/repositories"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func convertToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type LogEndUserUseCase struct {
	endUserRepository  domain_repositories.EndUserRepository
	metadataRepository domain_repositories.MetadataRepository
}

func NewLogEndUserUseCase(
	endUserRepository domain_repositories.EndUserRepository,
	metadataRepository domain_repositories.MetadataRepository,
) *LogEndUserUseCase {
	return &LogEndUserUseCase{
		endUserRepository:  endUserRepository,
		metadataRepository: metadataRepository,
	}
}

func (u *LogEndUserUseCase) Execute(command LogEndUserCommand) (*models.EndUser, error) {
	endUser, err := u.endUserRepository.FindByExternalID(command.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		endUser = models.CreateEndUser(uuid.New().String(), command.ID, command.FullName, command.Email, command.Picture, command.ProjectID)
		endUser, err = u.endUserRepository.Create(endUser)
		if err != nil {
			return nil, err
		}
	} else {
		now := time.Now().UTC()
		lastDate := endUser.LastLoginAt.UTC().Truncate(24 * time.Hour)
		today := now.Truncate(24 * time.Hour)
		yesterday := today.AddDate(0, 0, -1)

		switch {
		case lastDate.Equal(yesterday):
			endUser.CurrentStreak++
			if endUser.CurrentStreak > endUser.LongestStreak {
				endUser.LongestStreak = endUser.CurrentStreak
			}
		case lastDate.Before(yesterday):
			endUser.CurrentStreak = 1
		}

		endUser.LastLoginAt = now
		endUser, err = u.endUserRepository.Update(endUser)
		if err != nil {
			return nil, err
		}
	}

	if command.Metadata != nil {
		projectMetadataKeys, err := u.metadataRepository.ListProjectMetadataKey(command.ProjectID)
		if err != nil {
			return nil, err
		}

		var newEndUserMetadata []*models.EndUserMetadata
		var newProjectMetadataKeys []*models.ProjectMetadataKey
		var newMetadataValues []*models.MetadataValue

		existingMetadata := make(map[string]bool)

		for key, value := range command.Metadata {
			existingMetadata[key] = true
			stringValue := convertToString(value)

			if endUser.HasMetadataWithValue(key, stringValue) {
				continue
			}

			var projectMetadataKey *models.ProjectMetadataKey
			projectMetadataKeyFound := false
			for _, pmk := range *projectMetadataKeys {
				if pmk.HasMetadata(key) {
					projectMetadataKey = &pmk
					projectMetadataKeyFound = true
					break
				}
			}

			if projectMetadataKeyFound {
				if !projectMetadataKey.HasValue(stringValue) {
					metadataValue := models.CreateMetadataValue(uuid.New().String(), projectMetadataKey.ID, stringValue)
					newMetadataValues = append(newMetadataValues, metadataValue)
				}
			} else {
				// TODO Premium: Check project metadata key quota
				projectMetadataKey = models.CreateProjectMetadataKey(uuid.New().String(), command.ProjectID, key)
				newProjectMetadataKeys = append(newProjectMetadataKeys, projectMetadataKey)

				metadataValue := models.CreateMetadataValue(uuid.New().String(), projectMetadataKey.ID, stringValue)
				newMetadataValues = append(newMetadataValues, metadataValue)
			}

			existingMetadata := endUser.GetMetadata(key)
			if existingMetadata != nil {
				if existingMetadata.Value != stringValue {
					existingMetadata.Value = stringValue
					if err := u.metadataRepository.UpdateEndUserMetadata(existingMetadata); err != nil {
						return nil, err
					}
				}
			} else {
				metadata := models.CreateEndUserMetadata(endUser.ID, key, stringValue)
				newEndUserMetadata = append(newEndUserMetadata, metadata)
			}
		}

		for _, metadata := range endUser.EndUserMetadata {
			if !existingMetadata[metadata.KeyName] {
				if err := u.metadataRepository.DeleteEndUserMetadata(endUser.ID, metadata.KeyName); err != nil {
					return nil, err
				}
			}
		}

		if len(newProjectMetadataKeys) > 0 {
			if err := u.metadataRepository.BatchCreateProjectMetadataKeys(newProjectMetadataKeys); err != nil {
				return nil, err
			}
		}

		if len(newMetadataValues) > 0 {
			if err := u.metadataRepository.BatchCreateMetadataValues(newMetadataValues); err != nil {
				return nil, err
			}
		}

		if len(newEndUserMetadata) > 0 {
			if err := u.metadataRepository.BatchCreateEndUserMetadata(newEndUserMetadata); err != nil {
				return nil, err
			}
		}

		endUser, err = u.endUserRepository.Find(endUser.ID)
		if err != nil {
			return nil, err
		}
	}

	return endUser, nil
}
