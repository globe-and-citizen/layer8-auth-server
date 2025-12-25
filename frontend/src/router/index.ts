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
  ],
})

export default router
