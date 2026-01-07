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
        <InputField :placeholder="labels.name" :value="user.name"/>
        <InputField :placeholder="labels.redirect_uri" :value="user.redirect_uri"/>
        <InputField :placeholder="labels.backend_uri" :value="user.backend_uri"/>

        <div class="flex items-center space-x-2">
          <CopyField
            :placeholder="labels.id"
            :value="user.id"
            :copied="isCopied"
            @copy="copyToClipboard"
          />
        </div>

        <div class="flex items-center space-x-2">
          <CopyField
            :placeholder="labels.secret"
            :value="user.secret"
            :copied="isCopied"
            @copy="copyToClipboard"
          />
        </div>

        <div class="flex items-center space-x-2">
          <CopyField
            :placeholder="labels.x509_certificate"
            :value="user.x509_certificate"
            :copied="isCopied"
            @copy="copyToClipboard"
          />
          <input type="file" accept=".crt,.pem" @change="handleX509CertificateUpload"/>
        </div>
      </div>
    </div>

    <div :class="showToast ? 'opacity-100' : 'opacity-0 pointer-events-none'"
         class="fixed top-3 right-3 bg-green-500 text-white p-2 rounded-md transition-opacity ease-in-out duration-500 z-50">
      {{ toastMessage }}
    </div>
  </section>
</template>

<script setup>
import CopyField from "@/views/client/profile/components/CopyField.vue";
import InputField from "@/views/client/profile/components/ReadonlyField.vue";
import {onMounted, ref} from "vue";
import {ClientProfilePath, getAPI} from "@/api/paths.js";

const token = ref(localStorage.getItem("clientToken"))
const toastMessage = ref("")
const showToast = ref(false)

const labels = {
  name: "Name:",
  redirect_uri: "Redirect URI:",
  backend_uri: "Backend URI:",
  id: "UUID:",
  secret: "Secret:",
  x509_certificate: "X509 Certificate:",
}

const isCopied = ref();
const user = ref({
  id: "",
  name: "",
  secret: "",
  redirect_uri: "",
  backend_uri: "",
  x509_certificate: "",
})

const emit = defineEmits(["user-data"]);


// const onUpload = (e) => {
//   emit("upload-cert", e.target.files[0]);
// };

const showToastMessage = (msg) => {
  toastMessage.value = msg
  showToast.value = true
  setTimeout(() => (showToast.value = false), 3000)
}

const copyToClipboard = async (text) => {
  isCopied.value = text
  await navigator.clipboard.writeText(text)
  showToastMessage("Copied to clipboard")
}

const handleX509CertificateUpload = async (e) => {
  const cert = await e.target.files[0].text()
  user.value.x509_certificate = cert
  showToastMessage("Certificate uploaded")
}

onMounted(async () => {
  if (!token.value) {
    window.location.href = "/client-login"
    return
  }

  const userResp = await fetch(getAPI(ClientProfilePath),
    {
      // method: "GET",
      headers: {
        // ContentType: "application/json",
        Authorization: `Bearer ${token.value}`
      },
    })

  if (!userResp.ok) {
    showToastMessage("Failed to login")
    return
  }

  user.value = (await userResp.json()).data
  emit('user-data', user.value.name)
  console.log(user)
})

</script>
