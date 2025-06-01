<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useAuthStore } from "@/stores/authStore";
import { useRouter } from "vue-router";

const username = ref("");
const password = ref("");
const authStore = useAuthStore();
const router = useRouter();
const showPassword = ref(false);
const isConnectionError = ref(false);

const API_BASE_URL = authStore.API_BASE_URL;

onMounted(() => {
  if (authStore.isAuthenticated) {
    router.push("/");
  }
  testConnection();
});

async function testConnection() {
  try {
    const response = await fetch(`${API_BASE_URL}/api/auth/health`, {
      signal: AbortSignal.timeout(2000),
    }).catch(() => null);

    isConnectionError.value = !response || !response.ok;
  } catch (error) {
    isConnectionError.value = true;
  }
}

async function handleLogin() {
  if (!username.value || !password.value) return;
  await authStore.login(username.value, password.value);
}

function toggleShowPassword() {
  showPassword.value = !showPassword.value;
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 p-4"
  >
    <div
      class="w-full max-w-md p-6 md:p-8 space-y-6 md:space-y-8 bg-white rounded-xl shadow-lg"
    >
      <!-- Logo e título -->
      <div class="text-center">
        <h1
          class="text-2xl md:text-3xl font-bold text-indigo-700 flex items-center justify-center"
        >
          <span class="text-3xl md:text-4xl mr-2">📦</span> Estoque Simples
        </h1>
        <p class="mt-2 text-gray-600">Faça login para gerenciar seu estoque</p>
      </div>

      <!-- Aviso de erro de conexão -->
      <div
        v-if="isConnectionError"
        class="p-3 md:p-4 bg-red-100 text-red-800 rounded border border-red-200 flex items-start"
      >
        <span class="material-icons-outlined mr-2 text-red-600">error</span>
        <div>
          <p class="font-medium">Servidor indisponível</p>
          <p class="text-sm">
            Não foi possível conectar ao servidor. Verifique se o backend está
            em execução.
          </p>
        </div>
      </div>

      <!-- Formulário -->
      <form
        class="mt-6 md:mt-8 space-y-5 md:space-y-6"
        @submit.prevent="handleLogin"
      >
        <!-- Mensagem de erro -->
        <div
          v-if="authStore.authError"
          class="p-3 bg-red-100 text-red-700 rounded border border-red-200 flex items-center"
        >
          <span class="material-icons-outlined mr-2 text-red-600">error</span>
          {{ authStore.authError }}
        </div>

        <!-- Campo de usuário -->
        <div>
          <label
            for="username"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Nome de usuário
          </label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            :disabled="isConnectionError"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="Digite seu nome de usuário"
          />
        </div>

        <!-- Campo de senha -->
        <div>
          <label
            for="password"
            class="block text-sm font-medium text-gray-700 mb-1"
          >
            Senha
          </label>
          <div class="relative">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              required
              :disabled="isConnectionError"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
              placeholder="Digite sua senha"
            />
            <button
              type="button"
              @click="toggleShowPassword"
              class="absolute inset-y-0 right-0 px-3 flex items-center text-gray-500 hover:text-gray-700"
            >
              <span class="material-icons-outlined text-sm">
                {{ showPassword ? "visibility_off" : "visibility" }}
              </span>
            </button>
          </div>
        </div>

        <!-- Botão de login -->
        <div>
          <button
            type="submit"
            :disabled="isConnectionError || authStore.isLoading"
            class="w-full py-2.5 md:py-3 px-4 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 disabled:cursor-not-allowed text-white rounded-md transition flex items-center justify-center"
          >
            <span
              v-if="authStore.isLoading"
              class="material-icons-outlined animate-spin mr-2"
              >autorenew</span
            >
            <span>{{ authStore.isLoading ? "Entrando..." : "Entrar" }}</span>
          </button>
        </div>
      </form>

      <!-- Nota de informação -->
      <div class="mt-5 text-center text-xs text-gray-500">
        <p>Sistema requer autenticação para acesso.</p>
      </div>
    </div>
  </div>
</template>
