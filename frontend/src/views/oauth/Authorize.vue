<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo"/>
    <h1 class="heading">Layer8</h1>
    <div class="line"></div>

    <div class="body">
      <h2 class="center">
        Authorize <b>{{ clientName }}</b>
      </h2>

      <br/>

      <div class="box">
        <div class="box-item" v-for="s in scopes" :key="s.name">
          <span><input type="checkbox" checked disabled/></span>
          <span>{{ s.description }}</span>
        </div>
      </div>

      <br/>

      <form method="POST" id="submit" @submit="submit">
        <input type="hidden" name="decision" value="allow"/>

        <label
          style="display: flex; align-items: left; white-space: nowrap; align-self: flex-start;">
          <input
            type="checkbox"
            v-model="shareDisplayName"
            style="margin-right: 5px;"
          />
          <span style="font-size: 14px;">Share display name</span>
        </label>

        <label
          style="display: flex; align-items: left; white-space: nowrap; align-self: flex-start;">
          <input
            type="checkbox"
            v-model="shareIsEmailVerified"
            style="margin-right: 5px;"
          />
          <span style="font-size: 14px;">Share email verification data</span>
        </label>

        <label
          style="display: flex; align-items: left; white-space: nowrap; align-self: flex-start;">
          <input
            type="checkbox"
            v-model="shareColor"
            style="margin-right: 5px;"
          />
          <span style="font-size: 14px;">Share color</span>
        </label>

        <label
          style="display: flex; align-items: left; white-space: nowrap; align-self: flex-start;">
          <input
            type="checkbox"
            v-model="shareBio"
            style="margin-right: 5px;"
          />
          <span style="font-size: 14px;">Share bio</span>
        </label>

        <input type="submit" value="Authorize"/>
      </form>

      <br/>

      <div class="footer">
        <a class="cursor-pointer" @click="logout">Logout</a>
        | Layer8 © {{ getDate() }}
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue'
import {getAPI, OAuthGetAuthorizeContextPath, OAuthPostAuthorizeDecisionPath} from "@/api/paths.js";
import {useRoute, useRouter} from "vue-router";

const route = useRoute()
const router = useRouter()

const params = window.location.search
const queries = route.query
console.log(queries)
console.log(params)

const clientId = queries.client_id
const scopeParam = queries.scope || ''

const clientName = ref('')
const scopes = ref([])

const shareDisplayName = ref(false)
const shareIsEmailVerified = ref(false)
const shareColor = ref(false)
const shareBio = ref(false)

const getDate = () => new Date().getFullYear()

const logout = () => {
  document.cookie = 'token=; Max-Age=0; path=/'
  router.push('/oauth-login' + params)
}

onMounted(async () => {
  const res = await fetch(getAPI(OAuthGetAuthorizeContextPath) + `${params}`,
    {
      method: 'GET',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include'
    }
  )
  if (!res.ok) {
    await router.push('/oauth-login' + params)
    return
  }

  const data = await res.json()
  console.log(data)
  clientName.value = data.client_name
  scopes.value = data.scopes

  // clientName.value = "Layer8"
  // scopes.value = [
  //   {
  //     name: 'read:user',
  //     description: 'read anonymized information about your account'
  //   }
  // ]
})

const submit = async (e) => {
  e.preventDefault()

  const res = await fetch(getAPI(OAuthPostAuthorizeDecisionPath) + `${params}`,
    {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      redirect: "follow",
      body: JSON.stringify({
        client_id: clientId,
        scopes: scopeParam,
        share: {
          display_name: shareDisplayName.value,
          is_email_verified: shareIsEmailVerified.value,
          color: shareColor.value,
          bio: shareBio.value,
        },
        return_result: !!window.opener,
      }),
    })

  const data = await res.json()

  if (window.opener) {
    window.opener.postMessage(data, '*')
    window.close()
  } else {
    window.location.href = data.redirect_uri
  }
}
</script>
