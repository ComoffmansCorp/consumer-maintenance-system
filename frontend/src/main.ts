import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './style.css'
import App from './App.vue'
import router from './router'

const savedTheme = localStorage.getItem('cms.theme')
if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
  document.documentElement.classList.add('dark')
}

const app = createApp(App)
app.use(createPinia())
app.use(router)

// Wait for the initial navigation to resolve before mounting: otherwise
// App.vue's `route.meta.public` check runs against the pre-navigation
// route (meta = {}), briefly rendering AppShell on public pages like
// /login and firing authenticated-only requests with no token.
router.isReady().then(() => app.mount('#app'))
