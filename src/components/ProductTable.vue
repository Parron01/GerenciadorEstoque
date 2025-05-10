<script setup lang="ts">
import { useProductStore } from '@/stores/productStore'
import { useHistoryStore } from '@/stores/historyStore'
import { ref } from 'vue'
import { ProductChange, Product } from '@/models/product'
import { v4 as uuidv4 } from 'uuid'
import { useToast } from 'vue-toastification'

const productStore = useProductStore()
const historyStore = useHistoryStore()

// Toast instance
const toast = useToast()

// Estado para controlar o modo de edição
const isEditMode = ref(false)
// Estado para controlar o modo de adição de produto
const isAddProductMode = ref(false)

// Estado para controlar o diálogo de confirmação de exclusão
const showDeleteDialog = ref(false)
const productToDelete = ref<Product | null>(null)

// Novo produto
const newProduct = ref<Omit<Product, 'id'>>({
  name: '',
  unit: 'L',
  quantity: 0,
})

// Estado temporário para armazenar as quantidades durante a edição
const tempQuantities = ref<Record<string, number>>({})

// Inicializa as quantidades temporárias com os valores atuais
function initTempQuantities() {
  const quantities: Record<string, number> = {}
  productStore.products.forEach((product) => {
    quantities[product.id] = product.quantity
  })
  tempQuantities.value = quantities
}

// Ativa o modo de edição
function enableEditMode() {
  initTempQuantities()
  isEditMode.value = true
}

// Altera a quantidade temporária de um produto
function changeQuantity(id: string, delta: number) {
  tempQuantities.value[id] = Math.max(0, (tempQuantities.value[id] || 0) + delta)
}

// Atualiza diretamente a quantidade temporária
function updateQuantity(id: string, value: number) {
  tempQuantities.value[id] = Math.max(0, value)
}

// Confirma todas as alterações e cria entrada em lote no histórico
function confirmUpdates() {
  const changes: ProductChange[] = []

  // Identifica mudanças e atualiza o estoque
  productStore.products.forEach((product) => {
    const originalQuantity = product.quantity
    const newQuantity = tempQuantities.value[product.id] || 0

    if (originalQuantity !== newQuantity) {
      const delta = newQuantity - originalQuantity

      // Atualiza o estoque
      productStore.updateQuantity(product.id, delta)

      // Adiciona à lista de mudanças com informações de antes e depois
      changes.push({
        productId: product.id,
        productName: product.name, // Store the name directly
        action: delta > 0 ? 'add' : 'remove',
        quantityChanged: Math.abs(delta),
        quantityBefore: originalQuantity,
        quantityAfter: newQuantity,
      })
    }
  })

  // Se houver mudanças, adiciona uma única entrada em lote ao histórico
  if (changes.length > 0) {
    historyStore.addBatchEntry(changes)
    toast.success('Estoque atualizado com sucesso!', {
      icon: 'check_circle',
    })
  }

  // Desativa o modo de edição
  isEditMode.value = false
}

// Cancela a edição sem salvar
function cancelEdit() {
  isEditMode.value = false
}

// Abre o formulário para adicionar novo produto
function openAddProductForm() {
  if (!isEditMode.value) {
    newProduct.value = { name: '', unit: 'L', quantity: 0 }
    isAddProductMode.value = true
  }
}

// Cancela a adição de produto
function cancelAddProduct() {
  isAddProductMode.value = false
}

// Adiciona um novo produto
function addProduct() {
  if (!newProduct.value.name || newProduct.value.quantity < 0) {
    return
  }

  const product: Product = {
    id: uuidv4(),
    ...newProduct.value,
  }

  productStore.addProduct(product)

  // Adiciona entrada no histórico para o novo produto
  historyStore.addBatchEntry([
    {
      productId: product.id,
      productName: product.name, // Store the name directly
      action: 'add',
      quantityChanged: product.quantity,
      quantityBefore: 0,
      quantityAfter: product.quantity,
      isNewProduct: true, // Flag as a new product
    },
  ])

  // Exibe toast de sucesso
  toast.success(`Produto "${product.name}" adicionado com sucesso!`, {
    icon: 'add_circle',
  })

  isAddProductMode.value = false
}

// Remove um produto
function removeProduct(id: string) {
  if (!isEditMode.value) {
    const product = productStore.products.find((p) => p.id === id)
    if (product) {
      productToDelete.value = product
      showDeleteDialog.value = true
    }
  }
}

// Confirma a exclusão do produto
function confirmDelete() {
  if (productToDelete.value) {
    const productName = productToDelete.value.name

    // Registra a remoção no histórico
    historyStore.addBatchEntry([
      {
        productId: productToDelete.value.id,
        productName: productToDelete.value.name,
        action: 'remove',
        quantityChanged: productToDelete.value.quantity,
        quantityBefore: productToDelete.value.quantity,
        quantityAfter: 0,
        isProductRemoval: true,
      },
    ])

    productStore.removeProduct(productToDelete.value.id)

    // Exibe toast de remoção
    toast.info(`Produto "${productName}" removido do estoque`, {
      icon: 'delete',
    })

    closeDeleteDialog()
  }
}

// Fecha o diálogo de confirmação
function closeDeleteDialog() {
  showDeleteDialog.value = false
  productToDelete.value = null
}
</script>

<template>
  <div class="overflow-x-auto">
    <!-- Botões de ação -->
    <div class="flex justify-end mb-4 space-x-3" v-if="!isEditMode">
      <button
        @click="openAddProductForm"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition flex items-center"
      >
        <span class="material-icons-outlined mr-1">add_circle</span>
        Novo Produto
      </button>

      <button
        @click="enableEditMode"
        class="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition flex items-center"
      >
        <span class="material-icons-outlined mr-1">edit</span>
        Atualizar Dados
      </button>
    </div>

    <div class="flex justify-end mb-4 space-x-3" v-else>
      <button
        @click="cancelEdit"
        class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition flex items-center"
      >
        <span class="material-icons-outlined mr-1">cancel</span>
        Cancelar
      </button>

      <button
        @click="confirmUpdates"
        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition flex items-center"
      >
        <span class="material-icons-outlined mr-1">check_circle</span>
        Confirmar Atualizações
      </button>
    </div>

    <!-- Formulário para adicionar novo produto -->
    <div
      v-if="isAddProductMode && !isEditMode"
      class="bg-white p-4 rounded shadow mb-4 border border-gray-200"
    >
      <h2 class="text-xl font-bold mb-4">Adicionar Novo Produto</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Nome do Produto</label>
          <input
            v-model="newProduct.name"
            type="text"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Unidade</label>
          <select
            v-model="newProduct.unit"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="L">Litros (L)</option>
            <option value="kg">Quilogramas (kg)</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Quantidade Inicial</label>
          <input
            v-model.number="newProduct.quantity"
            type="number"
            min="0"
            class="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
      </div>
      <div class="mt-4 flex justify-end space-x-3">
        <button
          @click="cancelAddProduct"
          class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition"
        >
          Cancelar
        </button>
        <button
          @click="addProduct"
          class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition"
        >
          Adicionar Produto
        </button>
      </div>
    </div>

    <table class="min-w-full bg-white rounded shadow">
      <thead class="bg-indigo-600 text-white text-left">
        <tr>
          <th class="p-3">Produto</th>
          <th class="p-3">Quantidade</th>
          <th class="p-3">Unidade</th>
          <th v-if="isEditMode" class="p-3">Ações</th>
          <th v-if="!isEditMode" class="p-3 w-20">Remover</th>
        </tr>
      </thead>

      <tbody>
        <tr
          v-for="p in productStore.products"
          :key="p.id"
          class="even:bg-gray-50 hover:bg-gray-100"
        >
          <td class="p-3 font-medium">{{ p.name }}</td>
          <td class="p-3">
            <div v-if="!isEditMode">{{ p.quantity }}</div>
            <input
              v-else
              type="number"
              min="0"
              v-model.number="tempQuantities[p.id]"
              @input="updateQuantity(p.id, tempQuantities[p.id])"
              class="w-24 px-2 py-1 border border-gray-300 rounded text-center"
            />
          </td>
          <td class="p-3">{{ p.unit }}</td>
          <td v-if="isEditMode" class="p-3 space-x-1 flex">
            <button
              class="px-2 py-1 bg-red-500/90 hover:bg-red-600 text-white rounded text-xs"
              @click="changeQuantity(p.id, -10)"
            >
              -10
            </button>
            <button
              class="px-2 py-1 bg-red-500/90 hover:bg-red-600 text-white rounded text-xs"
              @click="changeQuantity(p.id, -1)"
            >
              -1
            </button>
            <button
              class="px-2 py-1 bg-green-500/90 hover:bg-green-600 text-white rounded text-xs"
              @click="changeQuantity(p.id, +1)"
            >
              +1
            </button>
            <button
              class="px-2 py-1 bg-green-500/90 hover:bg-green-600 text-white rounded text-xs"
              @click="changeQuantity(p.id, +10)"
            >
              +10
            </button>
          </td>
          <td v-if="!isEditMode" class="p-3">
            <button
              @click="removeProduct(p.id)"
              class="w-full px-2 py-1 bg-red-500 hover:bg-red-600 text-white rounded text-xs flex items-center justify-center"
            >
              <span class="material-icons-outlined text-sm">delete</span>
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Diálogo de confirmação de exclusão -->
    <div v-if="showDeleteDialog" class="fixed inset-0 flex items-center justify-center z-50">
      <!-- Overlay de fundo escurecido -->
      <div
        class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        @click="closeDeleteDialog"
      ></div>

      <!-- Conteúdo do diálogo -->
      <div
        class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 z-10 overflow-hidden transform transition-all"
      >
        <!-- Cabeçalho -->
        <div class="bg-gradient-to-r from-red-600 to-red-700 px-6 py-4">
          <h3 class="text-lg font-medium text-white flex items-center">
            <span class="material-icons-outlined mr-2">warning</span>
            Confirmar exclusão
          </h3>
        </div>

        <!-- Corpo da mensagem -->
        <div class="px-6 py-4">
          <p class="text-gray-700">
            Tem certeza que deseja remover o produto
            <span class="font-bold">{{ productToDelete?.name }}</span
            >?
          </p>
          <p class="text-sm text-gray-500 mt-2">Esta ação não pode ser desfeita.</p>
        </div>

        <!-- Botões de ação -->
        <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3">
          <button
            @click="closeDeleteDialog"
            class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition flex items-center"
          >
            <span class="material-icons-outlined mr-1 text-sm">close</span>
            Cancelar
          </button>
          <button
            @click="confirmDelete"
            class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition flex items-center"
          >
            <span class="material-icons-outlined mr-1 text-sm">delete</span>
            Excluir
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
