<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const isLoggedIn = computed(() => authStore.isAuthenticated || authStore.isLocalMode)
const userDisplayName = computed(() => {
  if (authStore.isAuthenticated) {
    return authStore.user || 'Usuário'
  }
  return authStore.isLocalMode ? 'Modo Local' : ''
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <!-- Wrapper de toda a aplicação -->
  <div class="min-h-screen flex flex-col bg-gradient-to-br from-gray-50 to-indigo-50 text-gray-800">
    <!-- Top Bar -->
    <header
      v-if="isLoggedIn"
      class="flex items-center justify-between px-6 md:px-12 h-16 shadow-lg bg-white sticky top-0 z-50"
    >
      <h1
        class="text-xl font-bold tracking-wide text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-700"
      >
        📦 Estoque Simples
      </h1>

      <!-- Navegação -->
      <nav class="space-x-6 font-medium flex items-center">
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

        <div class="border-l border-gray-300 h-6 mx-4"></div>

        <!-- User info and logout -->
        <div class="flex items-center">
          <span class="text-sm text-gray-600 mr-3">
            <span class="material-icons-outlined text-sm mr-1 align-text-bottom">
              {{ authStore.isAuthenticated ? 'person' : 'cloud_off' }}
            </span>
            {{ userDisplayName }}
          </span>
          <button
            @click="handleLogout"
            class="bg-gray-200 hover:bg-gray-300 text-gray-700 px-3 py-1 rounded text-sm flex items-center"
          >
            <span class="material-icons-outlined text-sm mr-1">logout</span>
            Sair
          </button>
        </div>
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
