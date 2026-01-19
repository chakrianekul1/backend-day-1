# Stage 1: Build
# Changed from 1.21 to 1.23 to match your go.mod
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Stage 2: Run
FROM alpine:latest
WORKDIR /root/
# Install security certificates so Go can make external calls if needed
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]