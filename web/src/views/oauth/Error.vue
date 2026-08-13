<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo" />

    <h1 class="heading">Layer8</h1>

    <div class="line"></div>

    <div class="body">
      <h2 class="center">Oops! We encountered some errors</h2>

      <div class="box">
        <div
          v-for="(err, index) in errors"
          :key="index"
          class="box-item"
        >
          <span>{{ err }}</span>
        </div>
      </div>

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
  bad_request: "Some error occurred during the request.",
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
  router.push("/oauth-login")
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

.container {
  width: 600px;
  max-width: 600px;
  min-width: 600px;

  padding: 30px;

  border: 5px solid #0000001f;
  border-radius: 20px;

  background: #fff;
}

.logo {
  display: block;

  width: 150px;

  margin: -110px auto 0;

  background: #fff;
  padding: 20px;
}

.heading {
  text-align: center;
  font-weight: bold;
  color: #484848;
  font-size: 2rem;
}

.center {
  text-align: center;
  font-size: 1.5rem;
  line-height: 1.4;
  margin-bottom: 24px;
}

.cursor-pointer {
  cursor: pointer;
}

.line {
  width: 100%;
  height: 2px;
  background-color: #0000001f;
  margin: 20px 0 30px;
}

.box {
  width: 100%;
  border: 1px solid #0000001f;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 24px;
}

.box-item {
  padding: 16px;
  border-bottom: 1px solid #0000001f;

  word-break: break-word;
  overflow-wrap: anywhere;
}

.box-item:last-child {
  border-bottom: none;
}

.footer {
  font-size: 13px;
  text-align: right;
  color: #0000006e;
}

.footer a {
  color: inherit;
}

/* Tablet */
@media (max-width: 768px) {
  .container {
    padding: 24px;
  }

  .heading {
    font-size: 1.75rem;
  }

  .center {
    font-size: 1.25rem;
  }
}

/* Mobile */
@media (max-width: 600px) {
  body {
    padding: 12px;
  }

  .container {
    padding: 20px;
    border-width: 2px;
    margin-top: 60px;
  }

  .logo {
    width: 120px;
    margin-top: -80px;
    padding: 12px;
  }

  .heading {
    font-size: 1.5rem;
  }

  .center {
    font-size: 1.125rem;
  }

  .box-item {
    padding: 14px;
    font-size: 14px;
  }

  .footer {
    text-align: center;
    line-height: 1.8;
  }
}
</style>
