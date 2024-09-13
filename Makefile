build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker-compose up --force-recreate --build -d

test:
	go test ./cart/internal/...

test-coverage:
	go test -cover ./cart/...

test-integration:
	go test ./cart/test/...
