package auth

import (
	"fmt"
	"net/http"
)

// CtxCheckAuth extracts and verifies the user information from a request context
// the returned struct contains user information about the logged-in user
func CtxCheckAuth(r *http.Request) (UserData, error) {

	var d UserData
	ctx := r.Context()

	val := ctx.Value(UserIsLoggedInKey)
	isLoggedIn, ok := val.(bool)
	if !ok || (ok && !isLoggedIn) {
		return d, ErrorUnauthorized{missingData: "isLoggedIn"}
	}

	val = ctx.Value(UserIdKey)
	userId, ok := val.(string)
	if !ok || (ok && userId == "") {
		return d, ErrorUnauthorized{missingData: "userId"}
	}

	d = UserData{
		UserId:          userId,
		IsAuthenticated: isLoggedIn,
	}
	return d, nil
}

type ErrorUnauthorized struct {
	missingData string
}

func (r ErrorUnauthorized) Error() string {
	return fmt.Sprintf("user login information not provided in request context: %s", r.missingData)
}
