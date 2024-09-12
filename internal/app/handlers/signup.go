package handlers

import (
	"log"
	"net/http"
	"regexp"
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
		repeatPw := r.FormValue("repeatPassword")

		if pw != repeatPw {
			log.Println("password do not match")
			http.Error(w, "Passwords do no match", http.StatusBadRequest)
			return
		}

		usrnRe := regexp.MustCompile(`^(([A-Za-z]+\w*)?){4,32}$`)
		if !usrnRe.MatchString(usrn) {
			http.Error(w, "Invalid username format", http.StatusBadRequest)
			return
		}

		pwReChar := regexp.MustCompile(`^[A-Za-z\d\-\_]{8,16}$`)
		pwReUpper := regexp.MustCompile(`[A-Z]`)
		pwReLower := regexp.MustCompile(`[a-z]`)
		pwReDigit := regexp.MustCompile(`\d`)
		if !(pwReChar.MatchString(pw) &&
			pwReUpper.MatchString(pw) &&
			pwReLower.MatchString(pw) &&
			pwReDigit.MatchString(pw)) {
			http.Error(w, "Invalid password format", http.StatusBadRequest)
			return
		}

		storeInst, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrStore := providers.NewUserProvider(storeInst)
		usr, err := usrStore.GetUser(usrn)
		if err != nil && err.Error() != "no rows" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if usr != nil {
			http.Error(w, "User already registered", http.StatusConflict)
			return
		}

		usr, err = usrStore.CreateUser(usrn, pw)
		if err != nil {
			log.Println("ERR", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		authTk, err := usr.NewAuthToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		refTk, err := usr.NewAuthToken()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "authenticationToken",
			Value:   string(authTk),
			Expires: time.Now().Add(time.Hour * 24),
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "refreshToken",
			Value:   string(refTk),
			Expires: time.Now().Add(time.Hour * 24),
		})
		w.Header().Set("HX-Redirect", "/chats")
		w.WriteHeader(http.StatusFound)
		return
	default:
		return
	}
}
