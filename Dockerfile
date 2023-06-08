FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY ./go-auth-backend/ .
RUN go mod download
RUN go build -ldflags='-s' -o=./bin/api ./cmd/api


FROM alpine:latest AS runner
WORKDIR /
COPY --from=builder /app/bin/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]