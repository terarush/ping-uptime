import './assets/globals.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// toggle dark mode based on system preference
const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
document.documentElement.classList.toggle('dark', mediaQuery.matches)
mediaQuery.addEventListener('change', (e) => {
  document.documentElement.classList.toggle('dark', e.matches)
})

const app = createApp(App)
app.use(router)
app.mount('#app')
