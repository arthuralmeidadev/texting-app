package handlers

import (
	"net/http"
	"texting-app/templates"
)

func Chats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		component := templates.Chats()
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		return
	}
}
