import { UserRole } from "@/api";
import { useFeedback } from "@/stores/feedback";
import { useUserStore } from "@/stores/user";
import { RouteRecordRaw } from "vue-router";
import { RouteNavDefinition } from ".";




function denyAccess(msg: string) {
  const { feedback } = useFeedback()
  feedback({ message: msg, type: "error" })
  return false
}

export function useGuards() {

  function guardAuth<T extends RouteRecordRaw>(route: T): T {
    return {
      ...route,
      beforeEnter: () => {
        const store = useUserStore()
        return store.isAuthenticated || { name: 'login', query: { redirect: route.path } }
      },
    }
  }

  function guardRole<T extends RouteRecordRaw>(role: UserRole, route: T): T {
    return {
      ...route,
      beforeEnter: async () => {
        const store = useUserStore()
        if (store.isAuthenticated) {
          await store.refreshAsNeeded()
          return store.isGranted(role) ? true : denyAccess(`Access requires ${role} privileges`)
        } else {
          return { name: 'login', query: { redirect: route.path } }
        }
      },
    }
  }

  function routeWithGuard(role: UserRole, route: RouteRecordRaw & { name: string }, nav: RouteNavDefinition) {
    return {
      ...guardRole(role, route),
      ...nav,
      granted: role
    }
  }

  return { guardAuth, guardRole, routeWithGuard }
}