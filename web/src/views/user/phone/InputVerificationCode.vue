<template>
  <div>
    <div class="mb-6">
      <label class="text-lg text-[#414141] mb-1 block">
        Input your verification code:
      </label>

      <input
        v-model="verificationCode"
        class="w-full bg-white rounded-md border border-[#EADFD8] py-2.5 px-3
               placeholder:text-[#414141] focus:outline-none"
        placeholder="Verification code"
      />
    </div>

    <button
      class="w-[70%] bg-[#4F80E1] rounded-lg text-white py-4 mb-12"
      @click="submitCode"
    >
      Submit
    </button>
  </div>
</template>

<script setup>
import { ref } from "vue";

const verificationCode = ref("");
const token = ref(localStorage.getItem("token"));

const submitCode = async () => {
  if (!verificationCode.value) {
    alert("Verification code is mandatory");
    return;
  }

  const response = await fetch(
    "[[ .ProxyURL ]]/api/v1/check-phone-number-verification-code",
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({
        verification_code: verificationCode.value,
      }),
    }
  );

  const result = await response.json();
  alert(result.message);

  window.location.href = "[[ .ProxyURL ]]/user";
};
</script>
