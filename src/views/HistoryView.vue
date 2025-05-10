<script setup lang="ts">
import HistoryList from '@/components/HistoryList.vue'
import { ref } from 'vue'

// Add state for selected filter
const selectedFilter = ref('all')
</script>

<template>
  <div class="container mx-auto px-4 py-8 max-w-6xl">
    <!-- Cabeçalho da página -->
    <header class="mb-8">
      <h1 class="text-3xl font-bold text-indigo-700 mb-2">Histórico de Alterações</h1>
      <p class="text-gray-600">Acompanhe todas as entradas e saídas do estoque</p>
      <div class="h-1 w-32 bg-gradient-to-r from-indigo-500 to-purple-600 mt-4 rounded-full"></div>
    </header>

    <!-- Seção de filtros (opcional) -->
    <section class="mb-8 bg-white p-4 rounded-lg shadow border border-gray-200">
      <div class="flex flex-wrap gap-4 items-end">
        <div class="flex-1 min-w-[200px]">
          <label class="block text-sm font-medium text-gray-700 mb-1">Período</label>
          <select
            v-model="selectedFilter"
            class="w-full border border-gray-300 rounded px-3 py-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="today">Hoje</option>
            <option value="week">Esta semana</option>
            <option value="month">Este mês</option>
            <option value="all">Todo o período</option>
          </select>
        </div>
        <button class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded transition">
          Filtrar
        </button>
      </div>
    </section>

    <!-- Seção da tabela de histórico -->
    <section class="bg-white rounded-xl shadow-lg overflow-hidden border border-gray-200">
      <div class="bg-gradient-to-r from-slate-700 to-slate-800 p-4">
        <h2 class="text-xl font-semibold text-white flex items-center">
          <span class="mr-2">📋</span>
          Registro de Movimentações
        </h2>
      </div>
      <div class="p-4">
        <HistoryList :filter-option="selectedFilter" />
      </div>
    </section>
  </div>
</template>
