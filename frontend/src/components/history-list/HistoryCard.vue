<script setup lang="ts">
import type {
  HistoryBatchGroup,
  ProductSummaryForBatch,
} from "@/models/history";
import type { ChangedField } from "@/models/product";
import { ref, computed } from "vue";
import {
  formatActionName,
  formatId,
  getActionBadgeClass,
  getQuantityChangeClass,
  formatQuantityChange,
  formatDateOnly,
} from "@/utils/formatters";

const props = defineProps<{
  batch: HistoryBatchGroup;
}>();

// Track which products are expanded
const expandedProducts = ref<Record<string, boolean>>({});

// Toggle product expansion
function toggleProduct(productId: string) {
  expandedProducts.value[productId] = !expandedProducts.value[productId];
}

// Group records by product for better organization
const recordsByProduct = computed(() => {
  const grouped: Record<string, any[]> = {};

  if (!props.batch.records) return grouped;

  // First group by product
  props.batch.records.forEach((record) => {
    let productId;

    if (record.entityType === "product") {
      productId = record.entityId;
    } else if (record.entityType === "lote" && record.details) {
      // For lote records, get the product ID from details or context
      productId =
        record.details.productId ||
        (record.productNameContext ? record.entityId : "unknown");
    }

    if (!productId) return;

    if (!grouped[productId]) {
      grouped[productId] = [];
    }

    grouped[productId].push(record);
  });

  return grouped;
});

// Get all unique product IDs from the batch
const productIds = computed(() => {
  return Object.keys(recordsByProduct.value);
});

// Helper function to safely access product name
function getProductName(productId: string): string {
  // First check product summaries
  if (props.batch.productSummaries && props.batch.productSummaries[productId]) {
    return props.batch.productSummaries[productId].productName;
  }

  // Otherwise try to find it in records
  const records = recordsByProduct.value[productId] || [];
  for (const record of records) {
    if (record.entityType === "product" && record.details?.productName) {
      return record.details.productName;
    }
    if (record.productNameContext) {
      return record.productNameContext;
    }
  }

  return `Produto ${formatId(productId)}`;
}

// Helper function to get current product quantity
function getProductCurrentQuantity(productId: string): number | null {
  // First check product summaries for after-batch quantity
  if (props.batch.productSummaries && props.batch.productSummaries[productId]) {
    return props.batch.productSummaries[productId].totalQuantityAfterBatch;
  }

  // Otherwise try to find it in records
  const records = recordsByProduct.value[productId] || [];
  for (const record of records) {
    if (record.productCurrentTotalQuantity !== undefined) {
      return record.productCurrentTotalQuantity;
    }
  }

  return null;
}

// Helper function to format N/A values as zero
function formatQuantityValue(value: any): string {
  if (value === undefined || value === null || value === "N/A") {
    return "0";
  }
  return value.toString();
}

// Helper function to calculate quantity difference
function calculateQuantityDifference(before: any, after: any): number {
  // Convert both values to numbers, treating undefined, null, N/A as 0
  const beforeNum =
    before === undefined || before === null || before === "N/A"
      ? 0
      : parseFloat(before.toString());

  const afterNum =
    after === undefined || after === null || after === "N/A"
      ? 0
      : parseFloat(after.toString());

  return afterNum - beforeNum;
}

// Helper function to get the difference in quantities, either from quantityChanged or calculated
function getQuantityDifference(record: any): number {
  if (record.details?.quantityChanged !== undefined) {
    return record.details.quantityChanged;
  }

  const before = record.details?.quantityBefore;
  const after = record.details?.quantityAfter;
  return calculateQuantityDifference(before, after);
}

// Get lote-specific records for a product
const getLoteRecordsForProduct = (productId: string) => {
  return (
    recordsByProduct.value[productId]?.filter((r) => r.entityType === "lote") ||
    []
  );
};

// Get product-specific records for a product
const getProductRecordsForProduct = (productId: string) => {
  return (
    recordsByProduct.value[productId]?.filter(
      (r) => r.entityType === "product"
    ) || []
  );
};
</script>

<template>
  <div class="sm:hidden">
    <!-- Group by product in mobile view too -->
    <div class="space-y-4">
      <div
        v-for="productId in productIds"
        :key="productId"
        class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
      >
        <!-- Product Header - Now always shows summary changes -->
        <div class="p-3">
          <div class="flex items-center justify-between">
            <div class="flex-grow">
              <div class="font-medium text-indigo-700 flex items-center">
                <span class="material-icons-outlined mr-1.5 text-indigo-500"
                  >inventory_2</span
                >
                {{ getProductName(productId) }}
              </div>

              <div class="text-xs text-gray-500 mt-0.5 flex items-center">
                ID: {{ formatId(productId) }}
              </div>
            </div>

            <div
              class="flex items-center cursor-pointer"
              @click="toggleProduct(productId)"
            >
              <span
                class="text-xs px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full mr-2"
              >
                {{ getLoteRecordsForProduct(productId).length }} lotes
              </span>
              <span
                class="material-icons-outlined text-indigo-600 transition-transform duration-200"
                :class="{ 'rotate-180': expandedProducts[productId] }"
              >
                expand_more
              </span>
            </div>
          </div>

          <!-- Product Summary - Always visible -->
          <div
            class="mt-2 p-3 bg-indigo-50 rounded-md border border-indigo-100"
          >
            <!-- Quantity change -->
            <div
              v-if="batch.productSummaries && batch.productSummaries[productId]"
              class="flex items-center justify-between mb-1.5"
            >
              <div class="text-sm text-gray-700">Quantidade:</div>
              <div class="flex items-center">
                <span class="text-sm text-gray-600">
                  {{
                    batch.productSummaries[
                      productId
                    ].totalQuantityBeforeBatch.toFixed(2)
                  }}
                </span>
                <span
                  class="material-icons-outlined text-gray-400 mx-1 text-xs"
                >
                  arrow_forward
                </span>
                <span class="text-sm font-medium">
                  {{
                    batch.productSummaries[
                      productId
                    ].totalQuantityAfterBatch.toFixed(2)
                  }}
                </span>
                <span
                  v-if="
                    batch.productSummaries[productId]
                      .netQuantityChangeInBatch !== 0
                  "
                  :class="
                    getQuantityChangeClass(
                      batch.productSummaries[productId].netQuantityChangeInBatch
                    )
                  "
                  class="ml-2 w-8 h-8 rounded-full flex items-center justify-center text-xs font-medium"
                >
                  {{
                    formatQuantityChange(
                      batch.productSummaries[productId].netQuantityChangeInBatch
                    )
                  }}
                </span>
              </div>
            </div>

            <!-- Product Updates - Always display changes to product details -->
            <template
              v-for="(record, idx) in getProductRecordsForProduct(productId)"
              :key="idx"
            >
              <!-- Show product details changes upfront -->
              <template v-if="record.details?.changedFields?.length">
                <div
                  v-for="(field, fidx) in record.details.changedFields"
                  :key="`${idx}-${fidx}`"
                  class="flex items-center justify-between mb-1.5"
                >
                  <div class="text-sm text-gray-700 capitalize">
                    {{ field.field.replace("_", " ") }}:
                  </div>
                  <div class="flex items-center">
                    <span class="text-sm text-gray-600">
                      {{ field.oldValue || "0" }}
                    </span>
                    <span
                      class="material-icons-outlined text-gray-400 mx-1 text-xs"
                    >
                      arrow_forward
                    </span>
                    <span class="text-sm font-medium">
                      {{ field.newValue || "0" }}
                    </span>
                  </div>
                </div>
              </template>
            </template>

            <!-- Special tags -->
            <div
              v-if="
                recordsByProduct[productId].some(
                  (r) => r.details?.isNewProduct || r.details?.isProductRemoval
                )
              "
              class="flex gap-1 mt-1"
            >
              <span
                v-if="
                  recordsByProduct[productId].some(
                    (r) => r.details?.isNewProduct
                  )
                "
                class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium bg-emerald-100 text-emerald-800 rounded-full"
              >
                <span class="material-icons-outlined text-xs mr-0.5"
                  >add_circle</span
                >
                Novo produto
              </span>
              <span
                v-if="
                  recordsByProduct[productId].some(
                    (r) => r.details?.isProductRemoval
                  )
                "
                class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium bg-red-100 text-red-800 rounded-full"
              >
                <span class="material-icons-outlined text-xs mr-0.5"
                  >delete</span
                >
                Produto removido
              </span>
            </div>
          </div>
        </div>

        <!-- Product Details (expanded) - Now only show Lote operations -->
        <div v-if="expandedProducts[productId]" class="p-3 pt-0 space-y-3">
          <!-- Individual lote operations -->
          <div
            v-for="(record, idx) in getLoteRecordsForProduct(productId)"
            :key="`${productId}-${idx}`"
            class="border-t border-gray-100 pt-2 mt-2 first:border-t-0 first:pt-0 first:mt-0"
          >
            <!-- Lote operation header - Made more prominent -->
            <div
              class="flex items-center justify-between bg-gray-50 p-2 rounded"
            >
              <div
                class="px-2 py-0.5 rounded-full text-xs font-medium flex items-center"
                :class="getActionBadgeClass(record.details?.action || '')"
              >
                <span class="material-icons-outlined mr-1 text-xs">
                  {{
                    record.details?.action?.includes("creat")
                      ? "add_circle"
                      : record.details?.action?.includes("delet")
                        ? "delete"
                        : "edit"
                  }}
                </span>
                {{ formatActionName(record.details?.action || "Alteração") }}
              </div>

              <div class="text-xs bg-gray-200 px-2 py-0.5 rounded-full">
                Lote {{ formatId(record.entityId) }}
              </div>
            </div>

            <!-- Record details -->

            <!-- Quantity changes -->
            <div
              v-if="
                record.details?.quantityBefore !== undefined ||
                record.details?.quantityAfter !== undefined
              "
              class="mt-1.5 bg-gray-50/80 p-2 rounded flex items-center text-sm"
            >
              <span class="material-icons-outlined text-amber-500 text-sm mr-1">
                inventory
              </span>
              <span class="text-xs text-gray-700 font-medium">Quantidade:</span>
              <span class="ml-1 mr-1">
                {{ formatQuantityValue(record.details?.quantityBefore) }}
              </span>
              <span class="material-icons-outlined text-gray-400 mx-1 text-xs">
                arrow_forward
              </span>
              <span class="font-medium">
                {{ formatQuantityValue(record.details?.quantityAfter) }}
              </span>

              <!-- Move change amount indicator right after quantity info -->
              <span
                class="ml-2 rounded-full w-6 h-6 flex items-center justify-center text-xs font-medium"
                :class="getQuantityChangeClass(getQuantityDifference(record))"
              >
                {{ formatQuantityChange(getQuantityDifference(record)) }}
              </span>
            </div>

            <!-- Lote details -->
            <div
              v-if="
                record.entityType === 'lote' &&
                (record.details?.dataValidadeOld ||
                  record.details?.dataValidadeNew ||
                  record.details?.dataValidade)
              "
              class="mt-1.5 bg-gray-50/80 p-2 rounded"
            >
              <div class="flex items-center text-xs">
                <span
                  class="material-icons-outlined text-green-600 text-sm mr-1"
                  >event</span
                >
                <span class="font-medium text-gray-700">Validade:</span>
              </div>

              <div class="mt-1 text-sm flex items-center">
                <template v-if="record.details.dataValidadeOld">
                  <span class="text-xs text-gray-500">
                    {{ formatDateOnly(record.details.dataValidadeOld) }}
                  </span>
                  <span
                    class="material-icons-outlined text-gray-400 mx-1 text-xs"
                  >
                    arrow_forward
                  </span>
                </template>
                <span class="text-xs font-medium">
                  {{
                    formatDateOnly(
                      record.details.dataValidadeNew ||
                        record.details.dataValidade
                    )
                  }}
                </span>
              </div>
            </div>
          </div>

          <!-- Empty lote operations message -->
          <div
            v-if="getLoteRecordsForProduct(productId).length === 0"
            class="text-center text-gray-500 p-3"
          >
            <span class="text-sm"
              >Não há operações de lotes neste registro</span
            >
          </div>
        </div>
      </div>

      <!-- Empty state -->
      <div
        v-if="!productIds || productIds.length === 0"
        class="p-4 text-center text-gray-500 bg-white rounded-lg shadow-sm border border-gray-200"
      >
        Nenhuma alteração encontrada neste lote.
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Utility classes specific to this component */
.tag-base {
  @apply px-2 py-0.5 rounded-full text-xs font-medium;
}
</style>
