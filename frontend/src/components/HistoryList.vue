<script setup lang="ts">
import { useHistoryStore } from "@/stores/historyStore";
import { computed, ref, watch } from "vue";
import type { HistoryBatchGroup } from "@/models/history";
import HistoryTable from "./history-list/HistoryTable.vue";
import HistoryCard from "./history-list/HistoryCard.vue";
import HistoryPagination from "./history-list/HistoryPagination.vue";

const props = defineProps<{
  filterOption: string;
}>();

const historyStore = useHistoryStore();

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
  firstDay.setHours(0, 0, 0, 0);
  lastDay.setHours(23, 59, 59, 999);
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

// New helper to check if batch contains records for filtered product
const matchesProductFilter = (batch: HistoryBatchGroup): boolean => {
  if (!historyStore.productFilter) return true;

  // Check if any record in batch references the filtered product
  if (batch.records && batch.records.length > 0) {
    // Check product-specific records
    const hasMatchingRecord = batch.records.some((record) => {
      // Direct match for product records
      if (
        record.entityType === "product" &&
        record.entityId === historyStore.productFilter
      ) {
        return true;
      }

      // Match for lote records that belong to the product
      if (
        record.entityType === "lote" &&
        record.details?.productId === historyStore.productFilter
      ) {
        return true;
      }

      return false;
    });

    if (hasMatchingRecord) return true;
  }

  // Check product summaries (which contain all affected products)
  if (
    batch.productSummaries &&
    historyStore.productFilter in batch.productSummaries
  ) {
    return true;
  }

  return false;
};

const filteredAndSortedBatches = computed((): HistoryBatchGroup[] => {
  if (!historyStore.groupedHistory || !historyStore.groupedHistory.groups)
    return [];

  return historyStore.groupedHistory.groups
    .filter((batch: HistoryBatchGroup) => {
      if (!batch.records || batch.records.length === 0) return false;

      // Apply date filter
      const batchDate = batch.createdAt;
      let matchesDateFilter = true;
      switch (props.filterOption) {
        case "today":
          matchesDateFilter = isToday(batchDate);
          break;
        case "week":
          matchesDateFilter = isThisWeek(batchDate);
          break;
        case "month":
          matchesDateFilter = isThisMonth(batchDate);
          break;
        default:
          matchesDateFilter = true;
      }

      // Apply product filter
      const productMatches = matchesProductFilter(batch);

      // Both filters must match
      return matchesDateFilter && productMatches;
    })
    .sort(
      (a, b) =>
        new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    );
});

const currentPage = computed(() => historyStore.currentPageForGrouped);
const totalPages = computed(() => historyStore.groupedHistory?.totalPages || 1);

const paginationInfo = computed(() => {
  const totalItems = historyStore.groupedHistory?.totalBatches || 0;
  if (totalItems === 0) return { showing: false, start: 0, end: 0, total: 0 };

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
});

function changePage(page: number) {
  historyStore.changeGroupedHistoryPage(page);
}

function formatBatchDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}
</script>

<template>
  <div>
    <div v-if="historyStore.isLoadingGroupedHistory" class="text-center p-8">
      <span class="material-icons-outlined animate-spin text-3xl"
        >hourglass_empty</span
      >
    </div>

    <div v-else-if="filteredAndSortedBatches.length > 0" class="space-y-6">
      <!-- Add filter results count -->
      <div
        class="text-sm text-gray-600 bg-gray-50 px-3 py-2 rounded"
        v-if="historyStore.productFilter || props.filterOption !== 'all'"
      >
        {{ filteredAndSortedBatches.length }} registro(s) encontrado(s)
      </div>

      <div
        v-for="batch in filteredAndSortedBatches"
        :key="batch.batchId"
        class="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden"
      >
        <div
          class="bg-gradient-to-r from-purple-600 to-indigo-600 text-white p-3"
        >
          <h3 class="font-medium flex items-center">
            <span class="material-icons-outlined mr-2">history</span>
            Alterações em {{ formatBatchDate(batch.createdAt) }}
          </h3>
        </div>

        <div class="p-0">
          <HistoryTable :batch="batch" />
          <HistoryCard :batch="batch" />
        </div>
      </div>
    </div>

    <div
      v-else-if="
        !historyStore.isLoadingGroupedHistory &&
        filteredAndSortedBatches.length === 0
      "
      class="bg-white p-8 rounded-lg shadow-md text-center text-gray-500"
    >
      <div class="flex flex-col items-center">
        <span class="material-icons-outlined text-5xl text-gray-300 mb-3"
          >history</span
        >
        <p class="text-lg">Nenhum registro de histórico encontrado.</p>
        <p
          class="text-sm text-gray-400"
          v-if="historyStore.productFilter || props.filterOption !== 'all'"
        >
          Tente ajustar os filtros de busca.
        </p>
        <p class="text-sm text-gray-400" v-else>
          Os registros aparecerão aqui quando você fizer alterações no estoque.
        </p>

        <!-- Clear filters button when no results found -->
        <button
          v-if="historyStore.productFilter || props.filterOption !== 'all'"
          @click="
            historyStore.clearProductFilter();
            $emit('update:filterOption', 'all');
          "
          class="mt-3 px-4 py-2 bg-indigo-100 text-indigo-700 rounded-lg hover:bg-indigo-200 transition-colors flex items-center"
        >
          <span class="material-icons-outlined mr-1">refresh</span>
          Limpar filtros
        </button>
      </div>
    </div>

    <!-- Only show pagination when not filtering -->
    <HistoryPagination
      v-if="
        !historyStore.isLoadingGroupedHistory &&
        (historyStore.groupedHistory?.totalBatches || 0) > 0 &&
        !historyStore.productFilter &&
        props.filterOption === 'all'
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
