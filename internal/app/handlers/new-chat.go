package handlers

import "net/http"

func NewChat(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        return
    case http.MethodPost:
        return
    default:
        return
    }
}
