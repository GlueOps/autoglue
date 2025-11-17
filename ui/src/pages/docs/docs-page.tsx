import { useEffect, useRef, useState, type FC } from "react"
import { useTheme } from "next-themes"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"

import "rapidoc"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx"

type RdThemeMode = "auto" | "light" | "dark"

export const DocsPage: FC = () => {
  const rdRef = useRef<any>(null)
  const { theme, systemTheme, setTheme } = useTheme()

  const [orgId, setOrgId] = useState("")
  const [rdThemeMode, setRdThemeMode] = useState<RdThemeMode>("auto")

  useEffect(() => {
    const stateSetter = () => {
      const stored = localStorage.getItem("autoglue.org")
      if (stored) setOrgId(stored)
    }

    stateSetter()
  }, [])

  useEffect(() => {
    const rd = rdRef.current
    if (!rd) return

    let effectiveTheme: "light" | "dark" = "light"
    if (rdThemeMode === "light") {
      effectiveTheme = "light"
    } else if (rdThemeMode === "dark") {
      effectiveTheme = "dark"
    } else {
      const appTheme = theme === "system" ? systemTheme : theme
      effectiveTheme = appTheme === "dark" ? "dark" : "light"
    }

    rd.setAttribute("theme", effectiveTheme)

    if (typeof window !== "undefined") {
      const defaultServer = `${window.location.origin}/api/v1`
      rd.setAttribute("default-api-server", defaultServer)
    }

    if (orgId) {
      rd.setAttribute("api-key-name", "X-ORG-ID")
      rd.setAttribute("api-key-location", "header")
      rd.setAttribute("api-key-value", orgId)
    } else {
      rd.removeAttribute("api-key-value")
    }
  }, [theme, systemTheme, rdThemeMode, orgId])

  const handleSaveOrg = () => {
    const trimmed = orgId.trim()
    localStorage.setItem("autoglue.org", trimmed)
    const rd = rdRef.current
    if (!rd) return

    if (trimmed) {
      rd.setAttribute("api-key-value", trimmed)
    } else {
      rd.removeAttribute("api-key-value")
    }
  }

  const handleResetOrg = () => {
    localStorage.removeItem("autoglue.org")
    setOrgId("")
    const rd = rdRef.current
    if (!rd) return
    rd.removeAttribute("api-key-value")
  }

  return (
    <div className="flex h-[100svh] flex-col">
      {/* Control bar */}
      <Card className="rounded-none border-b">
        <CardHeader className="py-3">
          <CardTitle className="flex flex-wrap items-center justify-between gap-4 text-base">
            <span>AutoGlue API Docs</span>
            <div className="flex items-center gap-2 text-xs">
              <div className="flex flex-wrap items-center gap-3 text-xs">
                {/* Theme selector */}
                <div className="flex items-center gap-2">
                  <span className="text-muted-foreground">Docs theme</span>
                  <Select
                    value={rdThemeMode}
                    onValueChange={(v) => {
                      const mode = v as RdThemeMode
                      setRdThemeMode(mode)

                      if (mode === "auto") {
                        setTheme("system")
                      } else {
                        setTheme(v)
                      }
                    }}
                  >
                    <SelectTrigger className="h-8 w-[120px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="auto">Auto (match app)</SelectItem>
                      <SelectItem value="light">Light</SelectItem>
                      <SelectItem value="dark">Dark</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Org ID controls */}
                <div className="flex items-center gap-2">
                  <span className="text-muted-foreground">Org ID (X-ORG-ID)</span>
                  <Input
                    className="h-8 w-80"
                    value={orgId}
                    onChange={(e) => setOrgId(e.target.value)}
                    placeholder="org_..."
                  />
                  <Button size="sm" onClick={handleSaveOrg}>
                    Save
                  </Button>
                  <Button size="sm" variant="outline" onClick={handleResetOrg}>
                    Reset
                  </Button>
                </div>
              </div>
            </div>
          </CardTitle>
        </CardHeader>
        <CardContent className="text-muted-foreground py-0 pb-2 text-xs">
          Requests from <code>&lt;rapi-doc&gt;</code> will include:
          <code className="ml-1">Cookie: ag_jwt=…</code> and{" "}
          <code className="ml-1">X-ORG-ID={orgId}</code>
          {!orgId && <> (set an Org ID above to send an X-ORG-ID header)</>}
        </CardContent>
      </Card>

      {/* @ts-expect-error ts-2339 */}
      <rapi-doc
        ref={rdRef}
        id="autoglue-docs"
        spec-url="/swagger/swagger.json"
        render-style="read"
        show-header="false"
        persist-auth="true"
        allow-advanced-search="true"
        schema-description-expanded="true"
        allow-schema-description-expand-toggle="false"
        allow-spec-file-download="true"
        allow-spec-file-load="false"
        allow-spec-url-load="false"
        allow-try="true"
        schema-style="tree"
        fetch-credentials="include"
      />
    </div>
  )
}