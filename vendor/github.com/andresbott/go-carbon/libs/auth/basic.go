package auth

import (
	"fmt"
	"net/http"
)

const (
	ActionLoginCheck = iota
	ActionLoginOk
	ActionLoginFailed
)

type Basic struct {
	User         UserLogin
	Message      string
	Redirect     string
	RedirectCode int
	logger       func(action int, user string)
}

func (auth *Basic) Middleware(next http.Handler) http.Handler {
	if auth.logger == nil {
		auth.logger = func(action int, user string) {}
	}
	if auth.Message == "" {
		auth.Message = "Authenticate"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		auth.logger(ActionLoginCheck, username)

		if ok {
			if auth.User.AllowLogin(username, password) {
				auth.logger(ActionLoginOk, username)
				next.ServeHTTP(w, r)
				return
			} else {
				auth.logger(ActionLoginFailed, username)
			}
		}

		if auth.Redirect != "" {
			http.Redirect(w, r, auth.Redirect, auth.RedirectCode)
			return
		}

		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, auth.Message))
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

type UserLogin interface {
	AllowLogin(user string, password string) bool
}
