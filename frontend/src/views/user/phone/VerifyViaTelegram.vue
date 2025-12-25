<template>
  <div>
    <div class="mb-6">
      <label class="text-lg text-[#414141] mb-1 block">
        Click the button below to verify your phone number via Telegram:
      </label>
    </div>

    <button
      class="w-[70%] bg-[#4F80E1] rounded-lg text-white py-4 mb-12"
      :disabled="loading"
      @click="verifyPhoneNumber"
    >
      {{ loading ? "Verifying..." : "Verify in Telegram" }}
    </button>
  </div>
</template>

<script setup>
import { ref } from "vue";

const emit = defineEmits(["verified-via-telegram"]);

const token = ref(localStorage.getItem("token"));
const loading = ref(false);

const verifyPhoneNumber = async () => {
  loading.value = true;

  try {
    const sessionRes = await fetch(
      "[[ .ProxyURL ]]/api/v1/generate-telegram-session-id",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token.value}`,
        },
      }
    );

    const sessionBody = await sessionRes.json();
    if (sessionRes.status !== 200) {
      alert("Invalid server response");
      return;
    }

    const sessionId = sessionBody.data.session_id;

    window.open(
      `https://t.me/Layer8PhoneNumberVerifierBot?start=${sessionId}`,
      "_blank",
      "noopener"
    );

    const response = await fetch(
      "[[ .ProxyURL ]]/api/v1/verify-phone-number-via-bot",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token.value}`,
        },
      }
    );

    if (response.status === 200) {
      emit("verified-via-telegram");
    } else {
      const result = await response.json();
      alert("Error: " + result.errors);
    }
  } catch (e) {
    console.error(e);
    alert("Verification failed");
  } finally {
    loading.value = false;
  }
};
</script>
