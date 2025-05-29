<script setup lang="ts">
import type { ProductHistory, ProductChange } from "@/models/product";
import type {
  ParsedHistoryRecord,
  HistoryBatchGroup, // Corrected type name
} from "@/models/history";
import type { LoteChangeDetails } from "@/models/lote";

const props = defineProps<{
  batch:
    | HistoryBatchGroup // Corrected type name
    | { batchId: string; createdAt: string; records: any[] }; // Keep local mode flexible
  isLocalMode: boolean;
}>();

function formatDateForDisplay(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}

// Type guards
function isParsedHistoryRecord(record: any): record is ParsedHistoryRecord {
  return (
    record &&
    typeof record.entityType === "string" &&
    typeof record.entityId === "string"
  );
}

function isProductHistory(record: any): record is ProductHistory {
  return (
    record && Array.isArray(record.changes) && typeof record.date === "string"
  );
}

// Helper function to safely access productName
function safeGetProductName(record: any): string {
  if (!record || !record.details) return "Unknown";

  // Check if it's a product record
  if (record.entityType === "product" && record.details) {
    return (
      (record.details as ProductChange).productName ||
      record.entityId ||
      "Unknown"
    );
  }

  // For other types or when missing data
  return record.productNameContext || record.entityId || "Unknown";
}

// Helper function to safely get action
function safeGetAction(record: any): string {
  if (!record || !record.details) return "unknown";
  // Ensure record.details.action is a string before calling replace
  const action = record.details.action;
  return (typeof action === "string" ? action : "unknown").replace("_", " ");
}

// Helper function to safely access record details properties
function safeGetQuantityBefore(record: any): string | number {
  if (!record || !record.details) return "N/A";
  return record.details.quantityBefore ?? "N/A";
}

function safeGetQuantityAfter(record: any): string | number {
  if (!record || !record.details) return "N/A";
  return record.details.quantityAfter ?? "N/A";
}

function safeGetQuantityChanged(record: any): string | number | null {
  if (!record || !record.details) return null;
  return record.details.quantityChanged ?? null;
}

// Additional safety check for lotes
function safeLoteDetails(record: any): any {
  if (!record || !record.details) return {};
  return record.details;
}
</script>

<template>
  <div class="sm:hidden">
    <!-- Para modo autenticado -->
    <div v-if="!isLocalMode && batch.records && batch.records.length > 0">
      <div
        v-for="record in batch.records"
        :key="record.id"
        class="bg-white border-b border-gray-200 last:border-0 p-3"
      >
        <!-- Cabeçalho do registro -->
        <div class="flex justify-between items-start mb-2">
          <div class="font-medium text-gray-800">
            {{ record.entityType === "product" ? "Produto" : "Lote" }}
          </div>
          <span
            :class="{
              'tag-green': ['add', 'created'].includes(safeGetAction(record)),
              'tag-red': ['remove', 'deleted'].includes(safeGetAction(record)),
              'tag-yellow': [
                'update',
                'updated',
                'product details updated',
              ].includes(safeGetAction(record)),
            }"
          >
            {{ safeGetAction(record) }}
          </span>
        </div>

        <!-- Detalhes do registro -->
        <div class="text-sm">
          <!-- Para produto -->
          <div v-if="record.entityType === 'product'">
            <div class="font-medium text-gray-700">
              {{ safeGetProductName(record) }}
            </div>
            <!-- Quantidade alterada -->
            <div
              v-if="
                record.details && // Ensure record.details exists
                (record.details.quantityBefore !== undefined ||
                  record.details.quantityAfter !== undefined)
              "
              class="text-gray-600 mt-1"
            >
              Qtd: {{ safeGetQuantityBefore(record) }} →
              {{ safeGetQuantityAfter(record) }}
              <span v-if="safeGetQuantityChanged(record)">
                ({{
                  (record.details?.action === "add" ? "+" : "-") + // Optional chaining for record.details.action
                  safeGetQuantityChanged(record)
                }})
              </span>
            </div>
            <!-- Campos alterados -->
            <div
              v-if="
                record.entityType === 'product' &&
                record.details && // Ensure record.details exists
                typeof (record.details as ProductChange).changedFields ===
                  'object' &&
                Array.isArray(
                  (record.details as ProductChange).changedFields
                ) &&
                (record.details as ProductChange).changedFields!.length > 0
              "
              class="mt-1 space-y-1"
            >
              <div
                v-for="(field, idx) in (record.details as ProductChange)
                  .changedFields"
                :key="idx"
                class="text-gray-600"
              >
                <span class="font-medium"
                  >{{ field.field.replace("_", " ") }}:</span
                >
                {{ field.oldValue }} → {{ field.newValue }}
              </div>
            </div>
            <!-- Tags especiais -->
            <div class="mt-2 flex flex-wrap gap-1">
              <span
                v-if="(record.details as ProductChange)?.isNewProduct"
                class="tag-new-sm"
                >Novo</span
              >
              <span
                v-if="(record.details as ProductChange)?.isProductRemoval"
                class="tag-removed-sm"
                >Removido</span
              >
            </div>
          </div>

          <!-- Para lote -->
          <div v-else>
            <div class="font-medium text-gray-700">
              ID: {{ record.entityId.substring(0, 8) }}
              <span
                v-if="record.productNameContext"
                class="ml-1 text-xs text-gray-500"
              >
                ({{ record.productNameContext }})
              </span>
            </div>
            <!-- Detalhes do lote -->
            <div class="space-y-1 mt-1">
              <div
                v-if="safeLoteDetails(record).quantityBefore !== undefined"
                class="text-gray-600"
              >
                Qtd:
                {{ safeLoteDetails(record).quantityBefore ?? "N/A" }} →
                {{ safeLoteDetails(record).quantityAfter ?? "N/A" }}
              </div>
              <div
                v-if="
                  safeLoteDetails(record).dataValidadeOld ||
                  safeLoteDetails(record).dataValidadeNew
                "
                class="text-gray-600"
              >
                Val:
                {{ safeLoteDetails(record).dataValidadeOld || "N/A" }}
                →
                {{
                  safeLoteDetails(record).dataValidadeNew ||
                  safeLoteDetails(record).dataValidade ||
                  "N/A"
                }}
              </div>
              <div
                v-else-if="safeLoteDetails(record).dataValidade"
                class="text-gray-600"
              >
                Val: {{ safeLoteDetails(record).dataValidade }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Para modo local -->
    <div v-else>
      <div
        v-for="(record, index) in batch.records"
        :key="index"
        class="bg-white border-b border-gray-200 last:border-0 p-3"
      >
        <div v-if="isProductHistory(record)">
          <div
            v-for="(change, changeIndex) in record.changes"
            :key="changeIndex"
            class="mb-3 last:mb-0 pb-3 last:pb-0 border-b border-gray-100 last:border-0"
          >
            <!-- Cabeçalho da mudança -->
            <div v-if="change" class="flex justify-between items-start mb-2">
              <div class="font-medium text-gray-800">
                {{ change.productName }}
              </div>
              <span
                :class="{
                  'tag-green':
                    change.action === 'add' || change.action === 'created',
                  'tag-red':
                    change.action === 'remove' || change.action === 'deleted',
                  'tag-yellow':
                    change.action === 'update' ||
                    change.action === 'updated' ||
                    change.action === 'product_details_updated',
                }"
              >
                {{ change.action.replace("_", " ") }}
              </span>
            </div>

            <!-- Detalhes da quantidade -->
            <div
              v-if="change && change.quantityBefore !== undefined"
              class="text-sm text-gray-600"
            >
              Qtd: {{ change.quantityBefore }} → {{ change.quantityAfter }}
              <span v-if="change.quantityChanged">
                ({{ change.action === "add" ? "+" : ""
                }}{{ change.quantityChanged }})
              </span>
            </div>

            <!-- Campos alterados -->
            <div
              v-if="
                change &&
                typeof change.changedFields === 'object' &&
                Array.isArray(change.changedFields) &&
                change.changedFields.length > 0
              "
              class="mt-1 text-sm space-y-1"
            >
              <div
                v-for="(field, fieldIdx) in change.changedFields"
                :key="fieldIdx"
                class="text-gray-600"
              >
                <span class="font-medium"
                  >{{ field.field.replace("_", " ") }}:</span
                >
                {{ field.oldValue }} → {{ field.newValue }}
              </div>
            </div>

            <!-- Tags especiais -->
            <div v-if="change" class="mt-2 flex flex-wrap gap-1">
              <span v-if="change?.isNewProduct" class="tag-new-sm">Novo</span>
              <span v-if="change?.isProductRemoval" class="tag-removed-sm"
                >Removido</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Estado vazio -->
    <div
      v-if="!batch.records || batch.records.length === 0"
      class="p-4 text-center text-gray-500"
    >
      Nenhuma alteração encontrada neste lote.
    </div>
  </div>
</template>

<style scoped>
.tag-green {
  @apply bg-green-100 text-green-800 px-2 py-0.5 rounded-full text-xs font-medium;
}
.tag-red {
  @apply bg-red-100 text-red-800 px-2 py-0.5 rounded-full text-xs font-medium;
}
.tag-yellow {
  @apply bg-amber-100 text-amber-800 px-2 py-0.5 rounded-full text-xs font-medium;
}
.tag-new-sm {
  @apply bg-green-500 text-white text-xs font-bold px-1.5 py-0.5 rounded;
}
.tag-removed-sm {
  @apply bg-red-500 text-white text-xs font-bold px-1.5 py-0.5 rounded;
}
</style>
