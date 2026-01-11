import {createApp} from 'vue'
import App from './App.vue'
import router from './router'

import '@/assets/styles/output.css'
import '@/assets/styles/base.css'
import '@/assets/styles/modal.css'
// import '@/assets/styles/oauth.css'

const app = createApp(App)
app.config.compilerOptions.isCustomElement = tag => tag === 'w3m-button'
app.use(router)
app.mount('#app')
