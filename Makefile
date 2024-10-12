# Deploy

.PHONY: run-local
run-local:
	@echo "Starting cart..."
	@go run ./cart/cmd/server & \
	echo "Starting loms..." && \
	go run ./loms/cmd/server & \
	echo "Starting HTTP-Gateway..." && \
	go run ./loms/cmd/gateway

DOCKER_YML=./docker-compose.yml
ENV_NAME="homework"

.PHONY: compose-up
compose-up:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} up -d

.PHONY: compose-down
compose-down:
	docker-compose -p ${ENV_NAME} -f -f ${DOCKER_YML} stop

.PHONY: compose-rm
compose-rm:
	docker-compose -p ${ENV_NAME} -f ${DOCKER_YML} rm -fvs

.PHONY: compose-rs
compose-rs:
	make compose-rm
	make compose-up

# Tests

CART := "./cart/internal/pkg/cart/..."
LOMS := "./loms/internal/service/loms/..."
LIST = $(shell go list ${CART} ${LOMS} | grep -v -e mock -e model -e benchmark)
INTEGRATION_LIST = $(shell go list ./cart/test/... ./loms/test/...)
TEST_BENCH = "./cart/internal/pkg/cart/repository/..."

.PHONY: test
test:
	go test $(LIST)

.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out $(LIST) && go tool cover -func=coverage.out

.PHONY: test-integration
test-integration:
	go test ${INTEGRATION_LIST}

.PHONY: test-bench
test-bench:
	go test -bench=. ${TEST_BENCH}

.PHONY: test-all
test-all: test test-coverage test-integration test-bench

.PHONY: gocyclo-get
gocyclo-get:
	$ go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

.PHONY: gocyclo-lint
gocyclo-lint:
	gocyclo -over 15 -avg -ignore "_test|mock" ./cart

.PHONY: gocognit-get
gocognit-get:
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

.PHONY: gocognit-lint
gocognit-lint:
	gocognit -over 15 -avg -ignore "_test|mock" ./cart
