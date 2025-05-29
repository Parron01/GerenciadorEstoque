<script setup lang="ts">
import HistoryList from "@/components/HistoryList.vue";
import { ref, onMounted, watch } from "vue";
import { useHistoryStore } from "@/stores/historyStore";
import { useAuthStore } from "@/stores/authStore";

const selectedFilter = ref("all"); // 'all', 'today', 'week', 'month'
const historyStore = useHistoryStore();
const authStore = useAuthStore();

onMounted(() => {
  // Initial fetch, HistoryList will use historyStore.groupedHistory
  historyStore.fetchGroupedHistory(1, historyStore.pageSizeForGrouped);
});

// Watch for auth mode changes to refetch history
watch(
  () => authStore.isLocalMode,
  () => {
    historyStore.fetchGroupedHistory(1, historyStore.pageSizeForGrouped);
  }
);
</script>

<template>
  <div class="container mx-auto p-4 sm:p-6 lg:p-8">
    <header class="mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Histórico de Alterações</h1>
      <p class="text-gray-600">
        Visualize todas as movimentações e alterações no estoque.
      </p>
    </header>

    <!-- Filtros -->
    <div class="mb-6 p-4 bg-white rounded-lg shadow">
      <label
        for="history-filter"
        class="block text-sm font-medium text-gray-700 mb-1"
        >Filtrar por:</label
      >
      <select
        id="history-filter"
        v-model="selectedFilter"
        class="w-full sm:w-auto px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
      >
        <option value="all">Todos</option>
        <option value="today">Hoje</option>
        <option value="week">Esta Semana</option>
        <option value="month">Este Mês</option>
      </select>
    </div>

    <HistoryList :filter-option="selectedFilter" />
  </div>
</template>
