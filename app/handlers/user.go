package handlrs

import (
	"encoding/json"
	"github.com/go-bumbu/userauth/handlers/sessionauth"
	"net/http"
)

type userStatus struct {
	User     string `json:"username"`
	LoggedIn bool   `json:"logged-in"`
}

func UserStatusHandler(session *sessionauth.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		data, err := session.GetSessData(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData := userStatus{
			User:     data.UserId,
			LoggedIn: data.IsAuthenticated,
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}
