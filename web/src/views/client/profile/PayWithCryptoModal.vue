<template>
  <div class="modal__dialog" :class="{ active }">
    <div class="modal__content">

      <h3 class="modal__header pb-3">
        <b>Pay for the traffic consumption using your crypto wallet</b>
      </h3>

      <table class="payment-table">
        <tbody>
        <tr>
          <td class="label p-3">Unpaid amount (USDC)</td>
          <td class="value readonly">
            <span class="m-4">{{ unpaidAmount }}</span>
          </td>
        </tr>

        <tr>
          <td class="label p-3">Payment amount (USDC)</td>
          <td class="value">
            <input
              class="table-input m-4 border"
              type="number"
              min="0.000001"
              step="any"
              :value="paymentAmount"
              @input="$emit('update:paymentAmount', $event.target.value)"
            />
          </td>
        </tr>
        </tbody>
      </table>

      <div class="modal__footer">
        <button
          class="btn-secondary hover:bg-[#4F80E1]"
          @click="$emit('cancel')"
        >
          Cancel
        </button>

        <button
          class="bg-[#4F80E1] text-white"
          :disabled="!paymentAmount || Number(paymentAmount) <= 0"
          @click="$emit('pay')"
        >
          Pay
        </button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  active: boolean
  unpaidAmount: string
  paymentAmount: string
  userID: string
}>()

defineEmits<{
  (e: 'cancel'): void
  (e: 'pay'): void
  (e: 'update:paymentAmount', value: string): void
}>()
</script>
