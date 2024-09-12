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
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		authTk := authCookie.Value
		jwtMngr := utils.NewJwtManager()
		usrn, err := jwtMngr.VerifyToken(authTk)
		if err != nil {
			refCookie, err := r.Cookie("refreshToken")
			if err != nil {
				http.Error(w, "Missing refresh token", http.StatusUnauthorized)
				return
			}

			refTk := refCookie.Value
			usrn, err = jwtMngr.VerifyToken(refTk)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}

		r.Header.Set("username", usrn)

		h(w, r)
	}
}
