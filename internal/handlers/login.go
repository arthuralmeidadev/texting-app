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
		component := templates.Login()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		store, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		userStore := providers.NewUserProvider(store)
		user, err := userStore.AuthUser(username, password)
		if err != nil {
			if err.Error() == "unauthorized" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		token, err := user.NewToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		authCookie := &http.Cookie{
			Name:    "authenticationToken",
			Value:   string(token),
			Expires: time.Now().Add(time.Hour * 2),
		}
		http.SetCookie(w, authCookie)
		http.Redirect(w, r, "/chats", http.StatusFound)
		return
	default:
		return
	}
}
