import type { User, LoginCredentials, RegisterCredentials } from './useApi'
import { useAuthApi, setAuthErrorHandler } from './useApi'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  loading: boolean
}

const state = ref<AuthState>({
  user: null,
  token: null,
  isAuthenticated: false,
  loading: false,
})

export function useAuth() {
  const router = useRouter()
  const authApi = useAuthApi()

  const initAuth = () => {
    if (import.meta.client) {
      const token = localStorage.getItem('auth_token')
      const userStr = localStorage.getItem('auth_user')

      if (token && userStr) {
        try {
          const user = JSON.parse(userStr)
          state.value.token = token
          state.value.user = user
          state.value.isAuthenticated = true
        } catch {
          clearAuth()
        }
      }
    }
  }

  const clearAuth = () => {
    state.value.user = null
    state.value.token = null
    state.value.isAuthenticated = false

    if (import.meta.client) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    }
  }

  const saveAuth = (token: string, user: User) => {
    state.value.token = token
    state.value.user = user
    state.value.isAuthenticated = true

    if (import.meta.client) {
      localStorage.setItem('auth_token', token)
      localStorage.setItem('auth_user', JSON.stringify(user))
    }
  }

  const login = async (credentials: LoginCredentials): Promise<{ success: boolean; error?: string }> => {
    state.value.loading = true

    try {
      const response = await authApi.login(credentials)

      if (response.success && response.data) {
        saveAuth(response.data.token, response.data.user)
        return { success: true }
      } else {
        return { success: false, error: response.error || 'Login failed' }
      }
    } catch (error: any) {
      return { success: false, error: error?.message || 'Network error' }
    } finally {
      state.value.loading = false
    }
  }

  const register = async (credentials: RegisterCredentials): Promise<{ success: boolean; error?: string }> => {
    state.value.loading = true

    try {
      const response = await authApi.register(credentials)

      if (response.success && response.data) {
        saveAuth(response.data.token, response.data.user)
        return { success: true }
      } else {
        return { success: false, error: response.error || 'Registration failed' }
      }
    } catch (error: any) {
      return { success: false, error: error?.message || 'Network error' }
    } finally {
      state.value.loading = false
    }
  }

  const logout = () => {
    clearAuth()
    router.push('/login')
  }

  const fetchCurrentUser = async (): Promise<boolean> => {
    if (!state.value.token) {
      return false
    }

    try {
      const response = await authApi.getCurrentUser()

      if (response.success && response.data) {
        state.value.user = response.data
        if (import.meta.client) {
          localStorage.setItem('auth_user', JSON.stringify(response.data))
        }
        return true
      } else {
        clearAuth()
        return false
      }
    } catch {
      clearAuth()
      return false
    }
  }

  const handleAuthError = () => {
    clearAuth()
    router.push('/login')
  }

  if (import.meta.client) {
    setAuthErrorHandler(handleAuthError)
  }

  return {
    user: computed(() => state.value.user),
    token: computed(() => state.value.token),
    isAuthenticated: computed(() => state.value.isAuthenticated),
    loading: computed(() => state.value.loading),
    initAuth,
    login,
    register,
    logout,
    fetchCurrentUser,
    clearAuth,
  }
}
