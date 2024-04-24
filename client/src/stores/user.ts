import { AccountService, ApiError } from "@/api"
import type { User } from "@/api"
import { defineStore } from "pinia"
import { ref, Ref } from "vue"


export const useUserStore = defineStore("user", () => {
  const user: Ref<User | undefined> = ref(undefined)
  const error: Ref<undefined | ApiError> = ref(undefined)
  const sessionToken: Ref<string | undefined> = ref(undefined)

  // Check whether a session is active
  getUser()

  async function getUser() {
    error.value = undefined
    await AccountService.currentUser()
      .then((res) => {
        if (res != undefined) {
          user.value = res.user
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

  return { user, error, token: sessionToken, getUser, logout }
})