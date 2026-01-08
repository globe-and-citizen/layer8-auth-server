<template>
  <div class="bg-[#E5E5E5] min-h-screen relative">
    <div id="app">
      <!-- Header -->
      <div
        class="py-3 px-4 md:px-7 mb-2 md:mb-4 bg-white rounded-b-2xl flex justify-between items-center">
        <div class="hidden md:block lg:block">
          <div class="flex space-x-4 items-center">
            <img
              class="w-12 h-12"
              src="@/assets/images/user_avatar.png"
            />
            <div>{{ userName }}</div>
          </div>
        </div>
        <div class="cursor-pointer block md:hidden lg:hidden" @click="showSidebar(true)">
          <svg width="26" height="26" viewBox="0 0 26 26" fill="none"
               xmlns="http://www.w3.org/2000/svg">
            <path
              d="M2 6C2 5.44772 2.44772 5 3 5H21C21.5523 5 22 5.44772 22 6C22 6.55228 21.5523 7 21 7H3C2.44772 7 2 6.55228 2 6Z"
              fill="currentColor"/>
            <path
              d="M2 12.0322C2 11.4799 2.44772 11.0322 3 11.0322H21C21.5523 11.0322 22 11.4799 22 12.0322C22 12.5845 21.5523 13.0322 21 13.0322H3C2.44772 13.0322 2 12.5845 2 12.0322Z"
              fill="currentColor"/>
            <path
              d="M3 17.0645C2.44772 17.0645 2 17.5122 2 18.0645C2 18.6167 2.44772 19.0645 3 19.0645H21C21.5523 19.0645 22 18.6167 22 18.0645C22 17.5122 21.5523 17.0645 21 17.0645H3Z"
              fill="currentColor"/>
          </svg>
        </div>
      </div>

      <!-- Mobile Sidebar Overlay -->
      <div
        @click="showSidebar(false)"
        v-show="sidebarShow"
        class="absolute block md:hidden lg:hidden top-0 left-0 h-dvh w-full bg-opacity-50 backdrop-blur-lg text-black"
      >
        <div class="w-[70%] h-dvh md:col-span-1 bg-white rounded-2xl px-5 py-2">
          <div>
            <span class="flex justify-end mb-4" @click="showSidebar(false)">&#x2715;</span>
            <div
              class="flex space-x-2.5 items-center pl-7 py-2.5 cursor-pointer mb-5 rounded-md bg-[#e4f6ff]">
              <img src="@/assets/images/icons/project-info.svg" alt="">
              <span class="font-medium text-[#2F80ED] text-base">Project info</span>
            </div>
            <div @click="logoutUser"
                 class="flex space-x-2.5 items-center pl-7 py-3 cursor-pointer mb-5">
              <img src="@/assets/images/icons/logout-icon.svg" alt="">
              <span class="font-medium text-black text-base">Log out</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Layout -->
      <div @click="showSidebar(false)" class="grid md:grid-cols-5 lg:grid-cols-5 gap-x-4 m-4">
        <div
          class="md:col-span-1 bg-white rounded-2xl px-2 md:px-5 py-3 md:py-6 hidden md:block lg:block">
          <div>
            <div
              class="flex space-x-2.5 items-center pl-2 md:pl-3 lg:pl-7 py-2.5 cursor-pointer mb-2 md:mb-5 rounded-md bg-[#e4f6ff]">
              <img src="@/assets/images/icons/project-info.svg" alt="">
              <span class="font-medium text-[#2F80ED] text-sm md:text-base">Project info</span>
            </div>
            <div @click="logoutUser"
                 class="flex space-x-2.5 items-center pl-2 md:pl-3 lg:pl-7 py-3 cursor-pointer mb-5">
              <img src="@/assets/images/icons/logout-icon.svg" alt="">
              <span class="font-medium text-black text-base">Log out</span>
            </div>
          </div>
        </div>

        <!-- Main Content -->
        <main class="md:col-span-4">
          <div class="font-bold text-2xl md:text-4xl text-[#2F80ED] text-center my-4 md:my-9">
            Welcome “{{ userName }}!”
            Client Portal
          </div>

          <div class="grid md:grid-cols-1 md:gap-x-4">
            <!-- User Data -->
            <UserDataSection
              @user-data=getUser
            />

            <!-- Statistics -->
            <UsageStatisticsSection/>
          </div>
        </main>
      </div>
    </div>
  </div>


</template>

<script setup lang="ts">
import {ref} from "vue"
import UsageStatisticsSection from "@/views/client/profile/UsageStatisticsSection.vue";
import UserDataSection from "@/views/client/profile/UserDataSection.vue";

const userName = ref("")
const sidebarShow = ref(false)

const showSidebar = (v: boolean) => {
  sidebarShow.value = v
  document.body.style.overflow = v ? "hidden" : "auto"
}

const logoutUser = () => {
  localStorage.removeItem("clientToken")
  window.location.href = "/"
}

const getUser = (text: string) => {
  userName.value = text
}
</script>

<style scoped>
.input {
  @apply border border-[#EADFD8] p-3 rounded-lg w-full font-medium;
}

/* NO 'scoped' attribute here so it hits the whole page */
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
  /* Applying the background color here ensures there's NO white gap at the bottom */
  background-color: #E5E5E5;
}

#app {
  min-height: 100%;
}
</style>
