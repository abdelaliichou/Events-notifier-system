package collections

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	services "middleware/example/internal/services/resources"
	"net/http"
)

// GetAllResources
// @Tags         resources
// @Summary      Get all resources
// @Description  Retrieve all resources from the service layer
// @Success      200 {array} models.Resource
// @Failure      500 "Error fetching resources"
// @Router       /resources [get]
func GetAllResources(w http.ResponseWriter, r *http.Request) {
	// Fetch all resources from the service layer
	resources, err := services.GetResources()
	if err != nil {
		logrus.Errorf("error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the list of resources
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(resources)
	_, _ = w.Write(body)
}
