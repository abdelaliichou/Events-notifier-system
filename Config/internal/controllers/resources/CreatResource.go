package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	services "middleware/example/internal/services/resources"
	"net/http"
)

// CreateResource
// @Tags         resources
// @Summary      Create a new resource
// @Description  This endpoint creates a new resource based on the request body.
// @Accept       json
// @Produce      json
// @Param        resource body models.Resource true "Resource Data"
// @Success      201 {object} models.Resource
// @Failure      400 "Invalid request body"
// @Failure      500 "Error creating resource"
// @Router       /resources [post]
func CreateResource(w http.ResponseWriter, r *http.Request) {

	var resource models.Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		logrus.Errorf("Invalid JSON: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newResource, err := services.CreateResource(resource)
	if err != nil {
		logrus.Errorf("Error creating resource: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	body, _ := json.Marshal(newResource)
	_, _ = w.Write(body)
}
