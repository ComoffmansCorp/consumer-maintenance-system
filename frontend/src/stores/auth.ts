import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { authApi, type CompanyRegistrationRequest, type LoginRequest } from '@/api/auth'
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
  tenantId?: number
  tenantCode?: string
  exp: number
}

export const useAuthStore = defineStore('auth', () => {
  const fullName = ref<string>(localStorage.getItem('cms.fullName') ?? '')
  const role = ref<Role | null>((localStorage.getItem('cms.role') as Role | null) ?? null)
  const tenantId = ref<number | null>(
    localStorage.getItem('cms.tenantId') ? Number(localStorage.getItem('cms.tenantId')) : null,
  )
  const tenantName = ref<string>(localStorage.getItem('cms.tenantName') ?? '')
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
    tenantId.value = auth.tenantId ?? null
    tenantName.value = auth.tenantName ?? ''
    localStorage.setItem('cms.fullName', auth.fullName)
    localStorage.setItem('cms.role', auth.role)
    if (auth.tenantId) localStorage.setItem('cms.tenantId', String(auth.tenantId))
    if (auth.tenantName) localStorage.setItem('cms.tenantName', auth.tenantName)

    const claims = decodeJwt<AccessClaims>(auth.accessToken)
    userId.value = claims?.userId ?? null
  }

  async function login(req: LoginRequest) {
    const auth = await authApi.login(req)
    applySession(auth)
  }

  async function registerCompany(req: CompanyRegistrationRequest) {
    const auth = await authApi.registerCompany(req)
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
    localStorage.removeItem('cms.tenantId')
    localStorage.removeItem('cms.tenantName')
    accessToken.value = null
    fullName.value = ''
    role.value = null
    tenantId.value = null
    tenantName.value = ''
    userId.value = null
  }

  return {
    fullName,
    role,
    tenantId,
    tenantName,
    userId,
    isAuthenticated,
    applySession,
    login,
    registerCompany,
    logout,
  }
})
