<script setup lang="ts">
import { useHistoryStore } from "@/stores/historyStore";
import { computed, ref, watch } from "vue";
import type { ProductHistory } from "@/models/product";
import type { ParsedHistoryRecord } from "@/models/history";
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

// Reset para primeira página quando o filtro muda
watch(
  () => props.filterOption,
  () => {
    currentPage.value = 1;
  }
);
</script>

<template>
  <div>
    <!-- Componente de tabela para desktop -->
    <HistoryTable :history="paginatedHistory" />

    <!-- Componente de cards para mobile -->
    <HistoryCard :history="paginatedHistory" />

    <!-- Componente de paginação -->
    <HistoryPagination
      v-if="filteredHistory.length > 0"
      :current-page="currentPage"
      :total-pages="totalPages"
      :pagination-info="paginationInfo"
      @prev-page="prevPage"
      @next-page="nextPage"
      @go-to-page="goToPage"
    />
  </div>
</template>
