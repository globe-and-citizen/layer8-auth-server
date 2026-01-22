<template>
  <div>
    <!-- Navbar -->
    <div class="user-container">
      <div class="bg-white flex justify-center items-center my-4">
        <img
          src="@/assets/images/L8Logo.png"
          alt="Layer8"
          width="250"
          height="535"
        />
      </div>
    </div>

    <!-- Body -->
    <div class="bg-[#F6F8FF] md:grid md:grid-cols-2">
      <div class="self-center py-4 mx-10 lg:mx-36">
        <div class="mb-6">
          <label class="text-lg text-[#414141] mb-1 block">
            Input your username:
          </label>
          <input
            v-model="username"
            class="w-full bg-white rounded-md border border-[#EADFD8]
                   py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
            placeholder="username"
          />
        </div>

        <div class="mb-6">
          <label class="text-lg text-[#414141] mb-1 block">
            Input your mnemonic:
          </label>
          <input
            v-model="mnemonicSentence"
            class="w-full bg-white rounded-md border border-[#EADFD8]
                   py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
            placeholder="12-word mnemonic"
          />
        </div>

        <div class="mb-6">
          <label class="text-lg text-[#414141] mb-1 block">
            Input new password:
          </label>
          <input
            v-model="newPassword"
            type="password"
            class="w-full bg-white rounded-md border border-[#EADFD8]
                   py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
            placeholder="new password"
          />
        </div>

        <div class="mb-12">
          <label class="text-lg text-[#414141] mb-1 block">
            Repeat the new password:
          </label>
          <input
            v-model="repeatedNewPassword"
            type="password"
            class="w-full bg-white rounded-md border border-[#EADFD8]
                   py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
            placeholder="repeat new password"
          />
        </div>

        <button
          class="w-[70%] bg-[#4F80E1] rounded-lg text-white py-4 mb-12"
          @click="resetPassword"
        >
          Reset
        </button>
      </div>

      <!-- Right image -->
      <div class="bg-white hidden md:flex items-center">
        <img src="@/assets/images/client-image.png" />
      </div>
    </div>

    <!-- Footer -->
    <div class="user-container">
      <div class="bg-white flex justify-between items-center my-8">
        <div>
          <img
            src="@/assets/images/L8Logo.png"
            alt="Layer8"
            class="mb-6 md:mb-12 h-[35px] md:h-[70px]"
          />
          <p class="font-bold text-sm md:text-base text-black">
            ©Layer8security 2023.
          </p>
        </div>

        <div>
          <div class="text-xl font-bold text-black mb-4">Contact</div>
          <ul class="font-medium text-sm md:text-base text-black">
            <li>Email: hi@layer8.com</li>
            <li>Client Support: support@layer8.com</li>
            <li>Phone number: 0371 525 777</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { mnemonic } from "@/utils/mnemonic.ts"
import { scram } from "@/utils/scram.ts"
import {getAPI, UserResetPasswordPath, UserResetPasswordPrecheckPath} from "@/api/paths.ts";

/**
 * State
 */
const username = ref("");
const mnemonicSentence = ref("");
const newPassword = ref("");
const repeatedNewPassword = ref("");

const messageToSign = "Sign-in with Layer8";

/**
 * Actions
 */
const resetPassword = async () => {
  if (
    !username.value ||
    !mnemonicSentence.value ||
    !newPassword.value ||
    !repeatedNewPassword.value
  ) {
    alert("All fields are mandatory");
    return;
  }

  if (newPassword.value !== repeatedNewPassword.value) {
    alert("Repeated password does not match");
    return;
  }

  const currMnemonic = mnemonicSentence.value.trim();

  // `mnemonic` is expected to be a global lib (same as original HTML)
  if (!mnemonic.isValid(currMnemonic)) {
    alert("The provided mnemonic is invalid");
    return;
  }

  const keyPair = mnemonic.getPrivateAndPublicKeys(currMnemonic);
  const signature = mnemonic.sign(keyPair.privateKey, messageToSign);

  try {
    const precheckResp = await fetch(
      getAPI(UserResetPasswordPrecheckPath),
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: username.value }),
      }
    );

    const precheckBody = await precheckResp.json();

    if (precheckResp.status !== 200) {
      alert("Error: " + precheckBody.message);
      return;
    }

    // `scram` is also expected as global lib
    const { data } = scram.keysHMAC(
      newPassword.value,
      precheckBody.data.salt,
      precheckBody.data.iteration_count
    );

    const resetResp = await fetch(
      getAPI(UserResetPasswordPath),
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: username.value,
          signature: Array.from(signature),
          stored_key: data.storedKey,
          server_key: data.serverKey,
        }),
      }
    );

    const resetBody = await resetResp.json();

    if (resetBody.is_success === true) {
      alert(resetBody.message);
      window.location.href = "/user-login";
    } else {
      console.error(resetBody.errors);
      alert("Error: " + resetBody.message);
    }
  } catch (error) {
    console.error(error);
    alert("Error happened");
  }
};
</script>
