package common_validator

import "net/http"

type CommonValidator struct{}

// mock middleware validator
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
