package collections

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"middleware/example/internal/models"
	"net/http"
)

// CtxResource extracts the Resource ID from the URL and adds it to the request context
func CtxResource(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resourceID, err := uuid.FromString(chi.URLParam(r, "id"))
		if err != nil {
			logrus.Errorf("Parsing error: %s", err.Error())
			customError := &models.CustomError{
				Message: fmt.Sprintf("Cannot parse id (%s) as UUID", chi.URLParam(r, "id")),
				Code:    http.StatusUnprocessableEntity,
			}
			w.WriteHeader(customError.Code)
			body, _ := json.Marshal(customError)
			_, _ = w.Write(body)
			return
		}

		// Store the resourceID in the request context
		ctx := context.WithValue(r.Context(), "resourceID", resourceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
