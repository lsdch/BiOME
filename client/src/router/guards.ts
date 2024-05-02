import { UserRole } from "@/api";
import { useUserStore } from "@/stores/user";
import { RouteRecordRaw } from "vue-router";


function denyAccess(msg: string) {
  console.warn(msg)
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