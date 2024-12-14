package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HttpErr struct {
	Error string
	Code  int
}

func JsonErrorHandler(err string, code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonError(w, err, code)
	})
}

func jsonError(w http.ResponseWriter, error string, code int) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	payload := HttpErr{
		Error: error,
		Code:  code,
	}
	byteErr, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, string(byteErr))
}
