<script setup lang="ts">
import { useHistoryStore } from '@/stores/historyStore'
import { computed } from 'vue'

// Define props before using them
const props = defineProps<{
  filterOption: string
}>()

const historyStore = useHistoryStore()

// Helper functions for date filtering
const isToday = (dateString: string): boolean => {
  const today = new Date()
  const date = new Date(dateString)
  return (
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()
  )
}

const isThisWeek = (dateString: string): boolean => {
  const date = new Date(dateString)
  const today = new Date()

  // Get the first day of the week (Sunday)
  const firstDay = new Date(today.getTime())
  const day = today.getDay()
  firstDay.setDate(today.getDate() - day)

  // Get the last day of the week (Saturday) - FIXED to use timestamp
  const lastDay = new Date(firstDay.getTime())
  lastDay.setDate(firstDay.getDate() + 6)

  // Check if the date is between first and last day of the week
  return date >= firstDay && date <= lastDay
}

const isThisMonth = (dateString: string): boolean => {
  const date = new Date(dateString)
  const today = new Date()
  return date.getMonth() === today.getMonth() && date.getFullYear() === today.getFullYear()
}

// Processa o histórico para filtro selecionado - agora usando nomes de produtos
// armazenados no próprio histórico
const rows = computed(() => {
  if (!historyStore.history) return []

  const filteredHistory = historyStore.history.filter((h) => {
    switch (props.filterOption) {
      case 'today':
        return isToday(h.date)
      case 'week':
        return isThisWeek(h.date)
      case 'month':
        return isThisMonth(h.date)
      default:
        return true // 'all' or any other value shows everything
    }
  })

  return filteredHistory
})
</script>

<template>
  <div class="overflow-x-auto">
    <table class="min-w-full bg-white rounded shadow">
      <thead class="bg-slate-700 text-white text-left">
        <tr>
          <th class="p-3 w-1/4">Data/Hora</th>
          <th class="p-3">Alterações</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="h in rows" :key="h.id" class="border-b border-gray-200 hover:bg-gray-50">
          <td class="p-3 align-top">{{ new Date(h.date).toLocaleString() }}</td>
          <td class="p-3">
            <div class="grid grid-cols-1 gap-2">
              <div
                v-for="(change, index) in h.changes"
                :key="index"
                class="flex flex-wrap items-center border-l-2 pl-3 mb-1 relative"
                :class="change.action === 'add' ? 'border-green-500' : 'border-red-500'"
              >
                <!-- Display product name from the stored history -->
                <span class="font-medium mr-2">{{ change.productName }}</span>
                <span
                  :class="
                    change.action === 'add'
                      ? 'bg-green-100 text-green-800'
                      : 'bg-red-100 text-red-800'
                  "
                  class="px-2 py-0.5 rounded-full text-xs font-medium"
                >
                  {{ change.action === 'add' ? '+' : '-' }}{{ change.quantityChanged }}
                </span>
                <span class="ml-2 text-gray-600 text-sm flex items-center">
                  <span class="text-gray-400 mx-1">De:</span>
                  <span class="font-medium">{{ change.quantityBefore }}</span>
                  <span class="mx-1">→</span>
                  <span class="text-gray-400 mr-1">Para:</span>
                  <span class="font-medium">{{ change.quantityAfter }}</span>
                </span>

                <!-- Tags for new product or product removal -->
                <span
                  v-if="change.isNewProduct"
                  class="absolute top-0 right-0 bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded"
                >
                  Novo produto
                </span>

                <span
                  v-if="change.isProductRemoval"
                  class="absolute top-0 right-0 bg-red-500 text-white text-xs font-bold px-2 py-0.5 rounded"
                >
                  Remoção de produto
                </span>
              </div>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <p v-if="rows.length === 0" class="text-gray-500 mt-4">
      Nenhuma alteração registrada no período selecionado.
    </p>
  </div>
</template>
