// src/config.ts
type RuntimeConfig = {
  BASE_API_URL: string;
  CONTRACT_ADDRESS: string;
  WALLET_PROJECT_ID: string;
};

declare global {
  interface Window {
    __APP_CONFIG__?: Partial<RuntimeConfig>;
  }
}

const runtime = window.__APP_CONFIG__;

export const config: RuntimeConfig = {
  BASE_API_URL:
    runtime?.BASE_API_URL ?? import.meta.env.VITE_BASE_API_URL,

  CONTRACT_ADDRESS:
    runtime?.CONTRACT_ADDRESS ??
    import.meta.env.VITE_CONTRACT_ADDRESS,

  WALLET_PROJECT_ID:
    runtime?.WALLET_PROJECT_ID ??
    import.meta.env.VITE_WALLET_PROJECT_ID,
};

// Optional safety check
Object.entries(config).forEach(([k, v]) => {
  if (!v && k !== 'BASE_API_URL') {
    throw new Error(`Missing config value: ${k}`);
  }
});
