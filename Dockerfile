# --- Étape 1: Build ---
# Utiliser une image Go 1.23 pour la compilation
FROM golang:1.23-alpine AS builder

# Le reste du fichier ne change pas
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server/main ./server

# --- Étape 2: Final ---
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /server/main .
EXPOSE 50051
CMD ["./main"]