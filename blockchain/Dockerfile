# Build stage
FROM golang:1.24.7-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main /app/cmd/server/main.go

# Run stage
FROM ubuntu:22.04
COPY --from=builder /app/main .
EXPOSE 3031

CMD ["./main"]