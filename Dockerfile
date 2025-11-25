# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /sistergo ./cmd/server

# Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /sistergo /sistergo
COPY .env /app/.env

EXPOSE 8080
ENTRYPOINT ["/sistergo"]
