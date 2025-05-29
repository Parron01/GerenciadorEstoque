<script setup lang="ts">
import type { ProductHistory, ProductChange } from "@/models/product";
import type { ParsedHistoryRecord } from "@/models/history";
import type { LoteChangeDetails } from "@/models/lote";

const props = defineProps<{
  history: (ProductHistory | ParsedHistoryRecord)[];
}>();

function formatDateForDisplay(dateStr: string): {
  short: string;
} {
  const date = new Date(dateStr);
  return {
    short:
      date.toLocaleDateString() +
      "\n" +
      date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }),
  };
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
  <div class="sm:hidden space-y-4">
    <div
      v-for="record in history"
      :key="record.id"
      class="bg-white rounded-lg shadow border border-gray-200 overflow-hidden"
    >
      <div class="bg-slate-700 text-white px-3 py-2 text-sm font-medium">
        {{ formatDateForDisplay(getDateFromRecord(record)).short }}
      </div>
      <div class="p-3">
        <!-- Authenticated Mode: ParsedHistoryRecord -->
        <div v-if="isParsedHistoryRecord(record)" class="space-y-2">
          <div
            v-if="record.entityType === 'product'"
            class="pb-2 mb-2 border-b last:border-0 last:mb-0 last:pb-0"
          >
            <div class="flex justify-between items-start mb-1">
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
                    ? 'tag-green'
                    : 'tag-red'
                "
                class="text-xs"
              >
                {{ (record.details as ProductChange).action }}
              </span>
            </div>
            <div class="text-sm text-gray-600">
              Qtd: {{ (record.details as ProductChange).quantityBefore }} →
              {{ (record.details as ProductChange).quantityAfter }}
            </div>
            <span
              v-if="(record.details as ProductChange).isNewProduct"
              class="tag-new-sm"
              >Novo</span
            >
            <span
              v-if="(record.details as ProductChange).isProductRemoval"
              class="tag-removed-sm"
              >Removido</span
            >
          </div>
          <div
            v-else-if="record.entityType === 'lote'"
            class="pb-2 mb-2 border-b last:border-0 last:mb-0 last:pb-0"
          >
            <div class="flex justify-between items-start mb-1">
              <span class="font-medium"
                >Lote: {{ record.entityId.substring(0, 8) }}</span
              >
              <span
                :class="
                  (record.details as LoteChangeDetails).action === 'created'
                    ? 'tag-green'
                    : (record.details as LoteChangeDetails).action === 'updated'
                      ? 'tag-yellow'
                      : 'tag-red'
                "
                class="text-xs"
              >
                {{ (record.details as LoteChangeDetails).action }}
              </span>
            </div>
            <div class="text-xs text-gray-500 mb-1">
              Prod:
              {{
                record.productNameContext ||
                (record.details as LoteChangeDetails).productId?.substring(0, 8)
              }}
            </div>
            <div class="text-sm text-gray-600">
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
                Val:
                {{
                  (record.details as LoteChangeDetails).dataValidadeOld || "N/A"
                }}
                →
                {{
                  (record.details as LoteChangeDetails).dataValidadeNew ||
                  (record.details as LoteChangeDetails).dataValidade ||
                  "N/A"
                }}
              </div>
              <div
                v-else-if="(record.details as LoteChangeDetails).dataValidade"
              >
                Val: {{ (record.details as LoteChangeDetails).dataValidade }}
              </div>
            </div>
          </div>
        </div>
        <!-- Local Mode: ProductHistory (batch changes) -->
        <div v-else-if="isProductHistory(record)">
          <div
            v-for="(change, index) in record.changes"
            :key="index"
            class="pb-3 mb-3 border-b border-gray-200 last:border-0 last:mb-0 last:pb-0 relative"
          >
            <div class="flex justify-between items-start mb-2">
              <span class="font-medium">{{ change.productName }}</span>
              <span
                :class="change.action === 'add' ? 'tag-green' : 'tag-red'"
                class="text-xs"
              >
                {{ change.action === "add" ? "+" : "-"
                }}{{ change.quantityChanged }}
              </span>
            </div>
            <div class="text-sm text-gray-600">
              De: {{ change.quantityBefore }} → Para: {{ change.quantityAfter }}
            </div>
            <span v-if="change.isNewProduct" class="tag-new-sm">Novo</span>
            <span v-if="change.isProductRemoval" class="tag-removed-sm"
              >Remoção</span
            >
          </div>
        </div>
      </div>
    </div>
    <div
      v-if="history.length === 0"
      class="bg-white rounded-lg text-center p-6 shadow border border-gray-200 text-gray-500"
    >
      Nenhuma alteração encontrada.
    </div>
  </div>
</template>

<style scoped>
.tag-green {
  @apply bg-green-100 text-green-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-red {
  @apply bg-red-100 text-red-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-yellow {
  @apply bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full font-medium;
}
.tag-new-sm {
  @apply bg-green-500 text-white text-xs font-bold px-1.5 py-0.5 rounded mr-1;
}
.tag-removed-sm {
  @apply bg-red-500 text-white text-xs font-bold px-1.5 py-0.5 rounded;
}
</style>
