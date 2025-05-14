# Dockerfile do Frontend - Vue.js (Vite)
# Stage 1: Build
FROM node:slim AS build
WORKDIR /app

# Instalando todas as dependências (prod + dev) apenas para o build
COPY package*.json ./
RUN npm install

# Copiando o código da aplicação
COPY . .

# Executando o build (typescript + vite)
RUN npm run build

# Stage 2: Imagem final simplificada (Apenas arquivos necessários)
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80

# Development Dockerfile for Frontend
FROM node:slim

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy the application code
COPY . .

# Expose Vite's default port
EXPOSE 5173

# Command to run the development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
