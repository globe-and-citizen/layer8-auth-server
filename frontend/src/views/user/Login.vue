<template>
  <div>
    <!-- Navbar -->
    <div id="navbar" class="user-container">
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
    <div id="body" class="bg-[#F6F8FF] md:grid md:grid-cols-2">
      <div class="self-center py-4 md:py-0 mx-10 lg:mx-36 md:mx-18">
        <h1 class="font-bold text-[#4F80E1] text-[40px] text-start mb-2">
          Login
        </h1>

        <p class="font-normal text-xl text-[#414141] text-start mb-12">
          Enter your email and password to login.
        </p>

        <div>
          <div class="mb-6">
            <label class="text-sm text-[#414141] mb-1 block">Username</label>
            <input
              v-model="loginUsername"
              class="w-full bg-white rounded-md border border-[#EADFD8] py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
              placeholder="Username"
            />
          </div>

          <div class="mb-12">
            <label class="text-sm text-[#414141] mb-1 block">Password</label>
            <input
              v-model="loginPassword"
              type="password"
              class="w-full bg-white rounded-md border border-[#EADFD8] py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
              placeholder="Password"
            />
          </div>

          <a
            href="/user-reset-password"
            class="text-sm text-[#414141] font-normal text-center block cursor-pointer"
          >
            Forgot your password?
          </a>

          <button
            class="w-full bg-[#4F80E1] rounded-lg text-center text-white py-4 mb-12"
            @click="loginUser"
          >
            Login
          </button>

          <a
            href="/user-register"
            class="text-sm text-[#414141] font-normal text-center block cursor-pointer"
          >
            Don't have an account?
            <span class="font-bold">Register</span>
          </a>
        </div>
      </div>

      <div class="bg-white hidden md:flex lg:flex items-center">
        <img src="@/assets/images/cyber-phone.png" />
      </div>
    </div>

    <!-- Footer -->
    <div id="footer" class="user-container">
      <div class="bg-white flex justify-between items-center my-8">
        <div>
          <img
            src="@/assets/images/L8Logo.png"
            alt="Layer8"
            class="mb-6 md:mb-12 h-[35px] md:w-full md:h-[70px]"
          />
          <p class="font-bold text-sm md:text-base text-black text-start">
            ©Layer8security 2023.
          </p>
        </div>

        <div>
          <div class="text-xl font-bold text-black text-end md:text-start mb-0 md:mb-4 lg:mb-6">
            Contact
          </div>
          <ul class="font-medium text-sm md:text-base text-black text-end md:text-start">
            <li>Email: hi@layer8.com</li>
            <li>Client Support: support@layer8.com</li>
            <li>Phone number: 0371 525 777</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <div
      :class="showToast ? 'opacity-100' : 'opacity-0 pointer-events-none'"
      class="fixed top-3 right-3 bg-red-500 text-white p-2 rounded-md transition-opacity ease-in-out duration-500 z-50"
    >
      {{ toastMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import {scram} from "@/utils/scram.ts"
import {getAPI, UserLoginPath, UserLoginPrecheckPath} from "@/api/paths.js";

const loginUsername = ref("");
const loginPassword = ref("");
const token = ref(localStorage.getItem("token"));
const showToast = ref(false);
const toastMessage = ref("");
const cNonce = ref("");

const isLoggedIn = computed(() => token.value !== null);

const showToastMessage = (message) => {
  toastMessage.value = message;
  showToast.value = true;
  setTimeout(() => {
    showToast.value = false;
  }, 3000);
};

const loginUser = async () => {
  try {
    if (!loginUsername.value || !loginPassword.value) {
      showToastMessage("Please enter a username and password!");
      return;
    }

    const cNonceBytes = new Uint8Array(32);
    crypto.getRandomValues(cNonceBytes);
    cNonce.value = btoa(String.fromCharCode(...cNonceBytes));

    const precheckRes = await fetch(
      getAPI(UserLoginPrecheckPath),
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: loginUsername.value,
          c_nonce: cNonce.value,
        }),
      }
    );

    if (precheckRes.status !== 200) {
      showToastMessage("Failed to login");
      return;
    }

    const precheckBody = await precheckRes.json();

    const { data } = scram.keysHMAC(
      loginPassword.value,
      precheckBody.data.salt,
      precheckBody.data.iteration_count
    );

    const clientKeyBytes = scram.hexStringToBytes(data.clientKey);

    const authMessage = `[n=${loginUsername.value},r=${cNonce.value},s=${precheckBody.data.salt},i=${precheckBody.data.iter_count},r=${precheckBody.data.nonce}]`;

    const clientSignature = scram.signatureHMAC(authMessage, data.storedKey);
    const clientProof = scram.bytesToHexString(
      scram.xorBytes(
        clientKeyBytes,
        scram.hexStringToBytes(clientSignature)
      )
    );

    const loginRes = await fetch(
      getAPI(UserLoginPath),
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          username: loginUsername.value,
          nonce: precheckBody.data.nonce,
          c_nonce: cNonce.value,
          client_proof: clientProof,
        }),
      }
    );

    const loginJSON = await loginRes.json();

    if (loginJSON.data?.server_signature) {
      const serverCheck = scram.signatureHMAC(authMessage, data.serverKey);
      if (serverCheck === loginJSON.data.server_signature) {
        token.value = loginJSON.data.token;
        localStorage.setItem("token", token.value);
        showToastMessage("Login successful!");
        window.location.href = "/user/profile";
      }
    } else {
      showToastMessage(loginJSON.message || "Login failed");
    }
  } catch (err) {
    console.error(err);
    showToastMessage("Login failed!");
  }
};
</script>
