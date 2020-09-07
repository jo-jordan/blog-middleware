package http

import (
	"net/http"
)

func Filter(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		inner.ServeHTTP(w, r)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	})
}
