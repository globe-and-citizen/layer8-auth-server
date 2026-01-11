<template>
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
    <div class="cursor-pointer block md:hidden lg:hidden" @click="emitShowSidebar(true)">
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

    <PayBalance :userID="userID"></PayBalance>
  </div>
</template>

<script setup lang="ts">
import PayBalance from "@/views/client/profile/PayBalance.vue";

/* =========================
   Props & Emits
   ========================= */

const props = defineProps<{
  userName: string
  userID: string
}>()

const emit = defineEmits<{
  (e: 'showSidebar', value: boolean): void
}>()

const emitShowSidebar = (show: boolean) => {
  emit('showSidebar', show)
}
</script>


<style scoped>

</style>
