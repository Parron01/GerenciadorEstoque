<script setup lang="ts">
import { useHistoryStore } from '@/stores/historyStore'
import { computed, ref, watch } from 'vue'

// Define props before using them
const props = defineProps<{
  filterOption: string
}>()

const historyStore = useHistoryStore()

// Paginação
const currentPage = ref(1)
const itemsPerPage = ref(6)

// Helper functions for date filtering
const isToday = (dateString: string): boolean => {
  const today = new Date()
  const date = new Date(dateString)
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  )
}

const isThisWeek = (dateString: string): boolean => {
  const date = new Date(dateString)
  const today = new Date()

  // Get the first day of the week (Sunday)
  const firstDay = new Date(today.getTime())
  const day = today.getDay()
  firstDay.setDate(today.getDate() - day)

  // Get the last day of the week (Saturday) - FIXED to use timestamp
  const lastDay = new Date(firstDay.getTime())
  lastDay.setDate(firstDay.getDate() + 6)

  // Check if the date is between first and last day of the week
  return date >= firstDay && date <= lastDay
}

const isThisMonth = (dateString: string): boolean => {
  const date = new Date(dateString)
  const today = new Date()
  return date.getMonth() === today.getMonth() && date.getFullYear() === today.getFullYear()
}

// Filtra o histórico com base na opção selecionada
const filteredHistory = computed(() => {
  if (!historyStore.history) return []

  return historyStore.history.filter((h) => {
    switch (props.filterOption) {
      case 'today':
        return isToday(h.date)
      case 'week':
        return isThisWeek(h.date)
      case 'month':
        return isThisMonth(h.date)
      default:
        return true // 'all' or any other value shows everything
    }
  })
})

// Total de páginas baseado no histórico filtrado
const totalPages = computed(() => {
  return Math.max(1, Math.ceil(filteredHistory.value.length / itemsPerPage.value))
})

// Registros a serem exibidos na página atual
const paginatedHistory = computed(() => {
  // Garantir que a página atual é válida
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value
  }

  const startIndex = (currentPage.value - 1) * itemsPerPage.value
  return filteredHistory.value.slice(startIndex, startIndex + itemsPerPage.value)
})

// Indicadores de paginação
const paginationInfo = computed(() => {
  const total = filteredHistory.value.length

  if (total === 0) {
    return { showing: false, start: 0, end: 0, total: 0 }
  }

  const start = Math.min((currentPage.value - 1) * itemsPerPage.value + 1, total)
  const end = Math.min(currentPage.value * itemsPerPage.value, total)

  return {
    showing: true,
    start,
    end,
    total,
  }
})

// Navegar para página anterior
function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

// Navegar para próxima página
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// Navegar para uma página específica
function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

// Gera array de números de página para exibição
const displayedPageNumbers = computed(() => {
  const total = totalPages.value
  const current = currentPage.value

  if (total <= 5) {
    // Se tivermos 5 ou menos páginas, exibir todas
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  // Estratégia para exibir páginas ao redor da página atual
  let pages = [1] // Sempre incluir a primeira página

  // Determinar o intervalo de páginas a serem exibidas
  let rangeStart = Math.max(2, current - 1)
  let rangeEnd = Math.min(total - 1, current + 1)

  // Ajustar para exibir sempre 3 números consecutivos (quando possível)
  if (current <= 2) {
    rangeEnd = Math.min(4, total - 1)
  } else if (current >= total - 1) {
    rangeStart = Math.max(2, total - 3)
  }

  // Adicionar elipses antes do intervalo, se necessário
  if (rangeStart > 2) {
    pages.push(-1) // -1 representa elipses
  }

  // Adicionar páginas do intervalo
  for (let i = rangeStart; i <= rangeEnd; i++) {
    pages.push(i)
  }

  // Adicionar elipses depois do intervalo, se necessário
  if (rangeEnd < total - 1) {
    pages.push(-2) // -2 representa elipses (usamos valor diferente para key única)
  }

  // Sempre incluir a última página
  if (total > 1) {
    pages.push(total)
  }

  return pages
})

// Formatar data para exibição responsiva
function formatDate(dateStr: string): { full: string; short: string } {
  const date = new Date(dateStr)
  return {
    full: date.toLocaleString(),
    short: date.toLocaleDateString() + '\n' + date.toLocaleTimeString().substring(0, 5),
  }
}

// Reset para primeira página quando o filtro muda
watch(
  () => props.filterOption,
  () => {
    currentPage.value = 1
  },
)

// Controle de visualização mobile
const isMobileView = ref(window.innerWidth < 640)

// Adiciona listener para redimensionamento da janela
window.addEventListener('resize', () => {
  isMobileView.value = window.innerWidth < 640
})
</script>

<template>
  <div>
    <!-- Tabela responsiva para telas médias e grandes -->
    <div class="hidden sm:block overflow-x-auto">
      <table class="min-w-full bg-white rounded shadow">
        <thead class="bg-slate-700 text-white text-left">
          <tr>
            <th class="p-3 w-1/4">Data/Hora</th>
            <th class="p-3">Alterações</th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="h in paginatedHistory"
            :key="h.id"
            class="border-b border-gray-200 hover:bg-gray-50"
          >
            <td class="p-3 align-top">{{ formatDate(h.date).full }}</td>
            <td class="p-3">
              <div class="grid grid-cols-1 gap-2">
                <div
                  v-for="(change, index) in h.changes"
                  :key="index"
                  class="flex flex-wrap items-center border-l-2 pl-3 mb-1 relative"
                  :class="change.action === 'add' ? 'border-green-500' : 'border-red-500'"
                >
                  <!-- Display product name from the stored history -->
                  <span class="font-medium mr-2">{{ change.productName }}</span>
                  <span
                    :class="
                      change.action === 'add'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    "
                    class="px-2 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ change.action === 'add' ? '+' : '-' }}{{ change.quantityChanged }}
                  </span>
                  <span class="ml-2 text-gray-600 text-sm flex items-center">
                    <span class="text-gray-400 mx-1">De:</span>
                    <span class="font-medium">{{ change.quantityBefore }}</span>
                    <span class="mx-1">→</span>
                    <span class="text-gray-400 mr-1">Para:</span>
                    <span class="font-medium">{{ change.quantityAfter }}</span>
                  </span>

                  <!-- Tags for new product or product removal -->
                  <span
                    v-if="change.isNewProduct"
                    class="absolute top-0 right-0 bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded"
                  >
                    Novo produto
                  </span>

                  <span
                    v-if="change.isProductRemoval"
                    class="absolute top-0 right-0 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded"
                  >
                    Remoção de produto
                  </span>
                </div>
              </div>
            </td>
          </tr>

          <!-- Linha de "sem dados" quando a tabela estiver vazia -->
          <tr v-if="paginatedHistory.length === 0">
            <td colspan="2" class="p-8 text-center text-gray-500">
              Nenhuma alteração encontrada no período selecionado.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Layout de cards para telas pequenas -->
    <div class="sm:hidden space-y-4">
      <div
        v-for="h in paginatedHistory"
        :key="h.id"
        class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden"
      >
        <div class="bg-slate-700 text-white px-3 py-2 text-sm font-medium">
          {{ formatDate(h.date).short }}
        </div>
        <div class="p-3">
          <div
            v-for="(change, index) in h.changes"
            :key="index"
            class="pb-3 mb-3 border-b border-gray-200 last:border-0 last:mb-0 last:pb-0 relative"
          >
            <div class="flex justify-between items-start mb-2">
              <span class="font-medium">{{ change.productName }}</span>
              <span
                :class="
                  change.action === 'add'
                    ? 'bg-green-100 text-green-800'
                    : 'bg-red-100 text-red-800'
                "
                class="px-2 py-0.5 rounded-full text-xs font-medium"
              >
                {{ change.action === 'add' ? '+' : '-' }}{{ change.quantityChanged }}
              </span>
            </div>

            <div class="text-sm text-gray-600 flex flex-wrap items-center">
              <span class="text-gray-400 mr-1">De:</span>
              <span class="font-medium mr-2">{{ change.quantityBefore }}</span>
              <span class="mr-2">→</span>
              <span class="text-gray-400 mr-1">Para:</span>
              <span class="font-medium">{{ change.quantityAfter }}</span>
            </div>

            <!-- Tags para mobile -->
            <div class="mt-2">
              <span
                v-if="change.isNewProduct"
                class="bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded mr-2"
              >
                Novo produto
              </span>
              <span
                v-if="change.isProductRemoval"
                class="bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded"
              >
                Remoção
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Mensagem quando não houver dados para mobile -->
      <div
        v-if="paginatedHistory.length === 0"
        class="bg-white rounded-lg text-center p-6 shadow border border-gray-200"
      >
        <div class="text-gray-500">Nenhuma alteração encontrada no período selecionado.</div>
      </div>
    </div>

    <!-- Controles de paginação modernizados -->
    <div
      v-if="filteredHistory.length > 0"
      class="mt-6 flex flex-col sm:flex-row justify-between items-center gap-4 bg-white p-3 rounded-lg shadow"
    >
      <!-- Informações sobre registros - apenas em telas maiores -->
      <div
        class="text-sm text-gray-500 sm:text-left text-center w-full sm:w-auto"
        v-if="paginationInfo.showing"
      >
        <span class="hidden sm:inline">Exibindo </span>
        <span class="font-medium">{{ paginationInfo.start }}-{{ paginationInfo.end }}</span>
        <span class="hidden sm:inline">de</span>
        <span class="sm:hidden">/</span>
        <span class="font-medium">{{ paginationInfo.total }}</span>
        <span class="hidden sm:inline"> registros</span>
      </div>

      <!-- Navegação de páginas modernizada -->
      <div class="flex justify-center items-center space-x-1.5">
        <!-- Botão anterior -->
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          :class="[
            currentPage === 1
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-white text-gray-600 hover:bg-indigo-50 hover:text-indigo-600',
          ]"
          class="flex items-center justify-center rounded-md border border-gray-200 w-10 h-10 transition-colors"
          aria-label="Página anterior"
        >
          <span class="material-icons-outlined text-lg">chevron_left</span>
        </button>

        <!-- Números de páginas para desktop -->
        <div class="hidden sm:flex space-x-1.5">
          <button
            v-for="pageNum in displayedPageNumbers"
            :key="pageNum"
            @click="pageNum > 0 ? goToPage(pageNum) : null"
            :class="[
              pageNum < 0
                ? 'bg-white border-transparent cursor-default'
                : pageNum === currentPage
                  ? 'bg-indigo-600 text-white border-indigo-600'
                  : 'bg-white text-gray-600 border-gray-200 hover:bg-indigo-50 hover:text-indigo-600',
            ]"
            class="flex items-center justify-center rounded-md border w-10 h-10 text-sm transition-colors"
            :disabled="pageNum < 0"
          >
            {{ pageNum > 0 ? pageNum : '...' }}
          </button>
        </div>

        <!-- Contador de página para mobile (atualizado para melhor visualização) -->
        <div
          class="sm:hidden flex items-center justify-center min-w-16 px-2 py-1 bg-indigo-50 text-indigo-700 rounded-md border border-indigo-100"
        >
          <span class="font-medium text-sm">{{ currentPage }} / {{ totalPages }}</span>
        </div>

        <!-- Botão próximo -->
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          :class="[
            currentPage >= totalPages
              ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
              : 'bg-white text-gray-600 hover:bg-indigo-50 hover:text-indigo-600',
          ]"
          class="flex items-center justify-center rounded-md border border-gray-200 w-10 h-10 transition-colors"
          aria-label="Próxima página"
        >
          <span class="material-icons-outlined text-lg">chevron_right</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Estilos específicos da paginação */
button:focus {
  outline: 2px solid rgba(79, 70, 229, 0.3);
  outline-offset: 2px;
}

/* Melhorar acessibilidade para toque em dispositivos móveis */
@media (max-width: 640px) {
  button {
    min-height: 40px;
    min-width: 40px;
  }
}
</style>
