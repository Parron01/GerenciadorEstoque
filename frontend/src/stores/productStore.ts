import { defineStore } from "pinia";
import { ref, watch, computed } from "vue";
import type { Product } from "@/models/product";
import type { Lote, LotePayload } from "@/models/lote"; // Import Lote models
import { useAuthStore } from "@/stores/authStore";
import { v4 as uuidv4 } from "uuid";
import { useToast } from "vue-toastification";
import * as api from "@/services/apiService"; // Import API service

// Storage key depends on auth mode
const getStorageKey = () => {
  const authStore = useAuthStore();
  return authStore.isLocalMode ? "estoque_produtos_local" : "estoque_produtos";
};

export const useProductStore = defineStore("product", () => {
  const authStore = useAuthStore();
  const products = ref<Product[]>([]);
  const toast = useToast();
  const isLoading = ref(false);

  // Demo data with lotes
  const demoLotesForProduct1: Lote[] = [
    {
      id: uuidv4(),
      productId: "demo-prod-1",
      quantity: 50,
      dataValidade: "2025-12-31",
      createdAt: new Date().toISOString(),
    },
    {
      id: uuidv4(),
      productId: "demo-prod-1",
      quantity: 70,
      dataValidade: "2026-06-30",
      createdAt: new Date().toISOString(),
    },
  ];
  const demoLotesForProduct2: Lote[] = [
    {
      id: uuidv4(),
      productId: "demo-prod-2",
      quantity: 45,
      dataValidade: "2025-08-15",
      createdAt: new Date().toISOString(),
    },
  ];

  const demoProducts: Product[] = [
    {
      id: "demo-prod-1",
      name: "Fertilizante NPK",
      unit: "kg",
      quantity: 120,
      lotes: demoLotesForProduct1,
    },
    {
      id: "demo-prod-2",
      name: "Herbicida Natural",
      unit: "L",
      quantity: 45,
      lotes: demoLotesForProduct2,
    },
    {
      id: uuidv4(),
      name: "Adubo Orgânico",
      unit: "kg",
      quantity: 200,
      lotes: [],
    },
    { id: uuidv4(), name: "Inseticida Biológico", unit: "L", quantity: 15 }, // No lotes explicitly
  ];

  async function fetchProductsFromApi() {
    if (authStore.isLocalMode) {
      loadFromStorage(); // Continue using local/demo for local mode
      return;
    }
    isLoading.value = true;
    try {
      products.value = await api.fetchProducts();
    } catch (error: any) {
      toast.error(`Erro ao buscar produtos: ${error.message}`);
      products.value = []; // Fallback to empty or handle as needed
    } finally {
      isLoading.value = false;
    }
  }

  function loadFromStorage() {
    if (authStore.isLocalMode) {
      const storageKey = getStorageKey();
      const data = localStorage.getItem(storageKey);
      if (data) {
        products.value = JSON.parse(data);
      } else {
        products.value = JSON.parse(JSON.stringify(demoProducts)); // Deep copy
        saveToStorage();
      }
    } else {
      // For authenticated mode, data is primarily fetched.
      // LocalStorage might be used as a cache, but API is source of truth.
      // For simplicity here, we rely on fetchProductsFromApi for auth mode.
      // If products are empty and auth mode, it implies fetch is needed.
    }
  }

  function saveToStorage() {
    if (authStore.isLocalMode) {
      const storageKey = getStorageKey();
      localStorage.setItem(storageKey, JSON.stringify(products.value));
    }
    // No saving to localStorage for authenticated mode here, API is source of truth.
  }

  // Computed quantity for a product
  const getProductQuantity = (product: Product): number => {
    if (product.lotes && product.lotes.length > 0) {
      return product.lotes.reduce((sum, lote) => sum + lote.quantity, 0);
    }
    return product.quantity;
  };

  // Update product quantity (used for products without lotes or initial set)
  function updateProductQuantity(id: string, newQuantity: number) {
    const product = products.value.find((p) => p.id === id);
    if (product && (!product.lotes || product.lotes.length === 0)) {
      product.quantity = Math.max(0, newQuantity);
      if (authStore.isLocalMode) saveToStorage();
      // In auth mode, if product quantity is directly updatable (no lotes),
      // an API call would be needed here. However, the prompt implies
      // product quantity is derived if lotes exist, and only name/unit are updatable for product.
    }
  }

  async function addProduct(productData: Omit<Product, "id" | "lotes">) {
    isLoading.value = true;
    try {
      if (authStore.isLocalMode) {
        const newProd: Product = { ...productData, id: uuidv4(), lotes: [] };
        products.value.push(newProd);
        saveToStorage();
        toast.success(`Produto "${newProd.name}" adicionado localmente.`);
      } else {
        // For API, quantity is initial. Lotes are added separately.
        const createdProduct = await api.createProductApi(productData);
        products.value.push(createdProduct); // API returns the created product
        toast.success(`Produto "${createdProduct.name}" criado.`);
      }
    } catch (error: any) {
      toast.error(`Erro ao adicionar produto: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  async function updateProductDetails(
    productId: string,
    details: Partial<Pick<Product, "name" | "unit">>
  ) {
    isLoading.value = true;
    try {
      const productIndex = products.value.findIndex((p) => p.id === productId);
      if (productIndex === -1)
        throw new Error("Produto não encontrado localmente.");

      if (authStore.isLocalMode) {
        Object.assign(products.value[productIndex], details);
        saveToStorage();
        toast.success(
          `Produto "${products.value[productIndex].name}" atualizado localmente.`
        );
      } else {
        const updatedProduct = await api.updateProductApi(productId, details);
        // The API returns the full product, potentially with updated lotes sum for quantity
        products.value[productIndex] = {
          ...products.value[productIndex],
          ...updatedProduct,
        };
        toast.success(`Produto "${updatedProduct.name}" atualizado.`);
      }
    } catch (error: any) {
      toast.error(`Erro ao atualizar produto: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  async function removeProduct(id: string) {
    isLoading.value = true;
    try {
      if (authStore.isLocalMode) {
        const index = products.value.findIndex((p) => p.id === id);
        if (index !== -1) {
          products.value.splice(index, 1);
          saveToStorage();
          toast.info("Produto removido localmente.");
        }
      } else {
        await api.deleteProductApi(id);
        products.value = products.value.filter((p) => p.id !== id);
        toast.info("Produto removido do servidor.");
      }
    } catch (error: any) {
      toast.error(`Erro ao remover produto: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  // Lote Management
  async function createLote(productId: string, loteData: LotePayload) {
    isLoading.value = true;
    try {
      const product = products.value.find((p) => p.id === productId);
      if (!product)
        throw new Error("Produto não encontrado para adicionar lote.");

      if (authStore.isLocalMode) {
        const newLote: Lote = {
          ...loteData,
          id: uuidv4(),
          productId,
          createdAt: new Date().toISOString(),
        };
        if (!product.lotes) product.lotes = [];
        product.lotes.push(newLote);
        // Update product quantity if it's derived from lotes in local mode too
        product.quantity = getProductQuantity(product);
        saveToStorage();
        toast.success("Lote adicionado localmente.");
      } else {
        const createdLote = await api.createLoteApi(productId, loteData);
        if (!product.lotes) product.lotes = [];
        product.lotes.push(createdLote);
        // Fetch the product again to get updated quantity and lotes list from server
        const updatedProduct = await api.fetchProductById(productId);
        const productIndex = products.value.findIndex(
          (p) => p.id === productId
        );
        if (productIndex !== -1) products.value[productIndex] = updatedProduct;

        toast.success("Lote criado no servidor.");
      }
    } catch (error: any) {
      toast.error(`Erro ao criar lote: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  async function updateLote(
    loteId: string,
    productId: string,
    loteData: LotePayload
  ) {
    isLoading.value = true;
    try {
      const product = products.value.find((p) => p.id === productId);
      if (!product || !product.lotes)
        throw new Error("Produto ou lotes não encontrados.");

      const loteIndex = product.lotes.findIndex((l) => l.id === loteId);
      if (loteIndex === -1) throw new Error("Lote não encontrado localmente.");

      if (authStore.isLocalMode) {
        Object.assign(product.lotes[loteIndex], loteData, {
          updatedAt: new Date().toISOString(),
        });
        product.quantity = getProductQuantity(product);
        saveToStorage();
        toast.success("Lote atualizado localmente.");
      } else {
        const updatedLote = await api.updateLoteApi(loteId, loteData);
        product.lotes[loteIndex] = updatedLote;
        // Fetch the product again to get updated quantity and lotes list from server
        const updatedProduct = await api.fetchProductById(productId);
        const productIndex = products.value.findIndex(
          (p) => p.id === productId
        );
        if (productIndex !== -1) products.value[productIndex] = updatedProduct;
        toast.success("Lote atualizado no servidor.");
      }
    } catch (error: any) {
      toast.error(`Erro ao atualizar lote: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  async function deleteLote(loteId: string, productId: string) {
    isLoading.value = true;
    try {
      const product = products.value.find((p) => p.id === productId);
      if (!product || !product.lotes)
        throw new Error("Produto ou lotes não encontrados.");

      if (authStore.isLocalMode) {
        product.lotes = product.lotes.filter((l) => l.id !== loteId);
        product.quantity = getProductQuantity(product);
        saveToStorage();
        toast.info("Lote removido localmente.");
      } else {
        await api.deleteLoteApi(loteId);
        product.lotes = product.lotes.filter((l) => l.id !== loteId);
        // Fetch the product again to get updated quantity and lotes list from server
        const updatedProduct = await api.fetchProductById(productId);
        const productIndex = products.value.findIndex(
          (p) => p.id === productId
        );
        if (productIndex !== -1) products.value[productIndex] = updatedProduct;
        toast.info("Lote removido do servidor.");
      }
    } catch (error: any) {
      toast.error(`Erro ao remover lote: ${error.message}`);
    } finally {
      isLoading.value = false;
    }
  }

  // Initial data load logic
  if (authStore.isAuthenticated) {
    fetchProductsFromApi();
  } else if (authStore.isLocalMode) {
    loadFromStorage();
  }

  watch(
    () => authStore.isLocalMode,
    (newIsLocalMode, oldIsLocalMode) => {
      if (newIsLocalMode) {
        loadFromStorage(); // Load demo/local data
      } else {
        // Switched from local to authenticated (e.g. after login)
        fetchProductsFromApi(); // Fetch server data
      }
    },
    { immediate: false }
  ); // immediate might cause issues if authStore is not fully ready

  watch(
    () => authStore.token,
    (newToken) => {
      if (newToken && !authStore.isLocalMode) {
        fetchProductsFromApi();
      } else if (!newToken && !authStore.isLocalMode) {
        products.value = []; // Clear products on logout if not in local mode
      }
    }
  );

  // Watch products for local mode saving (less critical for auth mode as API is truth)
  watch(
    products,
    () => {
      if (authStore.isLocalMode) {
        saveToStorage();
      }
    },
    { deep: true }
  );

  return {
    products,
    isLoading,
    fetchProductsFromApi, // Renamed for clarity
    addProduct,
    updateProductDetails, // Renamed for clarity
    removeProduct,
    updateProductQuantity, // For products without lotes
    getProductQuantity, // Helper to get correct quantity
    createLote,
    updateLote,
    deleteLote,
    loadFromStorage, // Keep for explicit reloads if needed
  };
});
