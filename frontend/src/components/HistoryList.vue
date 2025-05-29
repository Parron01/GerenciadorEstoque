<script setup lang="ts">
import { useHistoryStore } from "@/stores/historyStore";
import { computed, ref, watch } from "vue";
import type {
  HistoryBatchGroup, // Corrected type name
  ParsedHistoryRecord, // Keep if used by local mode transformation
  ProductHistory, // For local mode
} from "@/models/history";
import { useAuthStore } from "@/stores/authStore";
import HistoryTable from "./history-list/HistoryTable.vue";
import HistoryCard from "./history-list/HistoryCard.vue";
import HistoryPagination from "./history-list/HistoryPagination.vue";

// Define props before using them
const props = defineProps<{
  filterOption: string;
}>();

const historyStore = useHistoryStore();
const authStore = useAuthStore();

// itemsPerPage is now managed by historyStore.pageSizeForGrouped for API mode
// For local mode, we can define a local itemsPerPage if needed, or use a fixed one.
const localModeItemsPerPage = ref(4);

// Helper functions for date filtering
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
const filteredAndSortedBatches = computed((): HistoryBatchGroup[] => {
  if (authStore.isLocalMode) {
    // Use the getter that transforms local history for consistency
    const localGroups = historyStore.getLocalModeHistoryAsGroups;
    return localGroups
      .filter((batch: HistoryBatchGroup) => {
        const batchDate = batch.createdAt;
        switch (props.filterOption) {
          case "today":
            return isToday(batchDate);
          case "week":
            return isThisWeek(batchDate);
          case "month":
            return isThisMonth(batchDate);
          default:
            return true;
        }
      })
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );
  } else {
    // Authenticated mode: historyStore.groupedHistory already contains paginated groups from API
    // Client-side filtering is applied on the *current page* of batches fetched from the API.
    if (!historyStore.groupedHistory || !historyStore.groupedHistory.groups)
      return [];
    return historyStore.groupedHistory.groups // Corrected: filter on groups array
      .filter((batch: HistoryBatchGroup) => {
        // Added type for batch
        if (!batch.records || batch.records.length === 0) return false;
        const batchDate = batch.createdAt;
        switch (props.filterOption) {
          case "today":
            return isToday(batchDate);
          case "week":
            return isThisWeek(batchDate);
          case "month":
            return isThisMonth(batchDate);
          default:
            return true;
        }
      })
      .sort(
        (a, b) =>
          new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
      );
  }
});

// Pagination for local mode (since API mode pagination is handled by the store)
const localModeCurrentPage = ref(1);

const localModeTotalPages = computed(() => {
  if (!authStore.isLocalMode) return 1;
  return Math.max(
    1,
    Math.ceil(
      filteredAndSortedBatches.value.length / localModeItemsPerPage.value
    )
  );
});

const paginatedBatchesToDisplay = computed(() => {
  if (authStore.isLocalMode) {
    if (
      localModeCurrentPage.value > localModeTotalPages.value &&
      localModeTotalPages.value > 0
    ) {
      localModeCurrentPage.value = localModeTotalPages.value;
    } else if (localModeTotalPages.value === 0) {
      localModeCurrentPage.value = 1;
    }
    const startIndex =
      (localModeCurrentPage.value - 1) * localModeItemsPerPage.value;
    return filteredAndSortedBatches.value.slice(
      startIndex,
      startIndex + localModeItemsPerPage.value
    );
  }
  // For API mode, filteredAndSortedBatches already represents the current page's filtered data
  return filteredAndSortedBatches.value;
});

const currentPage = computed(() =>
  authStore.isLocalMode
    ? localModeCurrentPage.value
    : historyStore.currentPageForGrouped
);
const totalPages = computed(() =>
  authStore.isLocalMode
    ? localModeTotalPages.value
    : historyStore.groupedHistory?.totalPages || 1
);

const paginationInfo = computed(() => {
  if (authStore.isLocalMode) {
    const total = filteredAndSortedBatches.value.length;
    if (total === 0) return { showing: false, start: 0, end: 0, total: 0 };
    const start = Math.min(
      (localModeCurrentPage.value - 1) * localModeItemsPerPage.value + 1,
      total
    );
    const end = Math.min(
      localModeCurrentPage.value * localModeItemsPerPage.value,
      total
    );
    return { showing: true, start, end, total };
  } else {
    // API Mode
    const totalItems = historyStore.groupedHistory?.totalBatches || 0;
    if (totalItems === 0) return { showing: false, start: 0, end: 0, total: 0 };

    // If client-side filtering is active on the API-paginated data:
    // The 'total' should ideally reflect the total *unfiltered* batches from the API for pagination controls.
    // The 'start' and 'end' can reflect the currently *viewable* items if filtering reduces the count on the page.
    // However, for simplicity with API pagination, we'll use the store's pagination info.
    const pageSize = historyStore.pageSizeForGrouped;
    const apiCurrentPage = historyStore.currentPageForGrouped;
    const apiTotalBatches = historyStore.groupedHistory?.totalBatches || 0;
    const start = (apiCurrentPage - 1) * pageSize + 1;
    const end = Math.min(apiCurrentPage * pageSize, apiTotalBatches);

    return {
      showing: apiTotalBatches > 0,
      start: Math.min(start, apiTotalBatches),
      end: end,
      total: apiTotalBatches,
    };
  }
});

function changePage(page: number) {
  if (authStore.isLocalMode) {
    if (page >= 1 && page <= localModeTotalPages.value) {
      localModeCurrentPage.value = page;
    }
  } else {
    historyStore.changeGroupedHistoryPage(page);
  }
}

watch(
  () => props.filterOption,
  () => {
    if (authStore.isLocalMode) {
      localModeCurrentPage.value = 1; // Reset local pagination on filter change
    } else {
      // For API mode, if filtering is purely client-side on the current page,
      // no need to call changeGroupedHistoryPage unless we want to refetch page 1.
      // If the API supported filtering, we would call:
      // historyStore.fetchGroupedHistory(1, historyStore.pageSizeForGrouped, props.filterOption);
      // For now, client-side filtering on the current API page is fine.
      // If filteredAndSortedBatches becomes empty, it will show the empty state.
    }
  }
);

// Função de formatação da data para exibição
function formatBatchDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}
</script>

<template>
  <div>
    <div
      v-if="
        authStore.isLocalMode ? false : historyStore.isLoadingGroupedHistory
      "
      class="text-center p-8"
    >
      <!-- Loader ou indicador de carregamento pode ser adicionado aqui -->
      <span class="material-icons-outlined animate-spin text-3xl"
        >hourglass_empty</span
      >
    </div>

    <div v-else-if="paginatedBatchesToDisplay.length > 0" class="space-y-6">
      <!-- Para cada lote de histórico -->
      <div
        v-for="batch in paginatedBatchesToDisplay"
        :key="batch.batchId"
        class="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden"
      >
        <!-- Cabeçalho do lote com data/hora -->
        <div
          class="bg-gradient-to-r from-purple-600 to-indigo-600 text-white p-3"
        >
          <h3 class="font-medium flex items-center">
            <span class="material-icons-outlined mr-2">history</span>
            Alterações em {{ formatBatchDate(batch.createdAt) }}
            <span v-if="!authStore.isLocalMode" class="text-xs ml-2 opacity-75">
              (BatchID: {{ batch.batchId.substring(0, 8) }}...)
            </span>
          </h3>
        </div>

        <!-- Conteúdo dos registros do lote -->
        <div class="p-0">
          <!-- Desktop: Componente de tabela -->
          <HistoryTable :batch="batch" :is-local-mode="authStore.isLocalMode" />

          <!-- Mobile: Componente de cards -->
          <HistoryCard :batch="batch" :is-local-mode="authStore.isLocalMode" />
        </div>
      </div>
    </div>

    <!-- Indicação de vazio -->
    <div
      v-else-if="
        !(authStore.isLocalMode
          ? false
          : historyStore.isLoadingGroupedHistory) &&
        paginatedBatchesToDisplay.length === 0
      "
      class="bg-white p-8 rounded-lg shadow-md text-center text-gray-500"
    >
      <div class="flex flex-col items-center">
        <span class="material-icons-outlined text-5xl text-gray-300 mb-3"
          >history</span
        >
        <p class="text-lg">Nenhum registro de histórico encontrado.</p>
        <p class="text-sm text-gray-400">
          Os registros aparecerão aqui quando você fizer alterações no estoque.
        </p>
      </div>
    </div>

    <!-- Componente de paginação -->
    <HistoryPagination
      v-if="
        !(authStore.isLocalMode
          ? false
          : historyStore.isLoadingGroupedHistory) &&
        ((authStore.isLocalMode
          ? filteredAndSortedBatches.length
          : historyStore.groupedHistory?.totalBatches) || 0) > 0
      "
      :current-page="currentPage"
      :total-pages="totalPages"
      :pagination-info="paginationInfo"
      @prev-page="changePage(currentPage - 1)"
      @next-page="changePage(currentPage + 1)"
      @go-to-page="changePage"
    />
  </div>
</template>
