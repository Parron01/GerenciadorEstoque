<script setup lang="ts">
import HistoryList from "@/components/HistoryList.vue";
import { ref, onMounted, computed, watch, onUnmounted } from "vue";
import { useHistoryStore } from "@/stores/historyStore";
import { useProductStore } from "@/stores/productStore"; // Import product store

const selectedFilter = ref("all");
const historyStore = useHistoryStore();
const productStore = useProductStore(); // Add product store

// Product dropdown state
const isProductDropdownOpen = ref(false);
const productSearchQuery = ref("");
const selectedProductName = ref("");
const dropdownRef = ref<HTMLDivElement | null>(null);

// Update onMounted to also load products if needed
onMounted(() => {
  historyStore.fetchGroupedHistory(1, historyStore.pageSizeForGrouped);
  // Ensure products are loaded for the filter dropdown
  if (productStore.products.length === 0) {
    productStore.fetchProductsFromApi();
  }

  // Add click outside listener to close dropdown
  document.addEventListener("click", handleClickOutside);
});

onUnmounted(() => {
  // Clean up the event listener
  document.removeEventListener("click", handleClickOutside);
});

// Handle clicks outside the dropdown to close it
function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isProductDropdownOpen.value = false;
  }
}

// New computed property for sorted products
const sortedProducts = computed(() => {
  return [...productStore.products].sort((a, b) =>
    a.name.localeCompare(b.name)
  );
});

// Filtered products based on search query
const filteredDropdownProducts = computed(() => {
  if (!productSearchQuery.value) return sortedProducts.value;

  return sortedProducts.value.filter((product) =>
    product.name.toLowerCase().includes(productSearchQuery.value.toLowerCase())
  );
});

// Toggle dropdown visibility - simplified for reliability
function toggleProductDropdown() {
  // Simply toggle the dropdown state
  isProductDropdownOpen.value = !isProductDropdownOpen.value;

  // Reset search when opening
  if (isProductDropdownOpen.value) {
    productSearchQuery.value = "";
  }
}

// Select a product from dropdown
function selectProduct(productId: string, productName: string) {
  historyStore.setProductFilter(productId);
  selectedProductName.value = productName;
  isProductDropdownOpen.value = false;
}

// Update selectedProductName when productFilter changes
watch(
  () => historyStore.productFilter,
  (newFilter) => {
    if (!newFilter) {
      selectedProductName.value = "";
    } else {
      const product = productStore.products.find((p) => p.id === newFilter);
      selectedProductName.value = product?.name || "";
    }
  },
  { immediate: true }
);

// Function to handle filter updates
function updateFilter(value: string) {
  selectedFilter.value = value;
}

// Function to clear all filters
function clearAllFilters() {
  selectedFilter.value = "all";
  historyStore.clearProductFilter();
  selectedProductName.value = "";
}

// Check if any filter is active
const isAnyFilterActive = computed(() => {
  return selectedFilter.value !== "all" || historyStore.productFilter !== "";
});
</script>

<template>
  <div class="container mx-auto p-4 sm:p-6 lg:p-8">
    <header class="mb-6">
      <h1 class="text-3xl font-bold text-gray-800">Histórico de Alterações</h1>
      <p class="text-gray-600">
        Visualize todas as movimentações e alterações no estoque.
      </p>
    </header>

    <!-- Enhanced Filters -->
    <div class="mb-6 p-4 bg-white rounded-lg shadow border border-gray-200">
      <h3 class="text-lg font-medium text-gray-700 mb-3 flex items-center">
        <span class="material-icons-outlined mr-2 text-indigo-600"
          >filter_list</span
        >
        Filtrar Histórico
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Date Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Período</label
          >
          <div class="flex items-center space-x-2">
            <div
              class="flex items-center justify-center w-10 h-10 bg-indigo-100 rounded-lg"
            >
              <span class="material-icons-outlined text-indigo-600 text-sm"
                >event</span
              >
            </div>
            <div class="relative flex-1">
              <select
                v-model="selectedFilter"
                class="w-full px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow appearance-none pr-8"
              >
                <option value="all">Todos os períodos</option>
                <option value="today">Hoje</option>
                <option value="week">Esta Semana</option>
                <option value="month">Este Mês</option>
              </select>
              <span
                class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none text-gray-400"
              >
                <span class="material-icons-outlined text-sm">expand_more</span>
              </span>
            </div>
          </div>
        </div>

        <!-- Enhanced Product Filter with Custom Dropdown - FIXED -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1"
            >Produto</label
          >
          <div class="flex items-center space-x-2">
            <div
              class="flex items-center justify-center w-10 h-10 bg-indigo-100 rounded-lg"
            >
              <span class="material-icons-outlined text-indigo-600 text-sm"
                >inventory_2</span
              >
            </div>

            <!-- Custom dropdown container -->
            <div class="relative flex-1" ref="dropdownRef">
              <!-- Dropdown trigger button -->
              <button
                type="button"
                @click.stop="toggleProductDropdown"
                class="w-full px-3 py-2 border border-indigo-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-shadow cursor-pointer bg-white flex items-center justify-between"
                :class="{ 'ring-2 ring-indigo-500': isProductDropdownOpen }"
              >
                <span v-if="selectedProductName" class="text-gray-900">{{
                  selectedProductName
                }}</span>
                <span v-else class="text-gray-500">Todos os produtos</span>
                <span class="material-icons-outlined text-gray-400 text-sm">
                  {{ isProductDropdownOpen ? "expand_less" : "expand_more" }}
                </span>
              </button>

              <!-- Dropdown clear button -->
              <button
                v-if="historyStore.productFilter"
                type="button"
                @click.stop="
                  historyStore.clearProductFilter();
                  selectedProductName = '';
                "
                class="absolute inset-y-0 right-8 flex items-center pr-2 text-gray-400 hover:text-gray-600"
              >
                <span class="material-icons-outlined text-sm">close</span>
              </button>

              <!-- Dropdown menu -->
              <div
                v-show="isProductDropdownOpen"
                class="absolute z-10 mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg"
              >
                <!-- Search input -->
                <div class="p-2 border-b border-gray-200">
                  <div class="relative">
                    <span
                      class="absolute inset-y-0 left-0 pl-2 flex items-center text-gray-500"
                    >
                      <span class="material-icons-outlined text-sm"
                        >search</span
                      >
                    </span>
                    <input
                      id="product-dropdown-search"
                      v-model="productSearchQuery"
                      type="text"
                      class="w-full pl-8 pr-2 py-1.5 border border-gray-300 rounded text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                      placeholder="Buscar produto..."
                      @click.stop
                    />
                  </div>
                </div>

                <!-- Scrollable options list -->
                <div class="max-h-60 overflow-y-auto">
                  <!-- "All products" option -->
                  <div
                    @click.stop="selectProduct('', '')"
                    class="px-3 py-2 hover:bg-indigo-50 cursor-pointer flex items-center"
                    :class="{
                      'bg-indigo-100': historyStore.productFilter === '',
                    }"
                  >
                    <span class="text-gray-800">Todos os produtos</span>
                  </div>

                  <!-- Product options -->
                  <template v-if="filteredDropdownProducts.length > 0">
                    <div
                      v-for="product in filteredDropdownProducts"
                      :key="product.id"
                      @click.stop="selectProduct(product.id, product.name)"
                      class="px-3 py-2 hover:bg-indigo-50 cursor-pointer"
                      :class="{
                        'bg-indigo-100':
                          historyStore.productFilter === product.id,
                      }"
                    >
                      {{ product.name }}
                    </div>
                  </template>

                  <!-- No results message -->
                  <div
                    v-else-if="productSearchQuery"
                    class="px-3 py-3 text-gray-500 text-center text-sm italic"
                  >
                    Nenhum produto encontrado
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Clear Filters Button -->
        <div class="flex items-end">
          <button
            @click="clearAllFilters"
            class="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition text-sm shadow-sm hover:shadow flex items-center font-medium w-full justify-center"
            :disabled="!isAnyFilterActive"
            :class="{ 'opacity-50 cursor-not-allowed': !isAnyFilterActive }"
          >
            <span class="material-icons-outlined mr-1">clear</span>
            Limpar Filtros
          </button>
        </div>
      </div>

      <!-- Active Filters Display -->
      <div class="text-sm text-gray-600 mt-3" v-if="isAnyFilterActive">
        <span class="font-medium">Filtros ativos:</span>
        <span v-if="selectedFilter !== 'all'" class="ml-1">
          {{
            selectedFilter === "today"
              ? "Hoje"
              : selectedFilter === "week"
                ? "Esta semana"
                : selectedFilter === "month"
                  ? "Este mês"
                  : ""
          }}
        </span>
        <span v-if="selectedFilter !== 'all' && historyStore.productFilter"
          >,
        </span>
        <span v-if="selectedProductName" class="ml-1">
          {{ selectedProductName }}
        </span>
      </div>
    </div>

    <HistoryList
      :filter-option="selectedFilter"
      @update:filter-option="updateFilter"
    />
  </div>
</template>

<style scoped>
/* Remove default select arrow in some browsers */
select::-ms-expand {
  display: none;
}

/* Custom scrollbar for dropdown */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background-color: #f1f1f1;
  border-radius: 100px;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background-color: #d1d5db;
  border-radius: 100px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background-color: #9ca3af;
}
</style>
