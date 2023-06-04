FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -ldflags='-s' -o=./bin/api ./cmd/api


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/bin/api .
EXPOSE 8080
ENTRYPOINT ["./bin/api"]