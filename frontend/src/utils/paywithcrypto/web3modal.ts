import { createAppKit } from '@reown/appkit/vue'
import { WagmiAdapter } from '@reown/appkit-adapter-wagmi'
import { sepolia } from '@reown/appkit/networks'
// import { injected } from '@wagmi/vue/connectors'

// 1. Get projectId at https://cloud.reown.com
export const projectId = import.meta.env.VITE_WALLET_PROJECT_ID as `0x${string}`

// 2. Create Wagmi Adapter
export const networks = [sepolia]

export const wagmiAdapter = new WagmiAdapter({
  projectId,
  networks,
  // connectors: [injected()]
})

// 3. Create modal
createAppKit({
  adapters: [wagmiAdapter],
  networks: [sepolia],
  projectId,
  defaultNetwork: sepolia, // Ensures the app starts on Sepolia
  metadata: {
    name: 'Layer8 Traffic Usage Payment',
    description: 'Reverse Proxy Implementing the Layer8 protocol',
    url: window.location.origin,
    icons: ['https://avatars.githubusercontent.com/u/37784886']
  },
  features: {
    analytics: false // Disabling analytics can sometimes bypass AdBlocker triggers
  },
  enableReconnect: true
})

export const config = wagmiAdapter.wagmiConfig
