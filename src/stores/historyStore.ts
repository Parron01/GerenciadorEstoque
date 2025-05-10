import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { ProductHistory, ProductChange } from '@/models/product'
import { v4 as uuidv4 } from 'uuid'

const STORAGE_KEY = 'estoque_historico'
const API_URL = 'http://localhost:3000/api/history'

export const useHistoryStore = defineStore('history', () => {
  const history = ref<ProductHistory[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchHistory() {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(API_URL)
      if (!response.ok) throw new Error('Erro ao buscar histórico')

      history.value = await response.json()
      saveToStorage() // Manter backup local
    } catch (err) {
      console.error('Erro ao carregar histórico do servidor:', err)
      error.value = 'Falha ao carregar histórico. Usando dados locais.'
      loadFromStorage() // Fallback para dados locais
    } finally {
      isLoading.value = false
    }
  }

  function loadFromStorage() {
    const data = localStorage.getItem(STORAGE_KEY)
    history.value = data ? JSON.parse(data) : []
  }

  function saveToStorage() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history.value))
  }

  async function addBatchEntry(changes: ProductChange[]) {
    if (changes.length === 0) return

    const historyEntry: ProductHistory = {
      id: uuidv4(),
      date: new Date().toISOString(),
      changes,
    }

    try {
      const response = await fetch(API_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(historyEntry),
      })

      if (!response.ok) throw new Error('Falha ao registrar histórico')

      history.value.unshift(historyEntry)
      saveToStorage()
    } catch (err) {
      console.error('Erro ao adicionar histórico no servidor:', err)
      // Adiciona localmente mesmo se falhar no servidor
      history.value.unshift(historyEntry)
      saveToStorage()
    }
  }

  watch(history, saveToStorage, { deep: true })

  // Inicialização
  fetchHistory().catch(() => {
    loadFromStorage()
  })

  return {
    history,
    isLoading,
    error,
    addBatchEntry,
    fetchHistory,
  }
})
