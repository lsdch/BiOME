import { UserRole } from "@/api";
import { useFeedback } from "@/stores/feedback";
import { useUserStore } from "@/stores/user";
import { RouteRecordRaw } from "vue-router";




function denyAccess(msg: string) {
  const { feedback } = useFeedback()
  feedback({ message: msg, type: "error" })
  return false
}

export function useGuards() {
  const store = useUserStore()
  function guardRole(role: UserRole, route: RouteRecordRaw & { name: string }): RouteRecordRaw {
    return {
      ...route,
      beforeEnter: () => {
        if (store.isAuthenticated) {
          return store.isGranted(role) ? true : denyAccess(`Access requires ${role} privileges`)
        } else {
          return { name: 'login', query: { redirect: route.name } }
        }
      },
    }
  }

  return { guardRole }
}