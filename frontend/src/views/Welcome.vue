<template>
  <div class="welcome-container">
    <div
      class="flex flex-col lg:flex-row space-x-0 lg:space-x-7 justify-center h-dvh"
    >
      <div class="self-center mb-4">
        <h1
          class="text-center lg:text-start md:text-center md:mb-12 mb-4 md:text-[42px] text-3xl font-bold text-[#333333]"
        >
          Welcome page!
        </h1>

        <p
          class="text-center lg:text-start md:text-center font-normal text-sm text-[#333333] md:mb-11 mb-4"
        >
          Explore the incredible features and benefits our platform has to offer.
          From user-friendly interfaces to advanced tools, we're dedicated to
          making your experience extraordinary. To get started, please choose
          the appropriate login option below.
        </p>

        <div
          class="flex space-x-7 justify-center lg:justify-start md:justify-center"
        >
          <button
            @click="loginAsAUser"
            class="animate-bounce border border-black bg-white md:px-10 px-2 md:py-4 py-2 text-[#211F54] font-normal text-sm hover:bg-black hover:text-white"
          >
            Login as a User
          </button>

          <button
            @click="loginAsAClient"
            class="animate-bounce border border-black bg-white md:px-10 px-2 md:py-4 py-2 text-[#211F54] font-normal text-sm hover:bg-black hover:text-white"
          >
            Login as a Client
          </button>
        </div>
      </div>

      <img
        class="w-[264px] lg:w-[564px] md:w-[364px] h-[254px] md:h-[364px] lg:h-[564px] self-center"
        src="@/assets/images/team-working.png"
        alt="Team working"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const loginAsAUser = () => {
  router.push('/user-login')
}

const loginAsAClient = () => {
  router.push('/client-login')
}
</script>
