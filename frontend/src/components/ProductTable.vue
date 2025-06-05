<script setup lang="ts">
import { useProductStore } from "@/stores/productStore";
import { useHistoryStore } from "@/stores/historyStore";
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import type { Product } from "@/models/product";
import type { Lote, LotePayload } from "@/models/lote";
import { v4 as uuidv4 } from "uuid";
import { useToast } from "vue-toastification";
import AddLoteModal from "./AddLoteModal.vue";
import EditLoteModal from "./EditLoteModal.vue";
import ProductRow from "./product-table/ProductRow.vue";
import LoteDropdown from "./product-table/LoteDropdown.vue";
import type { ProductBatchContextPayload } from "@/models/history"; // Import new type

const productStore = useProductStore();
const historyStore = useHistoryStore();
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
const tempProductDetails = ref<
  Record<string, { name: string; unit: "L" | "kg" }>
>({});

// Store initial product states (name and total quantity) for batch context history
const initialProductStatesForBatch = ref<
  Record<string, { name: string; totalQuantity: number }>
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
  // If a product has no lotes, its quantity is assumed to be 0.
  return 0; // Was product.quantity, now strictly 0 if no lotes
}

function initTempStates() {
  productsBeforeEdit.value = JSON.parse(JSON.stringify(productStore.products));
  initialProductStatesForBatch.value = {}; // Clear previous batch states

  productStore.products.forEach((product) => {
    tempProductDetails.value[product.id] = {
      name: product.name,
      unit: product.unit,
    };
    // Capture initial state for product batch context history
    initialProductStatesForBatch.value[product.id] = {
      name: product.name,
      totalQuantity: getProductDisplayQuantity(product), // Calculated from lotes
    };
  });
}

function enableEditMode() {
  initTempStates();
  isEditMode.value = true;
}

const loteChangesTracking = ref<{
  created: { productId: string; loteData: LotePayload; localId: string }[];
  updated: {
    productId: string;
    loteId: string;
    loteData: LotePayload;
    originalLote: Lote;
  }[];
  deleted: { productId: string; loteId: string; originalLote: Lote }[];
}>({
  created: [],
  updated: [],
  deleted: [],
});

function closeDeleteDialog() {
  showDeleteDialog.value = false;
  productToDelete.value = null;
}

function openAddLote(productId: string) {
  if (!isEditMode.value) {
    toast.info("Ative o modo de edição para adicionar lotes.");
    return;
  }
  currentProductIdForLote.value = productId;
  showAddLoteModal.value = true;
}

function openEditLote(lote: Lote) {
  if (!isEditMode.value) {
    toast.info("Ative o modo de edição para editar lotes.");
    return;
  }
  currentLoteToEdit.value = JSON.parse(JSON.stringify(lote));
  showEditLoteModal.value = true;
}

function requestDeleteLote(loteId: string, productId: string) {
  if (!isEditMode.value) {
    toast.info("Ative o modo de edição para remover lotes.");
    return;
  }
  loteToDelete.value = { loteId, productId };
  showDeleteLoteDialog.value = true;
}

function closeDeleteLoteDialog() {
  showDeleteLoteDialog.value = false;
  loteToDelete.value = null;
}

const isComponentMounted = ref(true);

onBeforeUnmount(() => {
  isComponentMounted.value = false;
});

async function confirmUpdates() {
  const operationBatchId = uuidv4();
  let allApiCallsSuccessful = true;

  // Collect product IDs affected in this batch for context history later
  const affectedProductIdsInBatch = new Set<string>();
  loteChangesTracking.value.created.forEach((item) =>
    affectedProductIdsInBatch.add(item.productId)
  );
  loteChangesTracking.value.updated.forEach((item) =>
    affectedProductIdsInBatch.add(item.productId)
  );
  loteChangesTracking.value.deleted.forEach((item) =>
    affectedProductIdsInBatch.add(item.productId)
  );
  productStore.products.forEach((p) => {
    const originalProduct = productsBeforeEdit.value.find(
      (op) => op.id === p.id
    );
    if (originalProduct) {
      const editedName = tempProductDetails.value[p.id]?.name;
      const editedUnit = tempProductDetails.value[p.id]?.unit;
      if (
        (editedName && editedName !== originalProduct.name) ||
        (editedUnit && editedUnit !== originalProduct.unit)
      ) {
        affectedProductIdsInBatch.add(p.id);
      }
    }
  });

  // Delete lotes first
  for (const { productId, loteId } of loteChangesTracking.value.deleted) {
    try {
      if (!isComponentMounted.value) return;
      await productStore.deleteLote(loteId, productId, operationBatchId);
    } catch (e) {
      // The productStore.deleteLote action now handles its own specific toast on error.
      // We just need to mark that not all API calls were successful.
      if (!isComponentMounted.value) return;
      allApiCallsSuccessful = false;
      console.error(
        `[ProductTable] Error during productStore.deleteLote for lote ${loteId.substring(0, 6)}:`,
        e
      );
      // Removed toast.error here to avoid double toasting, store action handles it.
    }
  }

  // Update existing lotes
  for (const { productId, loteId, loteData } of loteChangesTracking.value
    .updated) {
    try {
      if (!isComponentMounted.value) return;
      await productStore.updateLote(
        loteId,
        productId,
        loteData,
        operationBatchId
      );
    } catch (e) {
      if (!isComponentMounted.value) return;
      allApiCallsSuccessful = false;
      toast.error(`Erro ao atualizar lote ${loteId.substring(0, 6)}: ${e}`);
    }
  }

  // Create new lotes
  for (const { productId, loteData } of loteChangesTracking.value.created) {
    try {
      if (!isComponentMounted.value) return;
      await productStore.createLote(productId, loteData, operationBatchId);
    } catch (e) {
      if (!isComponentMounted.value) return;
      allApiCallsSuccessful = false;
      toast.error(`Erro ao criar novo lote para ${productId}: ${e}`);
    }
  }

  // Update product details
  for (const product of productStore.products) {
    const originalProduct = productsBeforeEdit.value.find(
      (p) => p.id === product.id
    );
    if (!originalProduct) continue;

    const editedName = tempProductDetails.value[product.id]?.name;
    const editedUnit = tempProductDetails.value[product.id]?.unit;

    const productUpdatePayload: Partial<Pick<Product, "name" | "unit">> = {};
    let productDetailsChanged = false;

    if (editedName && editedName !== originalProduct.name) {
      productUpdatePayload.name = editedName;
      productDetailsChanged = true;
    }
    if (editedUnit && editedUnit !== originalProduct.unit) {
      productUpdatePayload.unit = editedUnit;
      productDetailsChanged = true;
    }

    if (productDetailsChanged) {
      try {
        if (!isComponentMounted.value) return;
        await productStore.updateProductDetails(
          product.id,
          productUpdatePayload,
          operationBatchId
        );
      } catch (e) {
        allApiCallsSuccessful = false;
        // Ensure product ID is added if update fails but was attempted
        affectedProductIdsInBatch.add(product.id); // Ensure affected ID is tracked
        toast.error(`Erro ao atualizar produto ${originalProduct.name}: ${e}`);
      }
    }
  }

  // After all individual operations, record product batch context history
  // This calculation happens *after* optimistic updates to productStore.products by lote CUD operations
  // and *before* fetching fresh data from the server.
  for (const productId of affectedProductIdsInBatch) {
    const initialProductState = initialProductStatesForBatch.value[productId];
    // Find the product in the current store state, which reflects optimistic updates
    const currentProductInStore = productStore.products.find(
      (p) => p.id === productId
    );

    if (initialProductState && currentProductInStore) {
      // Use the potentially updated name from tempProductDetails for the snapshot
      const finalName =
        tempProductDetails.value[productId]?.name || initialProductState.name;
      // Recalculate quantity based on the optimistically updated lotes in the store
      const finalQuantity = getProductDisplayQuantity(currentProductInStore);

      const productBatchContextData: ProductBatchContextPayload = {
        productId: productId,
        productNameSnapshot: finalName,
        quantityBeforeBatch: initialProductState.totalQuantity,
        quantityAfterBatch: finalQuantity,
      };

      try {
        if (!isComponentMounted.value) return;
        await historyStore.createProductBatchContextHistory(
          productBatchContextData,
          operationBatchId
        );
      } catch (e) {
        allApiCallsSuccessful = false;
        toast.error(
          `Erro ao registrar contexto de lote para produto ${finalName}: ${e}`
        );
      }
    } else if (initialProductState && !currentProductInStore) {
      // This case implies the product was deleted in this batch.
      // The 'product_deleted' history record will cover its state before deletion.
      // A product batch context might still be relevant if other changes happened before deletion,
      // but for simplicity, we'll rely on the product deletion history.
      // Or, we could create a context record indicating it was deleted.
      // For now, we only create context if the product still exists in the store post-optimistic updates.
      console.warn(
        `Product ${productId} was in initial states but not in current store; likely deleted.`
      );
    }
  }

  if (allApiCallsSuccessful) {
    toast.success("Alterações enviadas ao servidor com sucesso!");
  } else {
    toast.warning(
      "Algumas alterações falharam. Verifique os logs e tente novamente."
    );
  }

  await productStore.fetchProductsFromApi();
  await historyStore.refreshHistory();

  if (!isComponentMounted.value) return;
  isEditMode.value = false;
  productsBeforeEdit.value = [];
  loteChangesTracking.value = { created: [], updated: [], deleted: [] };
  expandedProducts.value = {};
  tempProductDetails.value = {};
  initialProductStatesForBatch.value = {}; // Clear for next batch
}

function cancelEdit() {
  productStore.fetchProductsFromApi();
  isEditMode.value = false;
  productsBeforeEdit.value = [];
  loteChangesTracking.value = { created: [], updated: [], deleted: [] };
  expandedProducts.value = {};
  initialProductStatesForBatch.value = {}; // Clear on cancel
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
  newProduct.value = { name: "", unit: "L", quantity: 0 };
}

async function addProductHandler() {
  if (!newProduct.value.name) {
    toast.error("Nome do produto é obrigatório.");
    return;
  }
  await productStore.addProduct(newProduct.value);
  await historyStore.refreshHistory();
  isAddProductMode.value = false;
  newProduct.value = { name: "", unit: "L", quantity: 0 };
}

function requestDeleteProduct(product: Product) {
  if (!isEditMode.value) {
    productToDelete.value = product;
    showDeleteDialog.value = true;
  }
}

async function confirmDeleteProduct() {
  if (productToDelete.value) {
    await productStore.removeProduct(productToDelete.value.id);
    await historyStore.refreshHistory();
    closeDeleteDialog();
  }
}

async function handleSaveLote(loteData: LotePayload) {
  if (!currentProductIdForLote.value || !isEditMode.value) return;

  const localLoteId = uuidv4();
  loteChangesTracking.value.created.push({
    productId: currentProductIdForLote.value,
    loteData: loteData,
    localId: localLoteId,
  });

  const product = productStore.products.find(
    (p) => p.id === currentProductIdForLote.value
  );
  if (product) {
    if (!product.lotes) product.lotes = [];
    product.lotes.push({
      ...loteData,
      id: localLoteId,
      productId: currentProductIdForLote.value,
      createdAt: new Date().toISOString(),
    });
  }
  toast.success("Lote adicionado e será salvo ao confirmar as atualizações.");
  showAddLoteModal.value = false;
  currentProductIdForLote.value = null;
}

async function handleUpdateLote(loteId: string, loteData: LotePayload) {
  if (!currentLoteToEdit.value || !isEditMode.value) return;
  const productId = currentLoteToEdit.value.productId;

  const productForOriginal = productStore.products.find(
    (p) => p.id === productId
  );
  const originalLote = productForOriginal?.lotes?.find((l) => l.id === loteId);
  if (!originalLote) {
    toast.error("Lote original não encontrado para atualização.");
    return;
  }

  loteChangesTracking.value.updated = loteChangesTracking.value.updated.filter(
    (t) => !(t.productId === productId && t.loteId === loteId)
  );
  loteChangesTracking.value.updated.push({
    productId,
    loteId,
    loteData,
    originalLote: JSON.parse(JSON.stringify(originalLote)),
  });

  const product = productStore.products.find((p) => p.id === productId);
  const loteToUpdate = product?.lotes?.find((l) => l.id === loteId);
  if (loteToUpdate) {
    Object.assign(loteToUpdate, loteData, {
      updatedAt: new Date().toISOString(),
    });
  }
  toast.success("Lote atualizado e será salvo ao confirmar as atualizações.");
  showEditLoteModal.value = false;
  currentLoteToEdit.value = null;
}

async function confirmDeleteLote() {
  if (!loteToDelete.value || !isEditMode.value) return;
  const { loteId, productId } = loteToDelete.value;

  const productForOriginal = productStore.products.find(
    (p) => p.id === productId
  );
  const originalLote = productForOriginal?.lotes?.find((l) => l.id === loteId);
  if (!originalLote) {
    toast.error("Lote original não encontrado para exclusão.");
    return;
  }

  loteChangesTracking.value.deleted.push({
    productId,
    loteId,
    originalLote: JSON.parse(JSON.stringify(originalLote)),
  });
  loteChangesTracking.value.created = loteChangesTracking.value.created.filter(
    (t) => !(t.productId === productId && t.localId === loteId)
  );
  loteChangesTracking.value.updated = loteChangesTracking.value.updated.filter(
    (t) => !(t.productId === productId && t.loteId === loteId)
  );

  const product = productStore.products.find((p) => p.id === productId);
  if (product && product.lotes) {
    const loteIndex = product.lotes.findIndex((l) => l.id === loteId);
    if (loteIndex !== -1) {
      product.lotes.splice(loteIndex, 1);
    }
  }
  toast.info("Lote removido e será finalizado ao confirmar as atualizações.");
  showDeleteLoteDialog.value = false;
  loteToDelete.value = null;
}

const products = computed(() => productStore.products);
const sortedProducts = computed(() => {
  return [...productStore.products].sort((a, b) =>
    a.name.localeCompare(b.name)
  );
});

onMounted(() => {
  productStore.initializeStore();
  historyStore.refreshHistory();
});
</script>

<template>
  <div>
    <!-- Action Buttons -->
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

    <!-- Add Product Form -->
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
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
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

    <!-- Products Table -->
    <div class="overflow-x-auto rounded-lg shadow-lg border border-gray-200">
      <table class="min-w-full bg-white">
        <thead
          class="bg-gradient-to-r from-indigo-600 to-indigo-800 text-white"
        >
          <tr>
            <th class="p-4 w-12 text-center">
              <span class="material-icons-outlined text-indigo-200"
                >expand_more</span
              >
            </th>
            <th class="p-4 text-left">Produto</th>
            <th class="p-4 text-left text-center">Qtd. Total</th>
            <th class="p-4 text-left text-center">Unidade</th>
            <th class="p-4 text-left text-center">Ações</th>
          </tr>
        </thead>
        <tbody>
          <template v-for="product in sortedProducts" :key="product.id">
            <ProductRow
              :product="product"
              :is-edit-mode="isEditMode"
              :expanded-products="expandedProducts"
              :temp-product-details="tempProductDetails"
              @toggle-product-lotes="toggleProductLotes"
              @request-delete="requestDeleteProduct"
            />
            <tr
              v-if="expandedProducts[product.id]"
              class="bg-gradient-to-r from-indigo-50/80 to-indigo-50/50"
            >
              <td></td>
              <!-- Empty cell for alignment under the expand icon -->
              <LoteDropdown
                :product="product"
                :is-edit-mode="isEditMode"
                @open-edit-lote="openEditLote"
                @request-delete-lote="requestDeleteLote"
                @open-add-lote="openAddLote"
              />
            </tr>
          </template>

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

    <!-- Delete Product Dialog -->
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

    <!-- Delete Lote Dialog -->
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
