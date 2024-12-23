include .env
export

LOCAL_BIN ?= $(CURDIR)/bin

.PHONY: run
run: ### run app
	go run cmd/simplecards/main.go

.PHONY: lint
lint: ### run linter
	golangci-lint run ./...

.PHONY: test
test: ### run tests
	go test ./...

.PHONY: migrate-up
migrate-up: ### run migrations
	bin/goose up

.PHONY: swag
swag:
	swag init -g internal/app/app.go

.PHONY: reqs
reqs: ### install binary deps to bin/
	GOBIN=$(LOCAL_BIN) go install go.uber.org/mock/mockgen@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
