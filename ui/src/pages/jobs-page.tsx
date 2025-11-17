import { useEffect, useState, type FC } from "react"
import { archerAdminApi } from "@/api/archer_admin"
import type { AdminListArcherJobsRequest } from "@/sdk"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Loader2, Plus, RefreshCw, Search, X } from "lucide-react"

import { cn } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Textarea } from "@/components/ui/textarea"

// Types (align with generated client camelCase)
type JobStatus =
  | "queued"
  | "running"
  | "succeeded"
  | "failed"
  | "canceled"
  | "retrying"
  | "scheduled"
type DtoJob = {
  id: string
  type: string
  queue: string
  status: JobStatus
  attempts: number
  maxAttempts?: number
  createdAt: string
  updatedAt?: string
  lastError?: string | null
  runAt?: string | null
  payload?: unknown
}
type DtoPageJob = {
  items: DtoJob[]
  total: number
  page: number
  pageSize: number
}
type QueueInfo = {
  name: string
  pending: number
  running: number
  failed: number
  scheduled: number
}

const STATUS: JobStatus[] = [
  "queued",
  "running",
  "succeeded",
  "failed",
  "canceled",
  "retrying",
  "scheduled",
]

const statusClass: Record<JobStatus, string> = {
  queued: "bg-amber-100 text-amber-800",
  running: "bg-sky-100 text-sky-800",
  succeeded: "bg-emerald-100 text-emerald-800",
  failed: "bg-red-100 text-red-800",
  canceled: "bg-zinc-200 text-zinc-700",
  retrying: "bg-orange-100 text-orange-800",
  scheduled: "bg-violet-100 text-violet-800",
}

function fmt(dt?: string | null) {
  if (!dt) return "—"
  const d = new Date(dt)
  return new Intl.DateTimeFormat(undefined, { dateStyle: "medium", timeStyle: "short" }).format(d)
}

// Small debounce hook for search input
function useDebounced<T>(value: T, ms = 300) {
  const [v, setV] = useState(value)
  useEffect(() => {
    const t = setTimeout(() => setV(value), ms)
    return () => clearTimeout(t)
  }, [value, ms])
  return v
}

export const JobsPage: FC = () => {
  const qc = useQueryClient()

  // Filters
  const [status, setStatus] = useState<string>("")
  const [queue, setQueue] = useState<string>("")
  const [q, setQ] = useState<string>("")
  const debouncedQ = useDebounced(q, 300)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(25)

  const key = ["archer", "jobs", { status, queue, q: debouncedQ, page, pageSize }]

  // Jobs query
  const jobsQ = useQuery({
    queryKey: key,
    queryFn: () =>
      archerAdminApi.listJobs({
        status: status,
        queue: queue || undefined,
        q: debouncedQ || undefined,
        page,
        pageSize,
      } as AdminListArcherJobsRequest) as Promise<DtoPageJob>,
    placeholderData: (prev) => prev,
    staleTime: 10_000,
  })

  // Queues summary (optional header)
  const queuesQ = useQuery({
    queryKey: ["archer", "queues"],
    queryFn: () => archerAdminApi.listQueues() as Promise<QueueInfo[]>,
    staleTime: 30_000,
  })

  // Mutations
  const enqueueM = useMutation({
    mutationFn: (body: {
      queue: string
      type: string
      payload?: object | undefined
      run_at?: string
    }) => archerAdminApi.enqueue(body),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["archer", "jobs"] }),
  })
  const retryM = useMutation({
    mutationFn: (id: string) => archerAdminApi.retryJob(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["archer", "jobs"] }),
  })
  const cancelM = useMutation({
    mutationFn: (id: string) => archerAdminApi.cancelJob(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["archer", "jobs"] }),
  })

  const busy = jobsQ.isFetching

  const data = jobsQ.data as DtoPageJob
  const totalPages = data ? Math.max(1, Math.ceil(data.total / data.pageSize)) : 1

  return (
    <div className="container mx-auto space-y-6 p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold">Archer Jobs</h1>
          <p className="text-muted-foreground text-sm">
            Inspect, enqueue, retry and cancel background jobs.
          </p>
        </div>
        <div className="flex gap-2">
          <EnqueueDialog
            onSubmit={(payload) => enqueueM.mutateAsync(payload)}
            submitting={enqueueM.isPending}
          />
          <Button
            variant="secondary"
            onClick={() => qc.invalidateQueries({ queryKey: ["archer", "jobs"] })}
            disabled={busy}
          >
            {busy ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <RefreshCw className="mr-2 h-4 w-4" />
            )}
            Refresh
          </Button>
        </div>
      </div>

      {/* Queue metrics (optional) */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {queuesQ.data?.map((q) => (
          <Card key={q.name}>
            <CardHeader>
              <CardTitle className="text-base">{q.name}</CardTitle>
            </CardHeader>
            <CardContent className="grid grid-cols-2 gap-2 text-sm">
              <Metric label="Pending" value={q.pending} />
              <Metric label="Running" value={q.running} />
              <Metric label="Failed" value={q.failed} />
              <Metric label="Scheduled" value={q.scheduled} />
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle>Filters</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-3 md:grid-cols-4">
          <div className="col-span-2 flex items-center gap-2">
            <Input
              placeholder="Search id, queue, error, payload…"
              value={q}
              onChange={(e) => {
                setQ(e.target.value)
                setPage(1)
              }}
              onKeyDown={(e) =>
                e.key === "Enter" && qc.invalidateQueries({ queryKey: ["archer", "jobs"] })
              }
            />
            {q && (
              <Button variant="ghost" size="icon" onClick={() => setQ("")}>
                <X className="h-4 w-4" />
              </Button>
            )}
            <Button onClick={() => qc.invalidateQueries({ queryKey: ["archer", "jobs"] })}>
              <Search className="mr-2 h-4 w-4" /> Search
            </Button>
          </div>
          <Select
            value={status || "all"}
            onValueChange={(v) => {
              setStatus(v === "all" ? "" : v)
              setPage(1)
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="All statuses" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All statuses</SelectItem>
              {STATUS.map((s) => (
                <SelectItem key={s} value={s}>
                  {s}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Input
            placeholder="Queue (optional)"
            value={queue}
            onChange={(e) => {
              setQueue(e.target.value)
              setPage(1)
            }}
          />
          <div className="flex items-center gap-2">
            <Label className="whitespace-nowrap">Page size</Label>
            <Select
              value={String(pageSize)}
              onValueChange={(v) => {
                setPageSize(Number(v))
                setPage(1)
              }}
            >
              <SelectTrigger className="w-[120px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {[10, 25, 50, 100].map((n) => (
                  <SelectItem key={n} value={String(n)}>
                    {n}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Table */}
      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Queue</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Attempts</TableHead>
                <TableHead>Run At</TableHead>
                <TableHead>Updated</TableHead>
                <TableHead className="pr-4 text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {jobsQ.isLoading && (
                <TableRow>
                  <TableCell colSpan={7} className="text-muted-foreground py-8 text-center">
                    Loading…
                  </TableCell>
                </TableRow>
              )}
              {jobsQ.isError && (
                <TableRow>
                  <TableCell colSpan={7} className="py-8 text-center text-red-600">
                    Failed to load jobs
                  </TableCell>
                </TableRow>
              )}
              {!jobsQ.isLoading && data && data.items.length === 0 && (
                <TableRow>
                  <TableCell colSpan={7} className="text-muted-foreground py-8 text-center">
                    No jobs match your filters.
                  </TableCell>
                </TableRow>
              )}
              {data?.items.map((j) => (
                <TableRow key={j.id}>
                  <TableCell>
                    <code className="text-xs">{j.id}</code>
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">{j.queue}</Badge>
                  </TableCell>
                  <TableCell>
                    <span className={cn("rounded-md px-2 py-0.5 text-xs", statusClass[j.status])}>
                      {j.status}
                    </span>
                  </TableCell>
                  <TableCell>
                    {j.maxAttempts ? `${j.attempts}/${j.maxAttempts}` : j.attempts}
                  </TableCell>
                  <TableCell>{fmt(j.runAt)}</TableCell>
                  <TableCell>{fmt(j.updatedAt ?? j.createdAt)}</TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      {(j.status === "failed" || j.status === "canceled") && (
                        <Button
                          size="sm"
                          variant="outline"
                          disabled={retryM.isPending}
                          onClick={() => retryM.mutate(j.id)}
                        >
                          Retry
                        </Button>
                      )}
                      {(j.status === "queued" ||
                        j.status === "running" ||
                        j.status === "scheduled") && (
                        <Button
                          size="sm"
                          variant="outline"
                          disabled={cancelM.isPending}
                          onClick={() => cancelM.mutate(j.id)}
                        >
                          Cancel
                        </Button>
                      )}
                      <DetailsButton job={j} />
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          {/* Pagination */}
          <div className="flex items-center justify-between border-t p-3 text-sm">
            <div>
              Page {page} of {totalPages} • {data?.total ?? 0} total
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                disabled={page <= 1 || jobsQ.isFetching}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                Prev
              </Button>
              <Button
                variant="outline"
                disabled={page >= totalPages || jobsQ.isFetching}
                onClick={() => setPage((p) => p + 1)}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div className="bg-muted/30 rounded-lg border p-3">
      <div className="text-muted-foreground text-xs">{label}</div>
      <div className="text-lg font-semibold">{value}</div>
    </div>
  )
}

function DetailsButton({ job }: { job: DtoJob }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm" variant="ghost">
          Details
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Job {job.id}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-3">
          {job.lastError && (
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Last error</CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="overflow-auto text-xs whitespace-pre-wrap">{job.lastError}</pre>
              </CardContent>
            </Card>
          )}
          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Payload</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="overflow-auto text-xs whitespace-pre-wrap">
                {JSON.stringify(job.payload, null, 2)}
              </pre>
            </CardContent>
          </Card>
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="secondary">Close</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

function EnqueueDialog({
  onSubmit,
  submitting,
}: {
  onSubmit: (body: {
    queue: string
    type: string
    payload?: object | undefined
    run_at?: string
  }) => Promise<unknown>
  submitting?: boolean
}) {
  const [open, setOpen] = useState(false)
  const [queue, setQueue] = useState("")
  const [type, setType] = useState("")
  const [payload, setPayload] = useState("{}")
  const [runAt, setRunAt] = useState("")

  const canSubmit = queue && type && !submitting

  async function handleSubmit() {
    const parsed = payload ? JSON.parse(payload) : undefined
    await onSubmit({ queue, type, payload: parsed, run_at: runAt || undefined })
    setOpen(false)
    setQueue("")
    setType("")
    setPayload("{}")
    setRunAt("")
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="mr-2 h-4 w-4" /> Enqueue
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Enqueue Job</DialogTitle>
        </DialogHeader>
        <div className="grid gap-3">
          <div className="grid gap-2">
            <Label>Queue</Label>
            <Input
              value={queue}
              onChange={(e) => setQueue(e.target.value)}
              placeholder="e.g. bootstrap_bastion"
            />
          </div>
          <div className="grid gap-2">
            <Label>Type</Label>
            <Input
              value={type}
              onChange={(e) => setType(e.target.value)}
              placeholder="e.g. bootstrap_bastion"
            />
          </div>
          <div className="grid gap-2">
            <Label>Payload (JSON)</Label>
            <Textarea
              value={payload}
              onChange={(e) => setPayload(e.target.value)}
              className="min-h-[120px] font-mono text-xs"
            />
          </div>
          <div className="grid gap-2">
            <Label>Run at (optional)</Label>
            <Input type="datetime-local" value={runAt} onChange={(e) => setRunAt(e.target.value)} />
          </div>
        </div>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="secondary">Cancel</Button>
          </DialogClose>
          <Button onClick={handleSubmit} disabled={!canSubmit}>
            {submitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
            Enqueue
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
