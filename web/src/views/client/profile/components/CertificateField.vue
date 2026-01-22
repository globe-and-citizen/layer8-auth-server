<template>
  <div class="flex items-center space-x-2">
    <input
      readonly
      class="bg-[#ECF4FD] border border-[#EADFD8] p-1 md:p-3 rounded-lg w-full font-medium"
      :value="value"
    />

    <button v-if="copied !== value" @click="$emit('copy', value)">
      📋
    </button>

    <input type="file" accept=".crt,.pem" @change="onUpload" />
  </div>
</template>

<script setup>
defineProps({
  value: String,
  copied: String,
})

const emit = defineEmits(["copy", "upload"])

const onUpload = (event) => {
  emit("upload", event.target.files[0])
}
</script>
