import { useCallback } from "react"
import { logoutEverywhere } from "@/auth/logout.ts"

export function useAuthActions() {
  const logout = useCallback(() => {
    return logoutEverywhere()
  }, [])

  return { logout }
}
