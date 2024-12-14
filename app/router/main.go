package router

import (
	_ "embed"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-bumbu/http/middleware"
	handlrs "github.com/go-bumbu/todo-app/app/handlers"
	"github.com/go-bumbu/userauth"
	"github.com/go-bumbu/userauth/handlers/sessionauth"
	"github.com/gorilla/mux"

	"github.com/go-bumbu/todo-app/app/spa"
	"github.com/go-bumbu/todo-app/internal/model/todolist"
)

type Cfg struct {
	Db             *gorm.DB
	SessionAuth    *sessionauth.Manager
	UserMngr       userauth.LoginHandler
	TodoListMngr   *todolist.Manager
	Logger         *slog.Logger
	ProductionMode bool
}

// MainAppHandler is the entrypoint http handler for the whole application
type MainAppHandler struct {
	router         *mux.Router
	db             *gorm.DB
	SessionAuth    *sessionauth.Manager
	userMngr       userauth.LoginHandler
	todoListMngr   *todolist.Manager
	logger         *slog.Logger
	productionMode bool
}

func (h *MainAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func New(cfg Cfg) (*MainAppHandler, error) {
	r := mux.NewRouter()
	app := MainAppHandler{
		router:       r,
		db:           cfg.Db,
		SessionAuth:  cfg.SessionAuth,
		userMngr:     cfg.UserMngr,
		logger:       cfg.Logger,
		todoListMngr: cfg.TodoListMngr,
	}

	// normally on real world project you would never add this middleware
	app.addDelayMiddleware(app.router)
	app.addPromMiddleware(app.router)

	app.attachUserAuth(app.router.PathPrefix("/auth").Subrouter())

	// add a handler for /api/v0, this includes authentication on tasks
	app.attachApiV0(app.router.PathPrefix("/api/v0").Subrouter())

	// attach another handler to /demo to showcase other use-cases
	app.attachDemo(app.router.Path("/demo").Subrouter())

	// protect /basic with basic auth
	app.attachBasicAuthProtected(app.router.Path("/basic").Subrouter())

	// add the spa to path /
	err := app.attachTodolistSpa(app.router.PathPrefix("/").Subrouter(), "/")
	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (h *MainAppHandler) attachTodolistSpa(r *mux.Router, path string) error {
	// if you want to serve the spa from the root, pass "/" to the spa handler and the path prefix
	// note that the SPA base and route needs to be adjusted accordingly
	spaHandler, err := spa.TodoApp(path)
	if err != nil {
		return err
	}
	r.Methods(http.MethodGet).PathPrefix(path).Handler(spaHandler)
	return nil
}

func (h *MainAppHandler) attachUserAuth(r *mux.Router) {

	//  LOGIN
	r.Path("/login").Methods(http.MethodPost).Handler(handlrs.UserLoginHandler(h.SessionAuth, h.userMngr))
	r.Path("/login").Methods(http.MethodOptions).Handler(
		http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}))
	// TODO add a basic form login here to the GET method
	r.Path("/login").HandlerFunc(StatusErr(http.StatusMethodNotAllowed))

	// LOGOUT
	r.Path("/logout").Handler(handlrs.UserLogoutHandler(h.SessionAuth))

	// STATUS
	r.Path("/status").Methods(http.MethodGet).Handler(handlrs.UserStatusHandler(h.SessionAuth))
	r.Path("/status").HandlerFunc(StatusErr(http.StatusMethodNotAllowed))

	// OPTIONS
	//r.Path("/user/options").Methods(http.MethodGet).Handler(handlers.StatusErr(http.StatusNotImplemented))
	r.Path("/user/options").HandlerFunc(StatusErr(http.StatusMethodNotAllowed))
}

func (h *MainAppHandler) addDelayMiddleware(r *mux.Router) {
	if !h.productionMode {
		throttle := middleware.ReqDelay{
			MinDelay: 1500 * time.Millisecond,
			MaxDelay: 3000 * time.Millisecond,
			On:       false,
		}
		r.Use(throttle.Delay)
	}

}

func (h *MainAppHandler) addPromMiddleware(r *mux.Router) {
	//add observability

	hist := middleware.NewHistogram("", nil, nil)
	r.Use(func(handler http.Handler) http.Handler {
		return middleware.PromLogMiddleware(handler, hist, h.logger)
	})

}

func StatusErr(status int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	}
}
