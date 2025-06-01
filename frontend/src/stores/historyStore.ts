import { defineStore } from "pinia";
import { fetchGroupedHistoryApi } from "@/services/apiService";
import type {
  PaginatedHistoryBatchGroups,
  HistoryBatchGroup,
} from "@/models/history";
import { useToast } from "vue-toastification";

export const useHistoryStore = defineStore("history", {
  state: () => ({
    groupedHistory: null as PaginatedHistoryBatchGroups | null,
    isLoadingGroupedHistory: false,
    currentPageForGrouped: 1,
    pageSizeForGrouped: 5,
  }),
  actions: {
    async fetchGroupedHistory(page: number = 1, pageSize?: number) {
      this.isLoadingGroupedHistory = true;
      try {
        const effectivePageSize = pageSize || this.pageSizeForGrouped;
        const data = await fetchGroupedHistoryApi(page, effectivePageSize);
        this.groupedHistory = data;
        this.currentPageForGrouped = page;
        if (data.groups.length === 0 && page > 1) {
          this.fetchGroupedHistory(Math.max(1, page - 1), effectivePageSize);
        }
      } catch (error) {
        console.error("Failed to fetch grouped history:", error);
        useToast().error("Falha ao carregar histórico do servidor.");
        this.groupedHistory = null;
      } finally {
        this.isLoadingGroupedHistory = false;
      }
    },

    async refreshHistory() {
      await this.fetchGroupedHistory(this.currentPageForGrouped);
    },

    async changeGroupedHistoryPage(page: number) {
      if (
        page < 1 ||
        (this.groupedHistory && page > this.groupedHistory.totalPages)
      ) {
        return;
      }
      await this.fetchGroupedHistory(page, this.pageSizeForGrouped);
    },

    clearHistory() {
      this.groupedHistory = null;
    },
  },
});
