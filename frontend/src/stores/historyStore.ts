import { defineStore } from "pinia";
import {
  fetchGroupedHistoryApi,
  createProductBatchContextHistoryApi, // Import the new API service function
} from "@/services/apiService";
import type {
  PaginatedHistoryBatchGroups,
  HistoryBatchGroup,
  ProductBatchContextPayload, // Import payload type
} from "@/models/history";
import { useToast } from "vue-toastification";

export const useHistoryStore = defineStore("history", {
  state: () => ({
    groupedHistory: null as PaginatedHistoryBatchGroups | null, // Renamed for clarity
    isLoadingGroupedHistory: false, // Renamed for clarity
    currentPageForGrouped: 1, // Renamed for clarity
    pageSizeForGrouped: 5, // Renamed for clarity
  }),
  // getters can remain similar, just ensure they point to the renamed state properties if needed
  getters: {
    historyGroups(): HistoryBatchGroup[] {
      return this.groupedHistory?.groups || [];
    },
    totalPages(): number {
      return this.groupedHistory?.totalPages || 0;
    },
    totalBatches(): number {
      return this.groupedHistory?.totalBatches || 0;
    },
  },
  actions: {
    async fetchGroupedHistory(page: number = 1, pageSize?: number) {
      // Renamed method
      this.isLoadingGroupedHistory = true;
      try {
        const effectivePageSize = pageSize || this.pageSizeForGrouped;
        const data = await fetchGroupedHistoryApi(page, effectivePageSize);
        this.groupedHistory = data;
        this.currentPageForGrouped = page;
        // Logic to handle empty page (e.g., after deleting last item on a page)
        if (data.groups.length === 0 && page > 1) {
          // Go to previous page if current page is empty and not the first page
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
      // Calls the renamed fetchGroupedHistory with current page and size
      await this.fetchGroupedHistory(
        this.currentPageForGrouped,
        this.pageSizeForGrouped
      );
    },

    async changeGroupedHistoryPage(page: number) {
      // Renamed method
      if (
        page < 1 ||
        (this.groupedHistory && page > this.groupedHistory.totalPages)
      ) {
        return; // Page out of bounds
      }
      // Calls the renamed fetchGroupedHistory
      await this.fetchGroupedHistory(page, this.pageSizeForGrouped);
    },

    clearHistory() {
      this.groupedHistory = null;
      this.currentPageForGrouped = 1; // Reset page
    },

    // New action to record product batch context
    async createProductBatchContextHistory(
      payload: ProductBatchContextPayload,
      operationBatchId: string
    ) {
      // const toast = useToast(); // Toasting handled by ProductTable for overall success/failure
      try {
        await createProductBatchContextHistoryApi(payload, operationBatchId);
        console.log(
          `[HistoryStore] Product batch context history recorded for product ${payload.productId} in batch ${operationBatchId}`
        );
      } catch (error) {
        console.error(
          `[HistoryStore] Failed to record product batch context history for product ${payload.productId}:`,
          error
        );
        // Re-throw the error so ProductTable.vue can handle it as part of allApiCallsSuccessful
        throw error;
      }
    },
  },
});
