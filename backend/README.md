# Projeto Estoque Simples - Backend

Este é o backend para o sistema de Gerenciamento de Estoque Simples, responsável pela lógica de negócios, persistência de dados e autenticação.

## Tecnologias Utilizadas

- **Node.js:** Ambiente de execução JavaScript.
- **Express.js:** Framework web para Node.js, usado para criar a API RESTful.
- **TypeScript:** Superset do JavaScript que adiciona tipagem estática.
- **SQLite (better-sqlite3):** Banco de dados relacional leve para armazenamento de dados.
- **JSON Web Tokens (JWT):** Para autenticação e autorização baseada em tokens.
- **bcrypt:** Para hashing seguro de senhas.
- **node-cron:** Para agendamento de tarefas (ex: backups).
- **cors:** Para habilitar Cross-Origin Resource Sharing.
- **dotenv:** Para gerenciamento de variáveis de ambiente.

## Funcionalidades Principais

- **API RESTful:**
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
    - `POST /api/history`: Cria um novo registro de histórico (requer autenticação, usado internamente ao modificar produtos).
  - **Health Check:**
    - `GET /api/auth/health`: Verifica o status do servidor.
- **Gerenciamento de Banco de Dados:**
  - Inicialização automática das tabelas (`users`, `products`, `history`) se não existirem.
  - Armazenamento dos dados em um arquivo SQLite (`src/database/inventory.sqlite`).
- **Segurança:**
  - Hashing de senhas com bcrypt.
  - Autenticação baseada em JWT para rotas protegidas.
- **Backups Automatizados:**
  - Agendamento semanal (Domingo às 03:00) para criar backups do banco de dados na pasta `src/backups/`.
- **Servir Arquivos Estáticos:**
  - Em modo de produção, pode servir os arquivos estáticos do frontend.

## Configuração e Execução

### Pré-requisitos

- Node.js (versão 22 ou superior recomendada)
- npm (geralmente vem com o Node.js)

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do diretório `backend` com as seguintes variáveis:

```env
PORT=3000
JWT_SECRET=sua_chave_secreta_jwt_aqui
# JWT_EXPIRATION=7d (opcional, padrão 7 dias)

ADMIN_USERNAME=seuusuario
ADMIN_PASSWORD=suasenha
```

### Instalação

1.  Navegue até o diretório `backend`:
    ```sh
    cd backend
    ```
2.  Instale as dependências:
    ```sh
    npm install
    ```

### Executando em Desenvolvimento

```sh
npm run dev
```

O servidor iniciará em `http://localhost:3000` (ou a porta definida em `.env`) com hot-reload via Nodemon e `tsx`.

## Scripts Disponíveis

- `npm run dev`: Inicia o servidor de desenvolvimento.
- `npm run build`: Compila o código TypeScript para JavaScript.
- `npm run start`: Inicia o servidor em modo de produção a partir dos arquivos compilados.
