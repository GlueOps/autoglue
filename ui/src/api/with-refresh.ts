// api/with-refresh.ts
import { authStore, type TokenPair } from "@/auth/store.ts"
import { API_BASE } from "@/sdkClient.ts"

let inflightRefresh: Promise<boolean> | null = null

async function doRefresh(): Promise<boolean> {
  const tokens = authStore.get()
  if (!tokens?.refresh_token) return false

  try {
    const res = await fetch(`${API_BASE}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refresh_token: tokens.refresh_token }),
    })

    if (!res.ok) return false

    const next = (await res.json()) as TokenPair
    authStore.set(next)
    return true
  } catch {
    return false
  }
}

async function refreshOnce(): Promise<boolean> {
  if (!inflightRefresh) {
    inflightRefresh = doRefresh().finally(() => {
      inflightRefresh = null
    })
  }
  return inflightRefresh
}

function isUnauthorized(err: any): boolean {
  return (
    err?.status === 401 ||
    err?.cause?.status === 401 ||
    err?.response?.status === 401 ||
    (err instanceof Response && err.status === 401)
  )
}

export async function withRefresh<T>(fn: () => Promise<T>): Promise<T> {
  // Optional: attempt a proactive refresh if close to expiry
  if (authStore.willExpireSoon?.(30)) {
    await refreshOnce()
  }

  try {
    return await fn()
  } catch (error) {
    if (!isUnauthorized(error)) throw error

    const ok = await refreshOnce()
    if (!ok) throw error

    return await fn()
  }
}
