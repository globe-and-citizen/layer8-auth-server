declare module '@/assets/js/crypto/web3modal.js' {
  export const WAGMI_CONFIG: any
  export function setupWeb3Modal(): void
}

declare module '@/assets/js/crypto/pay.js' {
  export function payBill(
    contractAddress: string,
    userID: string,
    amount: bigint
  ): Promise<string>
}
