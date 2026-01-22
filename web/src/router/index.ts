import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Welcome',
      component: () => import('@/views/Welcome.vue'),
    },
    {
      path: '/user-login',
      name: 'UserLogin',
      component: () => import('@/views/user/Login.vue'),
    },
    {
      path: '/user-register',
      name: 'UserRegister',
      component: () => import('@/views/user/Register.vue'),
    },
    {
      path: '/user-reset-password',
      name: 'UserResetPassword',
      component: () => import('@/views/user/ResetPassword.vue'),
    },
    {
      path: '/user/profile',
      name: 'UserProfile',
      component: () => import('@/views/user/Profile.vue'),
    },
    {
      path: '/user/verify-phone-number',
      name: 'UserVerifyPhoneNumber',
      component: () => import('@/views/user/VerifyPhone.vue'),
    },
    {
      path: '/user/verify-email',
      name: 'UserVerifyEmail',
      component: () => import('@/views/user/VerifyEmail.vue'),
    },
    {
      path: '/client-login',
      name: 'ClientLogin',
      component: () => import('@/views/client/Login.vue'),
    },
    {
      path: '/client-register',
      name: 'ClientRegister',
      component: () => import('@/views/client/Register.vue'),
    },
    {
      path: '/client/profile',
      name: 'ClientProfile',
      component: () => import('@/views/client/Profile.vue'),
    },
    {
      path: '/oauth-login',
      name: 'OAuthLogin',
      component: () => import('@/views/oauth/Login.vue'),
    },
    {
      // http://localhost:5001/authorize?client_id=e7fd0e02-8eff-4c91-8f9e-56d2680c371a&redirect_uri=http%3A%2F%2Flocalhost%3A5173%2Foauth2%2Fcallback&response_type=code&state=&scope=read%3Auser
      path: '/oauth/authorize',
      name: 'OAuthAuthorize',
      component: () => import('@/views/oauth/Authorize.vue'),
    },
    {
      path: '/oauth/error',
      name: 'OAuthError',
      component: () => import('@/views/oauth/Error.vue'),
    },
  ],
})

export default router
