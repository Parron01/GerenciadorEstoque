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

    async updateProductQuantity(
      productId: string,
      newQuantity: number,
      operationBatchId?: string
    ) {
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;

      const oldQuantity = product.quantity;
      product.quantity = newQuantity;

      try {
        await updateProductApi(
          productId,
          { quantity: newQuantity },
          operationBatchId
        );
      } catch (error) {
        console.error("Failed to update product quantity:", error);
        toast.error(`Falha ao atualizar quantidade do produto: ${error}`);
        product.quantity = oldQuantity;
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
      const product = this.products.find((p) => p.id === productId);
      if (!product) return;

      try {
        const createdLote = await createLoteApi(
          productId,
          loteData,
          operationBatchId
        );
        if (!product.lotes) product.lotes = [];
        product.lotes.push(createdLote);
      } catch (error) {
        console.error("Failed to create lote:", error);
        toast.error(`Falha ao criar lote: ${error}`);
      }
    },

    async updateLote(
      loteId: string,
      productId: string,
      loteData: LotePayload,
      operationBatchId?: string
    ) {
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      const lote = product?.lotes?.find((l) => l.id === loteId);
      if (!lote) return;

      const oldLoteData = { ...lote };
      Object.assign(lote, loteData, { updatedAt: new Date().toISOString() });

      try {
        const updatedLote = await updateLoteApi(
          loteId,
          loteData,
          operationBatchId
        );
        Object.assign(lote, updatedLote);
      } catch (error) {
        console.error("Failed to update lote:", error);
        toast.error(`Falha ao atualizar lote: ${error}`);
        Object.assign(lote, oldLoteData);
      }
    },

    async deleteLote(
      loteId: string,
      productId: string,
      operationBatchId?: string
    ) {
      const toast = useToast();
      const product = this.products.find((p) => p.id === productId);
      if (!product || !product.lotes) return;
      const loteIndex = product.lotes.findIndex((l) => l.id === loteId);
      if (loteIndex === -1) return;

      const removedLote = product.lotes[loteIndex];
      product.lotes.splice(loteIndex, 1);

      try {
        await deleteLoteApi(loteId, operationBatchId);
      } catch (error) {
        console.error("Failed to delete lote:", error);
        toast.error(`Falha ao remover lote: ${error}`);
        product.lotes.splice(loteIndex, 0, removedLote);
      }
    },
  },
});
