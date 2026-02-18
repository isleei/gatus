import { createApp } from 'vue'
import App from './App.vue'
import './index.css'
import router from './router'
import { initI18n } from './i18n'

initI18n()

createApp(App).use(router).mount('#app')
