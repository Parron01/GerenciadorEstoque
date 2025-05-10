<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <!-- Wrapper de toda a aplicação -->
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-gray-50 to-indigo-50 text-gray-800">
    <!-- Top Bar -->
    <header
      class="flex items-center justify-between px-6 md:px-12 h-16 shadow-lg bg-white sticky top-0 z-50"
    >
      <h1
        class="text-xl font-bold tracking-wide text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-700"
      >
        📦 Estoque Simples
      </h1>

      <!-- Navegação -->
      <nav class="space-x-6 font-medium">
        <RouterLink
          to="/"
          class="relative py-2 px-1 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:bg-indigo-600"
        >
          Estoque
        </RouterLink>

        <RouterLink
          to="/history"
          class="relative py-2 px-1 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:bg-indigo-600"
        >
          Histórico
        </RouterLink>
      </nav>
    </header>

    <!-- Conteúdo principal -->
    <main class="flex-1">
      <RouterView />
    </main>

    <!-- Footer -->
    <footer class="bg-indigo-900 text-indigo-200 text-center py-4 text-sm">
      <p>© 2025 Estoque Simples - Gerenciamento de inventário simplificado</p>
    </footer>
  </div>
</template>
