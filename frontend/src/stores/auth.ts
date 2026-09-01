import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi, type LoginRequest, type RegisterRequest } from '@/api/auth'
import { clearSession, getAccessToken, getRefreshToken, storeSession } from '@/api/client'
import type { AuthResponse, Role } from '@/types'

function decodeJwt<T>(token: string): T | null {
  try {
    const payload = token.split('.')[1]
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')))
  } catch {
    return null
  }
}

interface AccessClaims {
  userId: number
  role: Role
  exp: number
}

export const useAuthStore = defineStore('auth', () => {
  const fullName = ref<string>(localStorage.getItem('cms.fullName') ?? '')
  const role = ref<Role | null>((localStorage.getItem('cms.role') as Role | null) ?? null)
  const userId = ref<number | null>(null)

  // Kept as a ref (not read from localStorage inline in the computed below):
  // with `&&` short-circuiting, if the token starts out empty, `role.value`
  // would never be read on that first evaluation, so Vue would never track
  // it as a dependency and `isAuthenticated` would stay cached at `false`
  // forever, even after a successful login.
  const accessToken = ref<string | null>(getAccessToken())
  if (accessToken.value) {
    const claims = decodeJwt<AccessClaims>(accessToken.value)
    userId.value = claims?.userId ?? null
  }

  const isAuthenticated = computed(() => {
    const hasToken = !!accessToken.value
    const hasRole = !!role.value
    return hasToken && hasRole
  })

  function applySession(auth: AuthResponse) {
    storeSession(auth)
    accessToken.value = auth.accessToken
    fullName.value = auth.fullName
    role.value = auth.role
    localStorage.setItem('cms.fullName', auth.fullName)
    localStorage.setItem('cms.role', auth.role)

    const claims = decodeJwt<AccessClaims>(auth.accessToken)
    userId.value = claims?.userId ?? null
  }

  async function login(req: LoginRequest) {
    const auth = await authApi.login(req)
    applySession(auth)
  }

  async function registerClient(req: RegisterRequest) {
    const auth = await authApi.registerClient(req)
    applySession(auth)
  }

  async function registerMaster(req: RegisterRequest) {
    const auth = await authApi.registerMaster(req)
    applySession(auth)
  }

  async function logout() {
    const refreshToken = getRefreshToken()
    try {
      if (refreshToken) await authApi.logout(refreshToken)
    } catch {
      // best-effort — clear local session regardless
    }
    clearSession()
    localStorage.removeItem('cms.fullName')
    localStorage.removeItem('cms.role')
    accessToken.value = null
    fullName.value = ''
    role.value = null
    userId.value = null
  }

  return {
    fullName,
    role,
    userId,
    isAuthenticated,
    applySession,
    login,
    registerClient,
    registerMaster,
    logout,
  }
})
