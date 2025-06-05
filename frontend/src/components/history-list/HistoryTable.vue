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
  getActionColorClass,
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

// Helper to determine what kind of field change this is
function getFieldChangeType(field: any): { icon: string; colorClass: string } {
  if (!field || !field.field)
    return { icon: "error", colorClass: "text-gray-500" };

  const fieldName = field.field.toLowerCase();

  if (fieldName.includes("name")) {
    return { icon: "edit", colorClass: "text-blue-600" };
  }
  if (fieldName.includes("unit")) {
    return { icon: "straighten", colorClass: "text-purple-600" };
  }
  if (fieldName.includes("quantity")) {
    return { icon: "inventory", colorClass: "text-amber-600" };
  }
  if (
    fieldName.includes("data") ||
    fieldName.includes("date") ||
    fieldName.includes("validade")
  ) {
    return { icon: "event", colorClass: "text-green-600" };
  }

  return { icon: "settings", colorClass: "text-gray-600" };
}

// Get icon for record type
function getRecordTypeIcon(record: any): string {
  if (!record || !record.entityType) return "question_mark";

  if (record.entityType === "product") {
    return "inventory_2";
  }
  if (record.entityType === "lote") {
    return "inventory";
  }

  return "article";
}

function getChangeIcon(action: string): string {
  if (!action) return "help";

  const actionLower = action.toLowerCase();
  if (actionLower.includes("creat") || actionLower.includes("add")) {
    return "add_circle";
  }
  if (actionLower.includes("delet") || actionLower.includes("remov")) {
    return "delete";
  }
  if (actionLower.includes("updat")) {
    return "edit";
  }

  return "change_circle";
}

// Helper function to format N/A values as zero
function formatQuantityValue(value: any): string {
  if (value === undefined || value === null || value === "N/A") {
    return "0";
  }
  return value.toString();
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
</script>

<template>
  <div class="hidden sm:block">
    <!-- Product summaries are now integrated into each product's row -->
    <table class="min-w-full bg-white">
      <thead class="bg-slate-100 text-left text-gray-700">
        <tr>
          <th class="p-3 w-16"></th>
          <th class="p-3">Produto</th>
          <th class="p-3 text-right">Alterações</th>
        </tr>
      </thead>
      <tbody>
        <!-- Loop through product groups instead of individual records -->
        <template v-for="productId in productIds" :key="productId">
          <!-- Product header row (always visible) -->
          <tr class="border-b border-gray-200 hover:bg-gray-50">
            <td class="p-3 text-center">
              <span
                class="material-icons-outlined text-indigo-600 transition-transform duration-200 cursor-pointer"
                :class="{ 'rotate-90': expandedProducts[productId] }"
                @click="toggleProduct(productId)"
              >
                chevron_right
              </span>
            </td>
            <td class="p-3">
              <div class="font-medium flex items-center text-indigo-700">
                <span class="material-icons-outlined mr-1.5 text-indigo-500"
                  >inventory_2</span
                >
                {{ getProductName(productId) }}
                <span class="text-xs text-gray-500 ml-2"
                  >(ID: {{ formatId(productId) }})</span
                >
              </div>

              <!-- Summary Info - Always visible -->
              <div
                class="mt-2 p-2 bg-indigo-50 rounded-md border border-indigo-100"
              >
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-2">
                  <!-- Quantity change -->
                  <div
                    v-if="
                      batch.productSummaries &&
                      batch.productSummaries[productId]
                    "
                    class="flex items-center justify-between sm:justify-start sm:gap-2"
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
                        >arrow_forward</span
                      >
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
                            batch.productSummaries[productId]
                              .netQuantityChangeInBatch
                          )
                        "
                        class="ml-1 w-7 h-7 rounded-full flex items-center justify-center text-xs font-medium"
                      >
                        {{
                          formatQuantityChange(
                            batch.productSummaries[productId]
                              .netQuantityChangeInBatch
                          )
                        }}
                      </span>
                    </div>
                  </div>

                  <!-- Product Updates - Always display changes to product details -->
                  <template
                    v-for="(record, idx) in getProductRecordsForProduct(
                      productId
                    )"
                    :key="idx"
                  >
                    <!-- Show product details changes upfront -->
                    <template v-if="record.details?.changedFields?.length">
                      <div
                        v-for="(field, fidx) in record.details.changedFields"
                        :key="`${idx}-${fidx}`"
                        class="flex items-center justify-between sm:justify-start sm:gap-2"
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
                </div>

                <!-- Special tags -->
                <div
                  v-if="
                    recordsByProduct[productId].some(
                      (r) =>
                        r.details?.isNewProduct || r.details?.isProductRemoval
                    )
                  "
                  class="flex gap-1 mt-2"
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
            </td>
            <td class="p-3 text-right align-top">
              <div class="flex justify-end items-center space-x-1">
                <!-- Count of lote operations for this product -->
                <button
                  class="px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded-full hover:bg-indigo-100 transition-colors"
                  @click="toggleProduct(productId)"
                >
                  {{ getLoteRecordsForProduct(productId).length }} operações de
                  lote
                  <span
                    class="material-icons-outlined text-gray-500 text-xs ml-0.5"
                    :class="{ 'rotate-180': expandedProducts[productId] }"
                  >
                    expand_more
                  </span>
                </button>
              </div>
            </td>
          </tr>

          <!-- Expanded product details - Only show lote operations -->
          <template v-if="expandedProducts[productId]">
            <!-- Individual operation rows -->
            <tr
              v-for="(record, idx) in getLoteRecordsForProduct(productId)"
              :key="`${productId}-${idx}`"
              class="bg-gray-50/60 border-b border-gray-100 last:border-b-0"
            >
              <td></td>
              <td colspan="2" class="pl-8 pr-3 py-2">
                <div class="flex items-start">
                  <!-- Lote operation icon - More salient -->
                  <div
                    class="w-10 h-10 rounded-full flex items-center justify-center mr-3 mt-0.5"
                    :class="{
                      'bg-emerald-100':
                        record.details?.action?.includes('creat'),
                      'bg-red-100': record.details?.action?.includes('delet'),
                      'bg-amber-100': record.details?.action?.includes('updat'),
                      'bg-gray-100': !record.details?.action,
                    }"
                  >
                    <span
                      class="material-icons-outlined text-sm"
                      :class="{
                        'text-emerald-600':
                          record.details?.action?.includes('creat'),
                        'text-red-600':
                          record.details?.action?.includes('delet'),
                        'text-amber-600':
                          record.details?.action?.includes('updat'),
                        'text-gray-600': !record.details?.action,
                      }"
                    >
                      {{ getChangeIcon(record.details?.action) }}
                    </span>
                  </div>

                  <!-- Record content -->
                  <div class="flex-grow">
                    <!-- Lote operation header - More prominent -->
                    <div
                      class="flex items-center justify-between mb-2 bg-gray-100/80 p-2 rounded"
                    >
                      <div class="flex items-center">
                        <span
                          class="material-icons-outlined text-sm mr-1.5"
                          :class="
                            getActionColorClass(record.details?.action || '')
                          "
                        >
                          {{ getRecordTypeIcon(record) }}
                        </span>
                        <span
                          class="font-medium"
                          :class="
                            getActionColorClass(record.details?.action || '')
                          "
                        >
                          {{
                            formatActionName(
                              record.details?.action || "Alteração"
                            )
                          }}
                        </span>

                        <!-- Entity type badge -->
                        <span
                          class="ml-2 px-2 py-0.5 text-xs rounded-full bg-gray-200 text-gray-700"
                        >
                          Lote {{ formatId(record.entityId) }}
                        </span>
                      </div>
                    </div>

                    <!-- Record details -->

                    <!-- Quantity changes -->
                    <div
                      v-if="
                        record.details?.quantityBefore !== undefined ||
                        record.details?.quantityAfter !== undefined
                      "
                      class="mt-1.5 bg-gray-100/80 px-3 py-2 rounded flex items-center text-sm"
                    >
                      <span
                        class="material-icons-outlined text-amber-500 text-sm mr-2"
                        >inventory</span
                      >
                      <span class="text-gray-700 font-medium">Quantidade:</span>
                      <span class="ml-2 mr-2">
                        {{
                          formatQuantityValue(record.details?.quantityBefore)
                        }}
                      </span>
                      <span
                        class="material-icons-outlined text-gray-400 mx-1 text-sm"
                        >arrow_forward</span
                      >
                      <span class="font-medium">
                        {{ formatQuantityValue(record.details?.quantityAfter) }}
                      </span>

                      <!-- Move change amount indicator right after quantity info -->
                      <span
                        class="ml-2 rounded-full w-7 h-7 flex items-center justify-center text-xs font-medium"
                        :class="
                          getQuantityChangeClass(getQuantityDifference(record))
                        "
                      >
                        {{
                          formatQuantityChange(getQuantityDifference(record))
                        }}
                      </span>
                    </div>

                    <!-- Lote details for lotes -->
                    <div v-if="record.entityType === 'lote'" class="mt-1.5">
                      <!-- Data validade if present -->
                      <template
                        v-if="
                          record.details?.dataValidadeOld ||
                          record.details?.dataValidadeNew ||
                          record.details?.dataValidade
                        "
                      >
                        <div
                          class="bg-gray-100/80 px-3 py-2 rounded text-sm flex items-center"
                        >
                          <span
                            class="material-icons-outlined text-green-600 text-sm mr-2"
                            >event</span
                          >
                          <span class="text-gray-700 font-medium"
                            >Validade:</span
                          >

                          <template
                            v-if="
                              record.details.dataValidadeOld ||
                              record.details.dataValidadeNew
                            "
                          >
                            <span class="text-gray-500 ml-2">
                              {{
                                formatDateOnly(
                                  record.details.dataValidadeOld
                                ) || "N/A"
                              }}
                            </span>
                            <span
                              class="material-icons-outlined text-gray-400 mx-1 text-xs"
                              >arrow_forward</span
                            >
                            <span class="text-gray-900 font-medium">
                              {{
                                formatDateOnly(
                                  record.details.dataValidadeNew ||
                                    record.details.dataValidade
                                ) || "N/A"
                              }}
                            </span>
                          </template>
                          <span
                            class="text-gray-900 font-medium ml-2"
                            v-else-if="record.details.dataValidade"
                          >
                            {{ formatDateOnly(record.details.dataValidade) }}
                          </span>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </td>
            </tr>

            <!-- Empty lote operations message -->
            <tr v-if="getLoteRecordsForProduct(productId).length === 0">
              <td></td>
              <td colspan="2" class="text-center py-3 text-gray-500">
                <span class="text-sm"
                  >Não há operações de lotes neste registro</span
                >
              </td>
            </tr>
          </template>
        </template>

        <!-- Empty state -->
        <tr v-if="!productIds || productIds.length === 0">
          <td colspan="3" class="p-6 text-center text-gray-500">
            Nenhuma alteração encontrada neste lote.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.product-row {
  @apply cursor-pointer transition-colors duration-150 hover:bg-gray-50;
}
.product-row.expanded {
  @apply bg-indigo-50;
}

/* Additional styles for quantity changes */
.quantity-change-positive {
  @apply bg-emerald-100 text-emerald-800;
}
.quantity-change-negative {
  @apply bg-red-100 text-red-800;
}
</style>
