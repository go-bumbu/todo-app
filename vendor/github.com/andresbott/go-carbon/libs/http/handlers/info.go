package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// StaticInfo is a handler that will render a static map[string]any to an endpoint
// this is useful to expose static information like build version, or enabled/disabled features that won't change
// over the lifecycle of the application and indicate the clients how to behave, e.g. if authentication is disabled.
func StaticInfo(inf map[string]interface{}) (http.Handler, error) {
	_, err := json.Marshal(inf)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal the content: %v", err)
	}

	return &staticInfoHandler{inf}, nil
}

type staticInfoHandler struct {
	data map[string]interface{}
}

func (f *staticInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f.handleGet(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildtime"`
	Commit    string `json:"commit"`
}

func (f *staticInfoHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	respJson, err := json.Marshal(f.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(respJson)
}
