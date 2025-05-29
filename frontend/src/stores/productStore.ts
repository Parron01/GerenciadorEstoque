import { defineStore } from "pinia";
import type { Product } from "@/models/product";
import type { Lote, LotePayload } from "@/models/lote";
import {
  fetchProducts,
  createProductApi,
  updateProductApi,
  deleteProductApi,
  createLoteApi,
  updateLoteApi,
  deleteLoteApi,
  // fetchProductById, // Not explicitly used in actions below, but good to have
} from "@/services/apiService";
import { useAuthStore } from "./authStore";
import { useToast } from "vue-toastification";
import { v4 as uuidv4 } from "uuid";

const PRODUCTS_STORAGE_KEY = "inventoryProducts";

export const useProductStore = defineStore("product", {
  state: () => ({
    products: [] as Product[],
    isLoading: false,
  }),
  actions: {
    // Ensure saveToStorage is defined
    saveToStorage() {
      localStorage.setItem(PRODUCTS_STORAGE_KEY, JSON.stringify(this.products));
    },
    loadFromStorage() {
      const storedProducts = localStorage.getItem(PRODUCTS_STORAGE_KEY);
      if (storedProducts) {
        this.products = JSON.parse(storedProducts);
      } else {
        this.products = []; // Initialize as empty if nothing in storage
        // Optionally save empty state: this.saveToStorage();
      }
    },

    async fetchProductsFromApi() {
      this.isLoading = true;
      try {
        this.products = await fetchProducts();
      } catch (error) {
        console.error("Failed to fetch products from API:", error);
        useToast().error("Falha ao buscar produtos do servidor.");
        if (!useAuthStore().isLocalMode) {
          // Only load from storage if not in local mode and API fails
          this.products = []; // Or handle error differently
        }
      } finally {
        this.isLoading = false;
      }
    },

    // Ensure initializeStore is defined
    initializeStore() {
      const authStore = useAuthStore();
      if (authStore.isLocalMode) {
        this.loadFromStorage();
      } else {
        this.fetchProductsFromApi();
      }
    },

    async addProduct(
      productData: Omit<Product, "id" | "lotes"> & { quantity?: number }
      // No operationBatchId for standalone creation
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const newProductWithId: Product = {
        ...productData,
        id: uuidv4(),
        lotes: [],
        quantity: productData.quantity === undefined ? 0 : productData.quantity,
      };

      if (authStore.isLocalMode) {
        this.products.push(newProductWithId);
        this.saveToStorage();
        toast.success(
          `Produto "${newProductWithId.name}" adicionado localmente.`
        );
      } else {
        try {
          const createdProduct = await createProductApi(productData);
          // Instead of pushing, fetch all products again to get the latest state including the new one
          await this.fetchProductsFromApi();
          toast.success(`Produto "${createdProduct.name}" criado no servidor.`);
        } catch (error) {
          console.error("Failed to create product on API:", error);
          toast.error(`Falha ao criar produto: ${error}`);
        }
      }
    },

    async updateProductDetails(
      productId: string,
      details: Partial<Pick<Product, "name" | "unit">>,
      operationBatchId?: string // Corrected: Added optional
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;

      const oldDetails = { name: product.name, unit: product.unit };
      // Optimistic update for UI responsiveness
      Object.assign(product, details);

      if (authStore.isLocalMode) {
        this.saveToStorage();
      } else {
        try {
          await updateProductApi(productId, details, operationBatchId);
          // Success handled by confirmUpdates or specific UI feedback
        } catch (error) {
          console.error("Failed to update product details on API:", error);
          toast.error(`Falha ao atualizar detalhes do produto: ${error}`);
          Object.assign(product, oldDetails); // Revert optimistic update on error
        }
      }
    },

    async updateProductQuantity(
      productId: string,
      newQuantity: number,
      operationBatchId?: string // Corrected: Added optional
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;
      // This method should primarily be for products without lotes,
      // as quantity for products with lotes is derived.
      if (product.lotes && product.lotes.length > 0) {
        // console.warn("updateProductQuantity called for product with lotes. Quantity is derived.");
        // For products with lotes, this might update a 'base_quantity' if your backend/model supports it.
        // Otherwise, this change might be purely for local UI representation before backend recalculates.
      }

      const oldQuantity = product.quantity;
      product.quantity = newQuantity; // Optimistic update

      if (authStore.isLocalMode) {
        this.saveToStorage();
      } else {
        try {
          // The API call might update a 'base_quantity' or be ignored if quantity is purely derived.
          await updateProductApi(
            productId,
            { quantity: newQuantity },
            operationBatchId
          );
        } catch (error) {
          console.error("Failed to update product quantity on API:", error);
          toast.error(`Falha ao atualizar quantidade do produto: ${error}`);
          product.quantity = oldQuantity; // Revert
        }
      }
    },

    async removeProduct(
      productId: string /* No operationBatchId for standalone deletion */
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const productIndex = this.products.findIndex((p) => p.id === productId);
      if (productIndex === -1) return;

      const removedProduct = this.products[productIndex];
      this.products.splice(productIndex, 1); // Optimistic update

      if (authStore.isLocalMode) {
        this.saveToStorage();
        toast.info(`Produto "${removedProduct.name}" removido localmente.`);
      } else {
        try {
          await deleteProductApi(productId);
          toast.info(`Produto "${removedProduct.name}" removido do servidor.`);
          // No need to fetch all products again, optimistic update is usually fine for deletes.
          // If server might reject deletion, then a fetch might be needed on error.
        } catch (error) {
          console.error("Failed to delete product on API:", error);
          toast.error(`Falha ao remover produto: ${error}`);
          this.products.splice(productIndex, 0, removedProduct); // Revert
        }
      }
    },

    // Lote Actions
    async createLote(
      productId: string,
      loteData: LotePayload,
      operationBatchId?: string // Corrected: Added optional
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;

      if (authStore.isLocalMode) {
        const newLote: Lote = {
          ...loteData,
          id: uuidv4(),
          productId: productId,
          createdAt: new Date().toISOString(),
        };
        if (!product.lotes) product.lotes = [];
        product.lotes.push(newLote);
        this.saveToStorage();
      } else {
        try {
          const createdLote = await createLoteApi(
            productId,
            loteData,
            operationBatchId
          );
          if (!product.lotes) product.lotes = [];
          product.lotes.push(createdLote); // Add the server-confirmed lote
          // Optionally, refetch the product or all products to ensure full sync
          // await this.fetchProductsFromApi(); // Or fetchProductById(productId)
        } catch (error) {
          console.error("Failed to create lote on API:", error);
          toast.error(`Falha ao criar lote: ${error}`);
        }
      }
    },

    async updateLote(
      loteId: string,
      productId: string,
      loteData: LotePayload,
      operationBatchId?: string // Corrected: Added optional
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      const lote = product?.lotes?.find((l) => l.id === loteId);
      if (!lote) return;

      const oldLoteData = { ...lote };
      Object.assign(lote, loteData, { updatedAt: new Date().toISOString() }); // Optimistic

      if (authStore.isLocalMode) {
        this.saveToStorage();
      } else {
        try {
          const updatedLote = await updateLoteApi(
            loteId,
            loteData,
            operationBatchId
          );
          Object.assign(lote, updatedLote); // Update with server response
        } catch (error) {
          console.error("Failed to update lote on API:", error);
          toast.error(`Falha ao atualizar lote: ${error}`);
          Object.assign(lote, oldLoteData); // Revert
        }
      }
    },

    async deleteLote(
      loteId: string,
      productId: string,
      operationBatchId?: string // Corrected: Added optional
    ) {
      const authStore = useAuthStore();
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product || !product.lotes) return;
      const loteIndex = product.lotes.findIndex((l) => l.id === loteId);
      if (loteIndex === -1) return;

      const removedLote = product.lotes[loteIndex];
      product.lotes.splice(loteIndex, 1); // Optimistic

      if (authStore.isLocalMode) {
        this.saveToStorage();
      } else {
        try {
          await deleteLoteApi(loteId, operationBatchId);
        } catch (error) {
          console.error("Failed to delete lote on API:", error);
          toast.error(`Falha ao remover lote: ${error}`);
          product.lotes.splice(loteIndex, 0, removedLote); // Revert
        }
      }
    },
  },
});
