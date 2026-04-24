FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth-service ./cmd

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/auth-service /usr/local/bin/auth-service
# App currently requires godotenv.Load() at startup.
COPY .env ./.env

EXPOSE 8080

CMD ["auth-service"]