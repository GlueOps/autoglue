import { useEffect, useState } from "react"

import { api, ApiError } from "@/lib/api.ts"
import {
  EVT_ACTIVE_ORG_CHANGED,
  EVT_ORGS_CHANGED,
  getActiveOrgId,
  setActiveOrgId,
} from "@/lib/orgs-sync.ts"
import { Button } from "@/components/ui/button.tsx"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx"

type OrgLite = { id: string; name: string }

export const OrgSwitcher = () => {
  const [orgs, setOrgs] = useState<OrgLite[]>([])
  const [activeOrgId, setActiveOrgIdState] = useState<string | null>(null)

  async function fetchOrgs() {
    try {
      const data = await api.get<OrgLite[]>("/api/v1/orgs")
      setOrgs(data)
      if (!getActiveOrgId() && data.length > 0) {
        // default to first org if none selected yet
        setActiveOrgId(data[0].id)
        setActiveOrgIdState(data[0].id)
      }
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to load organizations"
      // optional: toast.error(msg);
      console.error(msg)
    }
  }

  useEffect(() => {
    // initial load
    setActiveOrgIdState(getActiveOrgId())
    void fetchOrgs()

    // cross-tab sync
    const onStorage = (e: StorageEvent) => {
      if (e.key === "active_org_id") setActiveOrgIdState(e.newValue)
    }
    window.addEventListener("storage", onStorage)

    // same-tab sync: active org + orgs list mutations
    const onActive = (e: Event) =>
      setActiveOrgIdState((e as CustomEvent<string | null>).detail ?? null)
    const onOrgs = () => void fetchOrgs()

    window.addEventListener(EVT_ACTIVE_ORG_CHANGED, onActive as EventListener)
    window.addEventListener(EVT_ORGS_CHANGED, onOrgs)

    return () => {
      window.removeEventListener("storage", onStorage)
      window.removeEventListener(EVT_ACTIVE_ORG_CHANGED, onActive as EventListener)
      window.removeEventListener(EVT_ORGS_CHANGED, onOrgs)
    }
  }, [])

  const switchOrg = (orgId: string) => {
    setActiveOrgId(orgId)
    setActiveOrgIdState(orgId)
  }

  const currentOrgName = orgs.find((o) => o.id === activeOrgId)?.name ?? "Select Org"

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="w-full justify-start">
          {currentOrgName}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-48">
        {orgs.length === 0 ? (
          <DropdownMenuItem disabled>No organizations</DropdownMenuItem>
        ) : (
          orgs.map((org) => (
            <DropdownMenuItem
              key={org.id}
              onClick={() => switchOrg(org.id)}
              className={org.id === activeOrgId ? "font-semibold" : undefined}
            >
              {org.name}
            </DropdownMenuItem>
          ))
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
