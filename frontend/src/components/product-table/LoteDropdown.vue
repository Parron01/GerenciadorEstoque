<script setup lang="ts">
import { computed } from "vue";
import type { Product } from "@/models/product";
import type { Lote } from "@/models/lote";

const props = defineProps<{
  product: Product;
  isEditMode: boolean;
}>();

const emit = defineEmits<{
  (e: "openEditLote", lote: Lote): void;
  (e: "requestDeleteLote", loteId: string, productId: string): void;
  (e: "openAddLote", productId: string): void;
}>();

function formatDate(dateString?: string) {
  if (!dateString) return "N/A";
  return new Date(dateString).toLocaleDateString();
}

function openEditLote(lote: Lote) {
  emit("openEditLote", lote);
}

function requestDeleteLote(loteId: string) {
  emit("requestDeleteLote", loteId, props.product.id);
}

function openAddLote() {
  emit("openAddLote", props.product.id);
}

const lotes = computed(() => props.product.lotes || []);
</script>

<template>
  <td colspan="4" class="py-3 px-4">
    <div class="rounded-lg border border-indigo-300 shadow-sm overflow-hidden">
      <!-- Lotes Header -->
      <div
        class="bg-gradient-to-r from-indigo-500 to-indigo-600 p-3 flex justify-between items-center"
      >
        <h3 class="font-medium text-white text-sm flex items-center">
          <span class="material-icons-outlined mr-1.5">inventory</span>
          Lotes de <span class="font-bold ml-1">{{ product.name }}</span>
        </h3>
        <span class="text-xs bg-white/20 px-2 py-0.5 rounded-full text-white">
          {{ lotes.length ? lotes.length + " lote(s)" : "Sem lotes" }}
        </span>
      </div>

      <!-- Lotes Content -->
      <div class="p-3 bg-white">
        <!-- Existing Lotes -->
        <div class="space-y-2 max-h-60 overflow-y-auto mb-3 pr-1">
          <div
            v-for="lote in lotes"
            :key="lote.id"
            class="p-3 border border-gray-200 rounded-lg bg-gray-50 flex justify-between items-center hover:bg-gray-100 transition-colors shadow-sm"
          >
            <div>
              <div class="flex items-baseline">
                <span class="font-bold text-lg text-indigo-700">{{
                  lote.quantity
                }}</span>
                <span class="text-gray-600 ml-1">{{ product.unit }}</span>
              </div>
              <div class="text-sm text-gray-600">
                <span class="font-medium text-gray-700">Validade:</span>
                <span class="font-medium">{{
                  formatDate(lote.dataValidade)
                }}</span>
              </div>
              <span class="text-gray-400 text-xs"
                >ID: {{ lote.id.substring(0, 6) }}</span
              >
            </div>

            <div class="flex gap-2">
              <button
                @click="openEditLote(lote)"
                class="btn-edit-enhanced"
                :disabled="!isEditMode"
                :class="{ 'opacity-50 cursor-not-allowed': !isEditMode }"
                title="Editar lote"
              >
                <span class="material-icons-outlined">edit</span>
              </button>
              <button
                @click="requestDeleteLote(lote.id)"
                class="btn-delete-enhanced"
                :disabled="!isEditMode"
                :class="{ 'opacity-50 cursor-not-allowed': !isEditMode }"
                title="Excluir lote"
              >
                <span class="material-icons-outlined">delete</span>
              </button>
            </div>
          </div>

          <!-- Empty state for no lotes -->
          <div
            v-if="lotes.length === 0"
            class="p-4 bg-gray-50 border border-gray-200 rounded-lg text-center text-gray-500 italic"
          >
            Nenhum lote cadastrado para este produto.
          </div>
        </div>

        <!-- Add New Lote Row -->
        <div
          @click="openAddLote"
          class="p-3 border-2 border-dashed border-indigo-300 rounded-lg bg-indigo-50 hover:bg-indigo-100 text-indigo-700 flex justify-center items-center gap-2 cursor-pointer transition-all hover:shadow-md"
          :class="{ 'opacity-60 cursor-not-allowed': !isEditMode }"
        >
          <span class="material-icons-outlined text-indigo-600"
            >add_circle</span
          >
          <span class="font-medium">Adicionar Novo Lote</span>
          <span v-if="!isEditMode" class="text-xs italic text-indigo-500">
            (Ative o modo de edição)
          </span>
        </div>
      </div>
    </div>
  </td>
</template>

<style scoped>
.btn-edit-enhanced {
  @apply p-2 bg-amber-100 hover:bg-amber-200 text-amber-700 rounded-lg transition-colors shadow-sm hover:shadow flex items-center justify-center;
}
.btn-delete-enhanced {
  @apply p-2 bg-red-100 hover:bg-red-200 text-red-700 rounded-lg transition-colors shadow-sm hover:shadow flex items-center justify-center;
}
</style>
