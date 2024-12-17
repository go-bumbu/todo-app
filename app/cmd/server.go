package cmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-bumbu/http/server"
	"github.com/go-bumbu/userauth"
	"github.com/go-bumbu/userauth/handlers/sessionauth"
	"github.com/go-bumbu/userauth/userstore/staticusers"
	"github.com/gorilla/securecookie"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-bumbu/todo-app/app/config"
	"github.com/go-bumbu/todo-app/app/logger"
	"github.com/go-bumbu/todo-app/app/metainfo"
	"github.com/go-bumbu/todo-app/app/router"
	"github.com/go-bumbu/todo-app/internal/model/todolist"
)

const dbFile = "carbon.db"

func serverCmd() *cobra.Command {
	var configFile = "./config.yaml"
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start a web server",
		Long:  "start a web server demonstrating the different features of the library",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(configFile)
		},
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", configFile, "config file")
	return cmd
}

func runServer(configFile string) error {
	cfg, err := config.Get(configFile)
	if err != nil {
		return err
	}
	_ = cfg
	// setup the logger
	l, err := logger.GetDefault(logger.GetLogLevel(cfg.Env.LogLevel))
	if err != nil {
		return err
	}

	l.Info("App startup",
		slog.String("component", "startup"),
		slog.String("version", metainfo.Version),
		slog.String("Build Date", metainfo.BuildTime),
		slog.String("commit", metainfo.ShaVer),
	)
	// print config messages delayed
	for _, m := range cfg.Msgs {
		if m.Level == "info" {
			l.Info(m.Msg, slog.String("component", "config"))
		} else {
			l.Debug(m.Msg, slog.String("component", "config"))
		}
	}

	// initialize DB
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		// TODO add slogger translation
		//Logger: zeroGorm.New(l.ZeroLog, zeroGorm.Cfg{IgnoreRecordNotFoundError: true}),
	})
	if err != nil {
		return err
	}

	userStore, err := getUserStore(cfg, l)
	if err != nil {
		return err
	}

	store, _ := sessionauth.NewFsStore("", securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	// create an instance of session auth
	sessionAuth, _ := sessionauth.New(sessionauth.Cfg{
		Store:         store,
		SessionDur:    time.Hour,       // time the user is logged in
		MaxSessionDur: 24 * time.Hour,  // time after the user is forced to re-login anyway
		MinWriteSpace: 2 * time.Minute, // throttle write operations on the session
	})

	// init task manager
	todoList, err := todolist.New(db)
	if err != nil {
		return fmt.Errorf("unable to create task manager :%v", err)
	}
	routerCfg := router.Cfg{
		Db:          db,
		SessionAuth: sessionAuth,
		UserMngr: userauth.LoginHandler{
			UserStore: userStore,
		},
		Logger:       l,
		TodoListMngr: todoList,
	}
	mainAppHandler, err := router.New(routerCfg)
	if err != nil {
		return fmt.Errorf("unable to create initialize main app handler:%v", err)
	}

	s, err := server.New(server.Cfg{
		Addr:    cfg.Server.Addr(),
		Handler: mainAppHandler,
		SkipObs: true,
		ObsAddr: cfg.Obs.Addr(),
		//ObsHandler: handlers.Observability(),
		Logger: func(msg string, isErr bool) {
			// TODO use slogger ?
			if isErr {
				l.Warn(msg, slog.String("component", "server"))
			} else {
				l.Info(msg, slog.String("component", "server"))
			}
		},
	})
	if err != nil {
		return err
	}

	return s.Start()

}

func getUserStore(cfg config.AppCfg, l *slog.Logger) (userauth.UserGetter, error) {
	var userGet userauth.UserGetter
	// load the correct user manager
	switch cfg.Auth.UserStore.StoreType {
	case "static":
		staticUsers := staticusers.Users{}
		for _, u := range cfg.Auth.UserStore.Users {
			staticUsers.Add(staticusers.User{
				Id:      u.Name,
				HashPw:  userauth.MustHashPw(u.Pw),
				Enabled: true,
			})
		}

		l.Debug("loading static users", slog.String("component", "users"),
			slog.Int("amount", len(staticUsers.Users)))
		userGet = &staticUsers

	case "file":

		if cfg.Auth.UserStore.FilePath == "" {
			return userGet, fmt.Errorf("no path for users file is empty")
		}
		users, err := staticusers.FromFile(cfg.Auth.UserStore.FilePath)
		if err != nil {
			return userGet, err
		}
		userGet = users
		l.Debug("loading users from file", slog.String("component", "users"),
			slog.Int("amount", len(users.Users)),
			slog.String("file", cfg.Auth.UserStore.FilePath))
	default:
		return userGet, fmt.Errorf("wrong user store in configuration, %s is not supported", cfg.Auth.UserStore.StoreType)
	}
	return userGet, nil
}
