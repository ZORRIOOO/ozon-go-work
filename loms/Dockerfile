FROM golang:1.22.5-alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/server.go

FROM scratch

COPY --from=builder /server /bin/server

ENV GIN_MODE=release

EXPOSE 8082

ENTRYPOINT ["/bin/server"]
