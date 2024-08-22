package handlers

import (
    "net/http"

    "texting-app/templates"
)

func Login(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username") 
    password := r.FormValue("password") 

    token, err := 

    templ := templates.Login()
    err := templ.Render(r.Context(), w)
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError) 
    }
}
