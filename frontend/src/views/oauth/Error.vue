<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo" />
    <h1 class="heading">Layer8</h1>
    <div class="line"></div>

    <div class="body">
      <h2 class="center">Oops! We encountered some errors</h2>
      <br />

      <div class="box">
        <div
          v-for="(err, index) in errors"
          :key="index"
          class="box-item"
        >
          <span>{{ err }}</span>
        </div>
      </div>

      <br />

      <div class="footer">
        <a class="cursor-pointer" @click="logout">Logout</a>
        | Layer8 &copy; {{ year }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import { useRoute, useRouter } from "vue-router"

const route = useRoute()
const router = useRouter()

/**
 * Maps ?opt=invalid_client,access_denied
 * to human-readable error messages
 * (same mapping as your Go code)
 */
const errorMap: Record<string, string> = {
  invalid_client: "The client is invalid.",
  access_denied: "The user denied the request.",
  server_error: "An error occurred on the server.",
  redirect_uri_mismatch:
    "The redirect uri does not match the client's redirect uri.",
}

const errors = computed(() => {
  const opt = (route.query.opt as string) || ""
  if (!opt) return []

  return opt
    .split(",")
    .map(code => errorMap[code])
    .filter(Boolean)
})

const year = new Date().getFullYear()

function logout() {
  // mirror original behavior
  document.cookie = "token=; Max-Age=0; path=/"
  router.push("/login")
}
</script>

