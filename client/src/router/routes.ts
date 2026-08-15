import { RouteRecordRaw } from 'vue-router'
import { useGuards } from './guards'

const { guardRole, guardAuth } = useGuards()

export const accountRoutes: Record<string, RouteRecordRaw> = {
  login: {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { title: 'Login' }
  },
  signup: {
    path: '/signup',
    name: 'signup',
    component: () => import('@/views/auth/SignUpView.vue'),
    meta: { title: 'Account request' }
  },
  pwdReset: {
    path: '/password-reset',
    name: 'password-reset',
    component: () => import('@/views/auth/PasswordResetView.vue'),
    meta: { title: 'Password reset' }
  },
  verifyEmail: {
    path: '/verify-email',
    name: 'verify-email',
    component: () => import('@/views/auth/EmailVerificationView.vue'),
    meta: { title: 'E-mail verification' }
  },
  account: guardAuth({
    path: '/account',
    name: 'account',
    component: () => import('@/views/accounts/AccountView.vue'),
    meta: { title: 'Account infos' }
  })
}

export default {
  settings: guardRole('admin', {
    label: 'Settings',
    icon: 'mdi-tools',
    path: '/settings/:category',
    name: 'app-settings',
    params: { category: 'instance' },
    component: () => import('@/features/settings/views/AdminSettings.vue'),
    props: true,
    meta: {
      title: 'Settings',
      drawer: {
        temporary: true
      }
    }
  }),
  ...accountRoutes
}
