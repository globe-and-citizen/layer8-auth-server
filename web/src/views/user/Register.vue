<template>
  <div>
    <!-- NAVBAR -->
    <div class="user-container">
      <div class="bg-white justify-center items-center my-4">
        <img src="@/assets/images/L8Logo.png" width="250" height="535"/>
      </div>
    </div>

    <!-- BODY -->
    <div class="md:grid md:grid-cols-2 min__width">
      <!-- LEFT -->
      <div class="mx-10 lg:mx-36">
        <h1 class="font-bold text-[#4F80E1] text-[40px] mb-2">Register</h1>
        <p class="text-xl text-[#414141] mb-5">Enter registration details.</p>

        <div>
          <div class="mb-3">
            <label class="text-sm text-[#414141] mb-1 block">Username</label>
            <input
                class="w-full rounded-md border border-[#EADFD8] py-2.5 px-3"
                v-model="registerUsername"
                placeholder="Username"
            />
          </div>

          <div class="mb-3">
            <label class="text-sm text-[#414141] mb-1 block">Password</label>
            <input
                class="w-full rounded-md border border-[#EADFD8] py-2.5 px-3"
                type="password"
                v-model="registerPassword"
                placeholder="Password"
            />
          </div>

          <div class="mb-12">
            <label class="text-sm text-[#414141] mb-1 block"
            >Confirm password</label
            >
            <input
                class="w-full rounded-md border border-[#EADFD8] py-2.5 px-3"
                type="password"
                v-model="confirmedPassword"
                placeholder="Confirm Password"
            />
          </div>

          <button
              class="w-full bg-[#4F80E1] rounded-lg text-white py-4 mb-4"
              @click="registerUser"
          >
            Register
          </button>

          <a
              href="/user-login-page"
              class="text-sm text-[#414141] block text-center"
          >
            Already have an account?
            <span class="font-bold">Login</span>
          </a>
        </div>
      </div>

      <!-- RIGHT -->
      <div class="bg-white hidden md:flex items-center">
        <img src="@/assets/images/cyber-computer.png"/>
      </div>
    </div>

    <!-- SUCCESS MODAL -->
    <div class="modal__dialog" :class="{ active: modalWindowActive }">
      <div class="modal__content">
        <h3 class="modal__header">
          Congratulations! You were registered successfully!
        </h3>

        <div class="modal__body">
          <div class="info-message">
            This is your 12-word recovery phrase:
          </div>

          <div
              class="toast-msg text-white rounded-md transition-opacity"
              :class="mnemonicCopied ? 'opacity-100' : 'opacity-0 pointer-events-none'"
          >
            Copied!
          </div>

          <div class="mnemonic_holder">
            <input
                class="input-mnemonic"
                readonly
                :value="currMnemonic"
            />
            <button @click="copyToClipboard">
              📋
            </button>
          </div>

          <div class="warning-msg">
            Save it somewhere safe and never share it with anybody!
          </div>
        </div>

        <div class="modal__footer">
          <button
              class="bg-[#4F80E1] rounded-lg text-white py-4 w-full"
              @click="backToLogin"
          >
            Got it!
          </button>
        </div>
      </div>
    </div>

    <div class="modal__overlay"></div>

    <!-- FOOTER -->
    <div class="user-container">
      <div class="bg-white flex justify-between items-center my-8">
        <div>
          <img
              src="@/assets/images/L8Logo.png"
              class="mb-6 h-[35px] md:h-[70px]"
          />
          <p class="font-bold text-sm">©Layer8security 2023.</p>
        </div>

        <div>
          <div class="text-xl font-bold mb-4">Contact</div>
          <ul class="text-sm">
            <li>Email: hi@layer8.com</li>
            <li>Client Support: support@layer8.com</li>
            <li>Phone number: 0371 525 777</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- TOAST -->
    <div
        class="fixed top-3 right-3 bg-red-500 text-white p-2 rounded-md transition-opacity z-50"
        :class="showToast ? 'opacity-100' : 'opacity-0 pointer-events-none'"
    >
      {{ toastMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue"
import {getAPI, UserRegisterPath, UserRegisterPrecheckPath} from "@/api/paths.ts";
import {mnemonic} from "@/utils/mnemonic.ts"
import {scram} from "@/utils/scram.ts"

const registerUsername = ref("")
const registerPassword = ref("")
const confirmedPassword = ref("")

const showToast = ref(false)
const toastMessage = ref("")

const currMnemonic = ref("")
const modalWindowActive = ref(false)
const mnemonicCopied = ref(false)

const copyToClipboard = async () => {
  await navigator.clipboard.writeText(currMnemonic.value)
  mnemonicCopied.value = true

  setTimeout(() => {
    mnemonicCopied.value = false
  }, 1500)
}

const backToLogin = () => {
  window.location.href = "/user-login"
}

const showToastMessage = (message: string) => {
  toastMessage.value = message
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

const registerUser = async () => {
  try {
    if (
        registerUsername.value === "" ||
        registerPassword.value === "" ||
        confirmedPassword.value === ""
    ) {
      showToastMessage("Please enter all details!")
      return
    }

    if (registerPassword.value !== confirmedPassword.value) {
      showToastMessage("Passwords do not match")
      return
    }

    // @ts-ignore – provided by bundled.js
    currMnemonic.value = mnemonic.generateBip39Mnemonic()

    // @ts-ignore – provided by bundled.js
    const keyPair = mnemonic.getPrivateAndPublicKeys(currMnemonic.value)

    const precheckResp = await fetch(getAPI(UserRegisterPrecheckPath), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: registerUsername.value,
      }),
    })

    const precheckBody = await precheckResp.json()
    if (precheckResp.status !== 201) {
      showToastMessage("Something went wrong!")
      return
    }

    // @ts-ignore – provided by scram-bundled.js
    const {data} = scram.keysHMAC(
        registerPassword.value,
        precheckBody.data.salt,
        precheckBody.data.iteration_count
    )

    const registerResp = await fetch(getAPI(UserRegisterPath), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: registerUsername.value,
        public_key: Array.from(keyPair.publicKey),
        stored_key: data.storedKey,
        server_key: data.serverKey,
      }),
    })

    const registerBody = await registerResp.json()

    if (registerResp.status === 201) {
      modalWindowActive.value = true
    } else if (registerBody.message) {
      showToastMessage(registerBody.message)
    } else {
      showToastMessage("Something went wrong!")
    }
  } catch (err) {
    console.error(err)
    showToastMessage("Registration failed!")
  }
}
</script>
