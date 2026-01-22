<template>
  <div id="app">
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
    <div class="bg-[#F6F8FF] md:grid md:grid-cols-2 min-h-screen">
      <!-- Left -->
      <div class="self-center py-4 mx-10 lg:mx-36">
        <!-- STEP 1: Input Email -->
        <template v-if="step === 'email'">
          <div class="mb-6">
            <label class="text-lg text-[#414141] mb-1 block">
              Input your email:
            </label>
            <input
              v-model="emailAddress"
              class="w-full bg-white rounded-md border border-[#EADFD8]
                     py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
              placeholder="Email"
            />
          </div>

          <button
            class="w-[70%] bg-[#4F80E1] rounded-lg text-white py-4 mb-12"
            @click="getVerificationCode"
          >
            Get code
          </button>
        </template>

        <!-- STEP 2: Input Verification Code -->
        <template v-else>
          <div class="mb-6">
            <label class="text-lg text-[#414141] mb-1 block">
              Input your verification code:
            </label>
            <input
              v-model="verificationCode"
              class="w-full bg-white rounded-md border border-[#EADFD8]
                     py-2.5 px-3 placeholder:text-[#414141] focus:outline-none"
              placeholder="Verification code"
            />
          </div>

          <button
            class="w-[70%] bg-[#4F80E1] rounded-lg text-white py-4 mb-12"
            @click="checkEmailVerificationCode"
          >
            Submit
          </button>
        </template>
      </div>

      <!-- Right -->
      <div class="bg-white hidden md:flex items-center">
        <img src="@/assets/images/client-image.png" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

type Step = "email" | "code";

const step = ref<Step>("email");
const emailAddress = ref("");
const verificationCode = ref("");
const token = ref<string | null>(localStorage.getItem("token"));

/**
 * STEP 1 — Request verification code
 */
const getVerificationCode = async () => {
  if (!emailAddress.value) {
    alert("Email address is mandatory");
    return;
  }

  localStorage.setItem("email", emailAddress.value);

  try {
    const response = await fetch(
      "[[ .ProxyURL ]]/api/v1/verify-email",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token.value}`,
        },
        body: JSON.stringify({
          email: emailAddress.value,
        }),
      }
    );

    const result = await response.json();

    if (response.status === 200) {
      step.value = "code"; // ✅ toggle view
    } else {
      alert("Error happened: " + result.errors);
    }
  } catch (error) {
    console.error(error);
  }
};

/**
 * STEP 2 — Verify code
 */
const checkEmailVerificationCode = async () => {
  if (!verificationCode.value) {
    alert("Verification code is mandatory");
    return;
  }

  const email = localStorage.getItem("email");

  const response = await fetch(
    "[[ .ProxyURL ]]/api/v1/check-email-verification-code",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({
        email,
        code: verificationCode.value,
      }),
    }
  );

  const result = await response.json();

  localStorage.removeItem("email");

  alert(result.message);
  window.location.href = "[[ .ProxyURL ]]/user";
};
</script>
