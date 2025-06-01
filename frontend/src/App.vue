<script setup lang="ts">
import { RouterLink, RouterView } from "vue-router";
import { useAuthStore } from "@/stores/authStore";
import { computed, ref } from "vue";
import { useRouter } from "vue-router";

const authStore = useAuthStore();
const router = useRouter();
const isMenuOpen = ref(false);

const isLoggedIn = computed(() => authStore.isAuthenticated);
const userDisplayName = computed(() => authStore.user || "Usuário");

const handleLogout = () => {
  authStore.logout();
};

const navigateAndCloseMenu = () => {
  if (isMenuOpen.value) {
    isMenuOpen.value = false;
  }
};
</script>

<template>
  <!-- Wrapper de toda a aplicação -->
  <div
    class="min-h-screen flex flex-col bg-gradient-to-br from-gray-50 to-indigo-50 text-gray-800"
  >
    <!-- Top Bar -->
    <header
      v-if="isLoggedIn"
      class="flex items-center justify-between px-4 md:px-6 lg:px-12 h-16 shadow-lg bg-white sticky top-0 z-50"
    >
      <h1
        class="text-xl font-bold tracking-wide text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-purple-700"
      >
        📦 Estoque Simples
      </h1>

      <!-- Botão de hamburguer para telas pequenas -->
      <button
        @click="isMenuOpen = !isMenuOpen"
        class="md:hidden text-gray-700 focus:outline-none"
      >
        <span class="material-icons-outlined text-2xl">{{
          isMenuOpen ? "close" : "menu"
        }}</span>
      </button>

      <!-- Navegação para telas médias e grandes -->
      <nav class="hidden md:flex space-x-6 font-medium items-center">
        <RouterLink
          to="/"
          class="relative py-2 px-1 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:bg-indigo-600"
          @click="navigateAndCloseMenu"
        >
          Estoque
        </RouterLink>

        <RouterLink
          to="/history"
          class="relative py-2 px-1 hover:text-indigo-600 transition-colors"
          active-class="text-indigo-600 after:absolute after:bottom-0 after:left-0 after:h-0.5 after:w-full after:bg-indigo-600"
          @click="navigateAndCloseMenu"
        >
          Histórico
        </RouterLink>

        <div class="border-l border-gray-300 h-6 mx-4"></div>

        <!-- User info and logout -->
        <div class="flex items-center">
          <span class="text-sm text-gray-600 mr-3">
            <span class="material-icons-outlined text-sm mr-1 align-text-bottom"
              >person</span
            >
            {{ userDisplayName }}
          </span>
          <button
            @click="handleLogout"
            class="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-sm flex items-center transition-colors"
          >
            <span class="material-icons-outlined text-sm mr-1">logout</span>
            Sair
          </button>
        </div>
      </nav>
    </header>

    <!-- Menu móvel -->
    <div
      v-if="isLoggedIn && isMenuOpen"
      class="md:hidden fixed top-16 left-0 right-0 z-40 bg-white shadow-lg border-t border-gray-100"
    >
      <div class="flex flex-col p-4 space-y-3">
        <RouterLink
          to="/"
          @click="navigateAndCloseMenu"
          class="py-2 px-4 rounded hover:bg-indigo-50 flex items-center"
          active-class="bg-indigo-100 text-indigo-700"
        >
          <span class="material-icons-outlined mr-2">inventory</span>
          Estoque
        </RouterLink>

        <RouterLink
          to="/history"
          @click="navigateAndCloseMenu"
          class="py-2 px-4 rounded hover:bg-indigo-50 flex items-center"
          active-class="bg-indigo-100 text-indigo-700"
        >
          <span class="material-icons-outlined mr-2">history</span>
          Histórico
        </RouterLink>

        <div class="border-t border-gray-200 my-2"></div>

        <div class="flex items-center px-4 py-2">
          <span class="material-icons-outlined text-gray-500 mr-2">person</span>
          <span class="text-gray-700">{{ userDisplayName }}</span>
        </div>

        <button
          @click="handleLogout"
          class="flex items-center py-2 px-4 bg-red-600 hover:bg-red-700 text-white rounded transition-colors w-full"
        >
          <span class="material-icons-outlined mr-2">logout</span>
          Sair
        </button>
      </div>
    </div>

    <!-- Conteúdo principal -->
    <main class="flex-1">
      <RouterView />
    </main>

    <!-- Footer -->
    <footer class="bg-indigo-900 text-indigo-200 text-center py-4 text-sm">
      <p>© 2025 Estoque Simples - Gerenciamento de inventário</p>
    </footer>
  </div>
</template>

<style>
/* Estilos personalizados para os toasts */
.custom-toast-container {
  z-index: 9999;
}

.custom-toast-style {
  border-radius: 12px !important;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
  padding: 12px 16px !important;
  display: flex !important;
  align-items: center !important;
  border-left: 5px solid !important;
  overflow: visible !important;
  position: relative !important;
}

/* Estilo base para toasts com ícones melhorados */
.Vue-Toastification__toast {
  min-height: 64px !important;
  padding-left: 60px !important;
}

/* Posicionamento dos ícones personalizados */
.Vue-Toastification__toast::before {
  content: "" !important;
  position: absolute !important;
  left: 16px !important;
  top: 50% !important;
  transform: translateY(-50%) !important;
  width: 32px !important;
  height: 32px !important;
  background-size: contain !important;
  background-repeat: no-repeat !important;
  background-position: center !important;
  border-radius: 50% !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  z-index: 10 !important;
  font-family: "Material Icons Outlined" !important;
  font-size: 20px !important;
  color: white !important;
  text-align: center !important;
}

/* Estilos específicos para cada tipo de toast */
.Vue-Toastification__toast--success {
  background-color: #ffffff !important;
  color: #0f766e !important;
  border-color: #10b981 !important;
}

.Vue-Toastification__toast--success::before {
  content: "check_circle" !important;
  background-color: #10b981 !important;
}

.Vue-Toastification__toast--error {
  background-color: #ffffff !important;
  color: #b91c1c !important;
  border-color: #ef4444 !important;
}

.Vue-Toastification__toast--error::before {
  content: "error" !important;
  background-color: #ef4444 !important;
}

.Vue-Toastification__toast--warning {
  background-color: #ffffff !important;
  color: #a16207 !important;
  border-color: #f59e0b !important;
}

.Vue-Toastification__toast--warning::before {
  content: "warning" !important;
  background-color: #f59e0b !important;
}

.Vue-Toastification__toast--info {
  background-color: #ffffff !important;
  color: #1e40af !important;
  border-color: #3b82f6 !important;
}

.Vue-Toastification__toast--info::before {
  content: "info" !important;
  background-color: #3b82f6 !important;
}

/* Estilização da barra de progresso */
.Vue-Toastification__progress-bar {
  bottom: 0 !important;
  height: 4px !important;
  opacity: 0.7 !important;
  background: linear-gradient(
    to right,
    rgba(0, 0, 0, 0.1),
    rgba(0, 0, 0, 0.2)
  ) !important;
}

/* Ajustes responsivos para mobile */
@media (max-width: 480px) {
  .Vue-Toastification__toast {
    margin-bottom: 8px !important;
    width: calc(100% - 16px) !important;
    padding: 10px 10px 10px 55px !important;
  }

  .Vue-Toastification__toast::before {
    left: 12px !important;
    width: 28px !important;
    height: 28px !important;
    font-size: 18px !important;
  }

  .custom-toast-style {
    border-radius: 8px !important;
    border-left-width: 4px !important;
  }
}
</style>
