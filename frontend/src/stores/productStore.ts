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
} from "@/services/apiService";
import { useToast } from "vue-toastification";

export const useProductStore = defineStore("product", {
  state: () => ({
    products: [] as Product[],
    isLoading: false,
  }),
  actions: {
    async fetchProductsFromApi() {
      this.isLoading = true;
      try {
        this.products = await fetchProducts();
      } catch (error) {
        console.error("Failed to fetch products from API:", error);
        useToast().error("Falha ao buscar produtos do servidor.");
        this.products = [];
      } finally {
        this.isLoading = false;
      }
    },

    initializeStore() {
      this.fetchProductsFromApi();
    },

    async addProduct(
      productData: Omit<Product, "id" | "lotes" | "quantity"> & {
        quantity?: number;
      }
    ) {
      const toast = useToast();
      try {
        await createProductApi(productData);
        await this.fetchProductsFromApi();
        toast.success(`Produto "${productData.name}" criado com sucesso.`);
      } catch (error) {
        console.error("Failed to create product:", error);
        toast.error(`Falha ao criar produto: ${error}`);
      }
    },

    async updateProductDetails(
      productId: string,
      details: Partial<Pick<Product, "name" | "unit">>,
      operationBatchId?: string
    ) {
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;

      const oldDetails = { name: product.name, unit: product.unit };
      Object.assign(product, details);

      try {
        await updateProductApi(productId, details, operationBatchId);
      } catch (error) {
        console.error("Failed to update product details:", error);
        toast.error(`Falha ao atualizar detalhes do produto: ${error}`);
        Object.assign(product, oldDetails);
      }
    },

    async removeProduct(productId: string) {
      const toast = useToast();
      const productIndex = this.products.findIndex((p) => p.id === productId);
      if (productIndex === -1) return;

      const removedProduct = this.products[productIndex];
      this.products.splice(productIndex, 1);

      try {
        await deleteProductApi(productId);
        toast.info(`Produto "${removedProduct.name}" removido com sucesso.`);
      } catch (error) {
        console.error("Failed to delete product:", error);
        toast.error(`Falha ao remover produto: ${error}`);
        this.products.splice(productIndex, 0, removedProduct);
      }
    },

    async createLote(
      productId: string,
      loteData: LotePayload,
      operationBatchId?: string
    ) {
      const toast = useToast();
      // Optimistic update is handled by ProductTable.vue's handleSaveLote.
      // This store action is now only responsible for the API call.
      // The productStore.products state will be fully synchronized by
      // fetchProductsFromApi() at the end of ProductTable.vue's confirmUpdates.

      try {
        // Call the API. We don't need to push to product.lotes here.
        await createLoteApi(productId, loteData, operationBatchId);
        // If createLoteApi returned the created lote, and we needed to update
        // the localId-based optimistic lote with the serverId, this is where it would happen,
        // but it requires more complex state management (passing localId to this action).
        // Relying on fetchProductsFromApi() for final sync is simpler.
      } catch (error) {
        console.error("Failed to create lote:", error);
        toast.error(`Falha ao criar lote: ${error}`);
        throw error; // Re-throw so confirmUpdates in ProductTable.vue knows this specific call failed
      }
    },

    async updateLote(
      loteId: string,
      productId: string,
      loteData: LotePayload,
      operationBatchId?: string
    ) {
      const toast = useToast();
      // Optimistic update to the specific lote object in productStore.products
      // is handled by ProductTable.vue's handleUpdateLote.
      // This store action is now only responsible for the API call.
      // The productStore.products state will be fully synchronized by
      // fetchProductsFromApi() at the end of ProductTable.vue's confirmUpdates.

      try {
        // Just call the API.
        await updateLoteApi(loteId, loteData, operationBatchId);
        // Similar to createLote, if selective update of the store instance was needed
        // without a full fetch, it would happen here using the server response.
      } catch (error) {
        console.error("Failed to update lote:", error);
        toast.error(`Falha ao atualizar lote: ${error}`);
        // Rollback of ProductTable's optimistic update will effectively happen
        // via fetchProductsFromApi() if the API call fails.
        throw error; // Re-throw so confirmUpdates in ProductTable.vue knows this specific call failed
      }
    },

    async deleteLote(
      loteId: string,
      productId: string, // Kept for context/logging, though not directly used by deleteLoteApi
      operationBatchId?: string
    ) {
      const toast = useToast();
      // No need to find product or lote in store's state here for the API call,
      // as the component has already handled the optimistic UI update
      // and the decision to delete is firm based on loteChangesTracking.

      try {
        console.log(
          `[ProductStore] Calling deleteLoteApi for loteId: ${loteId}, productId (for context): ${productId}, batchId: ${operationBatchId}`
        );
        await deleteLoteApi(loteId, operationBatchId);
        // Success can be noted, but the main success/failure summary is in ProductTable.vue
      } catch (error) {
        console.error(
          `[ProductStore] Failed to delete lote ${loteId} from API:`,
          error
        );
        // Toast the specific error here.
        toast.error(
          `Falha ao remover lote ${loteId.substring(0, 8)} do servidor: ${error}`
        );
        // Re-throw the error so confirmUpdates in ProductTable.vue knows this specific call failed
        // and can set allApiCallsSuccessful = false.
        throw error;
      }
    },
  },
});
