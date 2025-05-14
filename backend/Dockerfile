# Dockerfile do Backend - Node.js (TypeScript)
# Stage 1: Build
FROM node:22-slim AS build
WORKDIR /app

# Instalar dependências necessárias para compilação
RUN apt-get update && apt-get install -y \
    python3 \
    make \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copiar arquivos de configuração
COPY package*.json tsconfig.json ./

# Instalar dependências
RUN npm ci

# Copiar código fonte
COPY . .

# Compilar o projeto
RUN npm run build

# Stage 2: Produção
FROM node:22-slim
WORKDIR /app

# Copiar apenas os arquivos necessários
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/package.json ./

# Criar pasta do banco de dados
RUN mkdir -p dist/src/database

# Configurar ambiente
ENV NODE_ENV=production
EXPOSE 3000

# Comando para iniciar a aplicação
CMD ["node", "dist/src/index.js"]