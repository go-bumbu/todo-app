module github.com/go-bumbu/todo-app

go 1.23

toolchain go1.23.1

replace github.com/go-bumbu/userauth => ../userauth

replace github.com/go-bumbu/http => ../http

require (
	github.com/andresbott/go-carbon v0.1.0
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/go-bumbu/config v0.1.0
	github.com/go-bumbu/http v0.2.0
	github.com/go-bumbu/userauth v0.0.0-00010101000000-000000000000
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/securecookie v1.1.2
	github.com/mattn/go-isatty v0.0.20
	github.com/phsym/console-slog v0.3.1
	github.com/samber/slog-formatter v1.1.1
	github.com/spf13/cobra v1.8.1
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.61.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.13.1 // indirect
	github.com/samber/lo v1.47.0 // indirect
	github.com/samber/slog-multi v1.2.4 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)