
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .

RUN go build -o main /app/cmd/server/main.go

# Run stage
FROM alpine:3.19
COPY --from=builder /app/main .

CMD ["./main"]