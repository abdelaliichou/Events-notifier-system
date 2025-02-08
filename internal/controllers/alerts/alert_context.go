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

// CtxAlert extracts the Alert ID from the URL and adds it to the request context
func CtxAlert(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alertID, err := uuid.FromString(chi.URLParam(r, "id"))
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

		// Store the Alert ID in the request context
		ctx := context.WithValue(r.Context(), "alertID", alertID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
