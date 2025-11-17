import { memo, useMemo } from "react"
import { metaApi } from "@/api/footer"
import { useQuery } from "@tanstack/react-query"
import { Clipboard, ExternalLink, GitCommit, Info } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Separator } from "@/components/ui/separator"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

type VersionInfo = {
  built: string // ISO string or "unknown"
  builtBy: string
  commit: string
  go: string
  goArch: string
  goOS: string
  version: string
}

function shortCommit(c?: string) {
  return c && c !== "none" ? c.slice(0, 7) : "none"
}

function formatBuilt(built: string) {
  if (!built || built === "unknown") return "unknown"
  const d = new Date(built)
  return isNaN(+d) ? built : d.toLocaleString()
}

function asClipboardText(v?: VersionInfo) {
  if (!v) return ""
  return `v${v.version} (${shortCommit(v.commit)}) • built ${v.built} • ${v.go} ${v.goOS}/${v.goArch}`
}

export const Footer = memo(function Footer() {
  const footerQ = useQuery({
    queryKey: ["footer"],
    queryFn: () => metaApi.footer() as Promise<VersionInfo>,
    staleTime: 60_000,
    refetchOnWindowFocus: false,
  })

  const data = footerQ.data

  const copyText = useMemo(() => asClipboardText(data), [data])

  return (
    <footer className="bg-background text-muted-foreground w-full border-t px-3 py-2 text-xs sm:text-sm">
      <div className="mx-auto flex max-w-screen-2xl items-center justify-between">
        {/* Left: brand / copyright */}
        <div className="flex items-center gap-2 text-xs sm:text-sm">
          <span>© {new Date().getFullYear()} GlueOps</span>
          <Separator orientation="vertical" className="hidden h-4 sm:block" />
          <span className="hidden sm:block">All systems nominal.</span>
        </div>

        {/* Right: version/meta */}
        <div className="flex flex-wrap items-center gap-2 text-xs sm:text-sm">
          {footerQ.isLoading ? (
            <span className="animate-pulse">loading version…</span>
          ) : footerQ.isError ? (
            <span className="text-destructive">version unavailable</span>
          ) : data ? (
            <TooltipProvider>
              <div className="flex flex-wrap items-center gap-2">
                <Badge variant="secondary" className="font-mono">
                  {data.version}
                </Badge>

                <Tooltip>
                  <TooltipTrigger asChild>
                    <span className="inline-flex items-center gap-1">
                      <GitCommit className="h-3.5 w-3.5" />
                      <span className="font-mono">{shortCommit(data.commit)}</span>
                    </span>
                  </TooltipTrigger>
                  <TooltipContent side="top">
                    <div className="font-mono text-xs">{data.commit}</div>
                  </TooltipContent>
                </Tooltip>

                <Separator orientation="vertical" className="h-4" />

                <Tooltip>
                  <TooltipTrigger asChild>
                    <span className="inline-flex items-center gap-1">
                      <Info className="h-3.5 w-3.5" />
                      <span>{data.go}</span>
                    </span>
                  </TooltipTrigger>
                  <TooltipContent side="top">
                    <div className="font-mono text-xs">
                      {data.goOS}/{data.goArch}
                    </div>
                  </TooltipContent>
                </Tooltip>

                <Separator orientation="vertical" className="hidden h-4 sm:block" />

                <span className="hidden sm:inline">
                  built <span className="font-mono">{formatBuilt(data.built)}</span>
                </span>

                <Separator orientation="vertical" className="hidden h-4 sm:block" />

                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7"
                  onClick={() => {
                    navigator.clipboard?.writeText(copyText).catch(() => {})
                  }}
                  title="Copy version details"
                >
                  <Clipboard className="h-4 w-4" />
                </Button>

                <a
                  href="/api/v1/version"
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-xs underline-offset-4 hover:underline"
                  title="Open raw version JSON"
                >
                  JSON <ExternalLink className="h-3.5 w-3.5" />
                </a>
              </div>
            </TooltipProvider>
          ) : null}
        </div>
      </div>
    </footer>
  )
})
