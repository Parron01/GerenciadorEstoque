import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isProd = mode === 'production'

  return {
    plugins: [
      vue(),
      // Apenas usar o DevTools em desenvolvimento e sem os botões flutuantes
      !isProd &&
        vueDevTools({
          // Desabilitar os botões flutuantes
          appendTo: 'iframe', // Isso remove os botões flutuantes da interface principal
          toggleButtonVisibility: false, // Esconde o botão de toggle
        }),
    ].filter(Boolean),
    server: {
      host: '0.0.0.0',
      port: 5173,
      strictPort: true,
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
  }
})
