package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/resources"
	"net/http"
)

// UpdateResource
// @Tags         resources
// @Summary      Update an existing resource by ID
// @Description  This endpoint updates an existing resource with the provided details.
// @Param        id path string true "Resource UUID"
// @Param        resource body models.Resource true "Resource Data"
// @Success      200 "Resource updated successfully"
// @Failure      400 "Invalid request body"
// @Failure      404 "Resource not found"
// @Failure      500 "Error updating resource"
// @Router       /resources/{id} [put]
func UpdateResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceID, _ := ctx.Value("resourceID").(uuid.UUID)

	var resource models.Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		logrus.Errorf("Invalid JSON: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resource.Id = &resourceID // Ensure ID is set

	err = services.UpdateResource(&resource)
	if err != nil {
		logrus.Errorf("Error updating resource: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
