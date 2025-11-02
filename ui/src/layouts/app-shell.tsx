import { useEffect, useState } from "react"
import { meApi } from "@/api/me.ts"
import { orgStore } from "@/auth/org.ts"
import { authStore } from "@/auth/store.ts"
import { mainNav, orgNav, userNav } from "@/layouts/nav-config.ts"
import { OrgSwitcher } from "@/layouts/org-switcher.tsx"
import { Topbar } from "@/layouts/topbar.tsx"
import { NavLink, Outlet } from "react-router-dom"

import { cn } from "@/lib/utils.ts"
import { Button } from "@/components/ui/button.tsx"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from "@/components/ui/sidebar.tsx"

type Org = {
  id: string
  name: string
}

export const AppShell = () => {
  const [orgs, setOrgs] = useState<Org[]>([])

  useEffect(() => {
    let alive = true
    ;(async () => {
      try {
        const me = await meApi.getMe() // HandlersMeResponse
        const list = (me.organizations ?? []).map((o) => ({
          id: o.id,
          name: o.name ?? o.id,
        }))
        if (!alive) return
        setOrgs(list as Org[])

        // default selection if none
        if (!orgStore.get() && list.length > 0) {
          orgStore.set(list[0].id!)
        }
      } catch {
        // ignore; ProtectedRoute will handle auth
      }
    })()
    return () => {
      alive = false
    }
  }, [])

  return (
    <SidebarProvider defaultOpen>
      <Sidebar collapsible="icon" variant="floating">
        <SidebarHeader>
          <div className="px-2 py-2">
            <OrgSwitcher orgs={orgs} />
          </div>
        </SidebarHeader>

        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Navigation</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {mainNav.map((n) => (
                  <SidebarMenuItem key={n.to}>
                    <SidebarMenuButton asChild tooltip={n.label}>
                      <NavLink
                        to={n.to}
                        className={({ isActive }) =>
                          cn("flex items-center gap-2", isActive && "text-primary")
                        }
                      >
                        <n.icon className="h-4 w-4" />
                        <span>{n.label}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>

          <SidebarGroup>
            <SidebarGroupLabel>Organization</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {orgNav.map((n) => (
                  <SidebarMenuItem key={n.to}>
                    <SidebarMenuButton asChild tooltip={n.label}>
                      <NavLink
                        to={n.to}
                        className={({ isActive }) =>
                          cn("flex items-center gap-2", isActive && "text-primary")
                        }
                      >
                        <n.icon className="h-4 w-4" />
                        <span>{n.label}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>

          <SidebarGroup>
            <SidebarGroupLabel>User</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {userNav.map((n) => (
                  <SidebarMenuItem key={n.to}>
                    <SidebarMenuButton asChild tooltip={n.label}>
                      <NavLink
                        to={n.to}
                        className={({ isActive }) =>
                          cn("flex items-center gap-2", isActive && "text-primary")
                        }
                      >
                        <n.icon className="h-4 w-4" />
                        <span>{n.label}</span>
                      </NavLink>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>

        <SidebarFooter>
          <div className="px-2 py-2">
            <Button variant="ghost" size="sm" className="w-full" onClick={() => authStore.logout()}>
              Sign out
            </Button>
          </div>
        </SidebarFooter>
      </Sidebar>

      <SidebarInset className="min-h-screen">
        <Topbar />
        <main className="p-4">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  )
}
