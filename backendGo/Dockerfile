# Dockerfile
FROM golang:1.23-alpine AS build

WORKDIR /app

# Baixa módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/server

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /app

# Copia binário, .env e cria pasta de backups
COPY --from=build /app/api .
COPY .env .
RUN mkdir -p backups

EXPOSE 3000

CMD ["./api"]