package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	services "middleware/example/internal/services/resources"
	"net/http"
)

// DeleteResource
// @Tags         resources
// @Summary      Delete a resource by ID
// @Description  This endpoint deletes a resource using its unique ID
// @Param        id path string true "Resource UUID"
// @Success      200 {object} string
// @Failure      500 "Error deleting resource"
// @Router       /resources/{id} [delete]
func DeleteResource(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceID, _ := ctx.Value("resourceID").(uuid.UUID)

	// Call the service layer to delete the resource
	err := services.DeleteResource(resourceID)
	if err != nil {
		logrus.Errorf("Error deleting resource: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return success response (200 OK)
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal("Resource deleted successfully")
	_, _ = w.Write(body)
}
