package handlers

import (
	"net/http"
	"texting-app/internal/pkg/providers"
	"texting-app/internal/pkg/store"
	"texting-app/partials"
)

func FriendList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		storeInst, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrProvider := providers.NewUserProvider(storeInst)
		usrn := r.Header.Get("usrname")
		usrnFriends, err := usrProvider.GetUserFriends(usrn, 0)
		if err != nil {
			if err.Error() == "no rows" {
				http.Error(w, "404 not found", http.StatusNotFound)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

        if len(usrnFriends) == 0 {
            http.Error(w, "404 not found", http.StatusNotFound)
			return
        }

		friends := make([]string, len(usrnFriends))
		for i := 0; i < len(usrnFriends); i++ {
			friends = append(friends, usrnFriends[i].Username)
		}

		comp := partials.FriendList(friends)
		err = comp.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
    case http.MethodPost:
		storeInst, err := store.GetStore()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		usrProvider := providers.NewUserProvider(storeInst)
		usrn := r.Header.Get("username")
        recUsrn := r.FormValue("recUsrn")
        err = usrProvider.SendFriendRequest(usrn, recUsrn)
        if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
        }

        return
	default:
		return
	}
}
