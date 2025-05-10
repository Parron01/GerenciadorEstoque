import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { Product } from '@/models/product'

const STORAGE_KEY = 'estoque_produtos'

export const useProductStore = defineStore('product', () => {
  const products = ref<Product[]>([])

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

  function updateQuantity(id: string, delta: number) {
    const product = products.value.find((p) => p.id === id)
    if (product) {
      product.quantity += delta
      saveToStorage()
    }
  }

  function addProduct(product: Product) {
    products.value.push(product)
    saveToStorage()
  }

  function removeProduct(id: string) {
    const index = products.value.findIndex((p) => p.id === id)
    if (index !== -1) {
      products.value.splice(index, 1)
      saveToStorage()
    }
  }

  watch(products, saveToStorage, { deep: true })

  loadFromStorage()
  initializeDefaults()

  return {
    products,
    updateQuantity,
    addProduct,
    removeProduct,
  }
})
