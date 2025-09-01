import { api, API_BASE_URL } from "@/lib/api.ts"

export type MeUser = {
  id: string
  name?: string
  email?: string
  email_verified?: boolean
  role: "admin" | "user" | string
  created_at?: string
  updated_at?: string
}

export type MePayload = {
  user?: MeUser // preferred shape
  user_id?: MeUser // fallback (older shape)
  organization_id?: string | null
  org_role?: "admin" | "member" | string | null
  claims?: any
}

function getUser(me: MePayload | null | undefined): MeUser | undefined {
  return (me && (me.user || me.user_id)) as MeUser | undefined
}
export function isGlobalAdmin(me: MePayload | null | undefined): boolean {
  return getUser(me)?.role === "admin"
}
export function isOrgAdmin(me: MePayload | null | undefined): boolean {
  return (me?.org_role ?? "") === "admin"
}

export const authStore = {
  isAuthenticated(): boolean {
    return !!localStorage.getItem("access_token")
  },

  async login(email: string, password: string) {
    const data = await api.post<{ access_token: string; refresh_token: string }>(
      "/api/v1/auth/login",
      { email, password }
    )
    localStorage.setItem("access_token", data.access_token)
    localStorage.setItem("refresh_token", data.refresh_token)
  },

  async register(name: string, email: string, password: string) {
    await api.post("/api/v1/auth/register", { name, email, password })
  },

  async me() {
    return await api.get<MePayload>("/api/v1/auth/me")
  },

  async logout() {
    const rt = localStorage.getItem("refresh_token")
    if (rt) {
      try {
        await api.post("/api/v1/auth/logout", { refresh_token: rt })
      } catch {}
    }
    localStorage.removeItem("access_token")
    localStorage.removeItem("refresh_token")
  },

  async forgot(email: string) {
    await api.post("/api/v1/auth/password/forgot", { email })
  },

  async reset(token: string, new_password: string) {
    await api.post("/api/v1/auth/password/reset", { token, new_password })
  },

  async verify(token: string) {
    // GET with token query
    const res = await fetch(`${API_BASE_URL}/api/v1/auth/verify?token=${encodeURIComponent(token)}`)
    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg)
    }
  },
}
