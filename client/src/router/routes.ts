import { useGuards } from "./guards"

const { guardRole } = useGuards()

export default {
  settings: guardRole('Admin', {
    label: 'Settings',
    icon: 'mdi-tools',
    path: '/settings',
    name: "app-settings",
    component: () => import("@/views/settings/AdminSettings.vue")
  })
}