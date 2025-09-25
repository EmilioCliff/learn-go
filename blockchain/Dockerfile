# Build stage
FROM golang:1.24.7-alpine AS builder
WORKDIR /app
COPY . .

RUN go build -o main /app/cmd/server/main.go

# Run stage
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/main .
CMD ["./main"]