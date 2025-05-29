import { defineStore } from "pinia";
import {
  fetchHistoryApi,
  fetchGroupedHistoryApi, // Added
  // createHistoryBatchApi, // Kept for now, but less used by ProductTable
  // fetchHistoryBatchApi, // Kept for now
} from "@/services/apiService";
import type {
  ParsedHistoryRecord,
  ProductHistory,
  PaginatedHistoryBatchGroups, // Added
  HistoryBatchGroup, // Added
} from "@/models/history";
import { useAuthStore } from "./authStore";
import { useToast } from "vue-toastification";

const HISTORY_STORAGE_KEY = "inventoryHistory";

export const useHistoryStore = defineStore("history", {
  state: () => ({
    history: [] as ProductHistory[], // For local mode, or potentially ungrouped history
    groupedHistory: null as PaginatedHistoryBatchGroups | null, // For authenticated mode
    isLoading: false,
    isLoadingGroupedHistory: false, // Added
    currentPageForGrouped: 1, // Added for pagination of grouped history
    pageSizeForGrouped: 5, // Batches per page
  }),
  actions: {
    loadHistoryFromStorage() {
      const storedHistory = localStorage.getItem(HISTORY_STORAGE_KEY);
      if (storedHistory) {
        this.history = JSON.parse(storedHistory);
      }
    },
    saveHistoryToStorage() {
      localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(this.history));
    },

    // This was for local mode batching, where frontend constructed the history details.
    // It's less relevant now for authenticated mode.
    // For local mode, it can still be used to add to `this.history`.
    addBatchEntry(productChanges: any[]) {
      // For local mode, this structure might need to align with ProductHistory
      const newEntry: ProductHistory = {
        id: crypto.randomUUID(), // Or a more structured batch ID
        date: new Date().toISOString(),
        changes: productChanges,
        batchId: crypto.randomUUID(), // Ensure local entries also have a batchId
      };
      this.history.unshift(newEntry); // Add to the beginning
      this.saveHistoryToStorage();
    },

    async fetchHistory() {
      // This fetches ungrouped history, might be used for specific cases or admin views.
      // For the main HistoryView, fetchGroupedHistory will be used.
      const authStore = useAuthStore();
      if (authStore.isLocalMode) {
        this.loadHistoryFromStorage();
        return;
      }
      this.isLoading = true;
      try {
        // This is a placeholder if ungrouped history is ever needed directly.
        // const rawHistory = await fetchHistoryApi();
        // this.history = rawHistory; // Adapt if this.history structure changes
        // For now, main history view will use groupedHistory.
      } catch (error) {
        console.error("Failed to fetch history:", error);
        useToast().error("Falha ao carregar histórico do servidor.");
      } finally {
        this.isLoading = false;
      }
    },

    async fetchGroupedHistory(page: number = 1, pageSize?: number) {
      const authStore = useAuthStore();
      if (authStore.isLocalMode) {
        this.loadHistoryFromStorage(); // Local mode uses its own structure
        // Transform local history (ProductHistory[]) to PaginatedHistoryBatchGroups
        // This ensures the UI (HistoryList, HistoryTable, HistoryCard) can consistently
        // render history data regardless of whether it's from local storage or the API.
        const localGroups: HistoryBatchGroup[] = this.history.map(
          (localEntry) => ({
            batchId: localEntry.batchId || localEntry.id,
            createdAt: localEntry.date,
            records: localEntry.changes.map((change: any, index: number) => ({
              // Attempt to map ProductChange to ParsedHistoryRecord
              id: `${localEntry.id}-${index}`, // Create a unique ID for the sub-record
              entityType: change.loteId ? "lote" : "product", // Infer entity type
              entityId: change.productId || change.loteId,
              details: change, // The ProductChange object itself
              createdAt: localEntry.date,
              batchId: localEntry.batchId || localEntry.id,
              productNameContext: change.productName,
            })),
            recordCount: localEntry.changes.length,
          })
        );

        const totalLocalBatches = localGroups.length;
        const effectivePageSize = pageSize || this.pageSizeForGrouped;
        const totalLocalPages = Math.max(
          1,
          Math.ceil(totalLocalBatches / effectivePageSize)
        );
        const startIndex = (page - 1) * effectivePageSize;
        const paginatedLocalGroups = localGroups.slice(
          startIndex,
          startIndex + effectivePageSize
        );

        this.groupedHistory = {
          groups: paginatedLocalGroups,
          totalBatches: totalLocalBatches,
          page: page,
          pageSize: effectivePageSize,
          totalPages: totalLocalPages,
        };
        this.currentPageForGrouped = page;
        return;
      }

      this.isLoadingGroupedHistory = true;
      try {
        const effectivePageSize = pageSize || this.pageSizeForGrouped;
        // Fetches history data already grouped by BatchID from the backend.
        const data = await fetchGroupedHistoryApi(page, effectivePageSize);
        this.groupedHistory = data;
        this.currentPageForGrouped = page;
        if (data.groups.length === 0 && page > 1) {
          // If current page is empty and not the first page
          this.fetchGroupedHistory(Math.max(1, page - 1), effectivePageSize); // Go to previous page or first
        }
      } catch (error) {
        console.error("Failed to fetch grouped history:", error);
        useToast().error("Falha ao carregar histórico agrupado do servidor.");
        this.groupedHistory = null; // Clear on error
      } finally {
        this.isLoadingGroupedHistory = false;
      }
    },

    // refreshHistory now calls fetchGroupedHistory
    async refreshHistory() {
      await this.fetchGroupedHistory(this.currentPageForGrouped);
    },

    // Action to change page for grouped history
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
      this.history = [];
      this.groupedHistory = null;
      if (useAuthStore().isLocalMode) {
        localStorage.removeItem(HISTORY_STORAGE_KEY);
      }
    },
  },
  getters: {
    // Getter for local mode history if needed elsewhere, already transformed for grouped view
    getLocalModeHistoryAsGroups(state): HistoryBatchGroup[] {
      if (!useAuthStore().isLocalMode) return [];
      return state.history.map((localEntry) => ({
        batchId: localEntry.batchId || localEntry.id,
        createdAt: localEntry.date,
        records: localEntry.changes.map((change: any, index: number) => ({
          id: `${localEntry.id}-${index}`,
          entityType: change.loteId ? "lote" : "product",
          entityId: change.productId || change.loteId,
          details: change,
          createdAt: localEntry.date,
          batchId: localEntry.batchId || localEntry.id,
          productNameContext: change.productName,
        })),
        recordCount: localEntry.changes.length,
      }));
    },
  },
});
