import { useEffect, useMemo, useState, type ComponentType, type FC } from "react"
import { ChevronDown } from "lucide-react"
import { Link, useLocation } from "react-router-dom"

import { authStore, isGlobalAdmin, isOrgAdmin, type MePayload } from "@/lib/auth.ts"
import { Button } from "@/components/ui/button.tsx"
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx"
import { ModeToggle } from "@/components/mode-toggle.tsx"
import { OrgSwitcher } from "@/components/org-switcher.tsx"
import { items, type NavItem } from "@/components/sidebar/items.ts"

interface MenuItemProps {
  label: string
  icon: ComponentType<{ className?: string }>
  to?: string
  items?: MenuItemProps[]
}

function filterItems(items: NavItem[], isAdmin: boolean, isOrgAdminFlag: boolean): NavItem[] {
  return items
    .filter((it) => {
      if (it.requiresAdmin && !isAdmin) return false
      if (it.requiresOrgAdmin && !isOrgAdminFlag) return false
      return true
    })
    .map((it) => ({
      ...it,
      items: it.items ? filterItems(it.items, isAdmin, isOrgAdminFlag) : undefined,
    }))
    .filter((it) => !it.items || it.items.length > 0)
}

const MenuItem: FC<{ item: MenuItemProps }> = ({ item }) => {
  const location = useLocation()
  const Icon = item.icon

  if (item.to) {
    return (
      <Link
        to={item.to}
        className={`hover:bg-accent hover:text-accent-foreground flex items-center space-x-2 rounded-md px-4 py-2 text-sm ${location.pathname === item.to ? "bg-accent text-accent-foreground" : ""}`}
      >
        <Icon className="mr-4 h-4 w-4" />
        {item.label}
      </Link>
    )
  }

  if (item.items) {
    return (
      <Collapsible defaultOpen className="group/collapsible">
        <SidebarGroup>
          <SidebarGroupLabel asChild>
            <CollapsibleTrigger>
              <Icon className="mr-4 h-4 w-4" />
              {item.label}
              <ChevronDown className="ml-auto transition-transform group-data-[state=open]/collapsible:rotate-180" />
            </CollapsibleTrigger>
          </SidebarGroupLabel>
          <CollapsibleContent>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.items.map((subitem, index) => (
                  <SidebarMenuItem key={index}>
                    <SidebarMenuButton asChild>
                      <MenuItem item={subitem} />
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </CollapsibleContent>
        </SidebarGroup>
      </Collapsible>
    )
  }
  return null
}

export const DashboardSidebar = () => {
  const [me, setMe] = useState<MePayload | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let alive = true
    ;(async () => {
      try {
        const data = await authStore.me()
        if (!alive) return
        setMe(data)
      } catch {
        // ignore; unauthenticated users shouldn't be here anyway under ProtectedRoute
      } finally {
        // if (!alive) return
        setLoading(false)
      }
    })()
    return () => {
      alive = false
    }
  }, [])

  const visibleItems = useMemo(() => {
    const admin = isGlobalAdmin(me)
    const orgAdmin = isOrgAdmin(me)
    return filterItems(items, admin, orgAdmin)
  }, [me])

  if (loading) return <div className="p-6">Loading…</div>

  return (
    <Sidebar>
      <SidebarHeader className="flex items-center justify-between p-4">
        <h1 className="text-xl font-bold">AutoGlue</h1>
      </SidebarHeader>
      <SidebarContent>
        {visibleItems.map((item, i) => (
          <MenuItem item={item} key={i} />
        ))}
      </SidebarContent>
      <SidebarFooter className="space-y-2 p-4">
        <OrgSwitcher />
        <ModeToggle />
        <Button
          onClick={() => {
            localStorage.clear()
            window.location.reload()
          }}
          className="w-full"
        >
          Logout
        </Button>
      </SidebarFooter>
    </Sidebar>
  )
}
