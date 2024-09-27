FROM golang:1.22.5-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/server.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /gateway ./cmd/gateway/gateway.go

FROM alpine:3.18

RUN apk add --no-cache tini

COPY --from=builder /server /bin/server
COPY --from=builder /gateway /bin/gateway

ENV GIN_MODE=release
ENV DOCKER_PATH_ASSETS=/var/stock-data.json

EXPOSE 50051 8081

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/bin/sh", "-c", "/bin/server & /bin/gateway"]
