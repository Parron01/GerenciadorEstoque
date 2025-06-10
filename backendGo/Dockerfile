# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instale dependências de compilação
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the application with the same name as dev version
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Stage 2: Final image
FROM alpine:latest

# Instalar ferramentas necessárias (same as dev)
RUN apk add --no-cache postgresql-client ca-certificates tzdata bash curl

# Configure timezone
ENV TZ=America/Sao_Paulo

WORKDIR /app

# Copy wait script from builder stage
COPY --from=builder /app/scripts/wait-for-db.sh ./wait-for-db.sh
RUN chmod +x ./wait-for-db.sh

# Copy migrations from the builder stage
COPY --from=builder /app/migrations ./migrations

# Copy the built executable with same name as dev
COPY --from=builder /app/server .

# Set production mode
ENV GIN_MODE=release

# Expose port
EXPOSE 3000

# Use the wait script as entrypoint (same as dev)
ENTRYPOINT ["./wait-for-db.sh"]

# Command to run the executable (same path as dev)
CMD ["/app/server"]
