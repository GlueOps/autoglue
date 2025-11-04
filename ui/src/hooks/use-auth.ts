import { useSyncExternalStore } from "react"
import { authStore, type TokenPair } from "@/auth/store.ts"

export const useAuth = () => {
  const tokens = useSyncExternalStore<TokenPair | null>(
    (cb) => authStore.subscribe(cb),
    () => authStore.get(),
    () => authStore.get() // server snapshot (SSR)
  )

  return {
    tokens,
    authed: !!tokens?.access_token,
    isExpired: authStore.isExpired(),
    willExpireSoon: authStore.willExpireSoon(),
  }
}
