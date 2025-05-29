<script setup lang="ts">
import { ref, computed } from "vue";
import type { Product } from "@/models/product";
import type { Lote } from "@/models/lote";
import ProductQuantityEditor from "./ProductQuantityEditor.vue";
import LoteDropdown from "./LoteDropdown.vue";

const props = defineProps<{
  product: Product;
  isEditMode: boolean;
  expandedProducts: Record<string, boolean>;
  tempProductDetails: Record<string, { name: string; unit: "L" | "kg" }>;
  tempQuantities: Record<string, number>;
}>();

const emit = defineEmits<{
  (e: "toggleProductLotes", id: string): void;
  (e: "quantityChanged", id: string, delta: number): void;
  (e: "quantityUpdated", id: string, value: number): void;
  (e: "updateProductBaseQuantity", id: string, value: number): void;
  (e: "requestDelete", product: Product): void;
}>();

function getProductDisplayQuantity(product: Product): number {
  if (product.lotes && product.lotes.length > 0) {
    return product.lotes.reduce((sum, lote) => sum + lote.quantity, 0);
  }
  return product.quantity;
}

function toggleProductLotes(id: string) {
  emit("toggleProductLotes", id);
}

function updateProductQuantityDirectly(value: number) {
  emit("quantityUpdated", props.product.id, value);
}

function changeProductQuantity(delta: number) {
  emit("quantityChanged", props.product.id, delta);
}

function updateProductBaseQuantity(value: number) {
  emit("updateProductBaseQuantity", props.product.id, value);
}

function requestDelete() {
  emit("requestDelete", props.product);
}
</script>

<template>
  <tr
    class="border-b hover:bg-gray-50 transition-colors"
    :class="{ 'bg-indigo-50': expandedProducts[product.id] }"
  >
    <!-- Toggle Lotes button -->
    <td class="p-4 text-center">
      <button
        @click="toggleProductLotes(product.id)"
        class="w-8 h-8 rounded-full flex items-center justify-center transition-all duration-300 hover:bg-indigo-100"
        :class="{
          'bg-indigo-100 shadow-md': expandedProducts[product.id],
          'hover:shadow': !expandedProducts[product.id],
        }"
        title="Clique para expandir/recolher informações de lotes"
      >
        <span
          class="material-icons-outlined text-lg text-indigo-600 transition-transform duration-300"
          :class="{ 'rotate-90': expandedProducts[product.id] }"
        >
          chevron_right
        </span>
      </button>
    </td>

    <!-- Product Name -->
    <td class="p-4">
      <input
        v-if="isEditMode"
        type="text"
        v-model="tempProductDetails[product.id].name"
        class="input-field-enhanced w-full"
      />
      <span v-else class="font-medium text-gray-700">{{ product.name }}</span>
    </td>

    <!-- Product Unit -->
    <td class="p-4">
      <select
        v-if="isEditMode"
        v-model="tempProductDetails[product.id].unit"
        class="input-field-enhanced w-full"
      >
        <option value="L">L</option>
        <option value="kg">kg</option>
      </select>
      <span
        v-else
        class="px-2 py-1 bg-gray-100 rounded text-sm font-medium text-gray-700"
      >
        {{ product.unit }}
      </span>
    </td>

    <!-- Product Quantity -->
    <td class="p-4">
      <div v-if="isEditMode && (!product.lotes || product.lotes.length === 0)">
        <ProductQuantityEditor
          :quantity="tempQuantities[product.id]"
          @update-quantity="updateProductQuantityDirectly"
          @change-quantity="changeProductQuantity"
        />
      </div>
      <div
        v-else-if="isEditMode && product.lotes && product.lotes.length > 0"
        class="flex flex-col"
      >
        <span class="font-bold text-lg text-indigo-700">{{
          getProductDisplayQuantity(product)
        }}</span>

        <!-- Base quantity editor for products with lotes -->
        <div class="mt-2 pt-2 border-t border-dashed border-gray-200">
          <div class="text-xs text-gray-500 mb-1">
            <span class="font-medium text-indigo-600"
              >Qtd. Base do Produto:</span
            >
          </div>
          <div class="flex items-center space-x-2">
            <input
              type="number"
              min="0"
              :value="product.quantity"
              @input="
                updateProductBaseQuantity(
                  parseFloat(($event.target as HTMLInputElement).value)
                )
              "
              class="w-20 py-1 px-2 text-sm border border-indigo-300 rounded focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
            <span class="text-xs text-indigo-600 font-medium">{{
              product.unit
            }}</span>
          </div>
        </div>
      </div>
      <div v-else class="font-bold text-lg text-indigo-700">
        {{ getProductDisplayQuantity(product) }}
      </div>

      <div
        v-if="
          product.lotes &&
          product.lotes.length > 0 &&
          product.quantity !== getProductDisplayQuantity(product) &&
          !isEditMode
        "
        class="mt-1 p-1 bg-yellow-50 border border-yellow-300 rounded text-xs text-yellow-800 flex items-center"
      >
        <span class="material-icons-outlined text-xs mr-1 text-yellow-600"
          >warning</span
        >
        <span>
          Total lotes: {{ getProductDisplayQuantity(product) }} {{ product.unit
          }}<br />
          Base: {{ product.quantity }} {{ product.unit }}
          <span v-if="!isEditMode" class="text-xs italic">
            (Ative edição para corrigir)
          </span>
        </span>
      </div>
    </td>

    <!-- Product Actions -->
    <td class="p-4">
      <button
        v-if="!isEditMode"
        @click="requestDelete"
        class="btn-danger-enhanced flex items-center justify-center"
        title="Excluir produto"
      >
        <span class="material-icons-outlined text-sm">delete</span>
        <span class="ml-1">Excluir</span>
      </button>
      <div v-else class="flex items-center">
        <span
          v-if="!expandedProducts[product.id]"
          class="text-sm text-indigo-600 font-medium flex items-center"
          title="Clique na seta à esquerda para gerenciar lotes"
        >
          <span class="material-icons-outlined text-sm mr-1 animate-pulse"
            >arrow_back</span
          >
          Clique na seta para ver lotes
        </span>
        <span
          v-else
          class="text-sm text-indigo-600 font-medium flex items-center"
        >
          <span class="material-icons-outlined text-sm mr-1 text-indigo-500"
            >inventory</span
          >
          Gerenciando lotes...
        </span>
      </div>
    </td>
  </tr>
</template>

<style scoped>
.input-field-enhanced {
  @apply px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow;
}
.btn-danger-enhanced {
  @apply px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm shadow-sm hover:shadow flex items-center font-medium;
}
</style>
