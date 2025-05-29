<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  currentPage: number;
  totalPages: number;
  paginationInfo: {
    showing: boolean;
    start: number;
    end: number;
    total: number;
  };
}>();

const emit = defineEmits<{
  (e: "prevPage"): void;
  (e: "nextPage"): void;
  (e: "goToPage", page: number): void;
}>();

function prevPage() {
  emit("prevPage");
}

function nextPage() {
  emit("nextPage");
}

function goToPage(page: number) {
  emit("goToPage", page);
}

// Gera array de números de página para exibição
const displayedPageNumbers = computed(() => {
  const total = props.totalPages;
  const current = props.currentPage;
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
</script>

<template>
  <div
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
        :class="[currentPage >= totalPages ? 'btn-page-disabled' : 'btn-page']"
        aria-label="Próxima página"
      >
        <span class="material-icons-outlined text-lg">chevron_right</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
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
