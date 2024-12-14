package handlrs

import (
	"encoding/json"
	"github.com/go-bumbu/userauth"
	"github.com/go-bumbu/userauth/handlers/sessionauth"
	"net/http"
)

//type UserHandler struct {
//	user    auth.UserLogin
//	session *auth.SessionMgr
//}
//
//
//
//func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
//	var payload loginData
//
//	err := json.NewDecoder(r.Body).Decode(&payload)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if h.user.AllowLogin(payload.User, payload.Pw) {
//		err = h.session.Login(r, w, payload.User)
//		if err != nil {
//			http.Error(w, "internal error", http.StatusInternalServerError)
//			return
//		}
//	} else {
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
//
//}

type loginData struct {
	User string `json:"username"`
	Pw   string `json:"password"`
}

type userStatus struct {
	User     string `json:"username"`
	LoggedIn bool   `json:"logged-in"`
}

func UserStatusHandler(session *sessionauth.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := session.Read(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData := userStatus{
			User:     data.UserId,
			LoggedIn: data.IsAuthenticated,
		}
		err = json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

const anonymousUser = "anonymous"

type autDisabledStatus struct {
	AuthDisabled bool   `json:"auth-disabled,omitempty"`
	User         string `json:"user"`
}

func AuthDisabledHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData := autDisabledStatus{
			AuthDisabled: true,
			User:         anonymousUser,
		}
		err := json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func UserLogoutHandler(session *sessionauth.Manager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := session.LogoutUser(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		jsonData := userStatus{
			User:     "",
			LoggedIn: false,
		}

		err = json.NewEncoder(w).Encode(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})
}

func UserLoginHandler(session *sessionauth.Manager, user userauth.LoginHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload loginData
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		canlogin, err := user.CanLogin(payload.User, payload.Pw)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if canlogin {
			err = session.LoginUser(r, w, payload.User)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			// todo read user data...
			jsonData := userStatus{
				User:     payload.User,
				LoggedIn: true,
			}
			err = json.NewEncoder(w).Encode(jsonData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
