<template>
  <div class="container">
    <img src="@/assets/images/logo.png" alt="logo" class="logo"/>
    <h1 class="heading">Layer8</h1>
    <div class="line"></div>

    <div class="body">
      <h2 class="center">
        {{ isOIDC ? 'Sign in to' : 'Authorize' }}
        <b>{{ clientName }}</b>
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
        <div v-if="isOIDC" style="display: flex; flex-direction: column; align-items: flex-start;">
          <ConsentCheckbox label="Agree with all terms and conditions"
                           v-model="agreed"></ConsentCheckbox>
        </div>

        <div v-if="!isOIDC" style="display: flex; flex-direction: column; align-items: flex-start;">
          <ConsentCheckbox
            v-model="shareDisplayName"
            label="Share display name"
          />

          <ConsentCheckbox
            v-model="shareIsEmailVerified"
            label="Share email verification data"
          />

          <ConsentCheckbox
            v-model="shareColor"
            label="Share color"
          />

          <ConsentCheckbox
            v-model="shareBio"
            label="Share bio"
          />
        </div>
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
import ConsentCheckbox from "@/views/oauth/authorize/ConsentCheckbox.vue";

const route = useRoute()
const router = useRouter()

const params = window.location.search
const queries = route.query
console.log(queries)
console.log(params)

const clientId = queries.client_id
const scopeParam = queries.scope || ''
const isOIDC = scopeParam.split(' ').includes('openid')

const clientName = ref('')
const scopes = ref([])

const shareDisplayName = ref(false)
const shareIsEmailVerified = ref(false)
const shareColor = ref(false)
const shareBio = ref(false)
const agreed = ref(false)

const getDate = () => new Date().getFullYear()

const logout = () => {
  document.cookie = 'token=; Max-Age=0; path=/'
  router.push('/oauth-login' + params)
}

onMounted(async () => {
  try {
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
  } catch (e) {
    console.log(e)
  }
})

const submit = async (e) => {
  e.preventDefault()

  try {
    if (isOIDC && !agreed.value) {
      alert('You must agree with all terms and conditions to continue.')
      return
    }

    const res = await fetch(getAPI(OAuthPostAuthorizeDecisionPath) + `${params}`,
      {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
        redirect: "follow",
        body: JSON.stringify({
          client_id: clientId,
          scopes: scopeParam,
          oidc_agreed: agreed.value,
          share: {
            display_name: shareDisplayName.value,
            is_email_verified: shareIsEmailVerified.value,
            color: shareColor.value,
            bio: shareBio.value,
          },
          return_result: !!window.opener, // ??? what's this for?
        }),
      })

    const data = await res.json()

    if (window.opener) {
      window.opener.postMessage(data, '*')
      window.close()
    } else {
      window.location.href = data.redirect_uri
    }
  } catch (e) {
    console.log(e)
  }

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

.error {
  color: #ff0000b9;
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

.btn-primary {
  width: 100%;
  padding: 15px;
  margin: 10px 0;
  border: 1px solid #0000001f;
  outline: none;
  background-color: #000000;
  color: #fff;
  cursor: pointer;
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
