FROM golang:1.25-alpine AS builder

WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod tidy

COPY . .
RUN go build -o app.bin cmd/server/main.go

# -----
FROM golang:1.25-alpine AS runner
WORKDIR /app
COPY --from=builder /app/app.bin /app/app.bin
# COPY .env.example ./.env.example
# EXPOSE 8080
CMD ["./server"]

# -v ./env:/app.env -v ./firebase-key.json:/app/firebase-key.json