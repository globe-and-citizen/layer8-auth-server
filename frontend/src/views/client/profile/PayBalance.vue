<template>
  <div>
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
    v-model:paymentAmount="paymentAmount"
    @cancel="cancelModel"
    @pay="payWithCrypto"
  ></PayWithCryptoModal>
</template>

<script setup lang="ts">
import {onMounted, ref, watch} from "vue";
import {parseEther} from "viem";
import {ClientUnpaidAmountPath, getAPI} from "@/api/paths.ts";
import PayWithCryptoModal from "@/views/client/profile/PayWithCryptoModal.vue";

// Import the new hook
import {useTrafficContract} from "@/utils/paywithcrypto/usePaymentContract.ts";

const props = defineProps<{
  userID: string
}>()

/* =========================
   Reactive State
   ========================= */
const token = ref<string | null>(localStorage.getItem('clientToken'))
const modalWindowActive = ref(false)
const unpaidAmount = ref<bigint>(0n)
const paymentAmount = ref('')
const unpaidAmountETH = ref('')

const {
  connection,
  payTraffic,
  isPending,
  isSuccess,
  writeError,
  hash
} = useTrafficContract()

/* =========================
   Lifecycle & Data Fetching
   ========================= */
onMounted(async () => {
  if (!token.value) return
  await getOwingBalance()
})

async function getOwingBalance() {
  const unpaidAmountResponse = await fetch(
    getAPI(ClientUnpaidAmountPath),
    {
      headers: {Authorization: `Bearer ${token.value}`}
    }
  )

  if (unpaidAmountResponse.status !== 200) {
    alert('Failed to fetch unpaid amount, retry later')
    return
  }

  const unpaidAmountBody = await unpaidAmountResponse.json()
  unpaidAmountETH.value = unpaidAmountBody.data.balance_eth
}

/* =========================
   Payment Logic
   ========================= */
const cancelModel = async () => {
  modalWindowActive.value = false
}

const payWithCrypto = async () => {
  try {
    await payTraffic(props.userID, paymentAmount.value)
  } catch (error) {
    alert('Payment failed to initiate: ' + error)
  }
}

/* =========================
   Watchers for Feedback
   ========================= */
// Watch for successful transaction hash generation
watch(isSuccess, (newSuccess) => {
  if (newSuccess && hash.value) {
    modalWindowActive.value = false
    console.log('txhash:', hash.value)
    alert(
      'Transaction sent successfully! You can track it at txhash:' +
      hash.value
    )
  }
})

// Watch for contract errors
watch(writeError, (newError) => {
  if (newError) {
    alert('Transaction error: ' + (newError as any).shortMessage || newError.message)
  }
})
</script>

<style scoped>

</style>
