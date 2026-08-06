<template>
  <div>
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
  </div>
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
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.consent-checkbox {
  appearance: auto;

  width: 18px !important;
  height: 18px !important;

  min-width: 18px;
  min-height: 18px;

  margin: 0;
  padding: 0;

  flex-shrink: 0;
}

.consent-text {
  font-size: 14px;
}
</style>
