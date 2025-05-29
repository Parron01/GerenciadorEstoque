FROM golang:1.23-alpine AS build

# Instale dependências de compilação
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/server

FROM alpine:latest

# Instale ferramentas de cliente PostgreSQL e certificados CA
RUN apk add --no-cache postgresql-client ca-certificates tzdata

# Configure o fuso horário
ENV TZ=America/Sao_Paulo

WORKDIR /app

# Copie o binário compilado e as migrações
COPY --from=build /app/api .
COPY --from=build /app/migrations ./migrations

ENV GIN_MODE=release
EXPOSE 3000

CMD ["./api"]
