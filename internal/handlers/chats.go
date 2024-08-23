package handlers

import (
	"net/http"
	"texting-app/internal/pkg/utils"
	"texting-app/templates"
)

func Chats(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		authCookie, err := r.Cookie("authenticationToken")
		if err != nil {
			http.Error(w, "Missing authentication token", http.StatusUnauthorized)
			return
		}

		token := authCookie.Value
		jwtMngr := utils.NewJwtManager()
		if err := jwtMngr.VerifyToken(token); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		component := templates.Chats()
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		return
	}
}
