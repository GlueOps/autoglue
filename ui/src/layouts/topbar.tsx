import { useMemo } from "react"
import { Link, useLocation } from "react-router-dom"

import { useMe } from "@/hooks/use-me.ts"
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx"
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb.tsx"
import { Button } from "@/components/ui/button.tsx"
import { SidebarTrigger } from "@/components/ui/sidebar.tsx"

export const Topbar = () => {
  const loc = useLocation()
  const { data: me, isLoading } = useMe()

  const crumbs = useMemo(() => {
    const parts = loc.pathname.split("/").filter(Boolean)
    const acc: { to: string; label: string }[] = []
    let build = ""
    for (const p of parts) {
      build += `/${p}`
      acc.push({ to: build, label: p })
    }
    return acc
  }, [loc.pathname])

  const initials = useMemo(() => {
    if (!me) return "U"
    const name = me.display_name || me.primary_email || ""
    const parts = name.trim().split(/\s+/)
    if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
    if (parts.length === 1 && parts[0]) return parts[0][0]!.toUpperCase()
    return "U"
  }, [me])

  return (
    <div className="flex h-12 items-center gap-2 border-b px-3">
      <SidebarTrigger />
      <div className="flex-1">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <BreadcrumbLink asChild>
                <Link to="/">Home</Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            {crumbs.map((c, i) => (
              <span key={c.to} className="flex items-center">
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                  {i === crumbs.length - 1 ? (
                    <BreadcrumbPage className="capitalize">{c.label}</BreadcrumbPage>
                  ) : (
                    <BreadcrumbLink asChild>
                      <Link to={c.to} className="capitalize">
                        {c.label}
                      </Link>
                    </BreadcrumbLink>
                  )}
                </BreadcrumbItem>
              </span>
            ))}
          </BreadcrumbList>
        </Breadcrumb>
      </div>

      <Button variant="ghost" size="sm" asChild>
        <Link to="/me">{isLoading ? "…" : me?.display_name || "Profile"}</Link>
      </Button>

      <Avatar className="h-7 w-7">
        <AvatarFallback>{initials}</AvatarFallback>
      </Avatar>
    </div>
  )
}
