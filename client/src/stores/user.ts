import { AccountService, ApiError } from "@/api"
import type { User, UserRole } from "@/api"
import { orderedUserRoles } from "@/components/people/userRole"
import { defineStore } from "pinia"
import { computed, ref, Ref } from "vue"


export const useUserStore = defineStore("user", () => {
  const user: Ref<User | undefined> = ref(undefined)
  const error: Ref<undefined | ApiError> = ref(undefined)
  const sessionToken: Ref<string | undefined> = ref(undefined)

  async function getUser() {
    error.value = undefined
    await AccountService.currentUser()
      .then((res) => {
        if (res != undefined) {
          user.value = res.user
          console.info(`User ${user.value.identity.full_name} authenticated with role ${user.value.role}`)
          sessionToken.value = res.token
        }
      })
      .catch((reason) => {
        error.value = reason as ApiError
        user.value = undefined
        return undefined
      })
  }

  async function logout() {
    user.value = undefined
    await AccountService.logout()
  }

  const isAuthenticated = computed(() => {
    return user.value !== undefined
  })

  function isGranted(role: UserRole) {
    return user.value
      ? orderedUserRoles.indexOf(user.value.role) >= orderedUserRoles.indexOf(role)
      : false
  }

  return { user, error, token: sessionToken, getUser, logout, isGranted, isAuthenticated }
})