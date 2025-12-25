<template>
  <div class="modal__dialog" :class="{ active }">
    <div class="modal__content">
      <h3 class="modal__header">
        Pay for the traffic consumption using your crypto wallet!
      </h3>

      <div class="grid grid-cols-3 gap-4">
        <label>Unpaid amount (ETH)</label>
        <input :value="unpaidAmountETH" readonly />

        <label>Payment amount</label>
        <input
          :value="paymentAmount"
          @input="$emit('update:paymentAmount', $event.target.value)"
        />
      </div>

      <div class="modal__footer">
        <button class="close-button" @click="$emit('cancel')">Cancel</button>
        <button class="bg-[#4F80E1] text-white" @click="$emit('pay')">
          Pay
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  active: Boolean,
  unpaidAmountETH: String,
  paymentAmount: String,
});
</script>
