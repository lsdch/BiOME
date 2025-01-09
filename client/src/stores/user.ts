import { AccountService, AuthenticationResponse, ErrorModel, Meta, User, UserCredentials, UserRole } from "@/api"
import { useLocalStorage } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed, ref } from "vue"



export const useUserStore = defineStore("user", () => {
  const user = ref<User>()
  const error = ref<ErrorModel>()
  const session_expires = ref<Date>()
  const refresh_token = useLocalStorage<string | undefined>("refresh_token", undefined)

  function clearSession() {
    user.value = undefined
    refresh_token.value = undefined
    session_expires.value = undefined
  }

  function startSession(data: AuthenticationResponse) {
    user.value = data.user
    refresh_token.value = data.refresh_token
    session_expires.value = data.auth_token_expiration
  }

  async function login(credentials: UserCredentials) {
    const { data, error: err } = await AccountService.login({ body: credentials })
    error.value = err
    if (err) {
      clearSession()
      return err
    }
    startSession(data)
  }

  /**
   * Fetch currently authenticated user
   */
  async function getUser() {
    error.value = undefined
    await refreshAsNeeded()
    await AccountService.currentUser()
      .then(({ data, error: err, response }) => {
        if (err != undefined) {
          error.value = err
          user.value = undefined
          return
        }
        if (response.status === 204) {
          user.value = undefined
          return
        }
        user.value = data!.user
        console.info(`User ${user.value?.identity.full_name} authenticated with role ${user.value?.role}`)
      })
  }

  async function logout() {
    await AccountService.logout({ body: { refresh_token: refresh_token.value } })
    clearSession()
  }

  async function refresh() {
    if (!refresh_token.value) {
      console.warn("Attempt to refresh user session without a refresh token")
      return
    }
    const { data, error: err } = await AccountService.refreshSession({
      body: { refresh_token: refresh_token.value },
      priority: "high",
      headers: { noAuthRefresh: true }
    })

    error.value = err
    if (err) {
      clearSession()
      return
    }
    startSession(data)
  }

  async function refreshAsNeeded() {
    if (sessionExpired.value) {
      return await refresh()
    }
  }

  const isAuthenticated = computed(() => user.value !== undefined)

  const sessionExpired = computed(() =>
    session_expires.value === undefined ||
    (new Date() >= session_expires.value)
  )

  function isGranted(role: UserRole) {
    return user.value
      ? UserRole.isGranted(user.value, role)
      : false
  }

  function isOwner<
    Item extends { meta?: Meta }
  >(item: Item) {
    return user.value && item.meta?.created_by?.id === user.value.id
  }

  return {
    user, error, getUser,
    /**
     * Authenticate user and start a new session.
     * Session JWT is returned in the response, but is also saved in the cookies
     * and can be safely discarded.
     * Session refresh token is saved in local storage.
     */
    login,
    logout,
    refresh,
    refreshAsNeeded,
    isGranted,
    isOwner,
    isAuthenticated,
    sessionExpired
  }
})