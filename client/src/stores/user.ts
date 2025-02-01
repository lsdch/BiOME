import { AuthenticationResponse, Meta, User, UserCredentials, UserRole } from "@/api"
import { loginMutation, logoutMutation, refreshSessionMutation } from "@/api/gen/@tanstack/vue-query.gen"
import { client } from '@/api/gen/client.gen'
import { useMutation } from "@tanstack/vue-query"
import { until, useLocalStorage, useTimeoutPoll } from "@vueuse/core"
import { defineStore } from "pinia"
import { computed, ref } from "vue"


export const useUserStore = defineStore("user", () => {


  const user = ref<User>()
  const session_expires = ref<Date>()
  const refresh_token = useLocalStorage<string | undefined>("refresh_token", undefined)
  const isAuthenticated = computed(() => user.value !== undefined)

  // Session refresh using stored refresh token
  const { mutate: refreshSession, error: refreshError, isPending: refreshPending } = useMutation({
    ...refreshSessionMutation(),
    onSuccess: startSession,
    onError: clearSession
  })
  function refresh() {
    if (!refresh_token.value) {
      console.info("Attempt to refresh user session without a refresh token")
      return
    }
    return refreshSession({
      body: { refresh_token: refresh_token.value },
      priority: "high",
      // Prevent infinite refresh loop
      headers: { noAuthRefresh: true }
    })
  }
  const refreshState = computed(() => ({
    error: refreshError.value,
    pending: refreshPending.value
  }))

  // Intercept requests to refresh session if needed
  client.interceptors.request.use(async (request) => {
    if (isAuthenticated && sessionExpired() && !request.headers.has('noAuthRefresh')) {
      // Prevent concurrent refresh requests
      if (!refreshPending.value) {
        refresh()
      }
      await until(refreshPending).toBe(false)
    }
    return request
  })


  // Login
  const { mutate: mutateLogin, error: loginError, isPending: loginPending } = useMutation({
    ...loginMutation(),
    onSuccess: startSession,
    onError(error) {
      console.error("ERROR", error)
      clearSession()
    }
  })
  function login(credentials: UserCredentials) {
    return mutateLogin({ body: credentials })
  }
  const loginState = computed(() => ({
    error: loginError.value,
    pending: loginPending.value
  }))

  // Logout
  const { mutate: mutateLogout, error: logoutError, isPending: logoutPending } = useMutation({
    ...logoutMutation({ body: { refresh_token: refresh_token.value } }),
    onSuccess: clearSession
  })
  function logout() {
    return mutateLogout({ body: { refresh_token: refresh_token.value } })
  }
  const logoutState = computed(() => ({
    error: logoutError.value,
    pending: logoutPending.value
  }))


  function clearSession() {
    user.value = undefined
    refresh_token.value = undefined
    session_expires.value = undefined
  }

  function startSession(data: AuthenticationResponse) {
    user.value = data.user
    refresh_token.value = data.refresh_token
    session_expires.value = data.auth_token_expiration
    // Refresh session before it expires
    setTimeout(refresh, data.auth_token_expiration.getTime() - Date.now() - 30_000)
  }

  function sessionExpired() {
    return session_expires.value === undefined ||
      (new Date() >= session_expires.value)
  }


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
    /**
     * Currently authenticated user
     */
    user,
    /**
     * Authenticate user and start a new session.
     * Session JWT is returned in the response, but is also saved in the cookies
     * and can be safely discarded.
     * Session refresh token is saved in local storage.
     */
    login,
    /**
     * Login query state
     */
    loginState,
    /**
     * End the current session and clear all session data
     */
    logout,
    /**
     * Logout query state
     */
    logoutState,
    /**
     * Refresh the current session using the stored refresh token
     */
    refreshSession: refresh,
    /**
     * Refresh session query state
     */
    refreshState,
    /**
     * Checks if currently authenticated user has sufficient privileges
     */
    isGranted,
    /**
     * Checks whether the currently authenticated user is the owner of an item,
     * based on the item's metadata
     */
    isOwner,
    /**
     * Checks if the current session has expired
     */
    sessionExpired,
    /**
     * Indicates if the user is currently authenticated
     */
    isAuthenticated,
  }
})