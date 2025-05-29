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

- **Visualizar Produtos:** Exibe uma tabela com todos os produtos em estoque, incluindo nome, quantidade total (somada dos lotes, se existentes, ou base) e unidade.
- **Adicionar Novo Produto:**
  - Formulário para inserir nome, unidade (L ou kg) e quantidade inicial (se o produto não for gerenciado por lotes).
  - Ao adicionar, um registro é criado no histórico.
- **Gerenciamento de Lotes (para produtos aplicáveis):**
  - Expandir um produto na tabela para visualizar seus lotes.
  - **Adicionar Novo Lote:** Formulário para inserir quantidade e data de validade do lote.
  - **Editar Lote:** Modificar quantidade e data de validade de um lote existente.
  - **Remover Lote:** Excluir um lote.
  - Todas as operações de lote são registradas no histórico. A quantidade total do produto é automaticamente atualizada no backend com base na soma dos seus lotes.
- **Atualizar Detalhes e Quantidades:**
  - Modo de edição permite alterar nome, unidade dos produtos.
  - Para produtos _sem lotes_, permite alterar a quantidade diretamente.
  - Botões de atalho para incrementar/decrementar quantidades (+1, -1, +10, -10) para produtos sem lotes.
  - As alterações são salvas, gerando entradas no histórico. Em modo autenticado, múltiplas alterações podem ser agrupadas em um "batch" de histórico.
- **Remover Produto:**
  - Permite remover um produto do estoque.
  - Uma confirmação é solicitada antes da remoção.
  - A remoção é registrada no histórico.
- **Notificações:** Feedback visual (toasts) para ações como login, adição, atualização e remoção de produtos, e outros eventos importantes.

### Histórico de Alterações

- **Visualizar Histórico:** Exibe um registro de todas as movimentações de estoque (entradas, saídas, adições e remoções de produtos, criação/alteração/remoção de lotes, e atualizações de detalhes de produtos).
- **Detalhes da Alteração:** Para cada entrada no histórico, mostra:
  - Data e hora da alteração.
  - Identificador do "batch" da operação (agrupa múltiplas alterações feitas de uma vez).
  - Nome do produto e/ou ID do lote.
  - Tipo de ação (ex: `product_created`, `lote_updated`, `quantity_changed`).
  - Detalhes específicos da mudança (ex: quantidade antes/depois, campos alterados com valores antigos/novos).
- **Filtrar Histórico:** Permite filtrar as movimentações por período (ex: Hoje, Esta semana, Este mês).
- **Modo Autenticado:** O frontend envia as modificações de dados (criação/atualização/deleção de produtos/lotes) para o backend. Após a confirmação dessas operações, pode enviar um conjunto de entradas de histórico (`HistoryBatchInput`) para o endpoint `/api/history/batch` do backend, que agrupa essas entradas sob um `batchId` único. O backend também pode gerar registros de histórico individuais para cada operação de API.

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
