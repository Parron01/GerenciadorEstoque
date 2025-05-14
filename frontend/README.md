# Projeto Estoque Simples - Frontend

Este é o frontend para o sistema de Gerenciamento de Estoque Simples, desenvolvido para facilitar o controle de entradas e saídas de produtos, com uma interface de usuário interativa e responsiva.

## Tecnologias Utilizadas

- **Vue 3:** Framework progressivo para construção de interfaces de usuário.
- **Vite:** Ferramenta de build moderna e rápida para desenvolvimento frontend.
- **Pinia:** Gerenciador de estado intuitivo para Vue.js, usado para gerenciar produtos, histórico e estado de autenticação.
- **Vue Router:** Biblioteca de roteamento oficial para Vue.js, para navegação entre as páginas de Login, Estoque e Histórico.
- **Tailwind CSS:** Framework CSS utility-first para estilização rápida e customizável.
- **TypeScript:** Superset do JavaScript que adiciona tipagem estática.
- **Vue Toastification:** Para exibir notificações (toasts) de feedback ao usuário.
- **Axios (ou Fetch API):** Para realizar chamadas HTTP para o backend.
- **UUID:** Para geração de identificadores únicos no modo local.

## Funcionalidades

O sistema de frontend permite as seguintes operações:

### Autenticação de Usuário

- **Página de Login:** Permite que os usuários se autentiquem para acessar as funcionalidades do sistema.
- **Modo Local (Demonstração):**
  - Permite o uso da aplicação sem autenticação no backend, utilizando dados de demonstração e `localStorage` para persistência.
  - Útil para testes rápidos ou uso offline simplificado.
- **Modo Autenticado:**
  - Interage com o backend para carregar e salvar dados reais.
  - Requer login para acesso.

### Gerenciamento de Estoque

- **Visualizar Produtos:** Exibe uma tabela com todos os produtos em estoque, incluindo nome, quantidade e unidade.
- **Adicionar Novo Produto:**
  - Formulário para inserir nome, unidade (L ou kg) e quantidade inicial do novo produto.
  - Ao adicionar, um registro é criado no histórico.
- **Atualizar Quantidades:**
  - Modo de edição permite alterar a quantidade dos produtos existentes.
  - Botões de atalho para incrementar/decrementar quantidades (+1, -1, +10, -10).
  - As alterações são salvas, gerando entradas no histórico.
- **Remover Produto:**
  - Permite remover um produto do estoque.
  - Uma confirmação é solicitada antes da remoção.
  - A remoção é registrada no histórico.
- **Notificações:** Feedback visual (toasts) para ações como login, adição, atualização e remoção de produtos, e outros eventos importantes.

### Histórico de Alterações

- **Visualizar Histórico:** Exibe um registro de todas as movimentações de estoque (entradas, saídas, adições e remoções de produtos).
- **Detalhes da Alteração:** Para cada entrada no histórico, mostra:
  - Data e hora da alteração.
  - Nome do produto.
  - Tipo de ação (adição/remoção de quantidade, criação/remoção de produto).
  - Quantidade alterada.
  - Quantidade antes e depois da alteração.
- **Filtrar Histórico:** Permite filtrar as movimentações por período (ex: Hoje, Esta semana, Este mês).

### Persistência de Dados

- **Modo Local:** Os dados de produtos e o histórico de alterações são salvos localmente no navegador utilizando `localStorage`, garantindo que as informações persistam entre as sessões. Chaves de storage diferentes são usadas para modo local e modo autenticado.
- **Modo Autenticado:** Os dados são buscados e enviados para o backend, com `localStorage` podendo ser usado para cache ou preferências do usuário.

## Como Executar (Frontend)

### Pré-requisitos

- Node.js (versão slim ou superior recomendada para Docker, qualquer versão estável para desenvolvimento local)
- npm

### Instalação

1.  Navegue até o diretório `frontend`:
    ```sh
    cd frontend
    ```
2.  Instale as dependências:
    ```sh
    npm install
    ```

### Iniciar o Servidor de Desenvolvimento

```sh
npm run dev
```

O frontend estará acessível em `http://localhost:5173` (ou a porta configurada em [`vite.config.ts`](vite.config.ts)).

## Scripts Disponíveis

- `npm run dev`: Inicia o servidor de desenvolvimento com hot-reload.
- `npm run build`: Compila o projeto para produção (inclui type-checking). Os arquivos são gerados no diretório `dist/`.
