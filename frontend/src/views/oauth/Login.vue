<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo"/>
    <h1 class="heading">Layer8</h1>
    <div class="line"></div>

    <div class="body">
      <h2 class="center">Login</h2>
      <br/>

      <form @submit.prevent="submitLogin">
        <!-- keep next for backend compatibility -->
        <input type="hidden" name="next" :value="next"/>

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

        <input aria-required="true" type="submit" value="Login"/>

        <small v-if="error" class="error">{{ error }}</small>
      </form>

      <a href="/user-register">Don't have an account? Register</a>
      <br/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue"
import {useRoute, useRouter} from "vue-router"
import {getAPI, OAuthUserLoginPath, OAuthUserPrecheckLoginPath} from "@/api/paths.ts";
import scram from "@/utils/scram.ts";

const route = useRoute()
const router = useRouter()

const username = ref("")
const password = ref("")
const error = ref("")
const cNonce = ref("")
const token = ref("")

// equivalent of {{ .Next }}
const next = (route.query.next as string) || `/oauth/authorize${window.location.search}`

const submitLogin = async () => {
  try {
    if (!username.value || !password.value) {
      return;
    }

    const cNonceBytes = new Uint8Array(32);
    crypto.getRandomValues(cNonceBytes);
    cNonce.value = btoa(String.fromCharCode(...cNonceBytes));

    const precheckRes = await fetch(
      getAPI(OAuthUserPrecheckLoginPath),
      {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          username: username.value,
          c_nonce: cNonce.value,
        }),
      }
    );

    if (precheckRes.status !== 200) {
      return;
    }

    const precheckBody = await precheckRes.json();

    const {data} = scram.keysHMAC(
      password.value,
      precheckBody.data.salt,
      precheckBody.data.iteration_count
    );

    const clientKeyBytes = scram.hexStringToBytes(data.clientKey);

    const authMessage = `[n=${username.value},r=${cNonce.value},s=${precheckBody.data.salt},i=${precheckBody.data.iteration_count},r=${precheckBody.data.nonce}]`;

    const clientSignature = scram.signatureHMAC(authMessage, data.storedKey);
    const clientProof = scram.bytesToHexString(
      scram.xorBytes(
        clientKeyBytes,
        scram.hexStringToBytes(clientSignature)
      )
    );

    const loginRes = await fetch(
      getAPI(OAuthUserLoginPath + `${window.location.search}`),
      {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          username: username.value,
          nonce: precheckBody.data.nonce,
          c_nonce: cNonce.value,
          client_proof: clientProof,
        }),
      }
    );

    const loginJSON = await loginRes.json();

    if (loginJSON.data?.verifier) {
      const serverCheck = scram.signatureHMAC(authMessage, data.serverKey);
      if (serverCheck === loginJSON.data.verifier) {
        if (loginJSON.data?.redirect) {
          window.location.href = loginJSON.data.redirect
        } else {
          window.location.href = next
        }
      }
    } else {
      error.value = loginJSON.data?.error || "Login failed"
    }
  } catch (err) {
    console.error(err);
    error.value = "Network error";
  }
};

</script>
