<script setup lang="ts">
import ProductTable from '@/components/ProductTable.vue'
</script>

<template>
  <div class="container mx-auto px-4 py-8 max-w-6xl">
    <!-- Cabeçalho da página -->
    <header class="mb-8">
      <h1 class="text-3xl font-bold text-indigo-700 mb-2">Gerenciamento de Estoque</h1>
      <p class="text-gray-600">Controle entradas e saídas de produtos com facilidade</p>
      <div class="h-1 w-32 bg-gradient-to-r from-indigo-500 to-purple-600 mt-4 rounded-full"></div>
    </header>

    <!-- Seção da tabela de produtos -->
    <section class="bg-white rounded-xl shadow-lg overflow-hidden border border-gray-200">
      <div class="bg-gradient-to-r from-indigo-600 to-purple-700 p-4">
        <h2 class="text-xl font-semibold text-white flex items-center">
          <span class="material-icons-outlined mr-2">inventory</span>
          Produtos em Estoque
        </h2>
      </div>
      <div class="p-4">
        <ProductTable />
      </div>
    </section>
  </div>
</template>
