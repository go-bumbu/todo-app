COMMIT_SHA_SHORT ?= $(shell git rev-parse --short=12 HEAD)
PWD_DIR := ${CURDIR}

default: help

#==========================================================================================
##@ Testing
#==========================================================================================
test: ## run fast go tests
	@go test ./... -cover

lint: ## run go linter
	# depends on https://github.com/golangci/golangci-lint
	@golangci-lint run

benchmark: ## run go benchmarks
	@go test -run=^$$ -bench=. ./...

license-check: ## check for invalid licenses
	# depends on : https://github.com/elastic/go-licence-detector
	@go list -m -mod=readonly  -json all  | go-licence-detector -includeIndirect -validate -rules allowedLicenses.json

.PHONY: verify
verify: test license-check lint benchmark ## run all tests

cover-report: ## generate a coverage report
	go test -covermode=count -coverpkg=./... -coverprofile cover.out  ./...
	go tool cover -html cover.out -o cover.html
	open cover.html

#==========================================================================================
##@ Running
#==========================================================================================
run: ## start the GO service
	@APP_LOG_LEVEL="debug" go run main.go start -c zarf/appData/config.yaml

run-ui: package-ui run## build the UI and start the GO service

#==========================================================================================
##@ Building
#==========================================================================================
package-ui: build-ui ## build the web and copy into Go pacakge
	rm -rf ./app/spa/files/ui*
	mkdir -p ./app/spa/files/ui
	cp -r ./webui/dist/* ./app/spa/files/ui/
	touch ./app/spa/files/ui/.gitkeep
build-ui:
	@cd webui && \
	npm install && \
	export VITE_BASE="/ui" && \
	npm run build

build: package-ui ## use goreleaser to build
	@goreleaser build --auto-snapshot --clean

build-linux: package-ui ## use goreleaser to build a single linux target
	@goreleaser release --clean --auto-snapshot --skip publish


#==========================================================================================
##@ Release
#==========================================================================================

.PHONY: check-git-clean
check-git-clean: # check if git repo is clen
	@git diff --quiet

.PHONY: check-branch
check-branch:
	@current_branch=$$(git symbolic-ref --short HEAD) && \
	if [ "$$current_branch" != "main" ]; then \
		echo "Error: You are on branch '$$current_branch'. Please switch to 'main'."; \
		exit 1; \
	fi

check_env: # check for needed envs
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )


tag: check_env check-branch check-git-clean verify ## create a tag and push to git
	@git diff --quiet || ( echo 'git is in dirty state' ; exit 1 )
	@[ "${version}" ] || ( echo ">> version is not set, usage: make release version=\"v1.2.3\" "; exit 1 )
	@git tag -d $(version) || true
	@git tag -a $(version) -m "Release version: $(version)"
	@git push --delete origin $(version) || true
	@git push origin $(version) || true



#==========================================================================================
#  Help
#==========================================================================================
.PHONY: help
help: # Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
