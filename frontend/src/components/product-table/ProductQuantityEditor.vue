<script setup lang="ts">
const props = defineProps<{
  quantity: number;
  disabled?: boolean; // Keep disabled prop in case it's used for display-only elsewhere
}>();

const emit = defineEmits<{
  (e: "updateQuantity", value: number): void; // For direct input changes
}>();

function handleInput(event: Event) {
  const value = parseFloat((event.target as HTMLInputElement).value);
  if (!isNaN(value)) {
    emit("updateQuantity", Math.max(0, value));
  }
}
</script>

<template>
  <div>
    <input
      type="number"
      min="0"
      :value="quantity"
      @input="handleInput"
      class="w-28 input-field-enhanced text-center"
      :disabled="props.disabled"
    />
    <!-- Buttons removed -->
  </div>
</template>

<style scoped>
.input-field-enhanced {
  @apply px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow;
}
.input-field-enhanced:disabled {
  @apply bg-gray-100 cursor-not-allowed;
}
/* Button styles removed */
</style>
