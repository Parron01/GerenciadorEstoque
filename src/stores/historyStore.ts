import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { ProductHistory, ProductChange } from '@/models/product'
import { v4 as uuidv4 } from 'uuid'

const STORAGE_KEY = 'estoque_historico'

export const useHistoryStore = defineStore('history', () => {
  const history = ref<ProductHistory[]>([])

  function loadFromStorage() {
    const data = localStorage.getItem(STORAGE_KEY)
    history.value = data ? JSON.parse(data) : []
  }

  function saveToStorage() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history.value))
  }

  function addBatchEntry(changes: ProductChange[]) {
    if (changes.length === 0) return

    history.value.unshift({
      id: uuidv4(),
      date: new Date().toISOString(),
      changes,
    })
    saveToStorage()
  }

  watch(history, saveToStorage, { deep: true })

  loadFromStorage()

  return {
    history,
    addBatchEntry,
  }
})
