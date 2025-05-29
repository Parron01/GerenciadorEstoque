<script setup lang="ts">
import { useHistoryStore } from "@/stores/historyStore";
import { computed, ref, watch } from "vue";
import type { ProductHistory, ProductChange } from "@/models/product"; // For local mode
import type { ParsedHistoryRecord } from "@/models/history"; // For auth mode
import type { LoteChangeDetails } from "@/models/lote";
import { useAuthStore } from "@/stores/authStore";

// Define props before using them
const props = defineProps<{
  filterOption: string;
}>();

const historyStore = useHistoryStore();
const authStore = useAuthStore();

const currentPage = ref(1);
const itemsPerPage = ref(6);

// Helper functions for date filtering
const getDateFromRecord = (
  record: ProductHistory | ParsedHistoryRecord
): string => {
  return "createdAt" in record ? record.createdAt : record.date;
};

const isToday = (dateString: string): boolean => {
  const today = new Date();
  const date = new Date(dateString);
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
};

const isThisWeek = (dateString: string): boolean => {
  const date = new Date(dateString);
  const today = new Date();
  const firstDay = new Date(today.setDate(today.getDate() - today.getDay()));
  const lastDay = new Date(today.setDate(today.getDate() - today.getDay() + 6));
  firstDay.setHours(0, 0, 0, 0); // Start of the first day
  lastDay.setHours(23, 59, 59, 999); // End of the last day
  return date >= firstDay && date <= lastDay;
};

const isThisMonth = (dateString: string): boolean => {
  const date = new Date(dateString);
  const today = new Date();
  return (
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  );
};

// Filtra o histórico com base na opção selecionada
const filteredHistory = computed(() => {
  if (!historyStore.history) return [];

  return historyStore.history.filter((h_record) => {
    const recordDate = getDateFromRecord(h_record);
    switch (props.filterOption) {
      case "today":
        return isToday(recordDate);
      case "week":
        return isThisWeek(recordDate);
      case "month":
        return isThisMonth(recordDate);
      default:
        return true;
    }
  });
});

// Total de páginas baseado no histórico filtrado
const totalPages = computed(() => {
  return Math.max(
    1,
    Math.ceil(filteredHistory.value.length / itemsPerPage.value)
  );
});

// Registros a serem exibidos na página atual
const paginatedHistory = computed(() => {
  if (currentPage.value > totalPages.value && totalPages.value > 0) {
    // check totalPages > 0
    currentPage.value = totalPages.value;
  } else if (totalPages.value === 0) {
    currentPage.value = 1; // Reset to 1 if no pages
  }
  const startIndex = (currentPage.value - 1) * itemsPerPage.value;
  return filteredHistory.value.slice(
    startIndex,
    startIndex + itemsPerPage.value
  );
});

// Indicadores de paginação
const paginationInfo = computed(() => {
  const total = filteredHistory.value.length;
  if (total === 0) return { showing: false, start: 0, end: 0, total: 0 };
  const start = Math.min(
    (currentPage.value - 1) * itemsPerPage.value + 1,
    total
  );
  const end = Math.min(currentPage.value * itemsPerPage.value, total);
  return { showing: true, start, end, total };
});

// Navegar para página anterior
function prevPage() {
  if (currentPage.value > 1) currentPage.value--;
}

// Navegar para próxima página
function nextPage() {
  if (currentPage.value < totalPages.value) currentPage.value++;
}

// Navegar para uma página específica
function goToPage(page: number) {
  if (page >= 1 && page <= totalPages.value) currentPage.value = page;
}

// Gera array de números de página para exibição
const displayedPageNumbers = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1);
  let pages = [1];
  let rangeStart = Math.max(2, current - 1);
  let rangeEnd = Math.min(total - 1, current + 1);
  if (current <= 2) rangeEnd = Math.min(4, total - 1);
  else if (current >= total - 1) rangeStart = Math.max(2, total - 3);
  if (rangeStart > 2) pages.push(-1); // Ellipsis
  for (let i = rangeStart; i <= rangeEnd; i++) pages.push(i);
  if (rangeEnd < total - 1) pages.push(-2); // Ellipsis
  if (total > 1) pages.push(total);
  return pages;
});

// Reset para primeira página quando o filtro muda
watch(
  () => props.filterOption,
  () => {
    currentPage.value = 1;
  }
);

// Controle de visualização mobile
const isMobileView = ref(window.innerWidth < 640);
window.addEventListener("resize", () => {
  isMobileView.value = window.innerWidth < 640;
});

function formatDateForDisplay(dateStr: string): {
  full: string;
  short: string;
} {
  const date = new Date(dateStr);
  return {
    full: date.toLocaleString(),
    short:
      date.toLocaleDateString() +
      "\n" +
      date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
  };
}

// Type guard for ParsedHistoryRecord
function isParsedHistoryRecord(record: any): record is ParsedHistoryRecord {
  return (
    record &&
    typeof record.entityType === "string" &&
    typeof record.entityId === "string"
  );
}

// Type guard for ProductHistory (local mode)
function isProductHistory(record: any): record is ProductHistory {
  return (
    record && Array.isArray(record.changes) && typeof record.date === "string"
  );
}
</script>

<template>
  <div>
    <!-- Tabela responsiva para telas médias e grandes -->
    <div class="hidden sm:block overflow-x-auto">
      <table class="min-w-full bg-white rounded shadow">
        <thead class="bg-slate-700 text-white text-left">
          <tr>
            <th class="p-3 w-1/4">Data/Hora</th>
            <th class="p-3">Detalhes da Alteração</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="record in paginatedHistory"
            :key="record.id"
            class="border-b border-gray-200 hover:bg-gray-50"
          >
            <td class="p-3 align-top">
              {{ formatDateForDisplay(getDateFromRecord(record)).full }}
            </td>
            <td class="p-3">
              <!-- Authenticated Mode: ParsedHistoryRecord -->
              <div v-if="isParsedHistoryRecord(record)" class="space-y-2">
                <div
                  v-if="record.entityType === 'product'"
                  class="border-l-2 pl-3 border-blue-500"
                >
                  <span class="font-medium"
                    >Produto:
                    {{
                      (record.details as ProductChange).productName ||
                      record.entityId
                    }}</span
                  >
                  <span
                    :class="
                      (record.details as ProductChange).action === 'add' ||
                      (record.details as ProductChange).action === 'created'
                        ? 'text-green-700'
                        : 'text-red-700'
                    "
                    class="ml-2 text-sm"
                  >
                    Ação: {{ (record.details as ProductChange).action }}
                  </span>
                  <div class="text-xs text-gray-600">
                    Qtd:
                    {{ (record.details as ProductChange).quantityBefore }} →
                    {{ (record.details as ProductChange).quantityAfter }} ({{
                      (record.details as ProductChange).action === "add" ||
                      (record.details as ProductChange).action === "created"
                        ? "+"
                        : ""
                    }}{{ (record.details as ProductChange).quantityChanged }})
                  </div>
                  <span
                    v-if="(record.details as ProductChange).isNewProduct"
                    class="tag-new"
                    >Novo Produto</span
                  >
                  <span
                    v-if="(record.details as ProductChange).isProductRemoval"
                    class="tag-removed"
                    >Produto Removido</span
                  >
                </div>
                <div
                  v-else-if="record.entityType === 'lote'"
                  class="border-l-2 pl-3 border-purple-500"
                >
                  <span class="font-medium"
                    >Lote: {{ record.entityId.substring(0, 8) }}</span
                  >
                  <span class="text-xs text-gray-500 ml-1"
                    >(Produto:
                    {{
                      record.productNameContext ||
                      (
                        record.details as LoteChangeDetails
                      ).productId?.substring(0, 8)
                    }})</span
                  >
                  <span
                    :class="
                      (record.details as LoteChangeDetails).action === 'created'
                        ? 'text-green-700'
                        : (record.details as LoteChangeDetails).action ===
                            'updated'
                          ? 'text-yellow-700'
                          : 'text-red-700'
                    "
                    class="ml-2 text-sm"
                  >
                    Ação: {{ (record.details as LoteChangeDetails).action }}
                  </span>
                  <div class="text-xs text-gray-600">
                    <div
                      v-if="
                        (record.details as LoteChangeDetails).quantityBefore !==
                        undefined
                      "
                    >
                      Qtd:
                      {{ (record.details as LoteChangeDetails).quantityBefore }}
                      →
                      {{ (record.details as LoteChangeDetails).quantityAfter }}
                    </div>
                    <div
                      v-if="
                        (record.details as LoteChangeDetails).dataValidadeOld ||
                        (record.details as LoteChangeDetails).dataValidadeNew
                      "
                    >
                      Validade:
                      {{
                        (record.details as LoteChangeDetails).dataValidadeOld ||
                        "N/A"
                      }}
                      →
                      {{
                        (record.details as LoteChangeDetails).dataValidadeNew ||
                        (record.details as LoteChangeDetails).dataValidade ||
                        "N/A"
                      }}
                    </div>
                    <div
                      v-else-if="
                        (record.details as LoteChangeDetails).dataValidade
                      "
                    >
                      Validade:
                      {{ (record.details as LoteChangeDetails).dataValidade }}
                    </div>
                  </div>
                </div>
              </div>
              <!-- Local Mode: ProductHistory (batch changes) -->
              <div
                v-else-if="isProductHistory(record)"
                class="grid grid-cols-1 gap-2"
              >
                <div
                  v-for="(change, index) in record.changes"
                  :key="index"
                  class="flex flex-wrap items-center border-l-2 pl-3 mb-1 relative"
                  :class="
                    change.action === 'add'
                      ? 'border-green-500'
                      : 'border-red-500'
                  "
                >
                  <span class="font-medium mr-2">{{ change.productName }}</span>
                  <span
                    :class="
                      change.action === 'add'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    "
                    class="px-2 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ change.action === "add" ? "+" : "-"
                    }}{{ change.quantityChanged }}
                  </span>
                  <span class="ml-2 text-gray-600 text-sm flex items-center">
                    <span class="text-gray-400 mx-1">De:</span>
                    <span class="font-medium">{{ change.quantityBefore }}</span>
                    <span class="mx-1">→</span>
                    <span class="text-gray-400 mr-1">Para:</span>
                    <span class="font-medium">{{ change.quantityAfter }}</span>
                  </span>
                  <span v-if="change.isNewProduct" class="tag-new"
                    >Novo produto</span
                  >
                  <span v-if="change.isProductRemoval" class="tag-removed"
                    >Remoção de produto</span
                  >
                </div>
              </div>
            </td>
          </tr>
          <tr v-if="paginatedHistory.length === 0">
            <td colspan="2" class="p-8 text-center text-gray-500">
              Nenhuma alteração encontrada.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Layout de cards para telas pequenas -->
    <div class="sm:hidden space-y-4">
      <div
        v-for="record in paginatedHistory"
        :key="record.id"
        class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden"
      >
        <div class="bg-slate-700 text-white px-3 py-2 text-sm font-medium">
          {{ formatDateForDisplay(getDateFromRecord(record)).short }}
        </div>
        <div class="p-3">
          <!-- Authenticated Mode: ParsedHistoryRecord -->
          <div v-if="isParsedHistoryRecord(record)" class="space-y-2">
            <div
              v-if="record.entityType === 'product'"
              class="pb-2 mb-2 border-b last:border-0 last:mb-0 last:pb-0"
            >
              <div class="flex justify-between items-start mb-1">
                <span class="font-medium"
                  >Produto:
                  {{
                    (record.details as ProductChange).productName ||
                    record.entityId
                  }}</span
                >
                <span
                  :class="
                    (record.details as ProductChange).action === 'add' ||
                    (record.details as ProductChange).action === 'created'
                      ? 'tag-green'
                      : 'tag-red'
                  "
                  class="text-xs"
                >
                  {{ (record.details as ProductChange).action }}
                </span>
              </div>
              <div class="text-sm text-gray-600">
                Qtd: {{ (record.details as ProductChange).quantityBefore }} →
                {{ (record.details as ProductChange).quantityAfter }}
              </div>
              <span
                v-if="(record.details as ProductChange).isNewProduct"
                class="tag-new-sm"
                >Novo</span
              >
              <span
                v-if="(record.details as ProductChange).isProductRemoval"
                class="tag-removed-sm"
                >Removido</span
              >
            </div>
            <div
              v-else-if="record.entityType === 'lote'"
              class="pb-2 mb-2 border-b last:border-0 last:mb-0 last:pb-0"
            >
              <div class="flex justify-between items-start mb-1">
                <span class="font-medium"
                  >Lote: {{ record.entityId.substring(0, 8) }}</span
                >
                <span
                  :class="
                    (record.details as LoteChangeDetails).action === 'created'
                      ? 'tag-green'
                      : (record.details as LoteChangeDetails).action ===
                          'updated'
                        ? 'tag-yellow'
                        : 'tag-red'
                  "
                  class="text-xs"
                >
                  {{ (record.details as LoteChangeDetails).action }}
                </span>
              </div>
              <div class="text-xs text-gray-500 mb-1">
                Prod:
                {{
                  record.productNameContext ||
                  (record.details as LoteChangeDetails).productId?.substring(
                    0,
                    8
                  )
                }}
              </div>
              <div class="text-sm text-gray-600">
                <div
                  v-if="
                    (record.details as LoteChangeDetails).quantityBefore !==
                    undefined
                  "
                >
                  Qtd:
                  {{ (record.details as LoteChangeDetails).quantityBefore }} →
                  {{ (record.details as LoteChangeDetails).quantityAfter }}
                </div>
                <div
                  v-if="
                    (record.details as LoteChangeDetails).dataValidadeOld ||
                    (record.details as LoteChangeDetails).dataValidadeNew
                  "
                >
                  Val:
                  {{
                    (record.details as LoteChangeDetails).dataValidadeOld ||
                    "N/A"
                  }}
                  →
                  {{
                    (record.details as LoteChangeDetails).dataValidadeNew ||
                    (record.details as LoteChangeDetails).dataValidade ||
                    "N/A"
                  }}
                </div>
                <div
                  v-else-if="(record.details as LoteChangeDetails).dataValidade"
                >
                  Val: {{ (record.details as LoteChangeDetails).dataValidade }}
                </div>
              </div>
            </div>
          </div>
          <!-- Local Mode: ProductHistory (batch changes) -->
          <div v-else-if="isProductHistory(record)">
            <div
              v-for="(change, index) in record.changes"
              :key="index"
              class="pb-3 mb-3 border-b border-gray-200 last:border-0 last:mb-0 last:pb-0 relative"
            >
              <div class="flex justify-between items-start mb-2">
                <span class="font-medium">{{ change.productName }}</span>
                <span
                  :class="change.action === 'add' ? 'tag-green' : 'tag-red'"
                  class="text-xs"
                >
                  {{ change.action === "add" ? "+" : "-"
                  }}{{ change.quantityChanged }}
                </span>
              </div>
              <div class="text-sm text-gray-600">
                De: {{ change.quantityBefore }} → Para:
                {{ change.quantityAfter }}
              </div>
              <span v-if="change.isNewProduct" class="tag-new-sm">Novo</span>
              <span v-if="change.isProductRemoval" class="tag-removed-sm"
                >Remoção</span
              >
            </div>
          </div>
        </div>
      </div>
      <div
        v-if="paginatedHistory.length === 0"
        class="bg-white rounded-lg text-center p-6 shadow border border-gray-200 text-gray-500"
      >
        Nenhuma alteração encontrada.
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
        <span class="font-medium"
          >{{ paginationInfo.start }}-{{ paginationInfo.end }}</span
        >
        <span class="hidden sm:inline"> de </span>
        <span class="sm:hidden">/</span>
        <span class="font-medium">{{ paginationInfo.total }}</span>
        <span class="hidden sm:inline"> registros</span>
      </div>
      <div class="flex justify-center items-center space-x-1.5">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          :class="[currentPage === 1 ? 'btn-page-disabled' : 'btn-page']"
          aria-label="Página anterior"
        >
          <span class="material-icons-outlined text-lg">chevron_left</span>
        </button>
        <div class="hidden sm:flex space-x-1.5">
          <button
            v-for="pageNum in displayedPageNumbers"
            :key="pageNum"
            @click="pageNum > 0 ? goToPage(pageNum) : null"
            :class="[
              pageNum < 0
                ? 'btn-page-ellipsis'
                : pageNum === currentPage
                  ? 'btn-page-current'
                  : 'btn-page',
            ]"
            :disabled="pageNum < 0"
          >
            {{ pageNum > 0 ? pageNum : "..." }}
          </button>
        </div>
        <div
          class="sm:hidden flex items-center justify-center min-w-16 px-2 py-1 bg-indigo-50 text-indigo-700 rounded-md border border-indigo-100"
        >
          <span class="font-medium text-sm"
            >{{ currentPage }} / {{ totalPages }}</span
          >
        </div>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          :class="[
            currentPage >= totalPages ? 'btn-page-disabled' : 'btn-page',
          ]"
          aria-label="Próxima página"
        >
          <span class="material-icons-outlined text-lg">chevron_right</span>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tag-new {
  @apply absolute top-0 right-0 bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded;
}
.tag-removed {
  @apply absolute top-0 right-0 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded;
}
.tag-green {
  @apply bg-green-100 text-green-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-red {
  @apply bg-red-100 text-red-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-yellow {
  @apply bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-new-sm {
  @apply bg-green-500 text-white text-xs font-bold px-1.5 py-0.5 rounded mr-1;
}
.tag-removed-sm {
  @apply bg-red-500 text-white text-xs font-bold px-1.5 py-0.5 rounded;
}

.btn-page {
  @apply flex items-center justify-center rounded-md border border-gray-200 w-10 h-10 text-sm bg-white text-gray-600 hover:bg-indigo-50 hover:text-indigo-600 transition-colors;
}
.btn-page-current {
  @apply flex items-center justify-center rounded-md border w-10 h-10 text-sm bg-indigo-600 text-white border-indigo-600 transition-colors;
}
.btn-page-disabled {
  @apply flex items-center justify-center rounded-md border border-gray-200 w-10 h-10 text-sm bg-gray-100 text-gray-400 cursor-not-allowed transition-colors;
}
.btn-page-ellipsis {
  @apply flex items-center justify-center rounded-md border w-10 h-10 text-sm bg-white border-transparent cursor-default;
}

button:focus {
  outline: 2px solid rgba(79, 70, 229, 0.3);
  outline-offset: 2px;
}
@media (max-width: 640px) {
  button {
    min-height: 40px;
    min-width: 40px;
  }
}
</style>
