import { useEffect, useState, type ReactNode } from "react"
import { Navigate, Outlet, useLocation } from "react-router-dom"

import { authStore, isGlobalAdmin, type MePayload } from "@/lib/auth.ts"

type Props = { children?: ReactNode }

export function RequireAdmin({ children }: Props) {
  const [loading, setLoading] = useState(true)
  const [allowed, setAllowed] = useState(false)
  const location = useLocation()

  useEffect(() => {
    let alive = true
    ;(async () => {
      try {
        const me: MePayload = await authStore.me()
        if (!alive) return
        setAllowed(isGlobalAdmin(me))
      } catch {
        if (!alive) return
        setAllowed(false)
      } finally {
        setLoading(false)
        if (!alive) return
      }
    })()
    return () => {
      alive = false
    }
  }, [])

  if (loading) return null

  if (!allowed) return <Navigate to="/403" replace state={{ from: location }} />

  return children ? <>{children}</> : <Outlet />
}
