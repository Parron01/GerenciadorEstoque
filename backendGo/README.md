# Projeto Gerenciador de Estoque - Backend em Go

Este é o backend para o sistema de Gerenciamento de Estoque, reescrito em Go. O projeto é responsável pela lógica de negócios, persistência de dados e autenticação, utilizando PostgreSQL como banco de dados.

## Tecnologias Utilizadas

- **Go:** Linguagem de programação utilizada para o desenvolvimento do backend.
- **PostgreSQL:** Banco de dados relacional utilizado para armazenamento de dados.
- **Gin:** Framework web para Go, utilizado para criar a API RESTful.
- **Gorm:** ORM para Go, utilizado para interagir com o banco de dados.
- **JWT:** Para autenticação e autorização baseada em tokens.
- **dotenv:** Para gerenciamento de variáveis de ambiente.

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
│   │   └── postgres.go      # Conexão com o banco de dados PostgreSQL
│   ├── middleware
│   │   └── auth.go          # Middleware de autenticação
│   ├── models
│   │   └── models.go        # Modelos de dados
│   ├── routes
│   │   └── routes.go        # Configuração das rotas
│   └── utils
│       └── backup.go        # Funções utilitárias para backup
├── go.mod                    # Definição do módulo Go
├── go.sum                    # Checksums das dependências
├── Dockerfile                # Instruções para construir a imagem Docker
├── .env.example              # Exemplo de arquivo de variáveis de ambiente
└── README.md                 # Documentação do projeto
```

## Configuração e Execução

### Pré-requisitos

- Go (versão 1.18 ou superior recomendada)
- PostgreSQL (instalado e em execução)
- Dependências do Go (instaladas via `go mod tidy`)

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do diretório `backendGo` com as seguintes variáveis:

```env
PORT=8080
DATABASE_URL=postgres://usuario:senha@localhost:5432/nome_do_banco
JWT_SECRET=sua_chave_secreta_jwt_aqui
```

### Instalação

1. Navegue até o diretório `backendGo`:
   ```sh
   cd backendGo
   ```

2. Instale as dependências:
   ```sh
   go mod tidy
   ```

### Executando a Aplicação

Para iniciar o servidor, execute o seguinte comando:

```sh
go run cmd/server/main.go
```

O servidor estará disponível em `http://localhost:8080`.

### Endpoints da API

- **Autenticação:**
  - `POST /api/auth/login`: Login de usuário, retorna token JWT.
  - `GET /api/auth/verify`: Verifica a validade de um token JWT.

- **Produtos:**
  - `GET /api/products`: Lista todos os produtos.
  - `GET /api/products/:id`: Obtém um produto específico pelo ID.
  - `POST /api/products`: Cria um novo produto (requer autenticação).
  - `PUT /api/products/:id`: Atualiza um produto existente (requer autenticação).
  - `DELETE /api/products/:id`: Remove um produto (requer autenticação).

- **Histórico:**
  - `GET /api/history`: Lista todos os registros de histórico de alterações (requer autenticação).

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.