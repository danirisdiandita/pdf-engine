FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main cmd/server/main.go

FROM alpine:latest
COPY --from=builder /app/main /app/main
WORKDIR /app
CMD ["./main"]