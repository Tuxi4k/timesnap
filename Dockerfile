FROM golang:alpine AS builder

RUN adduser -D -u 10001 appuser

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init --ot json -g cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go
RUN mkdir -p sqlite

FROM scratch

WORKDIR /app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /build/main /app/main
COPY --from=builder /build/config.prod.yml /app/config.prod.yml
COPY --from=builder /build/docs /app/docs
COPY --chown=10001:10001 --from=builder /build/sqlite /app/sqlite

USER appuser

ENTRYPOINT ["./main"]