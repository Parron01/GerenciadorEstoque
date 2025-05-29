<script setup lang="ts">
import { useProductStore } from "@/stores/productStore";
import { useHistoryStore } from "@/stores/historyStore";
import { ref, computed } from "vue";
import type { ProductChange, Product } from "@/models/product";
import type { Lote, LotePayload } from "@/models/lote";
import { v4 as uuidv4 } from "uuid";
import { useToast } from "vue-toastification";
import AddLoteModal from "./AddLoteModal.vue";
import EditLoteModal from "./EditLoteModal.vue";
import { useAuthStore } from "@/stores/authStore";
import ProductRow from "./product-table/ProductRow.vue";
import LoteDropdown from "./product-table/LoteDropdown.vue";

const productStore = useProductStore();
const historyStore = useHistoryStore();
const authStore = useAuthStore();
const toast = useToast();

// Edit mode state
const isEditMode = ref(false);
const isAddProductMode = ref(false);
const showDeleteDialog = ref(false);
const productToDelete = ref<Product | null>(null);

// Product states
const newProduct = ref<Omit<Product, "id" | "lotes">>({
  name: "",
  unit: "L",
  quantity: 0,
});

// Temp states for edit mode
const tempQuantities = ref<Record<string, number>>({});
const tempProductDetails = ref<
  Record<string, { name: string; unit: "L" | "kg" }>
>({});

// Copy of original products for history tracking and cancellation
const productsBeforeEdit = ref<Product[]>([]);

// Lote Modals State
const showAddLoteModal = ref(false);
const currentProductIdForLote = ref<string | null>(null);
const showEditLoteModal = ref(false);
const currentLoteToEdit = ref<Lote | null>(null);
const showDeleteLoteDialog = ref(false);
const loteToDelete = ref<{ loteId: string; productId: string } | null>(null);

// State for expanded product lotes (accordion)
const expandedProducts = ref<Record<string, boolean>>({});

// Toggle accordion state for product lotes
function toggleProductLotes(productId: string) {
  expandedProducts.value[productId] = !expandedProducts.value[productId];
}

function getProductDisplayQuantity(product: Product): number {
  if (product.lotes && product.lotes.length > 0) {
    return product.lotes.reduce((sum, lote) => sum + lote.quantity, 0);
  }
  return product.quantity;
}

function initTempStates() {
  // Create deep copy of products for tracking changes and potential cancellation
  productsBeforeEdit.value = JSON.parse(JSON.stringify(productStore.products));

  // Initialize temp values for editable fields
  productStore.products.forEach((product) => {
    tempQuantities.value[product.id] = getProductDisplayQuantity(product);
    tempProductDetails.value[product.id] = {
      name: product.name,
      unit: product.unit,
    };
  });
}

function enableEditMode() {
  initTempStates();
  isEditMode.value = true;
}

function changeProductQuantity(id: string, delta: number) {
  // For products without lotes in edit mode
  const product = productStore.products.find((p) => p.id === id);
  if (product && (!product.lotes || product.lotes.length === 0)) {
    tempQuantities.value[id] = Math.max(
      0,
      (tempQuantities.value[id] || 0) + delta
    );
  }
}

function updateProductQuantityDirectly(id: string, value: number) {
  // For products without lotes in edit mode
  const product = productStore.products.find((p) => p.id === id);
  if (product && (!product.lotes || product.lotes.length === 0)) {
    tempQuantities.value[id] = Math.max(0, value);
  }
}

// Track all changes made during edit session to generate correct history
const loteChangesTracking = ref<{
  created: { productId: string; lote: Lote }[];
  updated: {
    productId: string;
    loteId: string;
    before: LotePayload;
    after: LotePayload;
  }[];
  deleted: { productId: string; loteId: string; lote: Lote }[];
}>({
  created: [],
  updated: [],
  deleted: [],
});

// Add new function to update base quantity of product
function updateProductBaseQuantity(productId: string, newBaseQuantity: number) {
  if (isEditMode.value) {
    const product = productStore.products.find((p) => p.id === productId);
    if (product) {
      // Store original base quantity for history tracking
      const originalBaseQuantity = product.quantity;

      // Update base quantity in temp state for history tracking
      if (!productsBeforeEdit.value.find((p) => p.id === productId)?.quantity) {
        // If not already stored, create a snapshot
        const originalProduct = JSON.parse(JSON.stringify(product));
        const existingIndex = productsBeforeEdit.value.findIndex(
          (p) => p.id === productId
        );
        if (existingIndex >= 0) {
          productsBeforeEdit.value[existingIndex] = originalProduct;
        } else {
          productsBeforeEdit.value.push(originalProduct);
        }
      }

      // Update the model
      product.quantity = Math.max(0, newBaseQuantity);

      toast.success(
        `Quantidade base do produto atualizada para ${newBaseQuantity}`
      );
    }
  }
}

async function confirmUpdates() {
  const productChangesBatch: ProductChange[] = [];

  // First collect all changes for history

  // 1. Product detail and quantity changes
  for (const product of productStore.products) {
    const originalProduct = productsBeforeEdit.value.find(
      (p) => p.id === product.id
    );
    if (!originalProduct) continue;

    const originalDisplayQuantity = getProductDisplayQuantity(originalProduct);
    const originalBaseQuantity = originalProduct.quantity;
    const editedName = tempProductDetails.value[product.id]?.name;
    const editedUnit = tempProductDetails.value[product.id]?.unit;

    // Handle base quantity changes (NEW)
    if (product.quantity !== originalBaseQuantity) {
      productChangesBatch.push({
        productId: product.id,
        productName: editedName || product.name,
        action: "product_base_quantity_updated",
        quantityBefore: originalBaseQuantity,
        quantityAfter: product.quantity,
        changedFields: [
          {
            field: "base_quantity",
            oldValue: originalBaseQuantity,
            newValue: product.quantity,
          },
        ],
      });
    }

    // Handle product detail changes (name, unit)
    if (
      (editedName && editedName !== originalProduct.name) ||
      (editedUnit && editedUnit !== originalProduct.unit)
    ) {
      // Add to history batch
      productChangesBatch.push({
        productId: product.id,
        productName: editedName || product.name,
        action: "product_details_updated",
        changedFields: [
          ...(editedName !== originalProduct.name
            ? [
                {
                  field: "name",
                  oldValue: originalProduct.name,
                  newValue: editedName,
                },
              ]
            : []),
          ...(editedUnit !== originalProduct.unit
            ? [
                {
                  field: "unit",
                  oldValue: originalProduct.unit,
                  newValue: editedUnit,
                },
              ]
            : []),
        ],
      });

      // Apply to product store
      await productStore.updateProductDetails(product.id, {
        name: editedName,
        unit: editedUnit,
      });
    }

    // Handle quantity changes for products WITHOUT lotes
    if (!product.lotes || product.lotes.length === 0) {
      const tempQty = tempQuantities.value[product.id];
      if (tempQty !== undefined && tempQty !== originalDisplayQuantity) {
        const delta = tempQty - originalDisplayQuantity;

        // Add to history batch
        productChangesBatch.push({
          productId: product.id,
          productName: editedName || product.name,
          action: delta > 0 ? "add" : "remove",
          quantityChanged: Math.abs(delta),
          quantityBefore: originalDisplayQuantity,
          quantityAfter: tempQty,
        });

        // Apply to product store
        productStore.updateProductQuantity(product.id, tempQty);
      }
    }
  }

  // 2. Lote changes

  // Created lotes
  for (const { productId, lote } of loteChangesTracking.value.created) {
    const product = productStore.products.find((p) => p.id === productId);
    if (!product) continue;

    productChangesBatch.push({
      productId: productId,
      productName: product.name,
      action: "lote_created",
      changedFields: [
        {
          field: "lote",
          loteId: lote.id,
          newValue: {
            quantity: lote.quantity,
            dataValidade: lote.dataValidade,
          },
        },
      ],
    });
  }

  // Updated lotes
  for (const { productId, loteId, before, after } of loteChangesTracking.value
    .updated) {
    const product = productStore.products.find((p) => p.id === productId);
    if (!product) continue;

    productChangesBatch.push({
      productId: productId,
      productName: product.name,
      action: "lote_updated",
      changedFields: [
        {
          field: "lote",
          loteId: loteId,
          oldValue: before,
          newValue: after,
        },
      ],
    });
  }

  // Deleted lotes
  for (const { productId, loteId, lote } of loteChangesTracking.value.deleted) {
    const product = productStore.products.find((p) => p.id === productId);
    if (!product) continue;

    productChangesBatch.push({
      productId: productId,
      productName: product.name,
      action: "lote_deleted",
      changedFields: [
        {
          field: "lote",
          loteId: loteId,
          oldValue: {
            quantity: lote.quantity,
            dataValidade: lote.dataValidade,
          },
        },
      ],
    });
  }

  // Create history entry from collected changes
  if (authStore.isLocalMode && productChangesBatch.length > 0) {
    historyStore.addBatchEntry(productChangesBatch);
    toast.success("Todas as alterações foram aplicadas e registradas!");
  } else if (!authStore.isLocalMode && productChangesBatch.length > 0) {
    // For auth mode, refresh history from server
    historyStore.refreshHistory();
    toast.success("Alterações enviadas ao servidor!");
  }

  // Refresh data if in auth mode
  if (!authStore.isLocalMode) {
    await productStore.fetchProductsFromApi();
    historyStore.refreshHistory();
  }

  // Reset states
  isEditMode.value = false;
  productsBeforeEdit.value = [];
  loteChangesTracking.value = { created: [], updated: [], deleted: [] };
  expandedProducts.value = {}; // Close all accordions
}

function cancelEdit() {
  // Revert to original state by loading from storage or API
  if (authStore.isLocalMode) {
    productStore.loadFromStorage();
  } else {
    productStore.fetchProductsFromApi();
  }

  // Reset all tracking and edit states
  isEditMode.value = false;
  productsBeforeEdit.value = [];
  loteChangesTracking.value = { created: [], updated: [], deleted: [] };
  expandedProducts.value = {}; // Close all accordions
  toast.info("Alterações canceladas.");
}

function openAddProductForm() {
  if (!isEditMode.value) {
    newProduct.value = { name: "", unit: "L", quantity: 0 };
    isAddProductMode.value = true;
  }
}

function cancelAddProduct() {
  isAddProductMode.value = false;
}

async function addProductHandler() {
  if (!newProduct.value.name || newProduct.value.quantity < 0) {
    toast.error("Nome do produto e quantidade válida são obrigatórios.");
    return;
  }
  await productStore.addProduct({ ...newProduct.value });
  if (authStore.isLocalMode) {
    historyStore.addBatchEntry([
      {
        productId:
          productStore.products.find((p) => p.name === newProduct.value.name)
            ?.id || uuidv4(), // Approximate ID for local
        productName: newProduct.value.name,
        action: "created",
        quantityChanged: newProduct.value.quantity,
        quantityBefore: 0,
        quantityAfter: newProduct.value.quantity,
        isNewProduct: true,
      },
    ]);
  } else {
    historyStore.refreshHistory(); // Backend handles history
  }
  isAddProductMode.value = false;
}

function requestDeleteProduct(product: Product) {
  if (!isEditMode.value) {
    productToDelete.value = product;
    showDeleteDialog.value = true;
  }
}

async function confirmDeleteProduct() {
  if (productToDelete.value) {
    const prodName = productToDelete.value.name;
    const prodId = productToDelete.value.id;
    const originalQuantity = getProductDisplayQuantity(productToDelete.value);

    await productStore.removeProduct(prodId);

    if (authStore.isLocalMode) {
      historyStore.addBatchEntry([
        {
          productId: prodId,
          productName: prodName,
          action: "deleted",
          quantityChanged: originalQuantity,
          quantityBefore: originalQuantity,
          quantityAfter: 0,
          isProductRemoval: true,
        },
      ]);
    } else {
      historyStore.refreshHistory(); // Backend handles history
    }
    toast.info(`Produto "${prodName}" removido.`);
    closeDeleteDialog();
  }
}

function closeDeleteDialog() {
  showDeleteDialog.value = false;
  productToDelete.value = null;
}

// Lote actions - only available in edit mode for consistent history
function openAddLote(productId: string) {
  if (!isEditMode.value) {
    toast.info(
      "Ative o modo de edição usando 'Atualizar Dados' para adicionar lotes."
    );
    return;
  }
  currentProductIdForLote.value = productId;
  showAddLoteModal.value = true;
}

async function handleSaveLote(loteData: LotePayload) {
  if (!currentProductIdForLote.value) return;

  if (isEditMode.value) {
    // In edit mode, track lote operations for batched history
    // Create lote in memory
    const newLote: Lote = {
      ...loteData,
      id: uuidv4(),
      productId: currentProductIdForLote.value,
      createdAt: new Date().toISOString(),
    };

    // Add to tracking for history
    loteChangesTracking.value.created.push({
      productId: currentProductIdForLote.value,
      lote: newLote,
    });

    // Find the product and add the lote to it
    const product = productStore.products.find(
      (p) => p.id === currentProductIdForLote.value
    );
    if (product) {
      if (!product.lotes) product.lotes = [];
      product.lotes.push(newLote);
    }

    toast.success("Lote adicionado e será salvo ao confirmar as atualizações.");
  } else {
    // Normal flow when not in edit mode (should never happen with UI restrictions)
    await productStore.createLote(currentProductIdForLote.value, loteData);
    if (!authStore.isLocalMode) historyStore.refreshHistory();
  }

  showAddLoteModal.value = false;
  currentProductIdForLote.value = null;
}

function openEditLote(lote: Lote) {
  if (!isEditMode.value) {
    toast.info(
      "Ative o modo de edição usando 'Atualizar Dados' para editar lotes."
    );
    return;
  }

  currentLoteToEdit.value = { ...lote }; // Pass a copy
  showEditLoteModal.value = true;
}

async function handleUpdateLote(loteId: string, loteData: LotePayload) {
  if (!currentLoteToEdit.value) return;
  const productId = currentLoteToEdit.value.productId;

  if (isEditMode.value) {
    // Find the original lote data
    const product = productStore.products.find((p) => p.id === productId);
    const lote = product?.lotes?.find((l) => l.id === loteId);

    if (lote) {
      // Save before state
      const beforeData: LotePayload = {
        quantity: lote.quantity,
        dataValidade: lote.dataValidade,
      };

      // Track the update
      loteChangesTracking.value.updated.push({
        productId,
        loteId,
        before: beforeData,
        after: loteData,
      });

      // Update the lote
      Object.assign(lote, loteData, {
        updatedAt: new Date().toISOString(),
      });

      toast.success(
        "Lote atualizado e será salvo ao confirmar as atualizações."
      );
    }
  } else {
    // Normal flow when not in edit mode (should never happen with UI restrictions)
    await productStore.updateLote(loteId, productId, loteData);
    if (!authStore.isLocalMode) historyStore.refreshHistory();
  }

  showEditLoteModal.value = false;
  currentLoteToEdit.value = null;
}

function requestDeleteLote(loteId: string, productId: string) {
  if (!isEditMode.value) {
    toast.info(
      "Ative o modo de edição usando 'Atualizar Dados' para remover lotes."
    );
    return;
  }

  loteToDelete.value = { loteId, productId };
  showDeleteLoteDialog.value = true;
}

async function confirmDeleteLote() {
  if (!loteToDelete.value) return;
  const { loteId, productId } = loteToDelete.value;

  if (isEditMode.value) {
    // Find the product and lote
    const product = productStore.products.find((p) => p.id === productId);
    const loteIndex = product?.lotes?.findIndex((l) => l.id === loteId) ?? -1;

    if (product && product.lotes && loteIndex >= 0) {
      // Save the lote for tracking
      const lote = product.lotes[loteIndex];

      // Add to tracking for history
      loteChangesTracking.value.deleted.push({
        productId,
        loteId,
        lote: { ...lote }, // Store a copy
      });

      // Remove from the product
      product.lotes.splice(loteIndex, 1);

      toast.info(
        "Lote removido e será finalizado ao confirmar as atualizações."
      );
    }
  } else {
    // Normal flow when not in edit mode (should never happen with UI restrictions)
    await productStore.deleteLote(loteId, productId);
    if (!authStore.isLocalMode) historyStore.refreshHistory();
  }

  showDeleteLoteDialog.value = false;
  loteToDelete.value = null;
}

function closeDeleteLoteDialog() {
  showDeleteLoteDialog.value = false;
  loteToDelete.value = null;
}

const products = computed(() => productStore.products);
</script>

<template>
  <div>
    <!-- Action Buttons - Enhanced with more vibrant colors and larger icons -->
    <div class="flex flex-col sm:flex-row justify-end mb-4 gap-3">
      <template v-if="!isEditMode">
        <button
          @click="openAddProductForm"
          class="px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center font-medium"
        >
          <span class="material-icons-outlined mr-1.5 text-xl">add_circle</span>
          Novo Produto
        </button>
        <button
          @click="enableEditMode"
          class="px-4 py-2.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center font-medium"
        >
          <span class="material-icons-outlined mr-1.5 text-xl">edit</span>
          Atualizar Dados
        </button>
      </template>
      <template v-else>
        <button
          @click="cancelEdit"
          class="px-4 py-2.5 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center font-medium"
        >
          <span class="material-icons-outlined mr-1.5 text-xl">cancel</span>
          Cancelar
        </button>
        <button
          @click="confirmUpdates"
          class="px-4 py-2.5 bg-emerald-600 text-white rounded-lg hover:bg-emerald-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center font-medium"
        >
          <span class="material-icons-outlined mr-1.5 text-xl"
            >check_circle</span
          >
          Confirmar Atualizações
        </button>
      </template>
    </div>

    <!-- Add Product Form with enhanced styling -->
    <div
      v-if="isAddProductMode && !isEditMode"
      class="bg-white p-5 rounded-lg shadow-lg mb-6 border-l-4 border-emerald-500"
    >
      <h2 class="text-xl font-bold mb-4 text-gray-800 flex items-center">
        <span class="material-icons-outlined text-emerald-500 mr-2"
          >add_box</span
        >
        Adicionar Novo Produto
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Nome</label
          >
          <input
            v-model="newProduct.name"
            type="text"
            class="w-full input-field"
            placeholder="Nome do produto"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Unidade</label
          >
          <select v-model="newProduct.unit" class="w-full input-field">
            <option value="L">Litros (L)</option>
            <option value="kg">Quilogramas (kg)</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Qtd. Inicial (se sem lotes)</label
          >
          <input
            v-model.number="newProduct.quantity"
            type="number"
            min="0"
            class="w-full input-field"
          />
        </div>
      </div>
      <div class="mt-5 flex justify-end space-x-3">
        <button @click="cancelAddProduct" class="btn-secondary">
          Cancelar
        </button>
        <button
          @click="addProductHandler"
          class="btn-primary bg-emerald-600 hover:bg-emerald-700"
        >
          <span class="material-icons-outlined mr-1">add</span>
          Adicionar Produto
        </button>
      </div>
    </div>

    <!-- Products Table with enhanced styling -->
    <div class="overflow-x-auto rounded-lg shadow-lg border border-gray-200">
      <table class="min-w-full bg-white">
        <thead
          class="bg-gradient-to-r from-indigo-600 to-indigo-800 text-white"
        >
          <tr>
            <th class="p-4 w-12 text-center">
              <!-- Enhanced toggle column header -->
              <span class="material-icons-outlined text-indigo-200"
                >expand_more</span
              >
            </th>
            <th class="p-4 text-left">Produto</th>
            <th class="p-4 text-left">Unidade</th>
            <th class="p-4 text-left">Qtd. Total</th>
            <th class="p-4 text-left">Ações</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="product in products" :key="product.id">
            <!-- Product Row Component -->
            <ProductRow
              :product="product"
              :is-edit-mode="isEditMode"
              :expanded-products="expandedProducts"
              :temp-product-details="tempProductDetails"
              :temp-quantities="tempQuantities"
              @toggle-product-lotes="toggleProductLotes"
              @quantity-changed="changeProductQuantity"
              @quantity-updated="updateProductQuantityDirectly"
              @update-product-base-quantity="updateProductBaseQuantity"
              @request-delete="requestDeleteProduct"
            />

            <!-- Lotes Dropdown Component (expanded row) -->
            <tr
              v-if="expandedProducts[product.id]"
              class="bg-gradient-to-r from-indigo-50/80 to-indigo-50/50"
            >
              <td></td>
              <LoteDropdown
                :product="product"
                :is-edit-mode="isEditMode"
                @open-edit-lote="openEditLote"
                @request-delete-lote="requestDeleteLote"
                @open-add-lote="openAddLote"
              />
            </tr>
          </template>

          <!-- Empty state with enhanced styling -->
          <tr v-if="products.length === 0">
            <td colspan="5" class="p-8 text-center">
              <div
                class="flex flex-col items-center justify-center text-gray-500"
              >
                <span
                  class="material-icons-outlined text-6xl text-gray-300 mb-2"
                  >inventory_2</span
                >
                <p class="text-lg">Nenhum produto encontrado.</p>
                <p class="text-sm text-gray-400">
                  Adicione um novo produto para começar.
                </p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modals -->
    <AddLoteModal
      :show="showAddLoteModal"
      :product-id="currentProductIdForLote || ''"
      @close="showAddLoteModal = false"
      @save="handleSaveLote"
    />
    <EditLoteModal
      :show="showEditLoteModal"
      :lote="currentLoteToEdit"
      @close="showEditLoteModal = false"
      @save="handleUpdateLote"
    />

    <!-- Delete Product Dialog with enhanced styling -->
    <div
      v-if="showDeleteDialog"
      class="fixed inset-0 flex items-center justify-center z-50"
    >
      <div
        class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm"
        @click="closeDeleteDialog"
      ></div>
      <div
        class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 z-10 p-6 border-l-4 border-red-500"
      >
        <h3 class="text-lg font-medium text-gray-900 flex items-center">
          <span class="material-icons-outlined text-red-500 mr-2">warning</span>
          Confirmar exclusão
        </h3>
        <p class="mt-3 text-gray-600">
          Deseja remover o produto "<strong class="font-medium">{{
            productToDelete?.name
          }}</strong
          >"?
        </p>
        <p class="mt-1 text-sm text-red-500">
          Esta ação não pode ser desfeita.
        </p>
        <div class="mt-4 flex justify-end space-x-3">
          <button @click="closeDeleteDialog" class="btn-secondary-enhanced">
            <span class="material-icons-outlined text-sm mr-1">cancel</span>
            Cancelar
          </button>
          <button @click="confirmDeleteProduct" class="btn-danger-enhanced">
            <span class="material-icons-outlined text-sm mr-1">delete</span>
            Excluir
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Lote Dialog with enhanced styling -->
    <div
      v-if="showDeleteLoteDialog"
      class="fixed inset-0 flex items-center justify-center z-50"
    >
      <div
        class="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm"
        @click="closeDeleteLoteDialog"
      ></div>
      <div
        class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4 z-10 p-6 border-l-4 border-red-500"
      >
        <h3 class="text-lg font-medium text-gray-900 flex items-center">
          <span class="material-icons-outlined text-red-500 mr-2">warning</span>
          Confirmar exclusão de Lote
        </h3>
        <p class="mt-3 text-gray-600">
          Deseja remover o lote ID
          <strong class="font-mono">{{
            loteToDelete?.loteId.substring(0, 8)
          }}</strong
          >?
        </p>
        <p class="mt-1 text-sm text-red-500">
          Esta ação não pode ser desfeita.
        </p>
        <div class="mt-4 flex justify-end space-x-3">
          <button @click="closeDeleteLoteDialog" class="btn-secondary-enhanced">
            <span class="material-icons-outlined text-sm mr-1">cancel</span>
            Cancelar
          </button>
          <button @click="confirmDeleteLote" class="btn-danger-enhanced">
            <span class="material-icons-outlined text-sm mr-1">delete</span>
            Excluir Lote
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Enhanced styles */
.input-field-enhanced {
  @apply px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow;
}
.btn-primary {
  @apply px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition text-sm shadow-sm hover:shadow flex items-center font-medium;
}
.btn-secondary {
  @apply px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm shadow-sm hover:shadow;
}
.btn-secondary-enhanced {
  @apply px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm shadow-sm hover:shadow flex items-center font-medium;
}
.btn-danger-enhanced {
  @apply px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition text-sm shadow-sm hover:shadow flex items-center font-medium;
}
.input-field {
  @apply px-3 py-2 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500;
}
</style>
