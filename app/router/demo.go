package router

import (
	"github.com/andresbott/go-carbon/libs/auth"
	"github.com/andresbott/go-carbon/libs/http/handlers"
	"github.com/go-bumbu/userauth/handlers/basicauth"

	"github.com/gorilla/mux"
)

func (h *MainAppHandler) attachDemo(r *mux.Router) {
	demoPage := handlers.SimpleText{
		Text: "Demo root page",
		Links: []handlers.Link{
			{Text: "Basic auth protected (demo:demo)", Url: "/basic"},
			{Text: "session Authentication protected page", Url: "/session"},
			//{
			//	Text: "session based login (demo:demo)",
			//	Url:  handlers2.sessionLogin,
			//},
			{Text: "Authenticaion", Child: []handlers.Link{
				{Text: "Status (/auth/status)", Url: "/auth/status"},
			}},
			{Text: "Json API", Child: []handlers.Link{
				{Text: "User options", Url: "/api/v0/user/options"},
			}},
			{Text: "Observability", Child: []handlers.Link{
				{Text: "metrics", Url: "http://localhost:9090/metrics"},
			}},
		},
	}
	r.Path("").Handler(demoPage)
}

func (h *MainAppHandler) attachBasicAuthProtected(r *mux.Router) {
	demoPage := handlers.SimpleText{
		Text: "Basic auth protected page",
		Links: []handlers.Link{
			{
				Text: "Back",
				Url:  "/demo",
			},
		},
	}
	basicAH := basicauth.NewHandler(h.userMngr, "", true, h.logger)
	// use the middleware to protect the page
	r.Use(basicAH.Middleware)
	r.Path("").Handler(demoPage)

}

// const SessionLogin = "/session-login"

func SessionProtected(r *mux.Router, session *auth.SessionMgr) error {
	pageHandler := handlers.SimpleText{
		Text: "Page protected by session auth",
		Links: []handlers.Link{
			{Text: "back to root", Url: "../"},
		},
	}

	ProtectedPage := session.Middleware(&pageHandler)
	r.Path("/session").Handler(ProtectedPage)

	return nil
}