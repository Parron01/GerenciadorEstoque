# Projeto Gerenciador de Estoque - Backend em Go

Este é o backend para o sistema de Gerenciamento de Estoque, desenvolvido em Go. O projeto é responsável pela lógica de negócios, persistência de dados, autenticação e backups automáticos, utilizando PostgreSQL como banco de dados.

## Tecnologias Utilizadas

- **Go 1.23+:** Linguagem de programação utilizada para o desenvolvimento do backend.
- **PostgreSQL:** Banco de dados relacional utilizado para armazenamento de dados.
- **Gin:** Framework web para Go, utilizado para criar a API RESTful.
- **JWT:** Para autenticação e autorização baseada em tokens.
- **golang-migrate/migrate:** Para gerenciamento e aplicação de migrações de banco de dados.
- **Cron:** Para agendamento de tarefas automáticas (backups semanais).
- **Docker:** Para containerização da aplicação.

## Funcionalidades Principais

- **Autenticação:** Login de usuário e verificação de token JWT.
- **Gerenciamento de Produtos:**
  - CRUD completo para produtos (ID, Nome, Unidade, Quantidade).
  - A quantidade de um produto é automaticamente atualizada por um gatilho no banco de dados (PostgreSQL) para refletir a soma das quantidades de seus lotes, se houver lotes.
- **Gerenciamento de Lotes (`product_lots`):**
  - CRUD completo para lotes (ID, ID do Produto, Quantidade, Data de Validade).
  - Lotes são associados a produtos. A criação, atualização ou exclusão de um lote dispara a atualização da quantidade total do produto pai.
- **Histórico de Alterações:**
  - Registro detalhado de todas as operações significativas (criação, atualização, exclusão de produtos e lotes).
  - Cada registro de histórico inclui `EntityType` (product, lote), `EntityID`, `ChangeDetails` (JSON com os detalhes da mudança), `CreatedAt`, e `BatchID`.
  - **Funcionamento do `BatchID`:**
    - **Geração Automática pelo Backend:** Para cada operação CRUD individual (ex: criar um lote, atualizar um produto), o backend automaticamente gera um registro de histórico.
    - **Agrupamento Iniciado pelo Cliente:** O cliente (frontend) pode gerar um UUID (ex: `X-Operation-Batch-ID`) e enviá-lo como um header HTTP em um conjunto de requisições de API (ex: múltiplas criações/atualizações de lotes e produtos feitas em uma única "transação" do usuário). O backend utilizará este `X-Operation-Batch-ID` como o `BatchID` para todos os registros de histórico gerados automaticamente por essas requisições específicas. Isso agrupa as alterações relacionadas.
    - **Fallback:** Se nenhum `X-Operation-Batch-ID` for fornecido pelo cliente, o backend geralmente define o `BatchID` do registro de histórico como o próprio `ID` do registro de histórico, tratando-o como uma operação individual.
    - **Endpoint de Batch Explícito (`POST /api/history/batch`):** Permite que um cliente envie um array de entradas de histórico pré-formatadas. O backend atribui um `BatchID` novo e único (gerado pelo servidor) a todas as entradas recebidas nesta requisição específica.
  - **Endpoints de Consulta de Histórico:**
    - `GET /api/history`: Retorna uma lista paginada de todos os registros de histórico.
    - `GET /api/history?batch_id={id}`: Retorna todos os registros de histórico associados a um `BatchID` específico.
    - `GET /api/history/batch/{batch_id}`: Similar ao anterior, focado em buscar um lote específico.
    - `GET /api/history/grouped`: **Novo endpoint** que retorna o histórico agrupado por `BatchID`. Cada grupo contém o `BatchID`, a data/hora da primeira entrada do lote, e todos os registros de histórico pertencentes àquele lote. Suporta paginação baseada nos lotes (batches).
    - `GET /api/history/{entity_type}/{entity_id}`: Retorna o histórico para uma entidade específica.
- **Backup:** Rotina de backup semanal do banco de dados (configurável via cron).

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
│   │   ├── product.go       # Controlador de produtos
│   │   └── lote_controller.go # Controlador de lotes de produtos
│   ├── database
│   │   └── database.go      # Conexão com o banco de dados PostgreSQL
│   ├── middleware
│   │   └── auth.go          # Middleware de autenticação
│   ├── models
│   │   └── models.go        # Modelos de dados (incluindo Lote)
│   ├── repository
│   │   ├── product_repository.go
│   │   ├── lote_repository.go
│   │   └── history_repository.go
│   ├── routes
│   │   └── routes.go        # Configuração das rotas
│   ├── service
│   │   ├── product_service.go
│   │   ├── lote_service.go
│   │   └── history_service.go
│   └── utils
│       └── backup.go        # Funções utilitárias para backup
├── backups                  # Diretório onde os backups são armazenados
├── migrations               # Arquivos de migração SQL
│   ├── 001_create_product_lots.sql
│   └── 002_update_history_table.sql
├── go.mod                   # Definição do módulo Go e dependências
├── go.sum                   # Checksums das dependências
├── Dockerfile               # Instruções para construir a imagem Docker de produção
├── Dockerfile.dev           # Instruções para construir a imagem Docker de desenvolvimento
├── docker-compose.yml       # Configuração para orquestração de containers
├── .env                     # Variáveis de ambiente
└── README.md                # Documentação do projeto
```

## Configuração e Execução

### Pré-requisitos

- Go 1.23 ou superior
- PostgreSQL
- Docker & Docker Compose (para execução via Docker)
- Utilitário pg_dump (para backups, geralmente incluído com PostgreSQL)

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

### Migrações de Banco de Dados

As migrações de banco de dados estão localizadas na pasta `migrations`. Elas são aplicadas automaticamente quando a aplicação backend é iniciada através do `docker-compose up` ou ao executar o `main.go` diretamente. A ferramenta `golang-migrate/migrate` é utilizada para este propósito.

### Executando com Docker (Recomendado)

O modo mais simples de executar o projeto é utilizando Docker Compose:

```sh
docker-compose up --build
```

O `--build` é recomendado na primeira vez ou após alterações nos Dockerfiles ou dependências.
Isso iniciará o banco de dados PostgreSQL, aplicará as migrações automaticamente e, em seguida, iniciará a API.

### Executando Localmente

1. Certifique-se de ter o PostgreSQL instalado e em execução
2. Instale as dependências:
   ```sh
   go mod download
   ```
3. Execute o servidor (as migrações serão aplicadas automaticamente):
   ```sh
   go run cmd/server/main.go
   ```

O servidor estará disponível em `http://localhost:3000`.

## Funcionalidades Principais

### Autenticação de Usuários

- Sistema de login seguro utilizando JWT
- Usuário admin criado automaticamente na inicialização
- Verificação de tokens para rotas protegidas

### Gerenciamento de Produtos e Lotes

- CRUD completo para produtos (criar, ler, atualizar, deletar).
- Produtos incluem ID, nome, unidade (L ou kg) e quantidade.
- Cada produto pode ser composto por múltiplos lotes.
- Cada lote possui ID, ID do produto, quantidade, data de validade.
- A quantidade total de um produto é automaticamente calculada como a soma das quantidades de seus lotes ativos (via gatilho no banco de dados).

### Histórico de Alterações

- Registro de todas as modificações em produtos e lotes.
- Armazenamento de alterações em formato JSON para flexibilidade.
- `EntityType` e `EntityID` nos registros de histórico indicam a qual entidade (produto ou lote) a alteração se refere.

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

- `GET /api/products`: Lista todos os produtos (incluindo seus lotes).
- `GET /api/products/:id`: Obtém um produto específico pelo ID (incluindo seus lotes).
- `POST /api/products`: Cria um novo produto (requer autenticação).
- `PUT /api/products/:id`: Atualiza um produto existente (requer autenticação).
- `DELETE /api/products/:id`: Remove um produto (e seus lotes associados) (requer autenticação).

### Lotes de Produtos

- `POST /api/products/:product_id/lotes`: Cria um novo lote para um produto específico (requer autenticação).
- `GET /api/products/:product_id/lotes`: Lista todos os lotes de um produto específico (requer autenticação).
- `PUT /api/lotes/:lote_id`: Atualiza um lote específico (requer autenticação).
- `DELETE /api/lotes/:lote_id`: Remove um lote específico (requer autenticação).

### Histórico

- `GET /api/history`: Lista todos os registros de histórico de alterações (requer autenticação).
  - Suporta query params `limit` e `offset` para paginação.
- `POST /api/history`: Adiciona um novo registro ao histórico (requer autenticação, geralmente usado internamente pelos serviços).
- `GET /api/history/:entity_type/:entity_id`: Lista registros de histórico para uma entidade específica (e.g., `/api/history/product/123` ou `/api/history/lote/abc`) (requer autenticação).
- `GET /api/history?batch_id={id}`: Retorna todos os registros de histórico associados a um `BatchID` específico.
- `GET /api/history/batch/{batch_id}`: Similar ao anterior, focado em buscar um lote específico.
- `GET /api/history/grouped`: **Novo endpoint** que retorna o histórico agrupado por `BatchID`. Cada grupo contém o `BatchID`, a data/hora da primeira entrada do lote, e todos os registros de histórico pertencentes àquele lote. Suporta paginação baseada nos lotes (batches).

## CORS

O servidor está configurado para aceitar requisições de qualquer origem (CORS habilitado), facilitando a integração com frontends em diferentes domínios durante o desenvolvimento.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.
