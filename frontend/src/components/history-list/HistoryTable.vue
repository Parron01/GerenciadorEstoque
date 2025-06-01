<script setup lang="ts">
import type {
  HistoryBatchGroup,
  ProductSummaryForBatch,
} from "@/models/history";
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
</script>

<template>
  <div class="hidden sm:block">
    <!-- Product summaries are now integrated into each product's expanded section -->

    <!-- Group records by product -->
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
          <tr
            class="border-b border-gray-200 hover:bg-gray-50 cursor-pointer"
            :class="{ 'bg-indigo-50': expandedProducts[productId] }"
            @click="toggleProduct(productId)"
          >
            <td class="p-3 text-center">
              <span
                class="material-icons-outlined text-indigo-600 transition-transform duration-200"
                :class="{ 'rotate-90': expandedProducts[productId] }"
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

              <div class="text-sm text-gray-600 mt-0.5">
                <template v-if="getProductCurrentQuantity(productId) !== null">
                  Quantidade atual:
                  <span class="font-medium">{{
                    getProductCurrentQuantity(productId)?.toFixed(2)
                  }}</span>
                </template>

                <!-- Show quantity change if available in product summary -->
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
                    class="ml-2 font-medium"
                  >
                    ({{
                      formatQuantityChange(
                        batch.productSummaries[productId]
                          .netQuantityChangeInBatch
                      )
                    }})
                  </span>
                </template>
              </div>
            </td>
            <td class="p-3 text-right">
              <div class="flex justify-end items-center space-x-1">
                <!-- Count of operations for this product -->
                <span
                  class="px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded-full"
                >
                  {{ recordsByProduct[productId].length }} operações
                </span>
              </div>
            </td>
          </tr>

          <!-- Expanded product details -->
          <template v-if="expandedProducts[productId]">
            <!-- Summary row if available -->
            <tr
              v-if="batch.productSummaries && batch.productSummaries[productId]"
              class="bg-indigo-50/60"
            >
              <td></td>
              <td colspan="2" class="p-3 border-b border-indigo-100">
                <div class="rounded-md bg-indigo-100/70 p-2 text-sm">
                  <div class="font-medium text-indigo-700 mb-1">
                    Resumo das alterações:
                  </div>
                  <div class="grid grid-cols-3 gap-2 text-gray-700">
                    <div>
                      <div class="text-xs text-gray-500">
                        Quantidade anterior:
                      </div>
                      <div class="font-medium">
                        {{
                          batch.productSummaries[
                            productId
                          ].totalQuantityBeforeBatch.toFixed(2)
                        }}
                      </div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500">Quantidade atual:</div>
                      <div class="font-medium">
                        {{
                          batch.productSummaries[
                            productId
                          ].totalQuantityAfterBatch.toFixed(2)
                        }}
                      </div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500">
                        Alteração líquida:
                      </div>
                      <div
                        class="font-medium"
                        :class="
                          getQuantityChangeClass(
                            batch.productSummaries[productId]
                              .netQuantityChangeInBatch
                          )
                        "
                      >
                        {{
                          formatQuantityChange(
                            batch.productSummaries[productId]
                              .netQuantityChangeInBatch
                          )
                        }}
                      </div>
                    </div>
                  </div>
                </div>
              </td>
            </tr>

            <!-- Individual operation rows -->
            <tr
              v-for="(record, idx) in recordsByProduct[productId]"
              :key="`${productId}-${idx}`"
              class="bg-gray-50/60 border-b border-gray-100 last:border-b-0"
            >
              <td></td>
              <td colspan="2" class="pl-8 pr-3 py-2">
                <div class="flex items-start">
                  <!-- Record icon based on type -->
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
                    <!-- Record header -->
                    <div class="flex items-center justify-between mb-1">
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
                          class="ml-2 px-1.5 py-0.5 text-xs rounded-full bg-gray-100 text-gray-600"
                        >
                          {{
                            record.entityType === "product" ? "Produto" : "Lote"
                          }}
                        </span>

                        <!-- Lote ID if applicable -->
                        <template v-if="record.entityType === 'lote'">
                          <span class="ml-1 text-xs text-gray-500">
                            {{ formatId(record.entityId) }}
                          </span>
                        </template>
                      </div>
                    </div>

                    <!-- Record details -->

                    <!-- Quantity changes -->
                    <div
                      v-if="
                        record.details?.quantityBefore !== undefined ||
                        record.details?.quantityAfter !== undefined
                      "
                      class="mt-1 rounded bg-gray-100/80 px-2.5 py-1.5 text-sm flex items-center"
                    >
                      <span
                        class="material-icons-outlined text-amber-500 text-sm mr-1"
                        >inventory</span
                      >
                      <span class="text-gray-600">Quantidade:</span>
                      <span class="ml-1 mr-1 font-medium">{{
                        record.details?.quantityBefore ?? "N/A"
                      }}</span>
                      <span
                        class="material-icons-outlined text-gray-400 mx-1 text-sm"
                        >arrow_forward</span
                      >
                      <span class="font-medium">{{
                        record.details?.quantityAfter ?? "N/A"
                      }}</span>

                      <!-- Change amount if available -->
                      <template
                        v-if="record.details?.quantityChanged !== undefined"
                      >
                        <span
                          class="ml-2 px-2 py-0.5 rounded-full text-xs font-medium"
                          :class="
                            getQuantityChangeClass(
                              record.details.quantityChanged
                            )
                          "
                        >
                          {{
                            formatQuantityChange(record.details.quantityChanged)
                          }}
                        </span>
                      </template>
                    </div>

                    <!-- Changed fields -->
                    <div
                      v-if="record.details?.changedFields?.length"
                      class="mt-1.5 space-y-1"
                    >
                      <div
                        v-for="(field, fieldIdx) in record.details
                          .changedFields"
                        :key="fieldIdx"
                        class="bg-gray-100/80 px-2.5 py-1.5 rounded text-sm flex items-center"
                      >
                        <!-- Field-specific icon -->
                        <span
                          class="material-icons-outlined text-sm mr-1.5"
                          :class="getFieldChangeType(field).colorClass"
                        >
                          {{ getFieldChangeType(field).icon }}
                        </span>

                        <span class="capitalize text-gray-600">
                          {{ field.field.replace("_", " ") }}:
                        </span>

                        <template
                          v-if="
                            field.oldValue !== undefined ||
                            field.newValue !== undefined
                          "
                        >
                          <span
                            class="text-gray-500 ml-1"
                            v-if="field.oldValue !== undefined"
                          >
                            {{ field.oldValue || "N/A" }}
                          </span>
                          <span
                            class="material-icons-outlined text-gray-400 mx-1 text-xs"
                            >arrow_forward</span
                          >
                          <span
                            class="text-gray-900 font-medium"
                            v-if="field.newValue !== undefined"
                          >
                            {{ field.newValue }}
                          </span>
                        </template>
                        <span
                          class="text-gray-900 font-medium ml-1"
                          v-else-if="field.newValue !== undefined"
                        >
                          {{ field.newValue }}
                        </span>
                      </div>
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
                          class="bg-gray-100/80 px-2.5 py-1.5 rounded text-sm flex items-center"
                        >
                          <span
                            class="material-icons-outlined text-green-600 text-sm mr-1.5"
                            >event</span
                          >
                          <span class="text-gray-600">Validade:</span>

                          <template
                            v-if="
                              record.details.dataValidadeOld ||
                              record.details.dataValidadeNew
                            "
                          >
                            <span class="text-gray-500 ml-1">
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
                            class="text-gray-900 font-medium ml-1"
                            v-else-if="record.details.dataValidade"
                          >
                            {{ formatDateOnly(record.details.dataValidade) }}
                          </span>
                        </div>
                      </template>
                    </div>

                    <!-- Special tags -->
                    <div
                      v-if="
                        record.details?.isNewProduct ||
                        record.details?.isProductRemoval
                      "
                      class="mt-1.5 flex gap-1"
                    >
                      <span
                        v-if="record.details.isNewProduct"
                        class="inline-flex items-center px-2 py-0.5 text-xs font-medium bg-emerald-100 text-emerald-800 rounded-full"
                      >
                        <span class="material-icons-outlined text-xs mr-0.5"
                          >add_circle</span
                        >
                        Novo
                      </span>
                      <span
                        v-if="record.details.isProductRemoval"
                        class="inline-flex items-center px-2 py-0.5 text-xs font-medium bg-red-100 text-red-800 rounded-full"
                      >
                        <span class="material-icons-outlined text-xs mr-0.5"
                          >delete</span
                        >
                        Removido
                      </span>
                    </div>
                  </div>
                </div>
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
</style>
