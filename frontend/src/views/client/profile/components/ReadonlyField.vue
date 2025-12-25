<template>
  <input
    class="bg-white border border-[#EADFD8] p-1 md:p-3 rounded-lg w-full font-medium"
    :value="value"
    readonly
  />
</template>

<script setup>
defineProps({ value: String })
</script>
