<template>
  <div v-if="Number(unpaidAmountETH) > 0">
    <div class="flex items-center gap-4">

        <span class="text-sm text-gray-700">
          Balance:
          <span class="font-semibold text-gray-900">
            {{ unpaidAmountETH }} ETH
          </span>
        </span>

      <div></div>

      <button
        @click="modalWindowActive = true"
        class="text-sm font-medium text-indigo-600
                 hover:text-indigo-700 border bg-[#E5E5E5]"
      >
        <span class="m-4">Pay Now</span>
      </button>
    </div>
  </div>

  <w3m-button class="pr-4"></w3m-button>

  <PayWithCryptoModal
    :active="modalWindowActive"
    :unpaidAmountETH="unpaidAmountETH"
    :paymentAmount="paymentAmount"
    @cancel="cancelModel"
    @pay="payWithCrypto"
  ></PayWithCryptoModal>
</template>

<script setup lang="ts">
import PayWithCryptoModal from "@/views/client/profile/PayWithCryptoModal.vue";
import {onMounted, ref} from "vue";
import {ClientUnpaidAmountPath, getAPI} from "@/api/paths.ts";

const props = defineProps<{
  userID: string
}>()

/* =========================
   Reactive State (UNCHANGED)
   ========================= */

const token = ref<string | null>(localStorage.getItem('clientToken'))
const modalWindowActive = ref(false)
const unpaidAmount = ref<bigint>(0n)
const paymentAmount = ref('')
const unpaidAmountETH = ref('')
const walletConnected = ref(false)

/* =========================
   External Modules (lazy-loaded)
   ========================= */

let reconnect: any
let watchConnections: any
let parseEther: any
let formatEther: any
let WAGMI_CONFIG: any
let setupWeb3Modal: any
let payBill: any

/* =========================
   Lifecycle
   ========================= */

onMounted(async () => {
  if (!token.value) {
    return
  }

  const [
    wagmi,
    viem,
    web3modal,
    pay,
  ] = await Promise.all([
    // @ts-expect-error – remote ESM import, resolved at runtime
    import(/* @vite-ignore */ 'https://esm.sh/@wagmi/core@2.11.6'),

    // @ts-expect-error – remote ESM import, resolved at runtime
    import(/* @vite-ignore */ 'https://esm.sh/viem@2.17.0'),

    import('@/assets/js/crypto/web3modal.js'),
    import('@/assets/js/crypto/pay.js'),
  ])

  reconnect = wagmi.reconnect
  watchConnections = wagmi.watchConnections

  parseEther = viem.parseEther
  formatEther = viem.formatEther

  WAGMI_CONFIG = web3modal.WAGMI_CONFIG
  setupWeb3Modal = web3modal.setupWeb3Modal

  payBill = pay.payBill

  reconnect(WAGMI_CONFIG)
  await checkWalletConnections()
  setupWeb3Modal()

  await getOwingBalance()
})

async function getOwingBalance() {
  const unpaidAmountResponse = await fetch(
    getAPI(ClientUnpaidAmountPath),
    {
      headers: {
        Authorization: `Bearer ${token.value}`
      }
    }
  )

  if (unpaidAmountResponse.status !== 200) {
    alert('Failed to fetch unpaid amount, retry later')
    return
  }

  const unpaidAmountBody = await unpaidAmountResponse.json()

  unpaidAmount.value = BigInt(unpaidAmountBody.data.unpaid_amount)
  unpaidAmountETH.value = formatEther(unpaidAmount.value)
}

/* =========================
   Wallet Logic
   ========================= */

const checkWalletConnections = async () => {
  watchConnections(WAGMI_CONFIG, {
    onChange(data: any[]) {
      walletConnected.value = data.length !== 0
    },
  })
}

/* =========================
   Payment Logic (UNCHANGED)
   ========================= */

const cancelModel = async () => {
  modalWindowActive.value = false
}

const payWithCrypto = async () => {
  try {
    const currentAmount = parseEther(paymentAmount.value)

    if (currentAmount < unpaidAmount.value) {
      alert(
        'Too small payment amount, at least ' +
        unpaidAmount.value +
        ' must be paid'
      )
      return
    }

    const transactionHash = await payBill(
      '[[ .SmartContractAddress ]]',
      props.userID,
      currentAmount
    )

    modalWindowActive.value = false

    alert(
      'Invoice was paid successfully! You can track your transaction at https://polygonscan.com/tx/' +
      transactionHash
    )
  } catch (error) {
    alert('Payment failed, error: ' + error)
  }
}
</script>

<style scoped>

</style>
