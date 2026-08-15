import { LoginResult, Meta, SessionTokens, User, UserCredentials, UserRole } from '@/api'
import {
  getCurrentUserOptions,
  loginMutation,
  logoutMutation,
  refreshSessionMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { client } from '@/api/gen/client.gen'
import { QueryClient, useMutation } from '@tanstack/vue-query'
import { until, useLocalStorage, useSessionStorage } from '@vueuse/core'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const refresh_token = useLocalStorage<string | undefined>('refresh_token', undefined)
  const session_id = useLocalStorage<string | undefined>('session_id', undefined)
  const user = useSessionStorage<User | undefined>('user', undefined)
  // Session expiration timestamp
  const session_expires = useSessionStorage<number | undefined>('session_expires', undefined)
  const authReady = ref<boolean>(false)
  const usePrivilege = useSessionStorage<UserRole | undefined>('usePrivilege', undefined)
  const authBootstrap = ref<Promise<void>>()
  const isAuthenticated = computed(() => user.value !== undefined)

  // Session refresh using stored refresh token
  const {
    mutateAsync: refreshSession,
    error: refreshError,
    isPending: refreshPending
  } = useMutation({
    ...refreshSessionMutation(),
    onSuccess(d) {
      useTokens(d)
    },
    onError: clearSession
  })

  function refresh() {
    if (!session_id.value) {
      return
    }
    if (!refresh_token.value) {
      console.info('Attempt to refresh user session without a refresh token')
      return
    }
    return refreshSession({
      priority: 'high',
      // Prevent infinite refresh loop
      headers: {
        noAuthRefresh: true,
        'X-Refresh-Token': refresh_token.value,
        'X-Session-ID': session_id.value
      }
    })
  }
  const refreshState = computed(() => ({
    error: refreshError.value,
    pending: refreshPending.value
  }))

  // Intercept requests to refresh session if needed
  client.interceptors.request.use(async (request) => {
    console.log('Intercepting request to check if session refresh is needed', sessionExpired())
    console.log('Refresh token:', refresh_token.value)
    if (!!refresh_token.value && sessionExpired() && !request.headers.has('noAuthRefresh')) {
      console.debug('Session expired, refreshing before request', request)
      // Prevent concurrent refresh requests
      if (!refreshPending.value) {
        await refresh()
      }
      await until(refreshPending).toBe(false)
    }
    return request
  })

  async function bootstrapAuth(client: QueryClient) {
    if (authReady.value) {
      return
    }
    console.debug('Bootstrapping auth state...')

    await client
      .fetchQuery(getCurrentUserOptions())
      .then((u) => {
        console.log('Bootstrapped auth state with user', u)
        user.value = u
        usePrivilege.value = u.role
        return
      })
      .catch((err) => {
        console.log('Using refresh token to bootstrap auth state', err)
        if (!authBootstrap.value) {
          authBootstrap.value = (async () => {
            if (refresh_token.value) {
              await refresh()
            } else {
              clearSession()
            }
          })()
        }
        return authBootstrap.value
      })
      .finally(() => {
        authReady.value = true
      })
  }

  // Login
  const {
    mutate: mutateLogin,
    error: loginError,
    isPending: loginPending
  } = useMutation({
    ...loginMutation(),
    onSuccess: startSession,
    onError(error) {
      console.error('ERROR', error)
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
  const {
    mutate: mutateLogout,
    error: logoutError,
    isPending: logoutPending
  } = useMutation({
    ...logoutMutation({ headers: { 'X-Refresh-Token': refresh_token.value } }),
    onSuccess: clearSession
  })
  function logout() {
    return mutateLogout({ headers: { 'X-Refresh-Token': refresh_token.value } })
  }
  const logoutState = computed(() => ({
    error: logoutError.value,
    pending: logoutPending.value
  }))

  function clearSession() {
    user.value = undefined
    refresh_token.value = undefined
    session_expires.value = undefined
    session_id.value = undefined
    usePrivilege.value = undefined
  }

  function startSession(data: LoginResult) {
    user.value = data.user
    usePrivilege.value = data.user.role
    useTokens(data.session)
  }

  function useTokens(tokens: SessionTokens) {
    refresh_token.value = tokens.refresh_token
    session_expires.value = tokens.auth_token_expiration.getTime()
    session_id.value = tokens.session_id

    setTimeout(
      refresh,
      Math.max(0, tokens.auth_token_expiration.getTime() - Date.now() - 30_000) // Refresh 30 seconds before expiration
    )
  }

  function sessionExpired() {
    return session_expires.value === undefined || Date.now() >= session_expires.value
  }

  function isGranted(role: UserRole) {
    return user.value
      ? UserRole.isGranted(usePrivilege.value ? { role: usePrivilege.value } : user.value, role)
      : false
  }

  function isOwner<Item extends { meta?: Meta }>(item: Item) {
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
    /**
     * Indicates whether the initial auth bootstrap has completed
     */
    authReady,
    /**
     * Bootstrap the current session state if possible
     */
    bootstrapAuth,
    /**
     * Current user role, used to determine UI privileges
     */
    usePrivilege
  }
})
