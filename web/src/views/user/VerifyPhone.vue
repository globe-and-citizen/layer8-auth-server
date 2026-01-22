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

    <!-- Shared background -->
    <div class="bg-[#F6F8FF] md:grid md:grid-cols-2 min-h-screen">
      <!-- Left -->
      <component
        :is="currentComponent"
        class="self-center py-4 mx-10 lg:mx-36"
        @verified-via-telegram="goToCodeInput"
      />

      <!-- Right -->
      <div class="bg-white hidden md:flex items-center">
        <img src="@/assets/images/client-image.png" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import VerifyViaTelegram from "@/views/user/phone/VerifyViaTelegram.vue";
import InputVerificationCode from "@/views/user/phone/InputVerificationCode.vue";

const step = ref("telegram"); // 'telegram' | 'code'

const currentComponent = computed(() =>
  step.value === "telegram"
    ? VerifyViaTelegram
    : InputVerificationCode
);

const goToCodeInput = () => {
  step.value = "code";
};
</script>
