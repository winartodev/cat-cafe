FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

RUN git clone https://github.com/golang-migrate/migrate /migrate-src && \
    cd /migrate-src && \
    git checkout v4.17.0 && \
    go build -tags 'postgres' -ldflags="-s -w" -o /usr/local/bin/migrate ./cmd/migrate

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/http

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8888

ENTRYPOINT ["/entrypoint.sh"]
CMD ["./main"]
