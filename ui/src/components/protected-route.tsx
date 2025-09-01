import { type ReactNode } from "react"
import { Navigate, Outlet, useLocation } from "react-router-dom"

import { authStore } from "@/lib/auth.ts"

export function ProtectedRoute({ children }: { children?: ReactNode }) {
  const location = useLocation()
  const isAuthed = authStore.isAuthenticated()
  if (!isAuthed) {
    return <Navigate to="/auth/login" state={{ from: location }} replace />
  }
  return children ? <>{children}</> : <Outlet />
}
