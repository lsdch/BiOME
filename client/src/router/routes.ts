import { useGuards } from "./guards"

const { guardRole, guardAuth } = useGuards()

export const accountRoutes = {
  login: {
    path: '/login',
    name: 'login',
    component: () => import('../views/auth/LoginView.vue'),
    // meta: { hideNavbar: true }
  },
  signup: {
    path: '/signup',
    name: 'signup',
    component: () => import('../views/auth/SignUpView.vue'),
    // meta: { hideNavbar: true }
  },
  pwdReset: {
    path: '/password-reset',
    name: 'password-reset',
    component: () => import('../views/auth/PasswordResetView.vue'),
  },
  verifyEmail: {
    path: '/verify-email',
    name: 'verify-email',
    component: () => import('../views/auth/EmailVerificationView.vue'),
  },
  account: guardAuth({
    path: "/account",
    name: "account",
    component: () => import("../views/AccountView.vue")
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
      drawer: {
        temporary: true
      }
    }
  }),
  ...accountRoutes
}