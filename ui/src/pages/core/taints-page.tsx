import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import {
  BadgeCheck,
  CircleSlash2,
  LinkIcon,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  Tags,
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx"

type NodePoolBrief = { id: string; name: string }
type Taint = {
  id: string
  key: string
  value?: string | null
  effect?: string | null
  node_groups?: NodePoolBrief[] // API uses "node_groups" for attached pools
}

const EFFECTS = ["NoSchedule", "PreferNoSchedule", "NoExecute"] as const

const CreateTaintSchema = z.object({
  key: z.string().trim().min(1, "Key is required").max(120, "Max 120 chars"),
  value: z.string().trim().optional(),
  effect: z.enum(EFFECTS),
  node_pool_ids: z.array(z.uuid()).optional().default([]),
})
type CreateTaintInput = z.input<typeof CreateTaintSchema>
type CreateTaintValues = z.output<typeof CreateTaintSchema>

const UpdateTaintSchema = z.object({
  key: z.string().trim().min(1, "Key is required").max(120).optional(),
  value: z.string().trim().optional(),
  effect: z.enum(EFFECTS as unknown as [string, ...string[]]).optional(),
})
type UpdateTaintValues = z.output<typeof UpdateTaintSchema>

const AttachPoolsSchema = z.object({
  node_pool_ids: z.array(z.string().uuid()).min(1, "Pick at least one node pool"),
})
type AttachPoolsValues = z.output<typeof AttachPoolsSchema>

function truncateMiddle(str?: string | null, keep = 12) {
  if (!str) return ""
  if (str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

function TaintBadge({ t }: { t: Pick<Taint, "key" | "value" | "effect"> }) {
  const label = `${t.key}${t.value ? `=${t.value}` : ""}${t.effect ? `:${t.effect}` : ""}`
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      <Tags className="mr-1 h-3 w-3" />
      {label}
    </Badge>
  )
}

export const TaintsPage = () => {
  const [loading, setLoading] = useState(true)
  const [err, setErr] = useState<string | null>(null)

  const [taints, setTaints] = useState<Taint[]>([])
  const [allPools, setAllPools] = useState<NodePoolBrief[]>([])

  const [q, setQ] = useState("")
  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Taint | null>(null)
  const [manageTarget, setManageTarget] = useState<Taint | null>(null)

  async function loadAll() {
    setLoading(true)
    setErr(null)
    try {
      // include attached node pools for quick display
      const [taintsData, poolsData] = await Promise.all([
        api.get<Taint[]>("/api/v1/taints?include=node_groups"),
        api.get<NodePoolBrief[]>("/api/v1/node-pools"),
      ])
      setTaints(taintsData || [])
      setAllPools(poolsData || [])

      if (manageTarget) {
        const refreshed = (taintsData || []).find((t) => t.id === manageTarget.id) || null
        setManageTarget(refreshed)
      }
      if (editTarget) {
        const refreshed = (taintsData || []).find((t) => t.id === editTarget.id) || null
        setEditTarget(refreshed)
      }
    } catch (e) {
      console.error(e)
      const msg = e instanceof ApiError ? e.message : "Failed to load taints or node pools"
      setErr(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadAll()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const filtered = useMemo(() => {
    const needle = q.trim().toLowerCase()
    if (!needle) return taints
    return taints.filter((t) => {
      const label =
        `${t.key}${t.value ? `=${t.value}` : ""}${t.effect ? `:${t.effect}` : ""}`.toLowerCase()
      const pools = (t.node_groups || []).some((p) => p.name.toLowerCase().includes(needle))
      return label.includes(needle) || pools
    })
  }, [taints, q])

  async function deleteTaint(id: string) {
    if (!confirm("Delete this taint? This cannot be undone.")) return
    await api.delete<void>(`/api/v1/taints/${id}`)
    await loadAll()
  }

  // -------- Create --------
  const createForm = useForm<CreateTaintInput, any, CreateTaintValues>({
    resolver: zodResolver(CreateTaintSchema),
    defaultValues: { key: "", value: "", effect: undefined, node_pool_ids: [] },
  })

  const submitCreate = async (values: CreateTaintValues) => {
    const payload: any = {
      key: values.key.trim(),
      effect: values.effect,
    }
    if (values.value) payload.value = values.value.trim()
    if (values.node_pool_ids && values.node_pool_ids.length > 0) {
      payload.node_pool_ids = values.node_pool_ids
    }
    await api.post("/api/v1/taints", payload)
    setCreateOpen(false)
    createForm.reset({ key: "", value: "", effect: undefined, node_pool_ids: [] })
    await loadAll()
  }

  // -------- Edit --------
  const editForm = useForm<UpdateTaintValues>({
    resolver: zodResolver(UpdateTaintSchema),
    defaultValues: {},
  })

  function openEdit(t: Taint) {
    setEditTarget(t)
    editForm.reset({ key: t.key, value: t.value || "", effect: (t.effect as any) || undefined })
  }

  const submitEdit = async (values: UpdateTaintValues) => {
    if (!editTarget) return
    const body: Record<string, any> = {}
    if (values.key !== undefined) body.key = values.key.trim()
    if (values.value !== undefined) body.value = values.value?.trim() ?? ""
    if (values.effect !== undefined) body.effect = values.effect
    await api.patch(`/api/v1/taints/${editTarget.id}`, body)
    setEditTarget(null)
    await loadAll()
  }

  // -------- Manage attached pools --------
  const attachForm = useForm<AttachPoolsValues>({
    resolver: zodResolver(AttachPoolsSchema),
    defaultValues: { node_pool_ids: [] },
  })

  function openManage(t: Taint) {
    setManageTarget(t)
    attachForm.reset({ node_pool_ids: [] })
  }

  const submitAttach = async (values: AttachPoolsValues) => {
    if (!manageTarget) return
    await api.post(`/api/v1/taints/${manageTarget.id}/node_pools`, {
      node_pool_ids: values.node_pool_ids,
    })
    attachForm.reset({ node_pool_ids: [] })
    await loadAll()
  }

  async function detachPool(poolId: string) {
    if (!manageTarget) return
    if (!confirm("Detach this taint from the node pool?")) return
    await api.delete(`/api/v1/taints/${manageTarget.id}/node_pools/${poolId}`)
    await loadAll()
  }

  const attachablePools = useMemo(() => {
    if (!manageTarget) return [] as NodePoolBrief[]
    const attachedIds = new Set((manageTarget.node_groups || []).map((p) => p.id))
    return allPools.filter((p) => !attachedIds.has(p.id))
  }, [manageTarget, allPools])

  // -------- UI --------
  if (loading) return <div className="p-6">Loading taints…</div>
  if (err) return <div className="p-6 text-red-500">{err}</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Taints</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Search taints or attached pools…"
              className="w-72 pl-8"
            />
          </div>

          <Button variant="outline" onClick={loadAll}>
            <RefreshCw className="mr-2 h-4 w-4" /> Refresh
          </Button>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" /> Create Taint
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create taint</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form onSubmit={createForm.handleSubmit(submitCreate)} className="space-y-4">
                  <FormField
                    control={createForm.control}
                    name="key"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Key</FormLabel>
                        <FormControl>
                          <Input placeholder="dedicated" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="value"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Value (optional)</FormLabel>
                        <FormControl>
                          <Input placeholder="gpu" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="effect"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Effect</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select effect" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {EFFECTS.map((e) => (
                              <SelectItem key={e} value={e}>
                                {e}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="node_pool_ids"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Attach to node pools (optional)</FormLabel>
                        <div className="max-h-56 space-y-2 overflow-auto rounded-xl border p-2">
                          {allPools.length === 0 && (
                            <div className="text-muted-foreground p-2 text-sm">
                              No node pools available
                            </div>
                          )}
                          {allPools.map((p) => {
                            const checked = field.value?.includes(p.id) || false
                            return (
                              <label
                                key={p.id}
                                className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1"
                              >
                                <Checkbox
                                  checked={checked}
                                  onCheckedChange={(v) => {
                                    const next = new Set(field.value || [])
                                    if (v === true) next.add(p.id)
                                    else next.delete(p.id)
                                    field.onChange(Array.from(next))
                                  }}
                                />
                                <div className="leading-tight">
                                  <div className="text-sm font-medium">{p.name}</div>
                                  <div className="text-muted-foreground text-xs">
                                    {truncateMiddle(p.id, 8)}
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
                <TableHead>Taint</TableHead>
                <TableHead>Attached Node Pools</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((t) => (
                <TableRow key={t.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      <TaintBadge t={t} />
                      <code className="text-muted-foreground text-xs">
                        {truncateMiddle(t.id, 6)}
                      </code>
                    </div>
                  </TableCell>

                  <TableCell>
                    <div className="mb-2 flex flex-wrap gap-2">
                      {(t.node_groups || []).slice(0, 6).map((p) => (
                        <Badge key={p.id} variant="outline" className="gap-1">
                          <BadgeCheck className="h-3 w-3" />
                          {p.name}
                        </Badge>
                      ))}

                      {(t.node_groups || []).length === 0 && (
                        <span className="text-muted-foreground">No node pools</span>
                      )}
                      {(t.node_groups || []).length > 6 && (
                        <span className="text-muted-foreground">
                          +{(t.node_groups || []).length - 6} more
                        </span>
                      )}
                    </div>

                    <Button variant="outline" size="sm" onClick={() => openManage(t)}>
                      <LinkIcon className="mr-2 h-4 w-4" /> Manage node pools
                    </Button>
                  </TableCell>

                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => openEdit(t)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit
                      </Button>

                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="destructive" size="sm">
                            <Trash className="mr-2 h-4 w-4" /> Delete
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => deleteTaint(t.id)}>
                            Confirm delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={3} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No taints match your search.
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
            <DialogTitle>Edit taint</DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form onSubmit={editForm.handleSubmit(submitEdit)} className="space-y-4">
              <FormField
                control={editForm.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Key</FormLabel>
                    <FormControl>
                      <Input placeholder="dedicated" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="value"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Value (optional)</FormLabel>
                    <FormControl>
                      <Input placeholder="gpu" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="effect"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Effect</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select effect" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {EFFECTS.map((e) => (
                          <SelectItem key={e} value={e}>
                            {e}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
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

      {/* Manage node pools dialog */}
      <Dialog open={!!manageTarget} onOpenChange={(o) => !o && setManageTarget(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              Manage pools for{" "}
              <span className="font-mono">
                {manageTarget
                  ? `${manageTarget.key}${manageTarget.value ? `=${manageTarget.value}` : ""}${manageTarget.effect ? `:${manageTarget.effect}` : ""}`
                  : ""}
              </span>
            </DialogTitle>
          </DialogHeader>

          {/* Attached pools */}
          <div className="space-y-3">
            <div className="text-sm font-medium">Attached node pools</div>
            <div className="overflow-hidden rounded-xl border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead className="w-[120px] text-right">Detach</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(manageTarget?.node_groups || []).map((p) => (
                    <TableRow key={p.id}>
                      <TableCell className="font-medium">{p.name}</TableCell>
                      <TableCell>
                        <div className="flex justify-end">
                          <Button variant="destructive" size="sm" onClick={() => detachPool(p.id)}>
                            <UnlinkIcon className="mr-2 h-4 w-4" /> Detach
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}

                  {(manageTarget?.node_groups || []).length === 0 && (
                    <TableRow>
                      <TableCell colSpan={2} className="text-muted-foreground py-8 text-center">
                        No node pools attached yet.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          </div>

          {/* Attach section */}
          <div className="pt-4">
            <Form {...attachForm}>
              <form onSubmit={attachForm.handleSubmit(submitAttach)} className="space-y-3">
                <FormField
                  control={attachForm.control}
                  name="node_pool_ids"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Attach more node pools</FormLabel>
                      <div className="grid max-h-64 grid-cols-1 gap-2 overflow-auto rounded-xl border p-2 md:grid-cols-2">
                        {attachablePools.length === 0 && (
                          <div className="text-muted-foreground p-2 text-sm">
                            No more node pools available to attach
                          </div>
                        )}
                        {attachablePools.map((p) => {
                          const checked = field.value?.includes(p.id) || false
                          return (
                            <label
                              key={p.id}
                              className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1"
                            >
                              <Checkbox
                                checked={checked}
                                onCheckedChange={(v) => {
                                  const next = new Set(field.value || [])
                                  if (v === true) next.add(p.id)
                                  else next.delete(p.id)
                                  field.onChange(Array.from(next))
                                }}
                              />
                              <div className="leading-tight">
                                <div className="text-sm font-medium">{p.name}</div>
                                <div className="text-muted-foreground text-xs">
                                  {truncateMiddle(p.id, 8)}
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
    </div>
  )
}
