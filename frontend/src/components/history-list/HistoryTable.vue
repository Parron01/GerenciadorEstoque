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

// Helper function to safely get action class
function getActionClass(record: any, actionType: string): string {
  if (!record || !record.details) return "";

  const action = (record.details.action || "").toLowerCase();

  if (actionType === "green" && ["add", "created"].includes(action)) {
    return "text-green-700";
  } else if (actionType === "red" && ["remove", "deleted"].includes(action)) {
    return "text-red-700";
  } else if (
    actionType === "amber" &&
    ["update", "updated", "product_details_updated"].includes(action)
  ) {
    return "text-amber-700";
  }

  return "";
}

// Helper function for safely accessing quantity properties
function safeQuantityBefore(record: any): any {
  return record?.details?.quantityBefore;
}

function safeQuantityAfter(record: any): any {
  return record?.details?.quantityAfter;
}

function safeQuantityChanged(record: any): any {
  return record?.details?.quantityChanged;
}

// Helper function for lote details
function safeLoteDetails(record: any): any {
  if (!record || !record.details) return {};
  return record.details;
}

// Helper function for action name
function getActionName(record: any): string {
  if (!record || !record.details) return "unknown";
  return ((record.details.action || "unknown") + "").replace("_", " ");
}
</script>

<template>
  <div class="hidden sm:block">
    <table class="min-w-full bg-white">
      <thead class="bg-slate-100 text-left text-gray-700">
        <tr>
          <th class="p-3 w-1/4">Entidade</th>
          <th class="p-3">Detalhes da Alteração</th>
        </tr>
      </thead>
      <tbody>
        <!-- Para registros no modo autenticado -->
        <template v-if="!isLocalMode && batch.records">
          <tr
            v-for="record in batch.records"
            :key="record.id"
            class="border-b border-gray-200 hover:bg-gray-50"
          >
            <td class="p-3 align-top">
              <div class="font-medium">
                {{ record.entityType === "product" ? "Produto" : "Lote" }}:
              </div>
              <div
                v-if="record.entityType === 'product'"
                class="text-sm text-gray-600"
              >
                {{ safeGetProductName(record) }}
              </div>
              <div v-else class="text-sm text-gray-600">
                ID: {{ record.entityId.substring(0, 8) }}...
                <div
                  v-if="record.productNameContext"
                  class="text-xs text-gray-500"
                >
                  Produto: {{ record.productNameContext }}
                </div>
              </div>
            </td>
            <td class="p-3">
              <!-- Detalhes do produto -->
              <div
                v-if="record.entityType === 'product' && record.details"
                class="border-l-2 pl-3 border-blue-500"
              >
                <div class="flex items-center">
                  <span
                    :class="{
                      'text-green-700': getActionClass(record, 'green'),
                      'text-red-700': getActionClass(record, 'red'),
                      'text-amber-700': getActionClass(record, 'amber'),
                    }"
                    class="font-medium"
                  >
                    {{ getActionName(record) }}
                  </span>
                </div>
                <!-- Quantidade alterada se aplicável -->
                <div
                  v-if="safeQuantityBefore(record) !== undefined"
                  class="text-sm text-gray-600"
                >
                  Qtd: {{ safeQuantityBefore(record) }} →
                  {{ safeQuantityAfter(record) }}
                  <span v-if="safeQuantityChanged(record)">
                    ({{
                      (record.details.action === "add" ? "+" : "-") +
                      safeQuantityChanged(record)
                    }})
                  </span>
                </div>
                <!-- Campos alterados -->
                <div
                  v-if="
                    record.entityType === 'product' &&
                    typeof (record.details as ProductChange).changedFields ===
                      'object' &&
                    Array.isArray(
                      (record.details as ProductChange).changedFields
                    ) &&
                    (record.details as ProductChange).changedFields!.length > 0
                  "
                  class="mt-1 text-sm space-y-1"
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
                <div class="mt-1 space-x-1">
                  <span
                    v-if="(record.details as ProductChange).isNewProduct"
                    class="inline-block px-2 py-0.5 text-xs font-medium bg-green-100 text-green-800 rounded-full"
                    >Novo Produto</span
                  >
                  <span
                    v-if="(record.details as ProductChange).isProductRemoval"
                    class="inline-block px-2 py-0.5 text-xs font-medium bg-red-100 text-red-800 rounded-full"
                    >Produto Removido</span
                  >
                </div>
              </div>

              <!-- Detalhes do lote -->
              <div
                v-else-if="record.entityType === 'lote'"
                class="border-l-2 pl-3 border-purple-500"
              >
                <div class="flex items-center">
                  <span
                    :class="{
                      'text-green-700':
                        (record.details as LoteChangeDetails).action ===
                        'created',
                      'text-red-700':
                        (record.details as LoteChangeDetails).action ===
                        'deleted',
                      'text-amber-700':
                        (record.details as LoteChangeDetails).action ===
                        'updated',
                    }"
                    class="font-medium"
                  >
                    {{ (record.details as LoteChangeDetails).action }}
                  </span>
                </div>
                <!-- Quantidade do lote -->
                <div class="text-sm text-gray-600 space-y-1">
                  <div
                    v-if="safeLoteDetails(record).quantityBefore !== undefined"
                  >
                    Qtd:
                    {{ safeLoteDetails(record).quantityBefore }} →
                    {{ safeLoteDetails(record).quantityAfter }}
                  </div>
                  <!-- Data de validade -->
                  <div
                    v-if="
                      (record.details as LoteChangeDetails).dataValidadeOld ||
                      (record.details as LoteChangeDetails).dataValidadeNew
                    "
                  >
                    Validade:
                    {{
                      (record.details as LoteChangeDetails).dataValidadeOld ||
                      "N/A"
                    }}
                    →
                    {{
                      (record.details as LoteChangeDetails).dataValidadeNew ||
                      (record.details as LoteChangeDetails).dataValidade ||
                      "N/A"
                    }}
                  </div>
                  <div
                    v-else-if="
                      (record.details as LoteChangeDetails).dataValidade
                    "
                  >
                    Validade:
                    {{ (record.details as LoteChangeDetails).dataValidade }}
                  </div>
                </div>
              </div>
            </td>
          </tr>
        </template>

        <!-- Para registros do modo local (ProductHistory) -->
        <template v-else>
          <tr
            v-for="(record, index) in batch.records"
            :key="index"
            class="border-b border-gray-200 hover:bg-gray-50"
          >
            <td class="p-3" colspan="2">
              <div v-if="isProductHistory(record)" class="space-y-3">
                <div
                  v-for="(change, changeIndex) in record.changes"
                  :key="changeIndex"
                  class="flex flex-wrap items-start border-l-2 pl-3 relative"
                  :class="{
                    'border-green-500':
                      change.action === 'add' || change.action === 'created',
                    'border-red-500':
                      change.action === 'remove' || change.action === 'deleted',
                    'border-amber-500':
                      change.action === 'update' ||
                      change.action === 'updated' ||
                      change.action === 'product_details_updated',
                  }"
                >
                  <!-- Cabeçalho da mudança -->
                  <div class="w-full mb-1">
                    <span class="font-medium mr-2">{{
                      change.productName
                    }}</span>
                    <span
                      :class="{
                        'bg-green-100 text-green-800':
                          change.action === 'add' ||
                          change.action === 'created',
                        'bg-red-100 text-red-800':
                          change.action === 'remove' ||
                          change.action === 'deleted',
                        'bg-amber-100 text-amber-800':
                          change.action === 'update' ||
                          change.action === 'updated' ||
                          change.action === 'product_details_updated',
                      }"
                      class="px-2 py-0.5 rounded-full text-xs font-medium"
                    >
                      {{ change.action.replace("_", " ") }}
                    </span>
                  </div>

                  <!-- Detalhes da quantidade -->
                  <div
                    v-if="change.quantityBefore !== undefined"
                    class="w-full text-gray-600 text-sm"
                  >
                    Qtd: {{ change.quantityBefore }} →
                    {{ change.quantityAfter }}
                    <span v-if="change.quantityChanged">
                      ({{ change.action === "add" ? "+" : ""
                      }}{{ change.quantityChanged }})
                    </span>
                  </div>

                  <!-- Campos alterados -->
                  <div
                    v-if="
                      typeof change.changedFields === 'object' &&
                      Array.isArray(change.changedFields) &&
                      change.changedFields!.length > 0
                    "
                    class="w-full mt-1 text-sm space-y-1"
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
                  <div class="mt-1 space-x-1">
                    <span
                      v-if="change.isNewProduct"
                      class="inline-block px-2 py-0.5 text-xs font-medium bg-green-100 text-green-800 rounded-full"
                      >Novo Produto</span
                    >
                    <span
                      v-if="change.isProductRemoval"
                      class="inline-block px-2 py-0.5 text-xs font-medium bg-red-100 text-red-800 rounded-full"
                      >Produto Removido</span
                    >
                  </div>
                </div>
              </div>
            </td>
          </tr>
        </template>

        <!-- Estado vazio -->
        <tr v-if="!batch.records || batch.records.length === 0">
          <td colspan="2" class="p-8 text-center text-gray-500">
            Nenhuma alteração encontrada neste lote.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
