<template>
  <div class="relative">
    <div class="client-container">
      <div class="grid grid-cols-1 md:grid-cols-2 bg-white h-dvh">
        <div class="self-center">
          <h2 class="font-bold text-3xl text-[#3751FE] md:mb-24 mb-10 animate-slideFromLeft">
            Client Portal
          </h2>
          <h1 class="font-bold text-4xl text-[#3751FE] mb-11 animate-slideFromLeft">
            Login
          </h1>

          <div class="mr-0 md:mr-16 lg:mr-28 animate-slideFromLeft">
            <!-- Username -->
            <div class="relative border border-[#C1BBBB]">
              <input
                v-model="username"
                @focus="isUsernameFocused = true"
                @blur="isUsernameFocused = false"
                type="text"
                class="w-full px-4 pt-10 pb-3 border-l-4 focus:border-blue-500 focus:outline-none text-lg text-[#3751FE]"
                placeholder=" "
              />
              <label
                class="absolute left-0 px-4 mt-6 transition-all duration-300 origin-0 text-[#636363] text-lg cursor-text"
                :class="{ '-top-4': isUsernameFocused || username }"
              >
                Username
              </label>
            </div>

            <!-- Password -->
            <div class="relative border border-[#C1BBBB] mb-9">
              <input
                v-model="password"
                @focus="isPasswordFocused = true"
                @blur="isPasswordFocused = false"
                type="password"
                class="w-full px-4 pt-10 pb-3 border-l-4 focus:border-blue-500 focus:outline-none text-lg text-[#3751FE]"
                placeholder=" "
              />
              <label
                class="absolute left-0 px-4 mt-6 transition-all duration-300 origin-0 text-[#636363] text-lg cursor-text"
                :class="{ '-top-4': isPasswordFocused || password }"
              >
                Password
              </label>
            </div>

            <button
              @click="loginClient"
              class="animate-bounce w-full py-4 border border-[#3751FE] text-[#3751FE] mb-7 hover:shadow-lg hover:text-white hover:bg-[#3751FE]"
            >
              Login
            </button>

            <a href="/client-register" class="text-sm text-[#414141]">
              Don’t have an account? <span class="font-bold">Register</span>
            </a>
          </div>
        </div>

        <div class="hidden"></div>
      </div>
    </div>

    <!-- Right image -->
    <div class="hidden md:flex items-center absolute bg-[#E5E5E5] right-0 top-0 pt-40 w-1/2 h-dvh">
      <img
        class="m-auto mt-10 animate-slideFromRight"
        src="@/assets/images/client-image.png"
      />
    </div>

    <!-- Toast -->
    <div
      class="fixed top-3 right-3 bg-red-500 text-white p-2 rounded-md transition-opacity duration-500 z-50"
      :class="showToast ? 'opacity-100' : 'opacity-0 pointer-events-none'"
    >
      {{ toastMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue"
import {scram} from "@/utils/scram.ts"
import {ClientLoginPath, ClientLoginPrecheckPath, getAPI} from "@/api/paths.ts";

const username = ref("")
const password = ref("")
const isUsernameFocused = ref(false)
const isPasswordFocused = ref(false)
const showToast = ref(false)
const toastMessage = ref("")
const cNonce = ref("")

const showToastMessage = (message: string) => {
  toastMessage.value = message
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

const loginClient = async () => {
  try {
    if (!username.value || !password.value) {
      showToastMessage("Please enter a username and password!")
      return
    }

    // scram is loaded globally
    // @ts-ignore
    cNonce.value = scram.generateCnonce()

    const precheckRes = await fetch(
      getAPI(ClientLoginPrecheckPath),
      {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          username: username.value,
          c_nonce: cNonce.value,
        }),
      }
    )

    if (!precheckRes.ok) {
      showToastMessage("Failed to login")
      return
    }

    const precheckBody = await precheckRes.json()

    // @ts-ignore
    const {data} = scram.keysHMAC(
      password.value,
      precheckBody.data.salt,
      precheckBody.data.iteration_count
    )

    // @ts-ignore
    const clientKeyBytes = scram.hexStringToBytes(data.clientKey)

    const authMessage = `[n=${username.value},r=${cNonce.value},s=${precheckBody.data.salt},i=${precheckBody.data.iteration_count},r=${precheckBody.data.nonce}]`

    // @ts-ignore
    const clientSignature = scram.signatureHMAC(authMessage, data.storedKey)
    // @ts-ignore
    const clientSignatureBytes = scram.hexStringToBytes(clientSignature)
    // @ts-ignore
    const clientProofBytes = scram.xorBytes(clientKeyBytes, clientSignatureBytes)
    // @ts-ignore
    const clientProof = scram.bytesToHexString(clientProofBytes)

    const loginRes = await fetch(getAPI(ClientLoginPath), {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        username: username.value,
        nonce: precheckBody.data.nonce,
        c_nonce: cNonce.value,
        client_proof: clientProof,
      }),
    })

    const loginBody = await loginRes.json()

    if (loginBody.data?.verifier) {
      // @ts-ignore
      const serverSigCheck = scram.signatureHMAC(
        authMessage,
        data.serverKey
      )

      if (serverSigCheck === loginBody.data.verifier) {
        localStorage.setItem("clientToken", loginBody.data.token)
        showToastMessage("Login successful!")
        window.location.href = "/client/profile"
      }
    } else {
      showToastMessage(loginBody.message || "Login failed")
    }
  } catch (err) {
    console.error(err)
    showToastMessage("Login failed!")
  }
}
</script>
