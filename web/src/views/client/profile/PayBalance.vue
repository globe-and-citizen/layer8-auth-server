<template>
  <div v-if="isLoggedIn">
    <div class="flex items-center gap-4">
      <span class="text-sm text-gray-700">
        Balance:
        <span class="font-semibold text-gray-900">
          {{ unpaidAmount }} USDC
        </span>
      </span>

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
    :unpaidAmount="unpaidAmount"
    :userID="props.userID"
    v-model:paymentAmount="paymentAmount"
    @cancel="cancelModal"
    @pay="payWithCrypto"
  />

  <CncPayCheckout
    :active="checkoutActive"
    :facture-id="props.userID"
    :amount="checkoutAmount"
    @close="closeCheckout"
    @success="onPaymentSuccess"
  />
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ClientUnpaidAmountPath, getAPI } from '@/api/paths.ts'
import PayWithCryptoModal from '@/views/client/profile/PayWithCryptoModal.vue'
import CncPayCheckout from '@/views/client/profile/CncPayCheckout.vue'

const props = defineProps<{
  userID: string
}>()

const token = ref<string | null>(localStorage.getItem('clientToken'))
const isLoggedIn = ref(!!token.value)

const modalWindowActive = ref(false)
const checkoutActive = ref(false)

const paymentAmount = ref('')
const checkoutAmount = ref('')
const unpaidAmount = ref('')

onMounted(async () => {
  if (!token.value) return
  await getOwingBalance()
})

async function getOwingBalance() {
  try {
    const response = await fetch(
      getAPI(ClientUnpaidAmountPath),
      {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      }
    )

    if (!response.ok) {
      throw new Error('Failed to fetch unpaid amount')
    }

    const body = await response.json()
    unpaidAmount.value = String(body.data.balance)
  } catch (error) {
    console.error(error)
    alert('Failed to fetch unpaid amount, retry later')
  }
}

function cancelModal() {
  modalWindowActive.value = false
}

function payWithCrypto() {
  const amount = Number(paymentAmount.value)

  if (!Number.isFinite(amount) || amount <= 0) {
    alert('Please enter a valid payment amount')
    return
  }

  // Preserve the amount used for this checkout.
  checkoutAmount.value = paymentAmount.value

  // Close the amount modal and open CNC Pay.
  modalWindowActive.value = false
  checkoutActive.value = true
}

function closeCheckout() {
  checkoutActive.value = false
}

async function onPaymentSuccess() {
  paymentAmount.value = ''
  await getOwingBalance()
}
</script>
