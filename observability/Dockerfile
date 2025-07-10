FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["/app/main"]