<script setup lang="ts">
import { ref, watch } from "vue";
import type { LotePayload } from "@/models/lote";
import { useToast } from "vue-toastification";

const props = defineProps<{
  show: boolean;
  productId: string;
}>();

const emit = defineEmits(["close", "save"]);

const toast = useToast();

const loteData = ref<LotePayload>({
  quantity: 0,
  dataValidade: "",
});

const today = new Date().toISOString().split("T")[0]; // YYYY-MM-DD format

watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      loteData.value = { quantity: 0, dataValidade: "" };
    }
  }
);

function adjustQuantity(amount: number) {
  const currentQuantity = Number(loteData.value.quantity) || 0;
  let newQuantity = currentQuantity + amount;
  newQuantity = Math.max(0.01, newQuantity); // Ensure minimum value is 0.01
  loteData.value.quantity = parseFloat(newQuantity.toFixed(2)); // Handle precision
}

function validateAndSave() {
  if (loteData.value.quantity <= 0) {
    toast.error("Quantidade do lote deve ser maior que zero.");
    return;
  }
  if (!loteData.value.dataValidade) {
    toast.error("Data de validade é obrigatória.");
    return;
  }
  if (loteData.value.dataValidade < today) {
    toast.error("Data de validade não pode ser anterior à data de hoje.");
    return;
  }
  emit("save", loteData.value);
}
</script>

<template>
  <div v-if="show" class="fixed inset-0 flex items-center justify-center z-50">
    <div
      class="fixed inset-0 bg-black bg-opacity-50"
      @click="$emit('close')"
    ></div>
    <div
      class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 z-10 p-6 space-y-4"
    >
      <h3 class="text-lg font-medium text-gray-900">Adicionar Novo Lote</h3>
      <div>
        <label
          for="lote-quantity"
          class="block text-sm font-medium text-gray-700 mb-1"
          >Quantidade</label
        >
        <div class="flex items-center space-x-1 mt-1">
          <button
            type="button"
            @click="adjustQuantity(-10)"
            class="btn-adjust-qty"
            aria-label="Diminuir 10 da quantidade"
          >
            -10
          </button>
          <button
            type="button"
            @click="adjustQuantity(-1)"
            class="btn-adjust-qty"
            aria-label="Diminuir 1 da quantidade"
          >
            -1
          </button>
          <input
            id="lote-quantity"
            type="number"
            v-model.number="loteData.quantity"
            min="0.01"
            step="0.01"
            class="input-qty flex-grow border-gray-300 focus:border-indigo-500 focus:ring-indigo-500"
            placeholder="0.00"
          />
          <button
            type="button"
            @click="adjustQuantity(1)"
            class="btn-adjust-qty"
            aria-label="Aumentar 1 na quantidade"
          >
            +1
          </button>
          <button
            type="button"
            @click="adjustQuantity(10)"
            class="btn-adjust-qty"
            aria-label="Aumentar 10 na quantidade"
          >
            +10
          </button>
        </div>
      </div>
      <div>
        <label
          for="lote-dataValidade"
          class="block text-sm font-medium text-gray-700"
          >Data de Validade</label
        >
        <input
          id="lote-dataValidade"
          type="date"
          v-model="loteData.dataValidade"
          :min="today"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
      <div class="flex justify-end space-x-3 pt-2">
        <button
          @click="$emit('close')"
          class="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400 transition"
        >
          Cancelar
        </button>
        <button
          @click="validateAndSave"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
        >
          Salvar Lote
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.input-qty {
  @apply text-center px-3 py-2 border rounded-md shadow-sm focus:outline-none sm:text-sm;
  /* Allow text selection */
  user-select: text;
  -webkit-user-select: text;
  -moz-user-select: text;
  -ms-user-select: text;
}

/* Hide default number input spinners */
.input-qty::-webkit-outer-spin-button,
.input-qty::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
.input-qty[type="number"] {
  -moz-appearance: textfield;
}

.btn-adjust-qty {
  @apply px-3 py-2 border border-gray-300 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-1 text-sm font-medium transition-colors;
  min-width: 40px; /* Ensure buttons have a decent tap size */
}
</style>
