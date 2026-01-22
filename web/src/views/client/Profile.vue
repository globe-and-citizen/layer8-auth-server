<template>
  <div class="bg-[#E5E5E5] min-h-screen relative">
    <div id="app">
      <!-- Header -->
      <TopBar
        :userName="user.name"
        :userID="user.id"
        @showSidebar="showSidebar"
      />

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
import TopBar from "@/views/client/profile/TopBar.vue";

const userName = ref("")
const sidebarShow = ref(false)
const user = ref({
  id: "",
  name: "",
  secret: "",
  redirect_uri: "",
  backend_uri: "",
  ntor_certificate: "",
})


const showSidebar = (v: boolean) => {
  sidebarShow.value = v
  document.body.style.overflow = v ? "hidden" : "auto"
}

const logoutUser = () => {
  localStorage.removeItem("clientToken")
  window.location.href = "/"
}

const getUser = (text: any) => {
  user.value = text
  userName.value = user.value.name
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
