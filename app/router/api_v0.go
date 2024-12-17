package router

import (
	"github.com/go-bumbu/todo-app/app/handlers"
	"github.com/go-bumbu/userauth/authenticator"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *MainAppHandler) attachApiV0(r *mux.Router) {
	// this sub router does enforce authentication
	authHandlers := []authenticator.AuthHandler{h.SessionAuth}
	auth := authenticator.New(authHandlers, h.logger, nil, nil)

	r.Use(auth.Middleware)
	h.attachApiTask(r)
}

func (h *MainAppHandler) attachApiTask(r *mux.Router) {
	// add tasks api
	th := handlrs.TodoListHandler{TaskManager: h.todoListMngr}
	r.Path("/tasks").Methods(http.MethodGet).Handler(th.List())
	r.Path("/task").Methods(http.MethodPost).Handler(th.Create())
	r.Path("/task/{ID}").Methods(http.MethodGet).Handler(th.Read())
	r.Path("/task/{ID}").Methods(http.MethodDelete).Handler(th.Delete())
	r.Path("/task/{ID}").Methods(http.MethodPut).Handler(th.Update())
}
