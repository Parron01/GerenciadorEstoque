import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

// Import Toast
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

// Toast options
const toastOptions = {
  position: 'top-right',
  timeout: 3000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: true,
  closeButton: 'button',
  icon: true,
  rtl: false,
  transition: 'Vue-Toastification__bounce',
  maxToasts: 3,
}

// Create app
const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(Toast, toastOptions)

app.mount('#app')
