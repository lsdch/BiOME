import { RouteRecord, RouteRecordRaw } from "vue-router"
import { useGuards } from "./guards"

const { guardRole, guardAuth } = useGuards()

export const accountRoutes: Record<string, RouteRecordRaw> = {
  login: {
    path: '/login',
    name: 'login',
    component: () => import('../views/auth/LoginView.vue'),
    meta: { subtitle: "Login" }
  },
  signup: {
    path: '/signup',
    name: 'signup',
    component: () => import('../views/auth/SignUpView.vue'),
    meta: { subtitle: "Account request" }
  },
  pwdReset: {
    path: '/password-reset',
    name: 'password-reset',
    component: () => import('../views/auth/PasswordResetView.vue'),
    meta: { subtitle: "Password reset" }
  },
  verifyEmail: {
    path: '/verify-email',
    name: 'verify-email',
    component: () => import('../views/auth/EmailVerificationView.vue'),
    meta: { subtitle: "E-mail verification" }
  },
  account: guardAuth({
    path: "/account",
    name: "account",
    component: () => import("../views/AccountView.vue"),
    meta: { subtitle: "Account infos" }
  }),
}

export default {
  settings: guardRole('Admin', {
    label: 'Settings',
    icon: 'mdi-tools',
    path: '/settings/:category',
    name: "app-settings",
    params: { category: "instance" },
    component: () => import("@/views/settings/AdminSettings.vue"),
    props: true,
    meta: {
      subtitle: "Settings",
      drawer: {
        temporary: true
      }
    }
  }),
  ...accountRoutes
}