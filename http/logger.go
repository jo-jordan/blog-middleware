package http

import (
	"log"
	"net/http"
	"time"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT")
		w.Header().Add("Access-Control-Allow-Headers", "content-type,x-token")

		if r.Method == "OPTIONS" {
			w.Header().Add("Access-Control-Max-Age", "2000")
		} else {
			w.Header().Add("Access-Control-Max-Age", "86400")
			//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			inner.ServeHTTP(w, r)

			log.Printf(
				"%s\t%s\t%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				name,
				time.Since(start),
				r.Header.Get("X-TOKEN"),
			)
		}
	})
}
