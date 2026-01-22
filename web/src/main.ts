import {createApp} from 'vue'
import App from './App.vue'
import router from './router'

import '@/assets/styles/output.css'
import '@/assets/styles/base.css'
import '@/assets/styles/modal.css'
// import '@/assets/styles/oauth.css'
import "@/utils/paywithcrypto/web3modal.ts"

import { WagmiPlugin } from '@wagmi/vue'
import { QueryClient, VueQueryPlugin } from '@tanstack/vue-query'
import { web3Config } from '@/utils/paywithcrypto/web3modal.ts' // Path to your file

const app = createApp(App)
// 1. Create a Query Client (Required for Wagmi hooks)
const queryClient = new QueryClient()

// 2. Provide the config to the whole app
app.use(WagmiPlugin, { config: web3Config })
app.use(VueQueryPlugin, { queryClient })
app.use(router)
app.mount('#app')
