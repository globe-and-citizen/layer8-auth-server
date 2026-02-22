<template>
  <label class="consent-label">
    <input
      type="checkbox"
      :checked="modelValue"
      @change="updateValue"
      class="consent-checkbox"
    />
    <span class="consent-text">
      {{ label }}
    </span>
  </label>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue'])

const updateValue = (e) => {
  emit('update:modelValue', e.target.checked)
}
</script>

<style scoped>
.consent-label {
  display: flex;
  white-space: nowrap;
  align-self: flex-start;
}

.consent-checkbox {
  margin-right: 5px;
}

.consent-text {
  font-size: 14px;
}
</style>
