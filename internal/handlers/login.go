package handlers

import (
	"net/http"
	"time"

	"texting-app/internal/pkg/auth"
	"texting-app/templates"
)

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templ := templates.Login()
		err := templ.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        return
    case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := auth.AuthenticateUser(username, password)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		token := user.GetToken()
		cookie := &http.Cookie{
			Name:    "authenticationToken",
			Value:   token,
			Expires: time.Now().Add(time.Hour * 2),
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/home", http.StatusFound)
        return
	default:
		return
	}
}
