package spa

import (
	"embed"
	handlers "github.com/go-bumbu/http/handlers/spa"

	"net/http"
)

//go:embed files/ui/*
var UiFiles embed.FS

func TodoApp(path string) (http.Handler, error) {
	return handlers.NewSpaHAndler(
		UiFiles,
		"files/ui",
		path,
	)
}
