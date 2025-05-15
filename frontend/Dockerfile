# --- Stage 1 : build ---
FROM node:slim AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
# Definindo a variável VITE_API_BASE_URL diretamente para o Vite
ENV VITE_API_BASE_URL="https://estoque.parron01.com"

RUN npm run build   # gera /app/dist

# --- Stage 2 : runtime ---
FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
# SPA: redireciona tudo p/ index.html
# nginx.conf direto no Dockerfile
RUN printf "server {\n\
  listen 80;\n\
  server_name _;\n\
  root /usr/share/nginx/html;\n\
  index index.html;\n\
  location / { try_files \$uri /index.html; }\n\
}\n" > /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
