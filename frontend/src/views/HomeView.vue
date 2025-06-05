<script setup lang="ts">
import ProductTable from "@/components/ProductTable.vue";
import { useProductStore } from "@/stores/productStore";
import { exportProductsToExcel } from "@/services/exportService"; // Import the service

const productStore = useProductStore();

function handleExportToExcel() {
  // Ensure products are loaded; you might want to add a loading check or fetch if empty
  if (productStore.products && productStore.products.length > 0) {
    exportProductsToExcel(productStore.products);
  } else {
    // Optionally, inform the user or fetch products if the store is empty
    alert(
      "Não há produtos para exportar ou os produtos ainda não foram carregados."
    );
  }
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-6xl">
    <!-- Cabeçalho da página -->
    <header class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-indigo-700 mb-2">
        Gerenciamento de Estoque
      </h1>
      <p class="text-gray-600">
        Controle entradas e saídas de produtos com facilidade
      </p>
      <div
        class="h-1 w-24 md:w-32 bg-gradient-to-r from-indigo-500 to-purple-600 mt-3 md:mt-4 rounded-full"
      ></div>
    </header>

    <!-- Botão de Exportar para Excel -->
    <div class="mb-4 flex justify-end">
      <button
        @click="handleExportToExcel"
        class="px-4 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-all shadow-md hover:shadow-lg flex items-center justify-center font-medium"
      >
        <span class="material-icons-outlined mr-1.5 text-xl"
          >file_download</span
        >
        Exportar para Excel
      </button>
    </div>

    <!-- Seção da tabela de produtos -->
    <section
      class="bg-white rounded-xl shadow-lg overflow-hidden border border-gray-200"
    >
      <div class="bg-gradient-to-r from-indigo-600 to-purple-700 p-3 md:p-4">
        <h2
          class="text-lg md:text-xl font-semibold text-white flex items-center"
        >
          <span class="material-icons-outlined mr-2">inventory</span>
          Produtos em Estoque
        </h2>
      </div>
      <div class="p-3 md:p-4">
        <ProductTable />
      </div>
    </section>
  </div>
</template>
