import "./assets/main.css";

import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

// Import Toast
import Toast, { POSITION, type PluginOptions } from "vue-toastification";
import "vue-toastification/dist/index.css";

// Toast options com configurações melhoradas e design moderno
const toastOptions: PluginOptions = {
  position: POSITION.TOP_RIGHT,
  timeout: 4000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: true,
  hideProgressBar: false,
  closeButton: "button",
  icon: false, // Desativamos os ícones padrão para usar nossos ícones personalizados via CSS
  rtl: false,
  transition: "Vue-Toastification__fade",
  maxToasts: 3,
  // Classes personalizadas para melhorar a estética
  toastClassName: "custom-toast-style",
  bodyClassName: "custom-toast-body",
  containerClassName: "custom-toast-container",
};

// Create app
const app = createApp(App);
app.use(createPinia());
app.use(router);
app.use(Toast, toastOptions);

app.mount("#app");
