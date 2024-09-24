build-all:
	docker compose up --force-recreate --build -d

run-all: build-all
	docker-compose start

LIST = $(shell go list ./cart/internal/pkg/cart/... | grep -v -e mock -e model -e benchmark)

test:
	go test $(LIST)

test-coverage:
	go test -coverprofile=coverage.out $(LIST) && go tool cover -func=coverage.out

test-integration:
	go test ./cart/test/...

test-bench:
	go test -bench=. ./cart/internal/pkg/cart/repository/...

test-all: test test-coverage test-integration test-bench

gocyclo-get:
	$ go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

gocyclo-lint:
	gocyclo -over 15 -avg -ignore "_test|mock" ./cart

gocognit-get:
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

gocognit-lint:
	gocognit -over 15 -avg -ignore "_test|mock" ./cart
