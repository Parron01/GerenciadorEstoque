<script setup lang="ts">
const props = defineProps<{
  quantity: number;
}>();

const emit = defineEmits<{
  (e: "updateQuantity", value: number): void;
  (e: "changeQuantity", delta: number): void;
}>();

function updateQuantity(value: number) {
  emit("updateQuantity", value);
}

function changeQuantity(delta: number) {
  emit("changeQuantity", delta);
}
</script>

<template>
  <div>
    <input
      type="number"
      min="0"
      :value="quantity"
      @input="
        updateQuantity(parseFloat(($event.target as HTMLInputElement).value))
      "
      class="w-28 input-field-enhanced text-center"
    />
    <div class="flex space-x-1 mt-2">
      <button @click="changeQuantity(-10)" class="btn-qty-enhanced">-10</button>
      <button @click="changeQuantity(-1)" class="btn-qty-enhanced">-1</button>
      <button
        @click="changeQuantity(1)"
        class="btn-qty-enhanced bg-emerald-500/90 hover:bg-emerald-600"
      >
        +1
      </button>
      <button
        @click="changeQuantity(10)"
        class="btn-qty-enhanced bg-emerald-500/90 hover:bg-emerald-600"
      >
        +10
      </button>
    </div>
  </div>
</template>

<style scoped>
.input-field-enhanced {
  @apply px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow;
}
.btn-qty-enhanced {
  @apply px-2 py-1 bg-red-500/90 hover:bg-red-600 text-white rounded font-bold shadow-sm hover:shadow transition-all;
}
</style>
