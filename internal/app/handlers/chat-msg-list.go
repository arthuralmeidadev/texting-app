package handlers

import (
	"net/http"

	"texting-app/partials"
)

func ChatMsgList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
        chatMsgList := make([]partials.ChatMsg, 0)
		component := partials.ChatMsgList(chatMsgList) 
		err := component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	default:
		return
	}
}
