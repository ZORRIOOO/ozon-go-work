build:
	docker build -t golang:latest .

run:
	docker run -p 8082:8082 golang:latest

build-dev:
	RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

run-dev:
	go run ./cmd/server/server.go
