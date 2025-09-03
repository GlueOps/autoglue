import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import {
  LinkIcon,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  ServerIcon,
  Tag,
  Trash,
  UnlinkIcon,
} from "lucide-react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { api, ApiError } from "@/lib/api.ts"
import { Badge } from "@/components/ui/badge.tsx"
import { Button } from "@/components/ui/button.tsx"
import { Checkbox } from "@/components/ui/checkbox.tsx"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx"

/* ----------------------------- Types & Schemas ---------------------------- */

type ServerBrief = {
  id: string
  hostname?: string | null
  ip?: string
  ip_address?: string
  role?: string
  status?: string
}

type LabelBrief = {
  id: string
  key: string
  value: string
}

type LabelWithPools = LabelBrief & {
  node_groups?: { id: string; name: string }[]
}

type NodePool = {
  id: string
  name: string
  servers?: ServerBrief[]
}

const CreatePoolSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
  server_ids: z.array(z.string().uuid()).optional().default([]),
})
type CreatePoolInput = z.input<typeof CreatePoolSchema>
type CreatePoolValues = z.output<typeof CreatePoolSchema>

const UpdatePoolSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
})
type UpdatePoolValues = z.output<typeof UpdatePoolSchema>

const AttachServersSchema = z.object({
  server_ids: z.array(z.string().uuid()).min(1, "Pick at least one server"),
})
type AttachServersValues = z.output<typeof AttachServersSchema>

const AttachLabelsSchema = z.object({
  label_ids: z.array(z.string().uuid()).min(1, "Pick at least one label"),
})
type AttachLabelsValues = z.output<typeof AttachLabelsSchema>

/* --------------------------------- Utils --------------------------------- */

function StatusBadge({ status }: { status?: string }) {
  const v =
    status === "ready"
      ? "default"
      : status === "provisioning"
        ? "secondary"
        : status === "failed"
          ? "destructive"
          : "outline"
  return (
    <Badge variant={v as any} className="capitalize">
      {status || "unknown"}
    </Badge>
  )
}

function truncateMiddle(str: string, keep = 12) {
  if (!str || str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

function serverLabel(s: ServerBrief) {
  const ip = s.ip || s.ip_address
  const name = s.hostname || ip || s.id
  const role = s.role ? ` · ${s.role}` : ""
  return `${name}${role}`
}

function labelKV(l: LabelBrief) {
  return `${l.key}=${l.value}`
}

/* --------------------------------- Page ---------------------------------- */

export const NodePoolPage = () => {
  const [loading, setLoading] = useState<boolean>(true)
  const [pools, setPools] = useState<NodePool[]>([])
  const [allServers, setAllServers] = useState<ServerBrief[]>([])

  // Pull labels with include=node_pools so we can map them to pools
  const [allLabels, setAllLabels] = useState<LabelWithPools[]>([])

  const [err, setErr] = useState<string | null>(null)
  const [q, setQ] = useState("")

  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<NodePool | null>(null)

  // Servers dialog state
  const [manageTarget, setManageTarget] = useState<NodePool | null>(null)

  // Labels dialog state
  const [manageLabelsTarget, setManageLabelsTarget] = useState<NodePool | null>(null)
  const [attachedLabels, setAttachedLabels] = useState<LabelBrief[]>([])
  const [labelsLoading, setLabelsLoading] = useState(false)
  const [labelsErr, setLabelsErr] = useState<string | null>(null)

  /* ------------------------------- Data Load ------------------------------ */

  async function loadAll() {
    setLoading(true)
    setErr(null)
    try {
      const [poolsData, serversData, labelsData] = await Promise.all([
        api.get<NodePool[]>("/api/v1/node-pools?include=servers"),
        api.get<ServerBrief[]>("/api/v1/servers"),
        api.get<LabelWithPools[]>("/api/v1/labels?include=node_pools"),
      ])
      setPools(poolsData || [])
      setAllServers(serversData || [])
      setAllLabels(labelsData || [])

      if (manageTarget) {
        const refreshed = (poolsData || []).find((p) => p.id === manageTarget.id) || null
        setManageTarget(refreshed)
      }
      if (editTarget) {
        const refreshed = (poolsData || []).find((p) => p.id === editTarget.id) || null
        setEditTarget(refreshed)
      }
      if (manageLabelsTarget) {
        await loadAttachedLabels(manageLabelsTarget.id)
      }
    } catch (e) {
      console.error(e)
      const msg = e instanceof ApiError ? e.message : "Failed to load node pools / servers / labels"
      setErr(msg)
    } finally {
      setLoading(false)
    }
  }

  async function loadAttachedLabels(poolId: string) {
    setLabelsLoading(true)
    setLabelsErr(null)
    try {
      const data = await api.get<LabelBrief[]>(`/api/v1/node-pools/${poolId}/labels`)
      setAttachedLabels(data || [])
    } catch (e) {
      console.error(e)
      const msg = e instanceof ApiError ? e.message : "Failed to load labels for pool"
      setLabelsErr(msg)
    } finally {
      setLabelsLoading(false)
    }
  }

  useEffect(() => {
    void loadAll()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  /* ---------------------------- Labels per Pool --------------------------- */

  // Build a quick lookup: poolId -> LabelBrief[]
  const labelsByPool = useMemo(() => {
    const map = new Map<string, LabelBrief[]>()
    for (const l of allLabels) {
      for (const ng of l.node_groups || []) {
        const arr = map.get(ng.id) || []
        arr.push({ id: l.id, key: l.key, value: l.value })
        map.set(ng.id, arr)
      }
    }
    return map
  }, [allLabels])

  /* -------------------------------- Filters ------------------------------- */

  const filtered = useMemo(() => {
    const needle = q.trim().toLowerCase()
    if (!needle) return pools
    return pools.filter(
      (p) =>
        p.name.toLowerCase().includes(needle) ||
        (p.servers || []).some(
          (s) =>
            (s.hostname || "").toLowerCase().includes(needle) ||
            (s.ip || s.ip_address || "").toLowerCase().includes(needle) ||
            (s.role || "").toLowerCase().includes(needle)
        ) ||
        (labelsByPool.get(p.id) || []).some(
          (l) =>
            l.key.toLowerCase().includes(needle) || (l.value || "").toLowerCase().includes(needle)
        )
    )
  }, [pools, q, labelsByPool])

  /* ------------------------------ Mutations ------------------------------- */

  async function deletePool(id: string) {
    if (!confirm("Delete this node pool? This cannot be undone.")) return
    await api.delete<void>(`/api/v1/node-pools/${id}`)
    await loadAll()
  }

  // Create Pool
  const createForm = useForm<CreatePoolInput, any, CreatePoolValues>({
    resolver: zodResolver(CreatePoolSchema),
    defaultValues: { name: "", server_ids: [] },
  })

  const submitCreate = async (values: CreatePoolValues) => {
    const payload: any = { name: values.name.trim() }
    if (values.server_ids && values.server_ids.length > 0) {
      payload.server_ids = values.server_ids
    }
    await api.post("/api/v1/node-pools", payload)
    setCreateOpen(false)
    createForm.reset({ name: "", server_ids: [] })
    await loadAll()
  }

  // Edit Pool
  const editForm = useForm<UpdatePoolValues>({
    resolver: zodResolver(UpdatePoolSchema),
    defaultValues: { name: "" },
  })

  function openEdit(p: NodePool) {
    setEditTarget(p)
    editForm.reset({ name: p.name })
  }

  const submitEdit = async (values: UpdatePoolValues) => {
    if (!editTarget) return
    await api.patch(`/api/v1/node-pools/${editTarget.id}`, { name: values.name.trim() })
    setEditTarget(null)
    await loadAll()
  }

  // Attach / Detach Servers
  const attachForm = useForm<AttachServersValues>({
    resolver: zodResolver(AttachServersSchema),
    defaultValues: { server_ids: [] },
  })

  function openManage(p: NodePool) {
    setManageTarget(p)
    attachForm.reset({ server_ids: [] })
  }

  const submitAttach = async (values: AttachServersValues) => {
    if (!manageTarget) return
    await api.post(`/api/v1/node-pools/${manageTarget.id}/servers`, {
      server_ids: values.server_ids,
    })
    attachForm.reset({ server_ids: [] })
    await loadAll()
  }

  async function detachServer(serverId: string) {
    if (!manageTarget) return
    if (!confirm("Detach this server from the pool?")) return
    await api.delete(`/api/v1/node-pools/${manageTarget.id}/servers/${serverId}`)
    await loadAll()
  }

  // Attach / Detach Labels
  const attachLabelsForm = useForm<AttachLabelsValues>({
    resolver: zodResolver(AttachLabelsSchema),
    defaultValues: { label_ids: [] },
  })

  function openManageLabels(p: NodePool) {
    setManageLabelsTarget(p)
    attachLabelsForm.reset({ label_ids: [] })
    void loadAttachedLabels(p.id)
  }

  const submitAttachLabels = async (values: AttachLabelsValues) => {
    if (!manageLabelsTarget) return
    await api.post(`/api/v1/node-pools/${manageLabelsTarget.id}/labels`, {
      label_ids: values.label_ids,
    })
    attachLabelsForm.reset({ label_ids: [] })
    await loadAttachedLabels(manageLabelsTarget.id)
    await loadAll() // refresh badges in table
  }

  async function detachLabel(labelId: string) {
    if (!manageLabelsTarget) return
    if (!confirm("Detach this label from the pool?")) return
    await api.delete(`/api/v1/node-pools/${manageLabelsTarget.id}/labels/${labelId}`)
    await loadAttachedLabels(manageLabelsTarget.id)
    await loadAll() // refresh badges in table
  }

  /* --------------------------------- Render -------------------------------- */

  if (loading) return <div className="p-6">Loading node pools…</div>
  if (err) return <div className="p-6 text-red-500">{err}</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Node Pools</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Search pools, servers, labels…"
              className="w-72 pl-8"
            />
          </div>

          <Button variant="outline" onClick={loadAll}>
            <RefreshCw className="mr-2 h-4 w-4" /> Refresh
          </Button>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" /> Create Pool
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create node pool</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form onSubmit={createForm.handleSubmit(submitCreate)} className="space-y-4">
                  <FormField
                    control={createForm.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Name</FormLabel>
                        <FormControl>
                          <Input placeholder="pool-workers-a" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="server_ids"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Initial servers (optional)</FormLabel>
                        <div className="max-h-56 space-y-2 overflow-auto rounded-xl border p-2">
                          {allServers.length === 0 && (
                            <div className="text-muted-foreground p-2 text-sm">
                              No servers available
                            </div>
                          )}
                          {allServers.map((s) => {
                            const checked = field.value?.includes(s.id) || false
                            return (
                              <label
                                key={s.id}
                                className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1"
                              >
                                <Checkbox
                                  checked={checked}
                                  onCheckedChange={(v) => {
                                    const next = new Set(field.value || [])
                                    if (v === true) next.add(s.id)
                                    else next.delete(s.id)
                                    field.onChange(Array.from(next))
                                  }}
                                />
                                <div className="leading-tight">
                                  <div className="text-sm font-medium">{serverLabel(s)}</div>
                                  <div className="text-muted-foreground text-xs">
                                    {truncateMiddle(s.id, 8)}
                                  </div>
                                </div>
                              </label>
                            )
                          })}
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createForm.formState.isSubmitting}>
                      {createForm.formState.isSubmitting ? "Creating…" : "Create"}
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <div className="bg-background overflow-hidden rounded-2xl border shadow-sm">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Servers</TableHead>
                <TableHead>Annotations</TableHead>
                <TableHead>Labels</TableHead>
                <TableHead>Taints</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((p) => {
                const labels = labelsByPool.get(p.id) || []
                return (
                  <TableRow key={p.id}>
                    <TableCell className="font-medium">{p.name}</TableCell>

                    {/* Servers cell */}
                    <TableCell>
                      <div className="flex flex-wrap gap-2">
                        {(p.servers || []).slice(0, 6).map((s) => (
                          <Badge key={s.id} variant="secondary" className="gap-1">
                            <ServerIcon className="h-3 w-3" />{" "}
                            {s.hostname || s.ip || s.ip_address || truncateMiddle(s.id, 6)}
                            <span className="ml-1">{s.role}</span>
                            {s.status && (
                              <span className="ml-1">
                                <StatusBadge status={s.status} />
                              </span>
                            )}
                          </Badge>
                        ))}
                        {(p.servers || []).length === 0 && (
                          <span className="text-muted-foreground">No servers</span>
                        )}
                        {(p.servers || []).length > 6 && (
                          <span className="text-muted-foreground">
                            +{(p.servers || []).length - 6} more
                          </span>
                        )}
                      </div>
                      <Button variant="outline" size="sm" onClick={() => openManage(p)}>
                        <LinkIcon className="mr-2 h-4 w-4" /> Manage servers
                      </Button>
                    </TableCell>

                    {/* Annotations placeholder */}
                    <TableCell>
                      <div className="flex flex-wrap gap-2">Annotations</div>
                      <Button variant="outline" size="sm" disabled>
                        <LinkIcon className="mr-2 h-4 w-4" /> Manage Annotations
                      </Button>
                    </TableCell>

                    {/* Labels cell with badges */}
                    <TableCell>
                      <div className="mb-2 flex flex-wrap gap-2">
                        {labels.slice(0, 6).map((l) => (
                          <Badge key={l.id} variant="outline" className="font-mono">
                            <Tag className="mr-1 h-3 w-3" />
                            {l.key}={l.value}
                          </Badge>
                        ))}
                        {labels.length === 0 && (
                          <span className="text-muted-foreground">No labels</span>
                        )}
                        {labels.length > 6 && (
                          <span className="text-muted-foreground">+{labels.length - 6} more</span>
                        )}
                      </div>
                      <Button variant="outline" size="sm" onClick={() => openManageLabels(p)}>
                        <Tag className="mr-2 h-4 w-4" /> Manage Labels
                      </Button>
                    </TableCell>

                    {/* Taints placeholder */}
                    <TableCell>
                      <div className="flex flex-wrap gap-2">Taints</div>
                      <Button variant="outline" size="sm" disabled>
                        <LinkIcon className="mr-2 h-4 w-4" /> Manage Taints
                      </Button>
                    </TableCell>

                    {/* Actions */}
                    <TableCell>
                      <div className="flex justify-end gap-2">
                        <Button variant="outline" size="sm" onClick={() => openEdit(p)}>
                          <Pencil className="mr-2 h-4 w-4" /> Edit
                        </Button>

                        <DropdownMenu>
                          <DropdownMenuTrigger asChild>
                            <Button variant="destructive" size="sm">
                              <Trash className="mr-2 h-4 w-4" /> Delete
                            </Button>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end">
                            <DropdownMenuItem onClick={() => deletePool(p.id)}>
                              Confirm delete
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={6} className="text-muted-foreground py-10 text-center">
                    No node pools match your search.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Edit dialog */}
      <Dialog open={!!editTarget} onOpenChange={(o) => !o && setEditTarget(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Edit node pool</DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form onSubmit={editForm.handleSubmit(submitEdit)} className="space-y-4">
              <FormField
                control={editForm.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="pool-workers-a" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button type="button" variant="outline" onClick={() => setEditTarget(null)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={editForm.formState.isSubmitting}>
                  {editForm.formState.isSubmitting ? "Saving…" : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Manage servers dialog */}
      <Dialog open={!!manageTarget} onOpenChange={(o) => !o && setManageTarget(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              Manage servers for <span className="font-mono">{manageTarget?.name}</span>
            </DialogTitle>
          </DialogHeader>

          {/* Attached servers list */}
          <div className="space-y-3">
            <div className="text-sm font-medium">Attached servers</div>
            <div className="overflow-hidden rounded-xl border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Server</TableHead>
                    <TableHead>IP</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="w-[120px] text-right">Detach</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(manageTarget?.servers || []).map((s) => (
                    <TableRow key={s.id}>
                      <TableCell className="font-medium">
                        {s.hostname || truncateMiddle(s.id, 8)}
                      </TableCell>
                      <TableCell>
                        <code className="font-mono text-sm">{s.ip || s.ip_address || "—"}</code>
                      </TableCell>
                      <TableCell className="capitalize">{s.role || "—"}</TableCell>
                      <TableCell>
                        <StatusBadge status={s.status} />
                      </TableCell>
                      <TableCell>
                        <div className="flex justify-end">
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => detachServer(s.id)}
                          >
                            <UnlinkIcon className="mr-2 h-4 w-4" /> Detach
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}

                  {(manageTarget?.servers || []).length === 0 && (
                    <TableRow>
                      <TableCell colSpan={5} className="text-muted-foreground py-8 text-center">
                        No servers attached yet.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          </div>

          {/* Attach servers */}
          <div className="pt-4">
            <Form {...attachForm}>
              <form onSubmit={attachForm.handleSubmit(submitAttach)} className="space-y-3">
                <FormField
                  control={attachForm.control}
                  name="server_ids"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Attach more servers</FormLabel>
                      <div className="grid max-h-64 grid-cols-1 gap-2 overflow-auto rounded-xl border p-2 md:grid-cols-2">
                        {/* options */}
                        {(() => {
                          const attachedIds = new Set(
                            (manageTarget?.servers || []).map((s) => s.id)
                          )
                          const attachableServers = allServers.filter((s) => !attachedIds.has(s.id))
                          if (attachableServers.length === 0) {
                            return (
                              <div className="text-muted-foreground p-2 text-sm">
                                No more servers available to attach
                              </div>
                            )
                          }
                          return attachableServers.map((s) => {
                            const checked = field.value?.includes(s.id) || false
                            return (
                              <label
                                key={s.id}
                                className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1"
                              >
                                <Checkbox
                                  checked={checked}
                                  onCheckedChange={(v) => {
                                    const next = new Set(field.value || [])
                                    if (v === true) next.add(s.id)
                                    else next.delete(s.id)
                                    field.onChange(Array.from(next))
                                  }}
                                />
                                <div className="leading-tight">
                                  <div className="text-sm font-medium">{serverLabel(s)}</div>
                                  <div className="text-muted-foreground text-xs">
                                    {truncateMiddle(s.id, 8)}
                                  </div>
                                </div>
                              </label>
                            )
                          })
                        })()}
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooter className="gap-2">
                  <Button type="submit" disabled={attachForm.formState.isSubmitting}>
                    <LinkIcon className="mr-2 h-4 w-4" />{" "}
                    {attachForm.formState.isSubmitting ? "Attaching…" : "Attach selected"}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </div>
        </DialogContent>
      </Dialog>

      {/* Manage labels dialog */}
      <Dialog open={!!manageLabelsTarget} onOpenChange={(o) => !o && setManageLabelsTarget(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              Manage labels for <span className="font-mono">{manageLabelsTarget?.name}</span>
            </DialogTitle>
          </DialogHeader>

          {/* Attached labels list */}
          <div className="space-y-3">
            <div className="text-sm font-medium">Attached labels</div>

            {labelsLoading ? (
              <div className="text-muted-foreground rounded-md border p-3 text-sm">Loading…</div>
            ) : labelsErr ? (
              <div className="rounded-md border p-3 text-sm text-red-500">{labelsErr}</div>
            ) : (
              <div className="overflow-hidden rounded-xl border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Key</TableHead>
                      <TableHead>Value</TableHead>
                      <TableHead className="w-[120px] text-right">Detach</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {attachedLabels.map((l) => (
                      <TableRow key={l.id}>
                        <TableCell className="font-mono text-sm">{l.key}</TableCell>
                        <TableCell className="font-mono text-sm">{l.value}</TableCell>
                        <TableCell>
                          <div className="flex justify-end">
                            <Button
                              variant="destructive"
                              size="sm"
                              onClick={() => detachLabel(l.id)}
                            >
                              <UnlinkIcon className="mr-2 h-4 w-4" /> Detach
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                    {attachedLabels.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={3} className="text-muted-foreground py-8 text-center">
                          No labels attached yet.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>
            )}
          </div>

          {/* Attach labels */}
          <div className="pt-4">
            <Form {...attachLabelsForm}>
              <form
                onSubmit={attachLabelsForm.handleSubmit(submitAttachLabels)}
                className="space-y-3"
              >
                <FormField
                  control={attachLabelsForm.control}
                  name="label_ids"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Attach more labels</FormLabel>
                      <div className="grid max-h-64 grid-cols-1 gap-2 overflow-auto rounded-xl border p-2 md:grid-cols-2">
                        {(() => {
                          const attachedIds = new Set(attachedLabels.map((l) => l.id))
                          const attachable = (allLabels as LabelBrief[]).filter(
                            (l) => !attachedIds.has(l.id)
                          )
                          if (attachable.length === 0) {
                            return (
                              <div className="text-muted-foreground p-2 text-sm">
                                No more labels available to attach
                              </div>
                            )
                          }
                          return attachable.map((l) => {
                            const checked = field.value?.includes(l.id) || false
                            return (
                              <label
                                key={l.id}
                                className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1"
                              >
                                <Checkbox
                                  checked={checked}
                                  onCheckedChange={(v) => {
                                    const next = new Set(field.value || [])
                                    if (v === true) next.add(l.id)
                                    else next.delete(l.id)
                                    field.onChange(Array.from(next))
                                  }}
                                />
                                <div className="leading-tight">
                                  <div className="text-sm font-medium">{labelKV(l)}</div>
                                  <div className="text-muted-foreground text-xs">
                                    {truncateMiddle(l.id, 8)}
                                  </div>
                                </div>
                              </label>
                            )
                          })
                        })()}
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooter className="gap-2">
                  <Button type="submit" disabled={attachLabelsForm.formState.isSubmitting}>
                    <LinkIcon className="mr-2 h-4 w-4" />{" "}
                    {attachLabelsForm.formState.isSubmitting ? "Attaching…" : "Attach selected"}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
