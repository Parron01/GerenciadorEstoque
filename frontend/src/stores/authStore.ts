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

  // Getters computados
  const isAuthenticated = computed(() => !!token.value);

  const API_BASE_URL =
    import.meta.env.VITE_API_BASE_URL || "http://localhost:3000";

  // Ações
  async function login(username: string, password: string) {
    authError.value = null;
    isLoading.value = true;

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
        signal: AbortSignal.timeout(10000), // 10 segundos de timeout
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || "Falha na autenticação");
      }

      token.value = data.token;
      user.value = data.user.username;

      // Salvar em localStorage
      localStorage.setItem("auth_token", data.token);
      localStorage.setItem("auth_user", data.user.username);

      // Toast de sucesso
      toast.success(`Bem-vindo, ${data.user.username}!`);

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

    // Limpar localStorage
    localStorage.removeItem("auth_token");
    localStorage.removeItem("auth_user");

    // Toast de informação
    toast.info("Você saiu do sistema");

    // Redirecionar para login
    router.push("/login");
  }

  async function verifyToken() {
    if (!token.value) return false;

    try {
      const response = await fetch(`${API_BASE_URL}/api/auth/verify`, {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
        signal: AbortSignal.timeout(5000), // 5 segundos de timeout
      });

      const data = await response.json();

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
    API_BASE_URL,
    login,
    logout,
    verifyToken,
  };
});
