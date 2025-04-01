package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"middleware/example/internal/services/resources"
	"net/http"
)

// GetResource
// @Tags         resources
// @Summary      Get a resource by ID
// @Description  Retrieve a resource using its unique ID
// @Param        id path string true "Resource UUID"
// @Success      200 {object} models.Resource
// @Failure      404 "Resource not found"
// @Failure      500 "Internal server error"
// @Router       /resources/{id} [get]
func GetResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceID, _ := ctx.Value("resourceID").(uuid.UUID)

	// Fetch the resource from the database
	resource, err := resources.GetResourceById(resourceID)
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		customError, isCustom := err.(*models.CustomError)
		if isCustom {
			w.WriteHeader(customError.Code)
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Respond with the resource data
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(resource)
	_, _ = w.Write(body)
}
