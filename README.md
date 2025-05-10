# Projeto Estoque Simples - Frontend

Este é o frontend para o sistema de Gerenciamento de Estoque Simples, desenvolvido para facilitar o controle de entradas e saídas de produtos.

## Tecnologias Utilizadas

- **Vue 3:** Framework progressivo para construção de interfaces de usuário.
- **Vite:** Ferramenta de build moderna e rápida para desenvolvimento frontend.
- **Pinia:** Gerenciador de estado intuitivo para Vue.js.
- **Vue Router:** Biblioteca de roteamento oficial para Vue.js.
- **Tailwind CSS:** Framework CSS utility-first para estilização rápida.
- **TypeScript:** Superset do JavaScript que adiciona tipagem estática.
- **Vue Toastification:** Para exibir notificações (toasts) de feedback ao usuário.
- **UUID:** Para geração de identificadores únicos.

## Funcionalidades

O sistema de frontend permite as seguintes operações:

### Gerenciamento de Estoque

- **Visualizar Produtos:** Exibe uma tabela com todos os produtos em estoque, incluindo nome, quantidade e unidade.
- **Adicionar Novo Produto:**
  - Formulário para inserir nome, unidade (L ou kg) e quantidade inicial do novo produto.
  - Ao adicionar, um registro é criado no histórico.
- **Atualizar Quantidades:**
  - Modo de edição permite alterar a quantidade dos produtos existentes.
  - Botões de atalho para incrementar/decrementar quantidades (+1, -1, +10, -10).
  - As alterações são salvas em lote, gerando uma única entrada no histórico com todas as modificações.
- **Remover Produto:**
  - Permite remover um produto do estoque.
  - Uma confirmação é solicitada antes da remoção.
  - A remoção é registrada no histórico.
- **Notificações:** Feedback visual (toasts) para ações como adição, atualização e remoção de produtos.

### Histórico de Alterações

- **Visualizar Histórico:** Exibe um registro de todas as movimentações de estoque (entradas, saídas, adições e remoções de produtos).
- **Detalhes da Alteração:** Para cada entrada no histórico, mostra:
  - Data e hora da alteração.
  - Nome do produto.
  - Tipo de ação (adição/remoção de quantidade).
  - Quantidade alterada.
  - Quantidade antes e depois da alteração.
  - Tags indicando se foi um "Novo produto" ou "Remoção de produto".
- **Filtrar Histórico:** Permite filtrar as movimentações por período:
  - Hoje
  - Esta semana
  - Este mês
  - Todo o período

### Persistência de Dados

- Os dados de produtos e o histórico de alterações são salvos localmente no navegador utilizando `localStorage`, garantindo que as informações persistam entre as sessões.

## Como Executar (Frontend)

1.  **Instalar dependências:**
    ```sh
    npm install
    ```
2.  **Iniciar o servidor de desenvolvimento:**
    ```sh
    npm run dev
    ```
    O frontend estará acessível em `http://localhost:5173` (ou a porta configurada em [`vite.config.ts`](vite.config.ts)).

## Scripts Disponíveis

- `npm run dev`: Inicia o servidor de desenvolvimento com hot-reload.
- `npm run build`: Compila o projeto para produção (inclui type-checking).
- `npm run preview`: Pré-visualiza a build de produção localmente.
- `npm run lint`: Executa o ESLint para verificar e corrigir problemas de código.
- `npm run format`: Formata o código usando Prettier.
