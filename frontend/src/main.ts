import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import '@/assets/styles/output.css'
import '@/assets/styles/base.css'
import '@/assets/styles/modal.css'

import "@/assets/js/scram-bundled.js"

const app = createApp(App)

app.use(router)

app.mount('#app')
