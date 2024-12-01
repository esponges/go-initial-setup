package common_validator

import "net/http"

// This is necessary because Go requires that middleware functions be methods on a struct in order to be used as handlers.
// If we just invoked common_validator.Validate directly,
// it would not work because Validate is not a standalone function that can be used as a handler.
type CommonValidator struct{}

func (c *CommonValidator) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the POST request has an unexpected header, return an error
		if r.Header.Get("Content-Type") != "application/json" && r.Method == "POST" {
			http.Error(w, "invalid content type", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
