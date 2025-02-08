package collections

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/resources"
	"net/http"
)

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
