package resources

import (
	"database/sql"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	repository "middleware/example/internal/repositories/resources"
	"net/http"
)

// GetResources retrieves all the resources from the database
func GetResources() ([]*models.Resource, error) {
	resources, err := repository.GetAllResources()
	if err != nil {
		logrus.Errorf("Error retrieving resources: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return resources, nil
}

// GetResourceById retrieves a resource by ID, handling errors properly
func GetResourceById(id uuid.UUID) (*models.Resource, error) {
	resource, err := repository.GetResourceById(id)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, &models.CustomError{
				Message: "Resource not found",
				Code:    http.StatusNotFound,
			}
		}
		logrus.Errorf("Error retrieving resource: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Something went wrong",
			Code:    http.StatusInternalServerError,
		}
	}

	return resource, nil
}

// CreateResource creat a new resource in our database
func CreateResource(resource models.Resource) (*models.Resource, error) {
	// Generate UUID for the new resource
	newID, err := uuid.NewV4()
	if err != nil {
		logrus.Errorf("Error generating UUID: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Failed to generate unique ID",
			Code:    http.StatusInternalServerError,
		}
	}
	resource.Id = &newID

	err = repository.CreateResource(resource)
	if err != nil {
		logrus.Errorf("Error creating resource: %s", err.Error())
		return nil, &models.CustomError{
			Message: "Error creating resource",
			Code:    http.StatusInternalServerError,
		}
	}

	return &resource, nil
}

// UpdateResource updates an existing resource
func UpdateResource(resource *models.Resource) error {
	err := repository.UpdateResource(resource)
	if err != nil {
		logrus.Errorf("Error updating resource: %s", err.Error())
		return &models.CustomError{
			Message: "Error updating the resource",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}

// DeleteResource deletes a resource by its ID
func DeleteResource(id uuid.UUID) error {
	err := repository.DeleteResource(id)
	if err != nil {
		logrus.Errorf("Error deleting resource: %s", err.Error())
		return &models.CustomError{
			Message: "Error deleting the resource",
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}
