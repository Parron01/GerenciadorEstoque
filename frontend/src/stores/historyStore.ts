import { defineStore } from "pinia";
import { ref, watch } from "vue";
import type { ProductHistory, ProductChange } from "@/models/product";
import type { ParsedHistoryRecord } from "@/models/history"; // Import new history model
import { v4 as uuidv4 } from "uuid";
import { useAuthStore } from "@/stores/authStore";
import { useToast } from "vue-toastification";
import * as api from "@/services/apiService"; // Import API service
import { useProductStore } from "./productStore";

// Storage key depends on auth mode
const getStorageKey = () => {
  const authStore = useAuthStore();
  return authStore.isLocalMode
    ? "estoque_historico_local"
    : "estoque_historico";
};

export const useHistoryStore = defineStore("history", () => {
  const authStore = useAuthStore();
  // History can be ProductHistory (local mode, batched product changes) or ParsedHistoryRecord (auth mode)
  const history = ref<(ProductHistory | ParsedHistoryRecord)[]>([]);
  const toast = useToast();
  const isLoading = ref(false);

  // Dados de demonstração para o modo local
  const createDemoHistory = (): ProductHistory[] => {
    const oneHourAgo = new Date(Date.now() - 3600000).toISOString();
    const yesterday = new Date(Date.now() - 86400000).toISOString();
    const lastWeek = new Date(Date.now() - 604800000).toISOString();

    return [
      {
        id: uuidv4(),
        date: oneHourAgo,
        changes: [
          {
            productId: "demo-prod-1", // Example ID from productStore demo data
            productName: "Fertilizante NPK",
            action: "add",
            quantityChanged: 20,
            quantityBefore: 100,
            quantityAfter: 120,
          },
        ],
      },
      {
        id: uuidv4(),
        date: yesterday,
        changes: [
          {
            productId: "demo-prod-2", // Example ID
            productName: "Herbicida Natural",
            action: "remove",
            quantityChanged: 5,
            quantityBefore: 50,
            quantityAfter: 45,
          },
        ],
      },
      // ... more demo entries
    ];
  };

  async function fetchHistoryFromApi() {
    if (authStore.isLocalMode) {
      loadFromStorage(); // Local mode uses localStorage
      return;
    }
    isLoading.value = true;
    try {
      const rawHistory = await api.fetchHistoryApi();
      const productStore = useProductStore(); // To get product names for context

      // Enrich lote history with product names if possible
      history.value = rawHistory.map((record) => {
        if (
          record.entityType === "lote" &&
          record.details &&
          (record.details as any).productId
        ) {
          const product = productStore.products.find(
            (p) => p.id === (record.details as any).productId
          );
          return {
            ...record,
            productNameContext: product?.name || "Produto Desconhecido",
          };
        }
        return record;
      });
    } catch (error: any) {
      toast.error(`Erro ao buscar histórico: ${error.message}`);
      history.value = [];
    } finally {
      isLoading.value = false;
    }
  }

  function loadFromStorage() {
    if (authStore.isLocalMode) {
      const storageKey = getStorageKey();
      const data = localStorage.getItem(storageKey);
      if (data) {
        history.value = JSON.parse(data);
      } else {
        history.value = createDemoHistory();
        saveToStorage();
      }
    } else {
      // For authenticated mode, data is fetched.
      // This function is primarily for local mode persistence.
    }
  }

  function saveToStorage() {
    if (authStore.isLocalMode) {
      const storageKey = getStorageKey();
      localStorage.setItem(storageKey, JSON.stringify(history.value));
    }
  }

  // This function is for local mode product batch changes
  function addBatchEntry(changes: ProductChange[]) {
    if (!authStore.isLocalMode || changes.length === 0) return;

    const newEntry: ProductHistory = {
      id: uuidv4(),
      date: new Date().toISOString(),
      changes,
    };
    history.value.unshift(newEntry);
    saveToStorage();
  }

  // Call this after Lote/Product operations in AUTH mode to refresh history
  function refreshHistory() {
    if (authStore.isAuthenticated) {
      fetchHistoryFromApi();
    }
  }

  // Initial data load logic
  if (authStore.isAuthenticated) {
    fetchHistoryFromApi();
  } else if (authStore.isLocalMode) {
    loadFromStorage();
  }

  watch(
    () => authStore.isLocalMode,
    (newIsLocalMode) => {
      if (newIsLocalMode) {
        loadFromStorage();
      } else {
        fetchHistoryFromApi();
      }
    },
    { immediate: false }
  );

  watch(
    () => authStore.token,
    (newToken) => {
      if (newToken && !authStore.isLocalMode) {
        fetchHistoryFromApi();
      } else if (!newToken && !authStore.isLocalMode) {
        history.value = [];
      }
    }
  );

  // Watch history for local mode saving
  watch(
    history,
    () => {
      if (authStore.isLocalMode) {
        saveToStorage();
      }
    },
    { deep: true }
  );

  return {
    history,
    isLoading,
    addBatchEntry, // For local mode product changes
    fetchHistoryFromApi, // For auth mode
    loadFromStorage, // For explicit reloads if needed
    refreshHistory, // To be called after CRUD operations in auth mode
  };
});
