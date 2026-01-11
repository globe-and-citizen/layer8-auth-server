<template>
  <div class="modal__dialog" :class="{ active }">
    <div class="modal__content">

      <!-- Header -->
      <h3 class="modal__header pb-3">
        <b>Pay for the traffic consumption using your crypto wallet</b>
      </h3>

      <!-- Amount table -->
      <table class="payment-table">
        <tbody>
        <tr>
          <td class="label p-3">Unpaid amount (ETH)</td>
          <td class="value readonly">
            <span class="m-4">{{ unpaidAmountETH }}</span>
          </td>
        </tr>

        <tr>
          <td class="label p-3">Payment amount</td>
          <td class="value">
            <input
              class="table-input m-4 border"
              :value="paymentAmount"
              @input="$emit('update:paymentAmount', $event.target.value)"
            />
          </td>
        </tr>
        </tbody>
      </table>

      <!-- Footer -->
      <div class="modal__footer">
        <button class="btn-secondary hover:bg-[#4F80E1]" @click="$emit('cancel')">Cancel</button>
        <button class="bg-[#4F80E1] text-white" @click="$emit('pay')">Pay</button>
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

defineEmits(["cancel", "pay", "update:paymentAmount"])
</script>
