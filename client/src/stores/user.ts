import type { ErrorModel, User, UserRole } from "@/api"
import { AccountService } from "@/api"
import { orderedUserRoles } from "@/components/people/userRole"
import { defineStore } from "pinia"
import { computed, ref } from "vue"


export const useUserStore = defineStore("user", () => {
  const user = ref<User>()
  const error = ref<ErrorModel>()
  const sessionToken = ref<string>()

  async function getUser() {
    error.value = undefined
    await AccountService.currentUser()
      .then(({ data, error: err, response }) => {
        if (err != undefined) {
          error.value = err
          user.value = undefined
          return
        }
        if (response.status === 204) {
          user.value = undefined
          sessionToken.value = undefined
          return
        }
        user.value = data!.user
        sessionToken.value = data!.token
        console.info(`User ${user.value?.identity.full_name} authenticated with role ${user.value?.role}`)
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