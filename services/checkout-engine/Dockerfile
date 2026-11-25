FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o checkout-engine ./cmd/api/main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/checkout-engine .
EXPOSE 8081
CMD ["./checkout-engine"]