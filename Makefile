build-all:
	docker compose up --force-recreate --build -d

run-all: build-all
	docker-compose start

run-local:
	@echo "Starting cart 1..."
	@go run ./cart/cmd/server & \
	echo "Starting loms 2..." && \
	go run ./loms/cmd/server & \
	echo "Starting HTTP-Gateway 3..." && \
	go run ./loms/cmd/gateway

CART := "./cart/internal/pkg/cart/..."
LOMS := "./loms/internal/service/loms/..."
LIST = $(shell go list ${CART} ${LOMS} | grep -v -e mock -e model -e benchmark)
INTEGRATION_LIST = $(shell go list ./cart/test/... ./loms/test/...)

test:
	go test $(LIST)

test-coverage:
	go test -coverprofile=coverage.out $(LIST) && go tool cover -func=coverage.out

test-integration:
	go test ${INTEGRATION_LIST}

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
