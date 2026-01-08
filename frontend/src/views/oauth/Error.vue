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

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  transition: all 0.3s ease-in-out;
}

body {
  font-family: monospace;
  font-size: 16px;
  line-height: 1.5;
  color: #333;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

a {
  text-decoration: none;
  color: #0000006e;
}

a:hover {
  color: #333;
}

.heading {
  text-align: center;
  font: bold;
  color: #484848;
}

.center {
  text-align: center;
}

.cursor-pointer {
  cursor: pointer;
}

.footer {
  font-size: 13px;
  text-align: right;
  color: #0000006e;
}

.line {
  width: 100%;
  height: 2px;
  background-color: #0000001f;
  margin: 20px 0;
}

.container {
  max-width: 600px;
  min-width: 600px;
  padding: 30px;
  border: 5px solid #0000001f;
  border-radius: 20px;
}

@media (max-width: 600px) {
  .container {
    width: 100%;
  }
}

.logo {
  width: 150px;
  margin-left: 50%;
  transform: translateX(-50%);
  margin-top: -110px;
  background-color: #fff;
  padding: 20px;
}

form {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

input {
  width: 100%;
  padding: 15px;
  margin: 10px 0;
  border: 1px solid #0000001f;
  outline: none;
}

input:focus {
  border: 1px solid #000000;
}

input[type="submit"] {
  background-color: #000000;
  color: #fff;
  cursor: pointer;
}

.box {
  width: 100%;
  border: 1px solid #0000001f;
  margin: 10px 0;
}

.box-item {
  padding: 20px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid #0000001f;
}

.box-item:last-child {
  border-bottom: none;
}

.box-item span:first-child {
  margin-right: 10px;
}
</style>

