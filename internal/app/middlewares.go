package app

import (
	"net/http"
	"texting-app/internal/pkg/utils"
)

func HTMXOnly(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			http.Error(w, "404 page not found", http.StatusNotFound)
			return
		}
		h(w, r)
	}
}

func EnsureAuth(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("authenticationToken")
		if err != nil {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		tk := authCookie.Value
		jwtMngr := utils.NewJwtManager()
        usrn, err := jwtMngr.VerifyToken(tk)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

        r.Header.Set("username", usrn)

		h(w, r)
	}
}
