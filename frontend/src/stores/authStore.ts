import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";

export const useAuthStore = defineStore("auth", () => {
  const router = useRouter();
  const toast = useToast();

  // Estado
  const token = ref<string | null>(localStorage.getItem("auth_token"));
  const user = ref<string | null>(localStorage.getItem("auth_user"));
  const authError = ref<string | null>(null);
  const isLoading = ref(false);
  const isLocalMode = ref(localStorage.getItem("local_mode") === "true");

  // Getters computados
  const isAuthenticated = computed(() => !!token.value && !isLocalMode.value);

  const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL || "http://localhost:3000";
  // Ações
  async function login(username: string, password: string) {
    authError.value = null;
    isLoading.value = true;
    isLocalMode.value = false;
    localStorage.removeItem("local_mode");

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
        // Adicionar timeout para evitar espera indefinida
        signal: AbortSignal.timeout(10000), // 10 segundos de timeout
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Falha na autenticação");
      }

      // Corrigido: extrair username do objeto user retornado pela API
      token.value = data.token;
      user.value = data.user.username; // <-- Modificação aqui

      // Salvar em localStorage
      localStorage.setItem("auth_token", data.token);
      localStorage.setItem("auth_user", data.user.username); // <-- Modificação aqui

      // Toast de sucesso
      toast.success(`Bem-vindo, ${data.user.username}!`, {
        // <-- Modificação aqui
        icon: "check_circle",
      });

      // Redirecionar para página inicial
      router.push("/");
      return true;
    } catch (error) {
      console.error("Erro de login:", error);
      authError.value =
        error instanceof Error
          ? error.message
          : "Erro desconhecido de autenticação";
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  function logout() {
    // Limpar dados de autenticação
    token.value = null;
    user.value = null;
    isLocalMode.value = false;

    // Limpar localStorage
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");
    localStorage.removeItem("local_mode");

    // Limpar dados locais dos produtos e histórico
    localStorage.removeItem("estoque_produtos_local");
    localStorage.removeItem("estoque_historico_local");

    // Toast de informação
    toast.info("Você saiu do sistema", {
      icon: "exit_to_app",
    });

    // Redirecionar para login
    router.push("/login");
  }

  function useLocalMode() {
    // Configurar para modo local apenas (sem autenticação)
    token.value = null;
    user.value = null;
    isLocalMode.value = true;

    // Limpar tokens de autenticação e definir modo local
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");
    localStorage.setItem("local_mode", "true");

    // Toast informativo
    toast.info("Usando modo local - dados de demonstração", {
      icon: "cloud_off",
    });

    // Redirecionar para home
    router.push("/");
  }

  // Verificar token
  async function verifyToken() {
    if (!token.value || isLocalMode.value) return false;

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/verify`, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
        // Adicionar timeout para evitar espera indefinida
        signal: AbortSignal.timeout(5000), // 5 segundos de timeout
      });

      const data = await response.json();

      // Se o token for inválido, faça logout
      if (!data.valid) {
        logout();
        return false;
      }

      return data.valid;
    } catch (error) {
      console.error("Erro ao verificar token:", error);
      return false;
    }
  }

  return {
    token,
    user,
    authError,
    isLoading,
    isAuthenticated,
    isLocalMode,
    login,
    logout,
    useLocalMode,
    verifyToken,
  };
});
