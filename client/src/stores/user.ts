import { defineStore } from "pinia"
import { ref, Ref } from "vue"
import type { User } from "@/api/models/User"
import { ApiError, PeopleService, AuthService } from "@/api"


export const useUserStore = defineStore("user", () => {
  const user: Ref<User | undefined> = ref(undefined)
  const error: Ref<undefined | ApiError> = ref(undefined)

  // Check whether a session is active
  getUser()

  async function getUser() {
    user.value = await PeopleService.currentUser()
      .catch((err: ApiError) => {
        error.value = err
        user.value = undefined
        return undefined
      })
  }

  async function logout() {
    user.value = undefined
    await AuthService.logout()
  }

  return { user, error, getUser, logout }
})