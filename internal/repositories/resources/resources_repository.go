package resources

import (
	"github.com/gofrs/uuid"
	"middleware/example/internal/helpers"
	"middleware/example/internal/models"
)

func GetAllResources() ([]*models.Resource, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db)

	query := models.GET_ALL_RESOURCES
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []*models.Resource
	for rows.Next() {
		var resource models.Resource
		if err := rows.Scan(&resource.Id, &resource.UcaID, &resource.Name); err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}

	return resources, nil
}

// GetResourceById fetches a resource by ID from the database
func GetResourceById(id uuid.UUID) (*models.Resource, error) {
	db, err := helpers.OpenDB()
	if err != nil {
		return nil, err
	}
	defer helpers.CloseDB(db) // Close DB connection after query

	var resource models.Resource
	query := models.GET_RESOURCE
	row := db.QueryRow(query, id.String())

	err = row.Scan(&resource.Id, &resource.UcaID, &resource.Name)
	if err != nil {
		return nil, err // Let the service layer handle errors
	}

	return &resource, nil
}

// CreateResource inserts a new resource into the database
func CreateResource(resource models.Resource) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := models.CREAT_RESOURCE
	_, err = db.Exec(query, resource.Id.String(), resource.UcaID, resource.Name)

	return err
}

// UpdateResource updates an existing resource in the database
func UpdateResource(resource *models.Resource) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	query := models.UPDATE_RESOURCE
	_, err = db.Exec(query, resource.UcaID, resource.Name, resource.Id.String())

	return err
}

// DeleteResource deletes an existing resource in the database
func DeleteResource(id uuid.UUID) error {
	db, err := helpers.OpenDB()
	if err != nil {
		return err
	}
	defer helpers.CloseDB(db)

	// Delete the resource based on the ID
	query := models.DELETE_RESOURCE
	result, err := db.Exec(query, id.String())
	if err != nil {
		return err
	}

	// Check if any row was affected (if no rows were deleted, the resource wasn't found)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return &models.CustomError{
			Message: "Resource not found",
			Code:    401,
		}
	}

	return nil
}
