# Projeto Gerenciador de Estoque - Backend em Go

Este é o backend para o sistema de Gerenciamento de Estoque, desenvolvido em Go. O projeto é responsável pela lógica de negócios, persistência de dados, autenticação e backups automáticos, utilizando PostgreSQL como banco de dados.

## Tecnologias Utilizadas

- **Go 1.23+:** Linguagem de programação utilizada para o desenvolvimento do backend.
- **PostgreSQL:** Banco de dados relacional utilizado para armazenamento de dados.
- **Gin:** Framework web para Go, utilizado para criar a API RESTful.
- **JWT:** Para autenticação e autorização baseada em tokens.
- **Cron:** Para agendamento de tarefas automáticas (backups semanais).
- **Docker:** Para containerização da aplicação.

## Estrutura do Projeto

```
backendGo
├── cmd
│   └── server
│       └── main.go          # Ponto de entrada da aplicação
├── internal
│   ├── config
│   │   └── config.go        # Configurações da aplicação
│   ├── controllers
│   │   ├── auth.go          # Controlador de autenticação
│   │   ├── history.go       # Controlador de histórico
│   │   └── product.go       # Controlador de produtos
│   ├── database
│   │   └── database.go      # Conexão com o banco de dados PostgreSQL
│   ├── middleware
│   │   └── auth.go          # Middleware de autenticação
│   ├── models
│   │   └── models.go        # Modelos de dados
│   ├── routes
│   │   └── routes.go        # Configuração das rotas
│   └── utils
│       └── backup.go        # Funções utilitárias para backup
├── backups                  # Diretório onde os backups são armazenados
├── go.mod                   # Definição do módulo Go e dependências
├── go.sum                   # Checksums das dependências
├── Dockerfile               # Instruções para construir a imagem Docker
├── docker-compose.yml       # Configuração para orquestração de containers
├── .env                     # Variáveis de ambiente
└── README.md                # Documentação do projeto
```

## Configuração e Execução

### Pré-requisitos

- Go 1.23 ou superior
- PostgreSQL
- Utilitário pg_dump (para backups)

### Variáveis de Ambiente

O arquivo `.env` na raiz do projeto deve conter:

```env
PORT=3000
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=inventory
JWT_SECRET=sua_chave_secreta_jwt_aqui
JWT_EXPIRATION=168h
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123
```

### Executando com Docker (Recomendado)

O modo mais simples de executar o projeto é utilizando Docker Compose:

```sh
docker-compose up
```

Isso iniciará tanto o banco de dados PostgreSQL quanto a API.

### Executando Localmente

1. Certifique-se de ter o PostgreSQL instalado e em execução
2. Instale as dependências:
   ```sh
   go mod download
   ```
3. Execute o servidor:
   ```sh
   go run cmd/server/main.go
   ```

O servidor estará disponível em `http://localhost:3000`.

## Funcionalidades Principais

### Autenticação de Usuários

- Sistema de login seguro utilizando JWT
- Usuário admin criado automaticamente na inicialização
- Verificação de tokens para rotas protegidas

### Gerenciamento de Produtos

- CRUD completo para produtos (criar, ler, atualizar, deletar)
- Produtos incluem ID, nome, unidade (L ou kg) e quantidade

### Histórico de Alterações

- Registro de todas as modificações no estoque
- Armazenamento de alterações em formato JSON para flexibilidade

### Sistema de Backup Automático

- Backups semanais automáticos (todo domingo às 3:00)
- Formato binário PostgreSQL para facilitar restauração
- Limpeza automática de backups com mais de 30 dias

## Endpoints da API

### Autenticação

- `POST /api/auth/login`: Login de usuário, retorna token JWT.
- `GET /api/auth/verify`: Verifica a validade de um token JWT.
- `GET /api/auth/health`: Verifica status de saúde do servidor.

### Produtos

- `GET /api/products`: Lista todos os produtos.
- `GET /api/products/:id`: Obtém um produto específico pelo ID.
- `POST /api/products`: Cria um novo produto (requer autenticação).
- `PUT /api/products/:id`: Atualiza um produto existente (requer autenticação).
- `DELETE /api/products/:id`: Remove um produto (requer autenticação).

### Histórico

- `GET /api/history`: Lista todos os registros de histórico de alterações (requer autenticação).
- `POST /api/history`: Adiciona um novo registro ao histórico (requer autenticação).

## CORS

O servidor está configurado para aceitar requisições de qualquer origem (CORS habilitado), facilitando a integração com frontends em diferentes domínios durante o desenvolvimento.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.
