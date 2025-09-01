import { api, API_BASE_URL } from "@/lib/api.ts"

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
    return await api.get<{ user_id: string; organization_id?: string; org_role?: string }>(
      "/api/v1/auth/me"
    )
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
