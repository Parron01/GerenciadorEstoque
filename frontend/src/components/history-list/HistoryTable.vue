<script setup lang="ts">
import type { ProductHistory, ProductChange } from "@/models/product";
import type { ParsedHistoryRecord } from "@/models/history";
import type { LoteChangeDetails } from "@/models/lote";

const props = defineProps<{
  history: (ProductHistory | ParsedHistoryRecord)[];
}>();

function formatDateForDisplay(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}

function getDateFromRecord(
  record: ProductHistory | ParsedHistoryRecord
): string {
  return "createdAt" in record ? record.createdAt : record.date;
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
</script>

<template>
  <div class="hidden sm:block overflow-x-auto">
    <table class="min-w-full bg-white rounded shadow">
      <thead class="bg-slate-700 text-white text-left">
        <tr>
          <th class="p-3 w-1/4">Data/Hora</th>
          <th class="p-3">Detalhes da Alteração</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="record in history"
          :key="record.id"
          class="border-b border-gray-200 hover:bg-gray-50"
        >
          <td class="p-3 align-top">
            {{ formatDateForDisplay(getDateFromRecord(record)) }}
          </td>
          <td class="p-3">
            <!-- Authenticated Mode: ParsedHistoryRecord -->
            <div v-if="isParsedHistoryRecord(record)" class="space-y-2">
              <div
                v-if="record.entityType === 'product'"
                class="border-l-2 pl-3 border-blue-500"
              >
                <span class="font-medium"
                  >Produto:
                  {{
                    (record.details as ProductChange).productName ||
                    record.entityId
                  }}</span
                >
                <span
                  :class="
                    (record.details as ProductChange).action === 'add' ||
                    (record.details as ProductChange).action === 'created'
                      ? 'text-green-700'
                      : 'text-red-700'
                  "
                  class="ml-2 text-sm"
                >
                  Ação: {{ (record.details as ProductChange).action }}
                </span>
                <div class="text-xs text-gray-600">
                  Qtd: {{ (record.details as ProductChange).quantityBefore }} →
                  {{ (record.details as ProductChange).quantityAfter }}
                  ({{
                    (record.details as ProductChange).action === "add" ||
                    (record.details as ProductChange).action === "created"
                      ? "+"
                      : ""
                  }}
                  {{ (record.details as ProductChange).quantityChanged }})
                </div>
                <span
                  v-if="(record.details as ProductChange).isNewProduct"
                  class="tag-new"
                  >Novo Produto</span
                >
                <span
                  v-if="(record.details as ProductChange).isProductRemoval"
                  class="tag-removed"
                  >Produto Removido</span
                >
              </div>
              <div
                v-else-if="record.entityType === 'lote'"
                class="border-l-2 pl-3 border-purple-500"
              >
                <span class="font-medium"
                  >Lote: {{ record.entityId.substring(0, 8) }}</span
                >
                <span class="text-xs text-gray-500 ml-1"
                  >(Produto:
                  {{
                    record.productNameContext ||
                    (record.details as LoteChangeDetails).productId?.substring(
                      0,
                      8
                    )
                  }})</span
                >
                <span
                  :class="
                    (record.details as LoteChangeDetails).action === 'created'
                      ? 'text-green-700'
                      : (record.details as LoteChangeDetails).action ===
                          'updated'
                        ? 'text-yellow-700'
                        : 'text-red-700'
                  "
                  class="ml-2 text-sm"
                >
                  Ação: {{ (record.details as LoteChangeDetails).action }}
                </span>
                <div class="text-xs text-gray-600">
                  <div
                    v-if="
                      (record.details as LoteChangeDetails).quantityBefore !==
                      undefined
                    "
                  >
                    Qtd:
                    {{ (record.details as LoteChangeDetails).quantityBefore }} →
                    {{ (record.details as LoteChangeDetails).quantityAfter }}
                  </div>
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
            </div>
            <!-- Local Mode: ProductHistory (batch changes) -->
            <div
              v-else-if="isProductHistory(record)"
              class="grid grid-cols-1 gap-2"
            >
              <div
                v-for="(change, index) in record.changes"
                :key="index"
                class="flex flex-wrap items-center border-l-2 pl-3 mb-1 relative"
                :class="
                  change.action === 'add'
                    ? 'border-green-500'
                    : 'border-red-500'
                "
              >
                <span class="font-medium mr-2">{{ change.productName }}</span>
                <span
                  :class="
                    change.action === 'add'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-red-100 text-red-800'
                  "
                  class="px-2 py-0.5 rounded-full text-xs font-medium"
                >
                  {{ change.action === "add" ? "+" : "-"
                  }}{{ change.quantityChanged }}
                </span>
                <span class="ml-2 text-gray-600 text-sm flex items-center">
                  <span class="text-gray-400 mx-1">De:</span>
                  <span class="font-medium">{{ change.quantityBefore }}</span>
                  <span class="mx-1">→</span>
                  <span class="text-gray-400 mr-1">Para:</span>
                  <span class="font-medium">{{ change.quantityAfter }}</span>
                </span>
                <span v-if="change.isNewProduct" class="tag-new"
                  >Novo produto</span
                >
                <span v-if="change.isProductRemoval" class="tag-removed"
                  >Remoção de produto</span
                >
              </div>
            </div>
          </td>
        </tr>
        <tr v-if="history.length === 0">
          <td colspan="2" class="p-8 text-center text-gray-500">
            Nenhuma alteração encontrada.
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.tag-new {
  @apply absolute top-0 right-0 bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded;
}
.tag-removed {
  @apply absolute top-0 right-0 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded;
}
</style>
