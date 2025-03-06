package collections

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// CtxAlert extracts the Event ID from the URL and adds it to the request context
func CtxAlert(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// in case we are working with id as uuid we do

		//eventID, err := uuid.FromString(chi.URLParam(r, "id"))

		//if err != nil {
		//	logrus.Errorf("Parsing error: %s", err.Error())
		//	customError := &models.CustomError{
		//		Message: fmt.Sprintf("Cannot parse id (%s) as UUID", chi.URLParam(r, "id")),
		//		Code:    http.StatusUnprocessableEntity,
		//	}
		//	w.WriteHeader(customError.Code)
		//	body, _ := json.Marshal(customError)
		//	_, _ = w.Write(body)
		//	return
		//}

		// Store the Event ID in the request context
		//ctx := context.WithValue(r.Context(), "eventID", eventID)
		//next.ServeHTTP(w, r.WithContext(ctx))

		// in my case we are working with uid as string

		eventID := chi.URLParam(r, "id")

		// Store the Event ID in the request context
		ctx := context.WithValue(r.Context(), "eventID", eventID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
