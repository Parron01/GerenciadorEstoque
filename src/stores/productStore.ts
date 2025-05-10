import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { Product } from '@/models/product'

const STORAGE_KEY = 'estoque_produtos'
const API_URL = 'http://localhost:3000/api/products'

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  async function fetchProducts() {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(API_URL)
      if (!response.ok) throw new Error('Erro ao buscar produtos')

      products.value = await response.json()
      saveToStorage() // Manter backup local
    } catch (err) {
      console.error('Erro ao carregar produtos do servidor:', err)
      error.value = 'Falha ao carregar produtos. Usando dados locais.'
      loadFromStorage() // Fallback para dados locais
    } finally {
      isLoading.value = false
    }
  }

  function loadFromStorage() {
    const data = localStorage.getItem(STORAGE_KEY)
    products.value = data ? JSON.parse(data) : []
  }

  function saveToStorage() {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(products.value))
  }

  function initializeDefaults() {
    if (products.value.length === 0) {
      products.value = [
        { id: '1', name: 'Alade', unit: 'L', quantity: 210 },
        { id: '2', name: 'Curbix', unit: 'L', quantity: 71 },
        { id: '3', name: 'Magnum', unit: 'kg', quantity: 110 },
        { id: '4', name: 'Instivo', unit: 'L', quantity: 3 },
        { id: '5', name: 'Kasumin', unit: 'L', quantity: 50 },
        { id: '6', name: 'Priori', unit: 'L', quantity: 33 },
      ]
      saveToStorage()
    }
  }

  async function updateQuantity(id: string, delta: number) {
    const product = products.value.find((p) => p.id === id)
    if (product) {
      const newQuantity = product.quantity + delta

      try {
        const response = await fetch(`${API_URL}/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ quantity: newQuantity }),
        })

        if (!response.ok) throw new Error('Falha ao atualizar produto')

        product.quantity = newQuantity
        saveToStorage()
      } catch (err) {
        console.error('Erro ao atualizar produto no servidor:', err)
        // Atualiza localmente mesmo se falhar no servidor
        product.quantity = newQuantity
        saveToStorage()
      }
    }
  }

  async function addProduct(product: Product) {
    try {
      const response = await fetch(API_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(product),
      })

      if (!response.ok) throw new Error('Falha ao adicionar produto')

      products.value.push(product)
      saveToStorage()
    } catch (err) {
      console.error('Erro ao adicionar produto no servidor:', err)
      // Adiciona localmente mesmo se falhar no servidor
      products.value.push(product)
      saveToStorage()
    }
  }

  async function removeProduct(id: string) {
    try {
      const response = await fetch(`${API_URL}/${id}`, {
        method: 'DELETE',
      })

      if (!response.ok) throw new Error('Falha ao remover produto')

      const index = products.value.findIndex((p) => p.id === id)
      if (index !== -1) {
        products.value.splice(index, 1)
        saveToStorage()
      }
    } catch (err) {
      console.error('Erro ao remover produto do servidor:', err)
      // Remove localmente mesmo se falhar no servidor
      const index = products.value.findIndex((p) => p.id === id)
      if (index !== -1) {
        products.value.splice(index, 1)
        saveToStorage()
      }
    }
  }

  watch(products, saveToStorage, { deep: true })

  // Inicialização
  fetchProducts().catch(() => {
    loadFromStorage()
    initializeDefaults()
  })

  return {
    products,
    isLoading,
    error,
    updateQuantity,
    addProduct,
    removeProduct,
    fetchProducts,
  }
})
