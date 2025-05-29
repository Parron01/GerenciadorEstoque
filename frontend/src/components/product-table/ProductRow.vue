<script setup lang="ts">
import { computed } from "vue";
import type { Product } from "@/models/product";

const props = defineProps<{
  product: Product;
  isEditMode: boolean;
  expandedProducts: Record<string, boolean>;
  tempProductDetails: Record<string, { name: string; unit: "L" | "kg" }>;
}>();

const emit = defineEmits<{
  (e: "requestDelete", product: Product): void;
  (e: "toggleProductLotes", productId: string): void;
}>();

const productHasLotes = computed(
  () => props.product.lotes && props.product.lotes.length > 0
);

// This function now always calculates quantity from lotes if they exist,
// or returns product.quantity (which should be 0 if it has no lotes and quantity is derived)
function getProductDisplayQuantity(product: Product): number {
  if (productHasLotes.value) {
    return product.lotes!.reduce((sum, lote) => sum + lote.quantity, 0);
  }
  // If a product has no lotes, its quantity should ideally be 0 according to the new rule.
  // The `product.quantity` field itself might be 0 from the backend.
  return product.quantity;
}

function toggleProductLotes() {
  emit("toggleProductLotes", props.product.id);
}

function requestDelete() {
  emit("requestDelete", props.product);
}
</script>

<template>
  <tr
    class="border-b border-gray-200 hover:bg-gray-50/50 transition-colors"
    :class="{ 'bg-indigo-50/30': expandedProducts[product.id] }"
  >
    <!-- Toggle Lotes button -->
    <td class="p-1 text-center">
      <button
        @click="toggleProductLotes"
        class="p-2 rounded-full hover:bg-indigo-100 text-indigo-600 transition-colors"
        :title="
          expandedProducts[product.id] ? 'Recolher lotes' : 'Expandir lotes'
        "
      >
        <span
          class="material-icons-outlined transition-transform duration-200"
          :class="{ 'rotate-180': expandedProducts[product.id] }"
        >
          expand_more
        </span>
      </button>
    </td>

    <!-- Product Name -->
    <td class="p-3">
      <input
        v-if="isEditMode"
        type="text"
        v-model="tempProductDetails[product.id].name"
        class="input-field-enhanced w-full"
      />
      <span v-else class="font-medium text-gray-800">{{ product.name }}</span>
    </td>

    <!-- Qtd. Total (New Order: 2nd data column) -->
    <td class="p-3 text-gray-700 text-center">
      {{ getProductDisplayQuantity(product) }}
    </td>

    <!-- Unidade (New Order: 3rd data column) -->
    <td class="p-3 text-center">
      <select
        v-if="isEditMode"
        v-model="tempProductDetails[product.id].unit"
        class="input-field-enhanced w-full"
      >
        <option value="L">L</option>
        <option value="kg">kg</option>
      </select>
      <span v-else>{{ product.unit }}</span>
    </td>

    <!-- Actions (New Order: 4th data column) -->
    <td class="p-3 text-center">
      <button
        @click="requestDelete"
        class="btn-delete-enhanced"
        :disabled="isEditMode"
        :class="{ 'opacity-50 cursor-not-allowed': isEditMode }"
        title="Excluir produto"
      >
        <span class="material-icons-outlined">delete</span>
      </button>
    </td>
  </tr>
</template>

<style scoped>
.input-field-enhanced {
  @apply px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow;
}
.btn-delete-enhanced {
  @apply px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm shadow-sm hover:shadow flex items-center font-medium;
}
</style>
