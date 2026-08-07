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

          <ConsentCheckbox
            v-model="shareLocation"
            label="Share location"
          />
        </div>
        <input type="submit" value="Authorize"/>
        <small v-if="error" class="error">{{ error }}</small>
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
const error = ref("");

const params = window.location.search
const queries = route.query

const clientId = queries.client_id
const scopeParam = queries.scope || ''
const isOIDC = scopeParam.split(' ').includes('openid')

const clientName = ref('')
const scopes = ref([])

const shareDisplayName = ref(false)
const shareIsEmailVerified = ref(false)
const shareColor = ref(false)
const shareBio = ref(false)
const shareLocation = ref(false)
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

    if (res.status >= 500) {
      await router.push("/oauth/error?opt=server_error");
      return
    }

    const data = await res.json()
    // console.log(data)
    clientName.value = data.client_name
    scopes.value = data.scopes
  } catch (e) {
    console.log(e)
    error.value = e.message
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
            location: shareLocation.value
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
  justify-content: center;
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
  border-color: #000000;
}

input[type="submit"] {
  background-color: #000000;
  color: #fff;
  cursor: pointer;
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
