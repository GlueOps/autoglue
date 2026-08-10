import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import type { DtoJobLogPage } from "@/sdk"
import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react"

const POLL_INTERVAL_MS = 1500
const MAX_RENDERED_LINES = 5000

type JobLogLine = {
  key: string
  text: string
  system: boolean
}

type Props = {
  /** Fetches one page from the cursor. Must be stable, or polling restarts. */
  fetchPage: (after: number) => Promise<DtoJobLogPage>
  /** Poll while true. Pass false to load once, e.g. for a finished run. */
  live?: boolean
  className?: string
  emptyMessage?: string
}

/**
 * Tails a background job's output.
 *
 * Chunks are appended by cursor rather than refetched, so a long bootstrap does
 * not re-download its whole transcript every poll. Polling stops when the API
 * reports `done`, which it withholds until the reader has caught up.
 */
export function JobLogViewer({
  fetchPage,
  live = true,
  className,
  emptyMessage = "No output yet.",
}: Props) {
  const [lines, setLines] = useState<JobLogLine[]>([])
  const [done, setDone] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [follow, setFollow] = useState(true)

  const cursorRef = useRef(0)
  const scrollRef = useRef<HTMLDivElement | null>(null)
  // Guards against two polls overlapping if one response is slow.
  const inFlightRef = useRef(false)

  const poll = useCallback(async () => {
    if (inFlightRef.current) return
    inFlightRef.current = true
    try {
      const page = await fetchPage(cursorRef.current)
      const items = page.items ?? []

      if (items.length > 0) {
        setLines((prev) => {
          const next = [...prev]
          for (const item of items) {
            const system = item.stream === "system"
            const body = item.chunk ?? ""
            // A chunk is a batch of output, not a line: split so the viewer can
            // cap what it renders without cutting mid-line.
            for (const text of body.split("\n")) {
              if (text === "" ) continue
              next.push({ key: `${item.id}-${next.length}`, text, system })
            }
          }
          return next.length > MAX_RENDERED_LINES
            ? next.slice(next.length - MAX_RENDERED_LINES)
            : next
        })
      }

      if (typeof page.next_cursor === "number") cursorRef.current = page.next_cursor
      if (page.done) setDone(true)
      setError(null)
    } catch (e: any) {
      setError(e?.message ?? "Failed to load logs")
    } finally {
      inFlightRef.current = false
    }
  }, [fetchPage])

  useEffect(() => {
    void poll()
    if (!live) return
    const t = setInterval(() => {
      if (!done) void poll()
    }, POLL_INTERVAL_MS)
    return () => clearInterval(t)
  }, [poll, live, done])

  // Stick to the bottom while following. useLayoutEffect so the jump happens
  // before paint rather than as a visible lurch.
  useLayoutEffect(() => {
    if (!follow) return
    const el = scrollRef.current
    if (el) el.scrollTop = el.scrollHeight
  }, [lines, follow])

  const onScroll = () => {
    const el = scrollRef.current
    if (!el) return
    // Re-arm following only when the user returns to the bottom themselves.
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 24
    setFollow(atBottom)
  }

  return (
    <div className={cn("flex flex-col gap-2", className)}>
      <div className="text-muted-foreground flex items-center gap-2 text-xs">
        <span
          className={cn(
            "inline-block size-2 rounded-full",
            done ? "bg-muted-foreground" : "animate-pulse bg-emerald-500",
          )}
          aria-hidden
        />
        <span>{done ? "Finished" : "Streaming"}</span>
        <span className="ml-auto tabular-nums">{lines.length} lines</span>
        {!follow && (
          <Button
            size="sm"
            variant="ghost"
            className="h-6 px-2 text-xs"
            onClick={() => setFollow(true)}
          >
            Jump to latest
          </Button>
        )}
      </div>

      {error && <p className="text-destructive text-xs">{error}</p>}

      <div
        ref={scrollRef}
        onScroll={onScroll}
        role="log"
        aria-live="polite"
        className="bg-muted/40 h-80 overflow-auto rounded-md border p-3 font-mono text-xs leading-relaxed"
      >
        {lines.length === 0 ? (
          <p className="text-muted-foreground">{emptyMessage}</p>
        ) : (
          lines.map((l) => (
            <div
              key={l.key}
              className={cn(
                "break-all whitespace-pre-wrap",
                l.system && "text-muted-foreground italic",
              )}
            >
              {l.text}
            </div>
          ))
        )}
      </div>
    </div>
  )
}
