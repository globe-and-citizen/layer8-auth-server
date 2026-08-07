<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo" />

    <h1 class="heading">Layer8</h1>

    <div class="line"></div>

    <div class="body">
      <h2 class="center">Login</h2>

      <form @submit.prevent="submitLogin">
        <!-- keep next for backend compatibility -->
        <input type="hidden" name="next" :value="next" />

        <input
          aria-required="true"
          type="text"
          v-model="username"
          name="username"
          id="username"
          placeholder="Username"
          required
        />

        <input
          aria-required="true"
          type="password"
          v-model="password"
          name="password"
          id="password"
          placeholder="Password"
          required
        />

        <input
          aria-required="true"
          type="submit"
          value="Login"
        />

        <small v-if="error" class="error">{{ error }}</small>
      </form>

      <a href="/user-register" class="register-link">
        Don't have an account? Register
      </a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRoute } from "vue-router";
import {
  getAPI,
  OAuthUserLoginPath,
  OAuthUserPrecheckLoginPath,
} from "@/api/paths.ts";
import scram from "@/utils/scram.ts";
import router from "@/router";

const route = useRoute();

const username = ref("");
const password = ref("");
const error = ref("");
const cNonce = ref("");

// equivalent of {{ .Next }}
const next =
  (route.query.next as string) ||
  `/oauth/authorize${window.location.search}`;

const submitLogin = async () => {
  try {
    if (!username.value || !password.value) {
      return;
    }

    const cNonceBytes = new Uint8Array(32);
    crypto.getRandomValues(cNonceBytes);
    cNonce.value = btoa(String.fromCharCode(...cNonceBytes));

    const precheckRes = await fetch(getAPI(OAuthUserPrecheckLoginPath), {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: username.value,
        c_nonce: cNonce.value,
      }),
    });

    if (precheckRes.status !== 200) {
      return;
    }

    const precheckBody = await precheckRes.json();

    const { data } = scram.keysHMAC(
      password.value,
      precheckBody.data.salt,
      precheckBody.data.iteration_count
    );

    const clientKeyBytes = scram.hexStringToBytes(data.clientKey);

    const authMessage = `[n=${username.value},r=${cNonce.value},s=${precheckBody.data.salt},i=${precheckBody.data.iteration_count},r=${precheckBody.data.nonce}]`;

    const clientSignature = scram.signatureHMAC(
      authMessage,
      data.storedKey
    );

    const clientProof = scram.bytesToHexString(
      scram.xorBytes(
        clientKeyBytes,
        scram.hexStringToBytes(clientSignature)
      )
    );

    const loginRes = await fetch(
      getAPI(OAuthUserLoginPath + window.location.search),
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username.value,
          nonce: precheckBody.data.nonce,
          c_nonce: cNonce.value,
          client_proof: clientProof,
        }),
      }
    );

    if (loginRes.status >= 500) {
      await router.push("/oauth/error?opt=server_error");
      return;
    }

    const loginJSON = await loginRes.json();

    if (loginJSON.data?.verifier) {
      const serverCheck = scram.signatureHMAC(
        authMessage,
        data.serverKey
      );

      if (serverCheck === loginJSON.data.verifier) {
        if (loginJSON.data?.redirect) {
          window.location.href = loginJSON.data.redirect;
        } else {
          window.location.href = next;
        }
      }
    } else {
      error.value = loginJSON.data?.error || "Login failed";
    }
  } catch (err) {
    console.error(err);
    error.value = "Network error" + err;
  }
};
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
}

.center {
  text-align: center;
  margin-bottom: 20px;
}

.line {
  width: 100%;
  height: 2px;
  background-color: #0000001f;
  margin: 20px 0;
}

.body {
  width: 100%;
}

form {
  display: flex;
  flex-direction: column;
  align-items: center;
}

input {
  width: 100%;

  padding: 15px;
  margin: 10px 0;

  font-size: 16px; /* prevents iOS zoom */

  border: 1px solid #0000001f;
  outline: none;
}

input:focus {
  border-color: #000;
}

input[type="submit"] {
  background: #000;
  color: #fff;
  cursor: pointer;
}

.register-link {
  display: block;
  text-align: center;
  margin-top: 10px;
}

.error {
  color: #ff0000b9;
  text-align: center;
}

.bold {
  font-weight: bold;
}

.cursor-pointer {
  cursor: pointer;
}

.footer {
  font-size: 13px;
  text-align: right;
  color: #0000006e;
}

.btn-primary {
  width: 100%;
  padding: 15px;
  margin: 10px 0;
  border: 1px solid #0000001f;
  outline: none;
  background: #000;
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

/* Mobile */
@media (max-width: 600px) {
  body {
    padding: 12px;
  }

  .container {
    width: 100%;
    min-width: 0;
    max-width: 100%;

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
    font-size: 1.75rem;
  }

  .center {
    font-size: 1.25rem;
    margin-bottom: 20px;
  }

  input {
    padding: 12px 14px;
  }
}
</style>
