package middlewares

import "net/http"

func HTMXOnly(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		h(w, r)
	}
}
