import { defineStore } from "pinia"
import { ref } from "vue"
import type { User } from "@/api/models/User"
import { Ref } from "vue"
import { PeopleService } from "@/api"

export const useUserStore = defineStore("user", () => {
  const jwt: Ref<string | undefined> = ref(undefined)
  const user: Ref<User | undefined> = ref(undefined)

  async function setToken(token: string) {
    jwt.value = token
    user.value = await PeopleService.currentUser().catch()
  }
})