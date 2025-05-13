import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { ProductHistory, ProductChange } from '@/models/product'
import { v4 as uuidv4 } from 'uuid'
import { useAuthStore } from '@/stores/authStore'

// Storage key depends on auth mode
const getStorageKey = () => {
  const authStore = useAuthStore()
  return authStore.isLocalMode ? 'estoque_historico_local' : 'estoque_historico'
}

export const useHistoryStore = defineStore('history', () => {
  const authStore = useAuthStore()
  const history = ref<ProductHistory[]>([])

  // Dados de demonstração para o modo local
  const createDemoHistory = (): ProductHistory[] => {
    const oneHourAgo = new Date(Date.now() - 3600000).toISOString()
    const yesterday = new Date(Date.now() - 86400000).toISOString()
    const lastWeek = new Date(Date.now() - 604800000).toISOString()

    return [
      {
        id: uuidv4(),
        date: oneHourAgo,
        changes: [
          {
            productId: '1',
            productName: 'Fertilizante NPK',
            action: 'add',
            quantityChanged: 20,
            quantityBefore: 100,
            quantityAfter: 120,
          },
        ],
      },
      {
        id: uuidv4(),
        date: yesterday,
        changes: [
          {
            productId: '2',
            productName: 'Herbicida Natural',
            action: 'remove',
            quantityChanged: 5,
            quantityBefore: 50,
            quantityAfter: 45,
          },
        ],
      },
      {
        id: uuidv4(),
        date: lastWeek,
        changes: [
          {
            productId: '3',
            productName: 'Adubo Orgânico',
            action: 'add',
            quantityChanged: 50,
            quantityBefore: 150,
            quantityAfter: 200,
            isNewProduct: true,
          },
        ],
      },
    ]
  }

  function loadFromStorage() {
    const storageKey = getStorageKey()
    const data = localStorage.getItem(storageKey)

    if (data) {
      history.value = JSON.parse(data)
    } else if (authStore.isLocalMode) {
      // Usar histórico de demonstração para modo local
      history.value = createDemoHistory()
      saveToStorage() // Salvar para persistência
    } else {
      // Inicializar com array vazio para modo autenticado sem dados
      history.value = []
    }
  }

  function saveToStorage() {
    const storageKey = getStorageKey()
    localStorage.setItem(storageKey, JSON.stringify(history.value))
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

  // Carregar dados ao inicializar o store
  loadFromStorage()

  // Watch for changes in auth mode to reload data
  watch(
    () => authStore.isLocalMode,
    () => {
      loadFromStorage()
    },
  )

  watch(history, saveToStorage, { deep: true })

  return {
    history,
    addBatchEntry,
    loadFromStorage,
  }
})
