# FROM golang:1.25-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod tidy

# COPY . .
# RUN go build -o app.bin cmd/server/main.go
# RUN go build -o /app/server ./cmd/server && ls -l /app

# -----
# FROM golang:1.25-alpine AS runner
# WORKDIR /app
# COPY --from=builder /app/app.bin /app/app.bin
# COPY .env.example ./.env.example
# EXPOSE 8080
# CMD ["/app/server"]

# -v ./env:/app.env -v ./firebase-key.json:/app/firebase-key.json

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server

FROM alpine:3.20 AS runner

WORKDIR /app

COPY --from=builder /app/server /app/server
COPY firebase-key.json /app/firebase-key.json

EXPOSE 8080

CMD ["/app/server"]