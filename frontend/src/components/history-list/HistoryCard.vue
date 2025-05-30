<script setup lang="ts">
import type { ProductHistory, ProductChange } from "@/models/product";
import type {
  ParsedHistoryRecord,
  HistoryBatchGroup,
  ProductSummaryForBatch,
} from "@/models/history";
import type { LoteChangeDetails } from "@/models/lote";
import { ref, computed } from "vue";
import {
  formatActionName,
  formatId,
  getActionColorClass,
  getActionBadgeClass,
  getQuantityChangeClass,
  formatQuantityChange,
  formatDateTime,
  formatDateOnly, // Import the new function
} from "@/utils/formatters";

const props = defineProps<{
  batch:
    | HistoryBatchGroup
    | {
        batchId: string;
        createdAt: string;
        records: any[];
        productSummaries?: Record<string, ProductSummaryForBatch>;
      };
  isLocalMode: boolean;
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
        <!-- Product Header -->
        <div
          class="p-3 flex items-center justify-between cursor-pointer"
          :class="{
            'bg-indigo-50 border-b border-indigo-100':
              expandedProducts[productId],
          }"
          @click="toggleProduct(productId)"
        >
          <div class="flex-grow">
            <div class="font-medium text-indigo-700 flex items-center">
              <span class="material-icons-outlined mr-1.5 text-indigo-500"
                >inventory_2</span
              >
              {{ getProductName(productId) }}
            </div>

            <div class="text-xs text-gray-500 mt-0.5 flex items-center">
              ID: {{ formatId(productId) }}

              <template v-if="getProductCurrentQuantity(productId) !== null">
                <span class="mx-1">•</span>
                <span>
                  Qtd: {{ getProductCurrentQuantity(productId)?.toFixed(2) }}
                </span>

                <!-- Show quantity change if available -->
                <template
                  v-if="
                    batch.productSummaries &&
                    batch.productSummaries[productId] &&
                    batch.productSummaries[productId]
                      .netQuantityChangeInBatch !== 0
                  "
                >
                  <span
                    :class="
                      getQuantityChangeClass(
                        batch.productSummaries[productId]
                          .netQuantityChangeInBatch
                      )
                    "
                    class="ml-1 font-medium"
                  >
                    ({{
                      formatQuantityChange(
                        batch.productSummaries[productId]
                          .netQuantityChangeInBatch
                      )
                    }})
                  </span>
                </template>
              </template>
            </div>
          </div>

          <div class="flex items-center">
            <span
              class="text-xs px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full mr-2"
            >
              {{ recordsByProduct[productId].length }}
            </span>
            <span
              class="material-icons-outlined text-indigo-600 transition-transform duration-200"
              :class="{ 'rotate-180': expandedProducts[productId] }"
            >
              expand_more
            </span>
          </div>
        </div>

        <!-- Product Details (expanded) -->
        <div v-if="expandedProducts[productId]" class="p-3 pt-0 space-y-3">
          <!-- Product summary if available -->
          <div
            v-if="batch.productSummaries && batch.productSummaries[productId]"
            class="p-2 bg-indigo-50 rounded-md mt-3 text-xs"
          >
            <div class="font-medium text-indigo-700 mb-1">
              Resumo das alterações
            </div>
            <div class="grid grid-cols-3 gap-1">
              <div class="px-1">
                <div class="text-gray-500">Anterior:</div>
                <div class="font-medium">
                  {{
                    batch.productSummaries[
                      productId
                    ].totalQuantityBeforeBatch.toFixed(2)
                  }}
                </div>
              </div>
              <div class="px-1">
                <div class="text-gray-500">Atual:</div>
                <div class="font-medium">
                  {{
                    batch.productSummaries[
                      productId
                    ].totalQuantityAfterBatch.toFixed(2)
                  }}
                </div>
              </div>
              <div class="px-1">
                <div class="text-gray-500">Alteração:</div>
                <div
                  class="font-medium"
                  :class="
                    getQuantityChangeClass(
                      batch.productSummaries[productId].netQuantityChangeInBatch
                    )
                  "
                >
                  {{
                    formatQuantityChange(
                      batch.productSummaries[productId].netQuantityChangeInBatch
                    )
                  }}
                </div>
              </div>
            </div>
          </div>

          <!-- Individual operations -->
          <div
            v-for="(record, idx) in recordsByProduct[productId]"
            :key="`${productId}-${idx}`"
            class="border-b border-gray-100 last:border-b-0 py-2"
          >
            <!-- Operation header -->
            <div class="flex items-center justify-between">
              <div
                class="px-2 py-0.5 rounded-full text-xs font-medium"
                :class="getActionBadgeClass(record.details?.action || '')"
              >
                {{ formatActionName(record.details?.action || "Alteração") }}
              </div>

              <div class="text-xs text-gray-500">
                {{ record.entityType === "product" ? "Produto" : "Lote" }}
                <template v-if="record.entityType === 'lote'">
                  {{ formatId(record.entityId) }}
                </template>
              </div>
            </div>

            <!-- Quantity changes -->
            <div
              v-if="
                record.details?.quantityBefore !== undefined ||
                record.details?.quantityAfter !== undefined
              "
              class="mt-1.5 bg-gray-50 p-2 rounded flex items-center text-sm"
            >
              <span class="material-icons-outlined text-amber-500 text-sm mr-1">
                inventory
              </span>
              <span class="text-xs text-gray-500">Qtd:</span>
              <span class="ml-1 text-xs">{{
                record.details?.quantityBefore ?? "N/A"
              }}</span>
              <span class="material-icons-outlined text-gray-400 mx-1 text-xs">
                arrow_forward
              </span>
              <span class="text-xs font-medium">{{
                record.details?.quantityAfter ?? "N/A"
              }}</span>

              <template v-if="record.details?.quantityChanged !== undefined">
                <span
                  class="ml-auto text-xs font-medium"
                  :class="
                    getQuantityChangeClass(record.details.quantityChanged)
                  "
                >
                  {{ formatQuantityChange(record.details.quantityChanged) }}
                </span>
              </template>
            </div>

            <!-- Changed fields -->
            <div
              v-if="record.details?.changedFields?.length"
              class="mt-1.5 space-y-1"
            >
              <div
                v-for="(field, fieldIdx) in record.details.changedFields"
                :key="fieldIdx"
                class="bg-gray-50 p-2 rounded"
              >
                <div class="flex items-center justify-between text-xs">
                  <span class="capitalize font-medium">{{
                    field.field.replace("_", " ")
                  }}</span>
                </div>

                <div class="mt-1 text-sm flex items-center">
                  <template v-if="field.oldValue !== undefined">
                    <span class="text-xs text-gray-500">{{
                      field.oldValue
                    }}</span>
                    <span
                      class="material-icons-outlined text-gray-400 mx-1 text-xs"
                    >
                      arrow_forward
                    </span>
                  </template>
                  <span class="text-xs font-medium">{{ field.newValue }}</span>
                </div>
              </div>
            </div>

            <!-- Lote details -->
            <div
              v-if="
                record.entityType === 'lote' &&
                (record.details?.dataValidadeOld ||
                  record.details?.dataValidadeNew ||
                  record.details?.dataValidade)
              "
              class="mt-1.5 bg-gray-50 p-2 rounded"
            >
              <div class="flex items-center justify-between text-xs">
                <span class="font-medium">Validade</span>
              </div>

              <div class="mt-1 text-sm flex items-center">
                <template v-if="record.details.dataValidadeOld">
                  <span class="text-xs text-gray-500">{{
                    formatDateOnly(record.details.dataValidadeOld)
                  }}</span>
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

            <!-- Special tags -->
            <div
              v-if="
                record.details?.isNewProduct || record.details?.isProductRemoval
              "
              class="flex gap-1 mt-1.5"
            >
              <span
                v-if="record.details.isNewProduct"
                class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium bg-emerald-100 text-emerald-800 rounded-full"
              >
                <span class="material-icons-outlined text-xs mr-0.5"
                  >add_circle</span
                >
                Novo
              </span>
              <span
                v-if="record.details.isProductRemoval"
                class="inline-flex items-center px-1.5 py-0.5 text-xs font-medium bg-red-100 text-red-800 rounded-full"
              >
                <span class="material-icons-outlined text-xs mr-0.5"
                  >delete</span
                >
                Removido
              </span>
            </div>
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
