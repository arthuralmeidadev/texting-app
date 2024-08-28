package handlers

import (
	"net/http"
	"time"

	"texting-app/internal/pkg/providers"
	"texting-app/internal/pkg/store"
	"texting-app/templates"
)

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		comp := templates.Login()
		err := comp.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	case http.MethodPost:
		usrn := r.FormValue("username")
		pw := r.FormValue("password")
		store, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		usrStore := providers.NewUserProvider(store)
		usr, err := usrStore.AuthUser(usrn, pw)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		tk, err := usr.NewToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authCookie := &http.Cookie{
			Name:    "authenticationToken",
			Value:   string(tk),
			Expires: time.Now().Add(time.Hour * 2),
		}
		http.SetCookie(w, authCookie)
		w.Header().Set("HX-Redirect", "/chats")
		w.WriteHeader(http.StatusFound)
		return
	default:
		return
	}
}
