package handlers

import (
	"net/http"
	"texting-app/internal/pkg/providers"
	"texting-app/internal/pkg/store"
	"texting-app/templates"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		comp := templates.Signup()
		err := comp.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case http.MethodPost:
		usrn := r.FormValue("username")
		pw := r.FormValue("password")
		repeatPw := r.FormValue("repatedPassword")

		if pw != repeatPw {
			http.Error(w, "Passwords do no match", http.StatusBadRequest)
		}

		store, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		usrStore := providers.NewUserProvider(store)
		usr, err := usrStore.GetUser(usrn)
		if err != nil && err.Error() != "ErrNoRows" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if usr != nil {
			http.Error(w, "User already registered", http.StatusConflict)
			return
		}

		usr, err = usrStore.CreateUser(usrn, pw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
