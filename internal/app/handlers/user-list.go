package handlers

import (
	"net/http"
	"strconv"
	"texting-app/internal/pkg/providers"
	"texting-app/internal/pkg/store"
	"texting-app/partials"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		storeInst, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrn := r.URL.Query().Get("username")
		if len(usrn) == 0 {
			comp := partials.UserList(make([]string, 0))
			err = comp.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrProvider := providers.NewUserProvider(storeInst)

		forUsr := r.Header.Get("username")
		usrns, err := usrProvider.FindUser(forUsr, usrn, uint(offset))
		if err != nil {
			if err.Error() == "no rows" {
				http.Error(w, "404 Not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		comp := partials.UserList(usrns)
		err = comp.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		return
	default:
		return
	}
}
