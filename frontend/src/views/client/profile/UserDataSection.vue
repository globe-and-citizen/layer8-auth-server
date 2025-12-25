<template>
  <section class="bg-white rounded-2xl py-3 md:py-4 px-4 md:px-6 mb-6 md:mb-0">
    <h1 class="font-medium text-lg md:text-xl text-black mb-2 md:mb-4">Your data</h1>
    <p class="font-normal text-sm md:text-base text-[#8E8E93] mb-6 md:mb-12">
      Your product data to use on your own project
    </p>

    <div class="grid grid-cols-3">
      <ul class="col-span-1 flex flex-col space-y-6 my-auto">
        <li class="font-bold md:text-xl text-sm text-black">Name:</li>
        <li class="font-bold md:text-xl text-sm text-black">Redirect URI:</li>
        <li class="font-bold md:text-xl text-sm text-black">Backend URI:</li>
        <li class="font-bold md:text-xl text-sm text-black">UUID:</li>
        <li class="font-bold md:text-xl text-sm text-black">Secret:</li>
        <li class="font-bold md:text-xl text-sm text-black">X509 Certificate:</li>
      </ul>

      <div class="col-span-2 flex flex-col space-y-2">
        <InputField :value="user.name"/>
        <InputField :value="user.redirect_uri"/>
        <InputField :value="user.backend_uri"/>

        <div class="flex items-center space-x-2">
          <CopyField :value="user.id" :copied="isCopied" @copy="$emit('copy', user.id)"/>
        </div>

        <div class="flex items-center space-x-2">
          <CopyField :value="user.secret" :copied="isCopied" @copy="$emit('copy', user.secret)"/>
        </div>

        <div class="flex items-center space-x-2">
          <CopyField
            :value="user.x509_certificate"
            :copied="isCopied"
            @copy="$emit('copy', user.x509_certificate)"
          />
          <input type="file" accept=".crt,.pem" @change="onUpload"/>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import CopyField from "@/views/client/profile/components/CopyField.vue";
import InputField from "@/views/client/profile/components/ReadonlyField.vue";

const props = defineProps({
  user: Object,
  isCopied: String,
});

const emit = defineEmits(["copy", "upload-cert"]);

const labels = [
  "Name:",
  "Redirect URI:",
  "Backend URI:",
  "UUID:",
  "Secret:",
  "X509 Certificate:",
];

const onUpload = (e) => {
  emit("upload-cert", e.target.files[0]);
};
</script>
