import { defineStore } from "pinia";
import { ref, watch } from "vue";
import type { Product } from "@/models/product";
import { useAuthStore } from "@/stores/authStore";
import { v4 as uuidv4 } from "uuid";

// Storage key depends on auth mode
const getStorageKey = () => {
  const authStore = useAuthStore();
  return authStore.isLocalMode ? "estoque_produtos_local" : "estoque_produtos";
};

export const useProductStore = defineStore("product", () => {
  const authStore = useAuthStore();
  const products = ref<Product[]>([]);

  // Dados de demonstração para modo local
  const demoProducts: Product[] = [
    { id: uuidv4(), name: "Fertilizante NPK", unit: "kg", quantity: 120 },
    { id: uuidv4(), name: "Herbicida Natural", unit: "L", quantity: 45 },
    { id: uuidv4(), name: "Adubo Orgânico", unit: "kg", quantity: 200 },
    { id: uuidv4(), name: "Inseticida Biológico", unit: "L", quantity: 15 },
    { id: uuidv4(), name: "Calcário Agrícola", unit: "kg", quantity: 90 },
    { id: uuidv4(), name: "Óleo de Neem", unit: "L", quantity: 8 },
  ];

  function loadFromStorage() {
    // Determinar qual fonte de dados usar com base no modo de autenticação
    const storageKey = getStorageKey();
    const data = localStorage.getItem(storageKey);

    if (data) {
      products.value = JSON.parse(data);
    } else if (authStore.isLocalMode) {
      // Usar dados de demonstração para modo local se não houver dados salvos
      products.value = [...demoProducts];
      saveToStorage(); // Salvar para persistência
    } else {
      // Usar produtos padrão para modo autenticado sem dados salvos
      products.value = [
        { id: "1", name: "Alade", unit: "L", quantity: 210 },
        { id: "2", name: "Curbix", unit: "L", quantity: 71 },
        { id: "3", name: "Magnum", unit: "kg", quantity: 110 },
        { id: "4", name: "Instivo", unit: "L", quantity: 3 },
        { id: "5", name: "Kasumin", unit: "L", quantity: 50 },
        { id: "6", name: "Priori", unit: "L", quantity: 33 },
      ];
      saveToStorage();
    }
  }

  function saveToStorage() {
    const storageKey = getStorageKey();
    localStorage.setItem(storageKey, JSON.stringify(products.value));
  }

  function updateQuantity(id: string, delta: number) {
    const product = products.value.find((p) => p.id === id);
    if (product) {
      product.quantity += delta;
      saveToStorage();
    }
  }

  function addProduct(product: Product) {
    products.value.push(product);
    saveToStorage();
  }

  function removeProduct(id: string) {
    const index = products.value.findIndex((p) => p.id === id);
    if (index !== -1) {
      products.value.splice(index, 1);
      saveToStorage();
    }
  }

  // Carregar dados ao inicializar o store
  loadFromStorage();

  // Watch for changes in auth mode to reload data
  watch(
    () => authStore.isLocalMode,
    () => {
      loadFromStorage();
    }
  );

  watch(products, saveToStorage, { deep: true });

  return {
    products,
    updateQuantity,
    addProduct,
    removeProduct,
    loadFromStorage,
  };
});
