import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'
import { router } from './router'
import { i18n } from './i18n'
import './style.css'

const app = createApp(App)
app.config.errorHandler = (err, _instance, info) => {
  console.error('[Vue Error]', err, info)
}
app.use(createPinia())
app.use(router)
app.use(naive)
app.use(i18n)
app.mount('#app')
