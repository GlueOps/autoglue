import { useEffect, useMemo, useState } from "react"
import { annotationsApi } from "@/api/annotations"
import { labelsApi } from "@/api/labels"
import { canAttachToPool, nodePoolsApi } from "@/api/node_pools"
import { serversApi } from "@/api/servers"
import { taintsApi } from "@/api/taints"
import type { DtoNodePoolResponse } from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Ellipsis, LinkIcon, Pencil, Plus, Search, ServerIcon, Trash2 } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
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

// ---------------------------------------------
// Helpers & shared UI
// ---------------------------------------------

const ROLE_OPTIONS = ["master", "worker"] as const
type Role = (typeof ROLE_OPTIONS)[number]

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

// ---------------------------------------------
// Reusable manage-many dialog
// ---------------------------------------------

type AssocItem = {
  id: string
  // common searchable/display fields
  name?: string
  key?: string
  value?: string
  effect?: string
  role?: string
  status?: string
  hostname?: string
  private_ip_address?: string
  public_ip_address?: string
}

function fuzzyIncludes(hay: string | undefined, q: string) {
  return (hay ?? "").toLowerCase().includes(q)
}

function ManageManyDialog(props: {
  open: boolean
  title: string
  onOpenChange: (v: boolean) => void
  items: AssocItem[]
  initialSelectedIds: Set<string>
  onSave: (diff: { toAttach: string[]; toDetach: string[] }) => Promise<void> | void
  columns: { header: string; render: (item: AssocItem) => React.ReactNode }[]
  allowItem?: (item: AssocItem) => boolean
}) {
  const { open, title, onOpenChange, items, initialSelectedIds, onSave, columns, allowItem } = props
  const [q, setQ] = useState("")
  const [selected, setSelected] = useState<Set<string>>(new Set(initialSelectedIds))
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    setSelected(new Set(initialSelectedIds))
    setQ("")
  }, [initialSelectedIds, open])

  const filtered = useMemo(() => {
    const qq = q.trim().toLowerCase()
    return items.filter((it) => {
      if (allowItem && !allowItem(it)) return false
      if (!qq) return true
      return (
        fuzzyIncludes(it.name, qq) ||
        fuzzyIncludes(it.key, qq) ||
        fuzzyIncludes(it.value, qq) ||
        fuzzyIncludes(it.effect, qq) ||
        fuzzyIncludes(it.hostname, qq) ||
        fuzzyIncludes(it.private_ip_address, qq) ||
        fuzzyIncludes(it.public_ip_address, qq) ||
        fuzzyIncludes(it.role, qq) ||
        fuzzyIncludes(it.status, qq)
      )
    })
  }, [items, q, allowItem])

  const initial = initialSelectedIds
  const changed =
    Array.from(selected).some((id) => !initial.has(id)) ||
    Array.from(initial).some((id) => !selected.has(id))

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-3xl">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>

        <div className="space-y-3">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="Search…"
              className="pl-8"
            />
          </div>

          <div className="max-h-[50vh] overflow-auto rounded border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[48px]"></TableHead>
                  {columns.map((c, i) => (
                    <TableHead key={i}>{c.header}</TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((it) => {
                  const id = it.id
                  const checked = selected.has(id)
                  return (
                    <TableRow key={id}>
                      <TableCell className="text-center align-middle">
                        <input
                          type="checkbox"
                          className="h-4 w-4"
                          checked={checked}
                          onChange={(e) => {
                            const next = new Set(selected)
                            if (e.target.checked) next.add(id)
                            else next.delete(id)
                            setSelected(next)
                          }}
                        />
                      </TableCell>
                      {columns.map((c, i) => (
                        <TableCell key={i}>{c.render(it)}</TableCell>
                      ))}
                    </TableRow>
                  )
                })}
                {filtered.length === 0 && (
                  <TableRow>
                    <TableCell
                      colSpan={1 + columns.length}
                      className="text-muted-foreground py-8 text-center"
                    >
                      No items found.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          <div className="text-muted-foreground text-sm">
            Selected: <span className="text-foreground font-medium">{selected.size}</span>
          </div>
        </div>

        <DialogFooter className="gap-2">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={saving}
          >
            Close
          </Button>
          <Button
            onClick={async () => {
              const toAttach: string[] = []
              const toDetach: string[] = []
              for (const id of selected) if (!initial.has(id)) toAttach.push(id)
              for (const id of initial) if (!selected.has(id)) toDetach.push(id)
              try {
                setSaving(true)
                await onSave({ toAttach, toDetach })
                onOpenChange(false)
              } finally {
                setSaving(false)
              }
            }}
            disabled={saving || !changed}
          >
            {saving ? "Saving…" : "Save changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

// ---------------------------------------------
// Page
// ---------------------------------------------

const createNodePoolSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
  role: z.enum(ROLE_OPTIONS),
})
type CreateNodePoolValues = z.infer<typeof createNodePoolSchema>

const updateNodePoolSchema = createNodePoolSchema.partial()
type UpdateNodePoolValues = z.infer<typeof updateNodePoolSchema>

export function NodePoolsPage() {
  const [filter, setFilter] = useState("")
  const [createOpen, setCreateOpen] = useState(false)
  const [updateOpen, setUpdateOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [editing, setEditing] = useState<DtoNodePoolResponse | null>(null)
  const [deleting, setDeleting] = useState<DtoNodePoolResponse | null>(null)

  // Manage Servers dialog state
  const [manageOpen, setManageOpen] = useState(false)
  const [managePool, setManagePool] = useState<DtoNodePoolResponse | null>(null)
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [initialSelected, setInitialSelected] = useState<Set<string>>(new Set())
  const [serverFilter, setServerFilter] = useState("")

  // Manage Labels / Annotations / Taints dialog state
  const [manageLabelsOpen, setManageLabelsOpen] = useState(false)
  const [manageAnnotationsOpen, setManageAnnotationsOpen] = useState(false)
  const [manageTaintsOpen, setManageTaintsOpen] = useState(false)
  const [manageLATPool, setManageLATPool] = useState<DtoNodePoolResponse | null>(null)

  const [labelsInitial, setLabelsInitial] = useState<Set<string>>(new Set())
  const [annotationsInitial, setAnnotationsInitial] = useState<Set<string>>(new Set())
  const [taintsInitial, setTaintsInitial] = useState<Set<string>>(new Set())

  const qc = useQueryClient()

  // Queries
  const nodePoolQ = useQuery({
    queryKey: ["node-pools"],
    queryFn: () => nodePoolsApi.listNodePools(),
  })

  const serverQ = useQuery({
    queryKey: ["servers"],
    queryFn: () => serversApi.listServers(),
  })

  const annotationQ = useQuery({
    queryKey: ["annotations"],
    queryFn: () => annotationsApi.listAnnotations(),
  })

  const labelQ = useQuery({
    queryKey: ["labels"],
    queryFn: () => labelsApi.listLabels(),
  })

  const taintQ = useQuery({
    queryKey: ["taints"],
    queryFn: () => taintsApi.listTaints(),
  })

  // --- Create
  const createForm = useForm<CreateNodePoolValues>({
    resolver: zodResolver(createNodePoolSchema),
    defaultValues: { name: "", role: "worker" },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateNodePoolValues) => nodePoolsApi.createNodePool(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["node-pools"] })
      createForm.reset({ name: "", role: "worker" })
      setCreateOpen(false)
      toast.success("Node pool created.")
    },
    onError: (err: any) => toast.error(err?.message ?? "Unable to create node pool."),
  })

  const onCreateSubmit = (values: CreateNodePoolValues) => createMut.mutate(values)

  // --- Update
  const updateForm = useForm<UpdateNodePoolValues>({
    resolver: zodResolver(updateNodePoolSchema),
    defaultValues: { name: undefined, role: undefined },
  })

  useEffect(() => {
    if (editing) {
      updateForm.reset({ name: editing.name, role: editing.role as Role })
    } else {
      updateForm.reset({ name: undefined, role: undefined })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editing])

  const updateMut = useMutation({
    mutationFn: async (values: UpdateNodePoolValues) => {
      if (!editing) return
      const body: UpdateNodePoolValues = {}
      if (values.name !== editing.name) body.name = values.name
      if ((values.role as string) !== editing.role) body.role = values.role
      return await nodePoolsApi.updateNodePool(editing.id as unknown as string, body)
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["node-pools"] })
      setUpdateOpen(false)
      setEditing(null)
      toast.success("Node pool updated.")
    },
    onError: (err: any) => toast.error(err?.message ?? "Unable to update node pool."),
  })

  const onUpdateSubmit = (values: UpdateNodePoolValues) => updateMut.mutate(values)

  // --- Delete
  const deleteMut = useMutation({
    mutationFn: async () => {
      if (!deleting) return
      await nodePoolsApi.deleteNodePool(deleting.id as unknown as string)
    },
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["node-pools"] })
      setDeleteOpen(false)
      setDeleting(null)
      toast.success("Node pool deleted.")
    },
    onError: (err: any) => toast.error(err?.message ?? "Unable to delete node pool."),
  })

  // --- Filter table
  const filtered = useMemo(() => {
    const data = (nodePoolQ.data ?? []) as DtoNodePoolResponse[]
    const q = filter.trim().toLowerCase()
    return q
      ? data.filter(
          (p) => p.name?.toLowerCase().includes(q) || (p.role as string)?.toLowerCase().includes(q)
        )
      : data
  }, [filter, nodePoolQ.data])

  if (nodePoolQ.isLoading) return <div className="p-6">Loading node pools…</div>
  if (nodePoolQ.error)
    return (
      <div className="p-6 text-red-500">
        Error loading node pools.
        <pre className="bg-muted mt-3 rounded p-3 text-xs">
          {JSON.stringify(nodePoolQ.error, null, 2)}
        </pre>
      </div>
    )

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Node Pools</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search node pools"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Node Pool
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create Node Pool</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form className="space-y-4" onSubmit={createForm.handleSubmit(onCreateSubmit)}>
                  <FormField
                    control={createForm.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Name</FormLabel>
                        <FormControl>
                          <Input placeholder="master-pool" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="role"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Role</FormLabel>
                        <Select
                          onValueChange={(v) =>
                            createForm.setValue("role", v as Role, {
                              shouldDirty: true,
                              shouldValidate: true,
                            })
                          }
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select role" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectItem value="master">master</SelectItem>
                            <SelectItem value="worker">worker</SelectItem>
                          </SelectContent>
                        </Select>
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
                <TableHead>Role</TableHead>
                <TableHead>Servers</TableHead>
                <TableHead>Annotations</TableHead>
                <TableHead>Labels</TableHead>
                <TableHead>Taints</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((p) => {
                const serverCount = Array.isArray(p.servers) ? p.servers.length : 0
                return (
                  <TableRow key={p.id}>
                    <TableCell className="font-medium">{p.name}</TableCell>
                    <TableCell className="font-medium">{p.role}</TableCell>

                    {/* Servers cell */}
                    <TableCell>
                      <div className="flex flex-wrap items-center gap-2">
                        {(p.servers || []).slice(0, 6).map((s) => (
                          <Badge key={s.id} variant="secondary" className="gap-1">
                            <ServerIcon className="h-3 w-3" />
                            {s.hostname || s.private_ip_address}
                            <span className="ml-1">{s.role}</span>
                            {s.status && (
                              <span className="ml-1">
                                <StatusBadge status={s.status} />
                              </span>
                            )}
                          </Badge>
                        ))}
                        {serverCount === 0 && (
                          <span className="text-muted-foreground">No servers</span>
                        )}
                        {serverCount > 6 && (
                          <span className="text-muted-foreground">+{serverCount - 6} more</span>
                        )}

                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setManagePool(p)
                            const ids = new Set((p.servers || []).map((s) => s.id as string))
                            setSelected(new Set(ids))
                            setInitialSelected(new Set(ids))
                            setServerFilter("")
                            setManageOpen(true)
                          }}
                        >
                          <LinkIcon className="mr-2 h-4 w-4" />
                          Manage Servers
                        </Button>
                      </div>
                    </TableCell>

                    {/* Annotations */}
                    <TableCell>
                      <div className="flex flex-wrap items-center gap-2">
                        {(p.annotations || []).slice(0, 6).map((a: any) => (
                          <Badge key={a.id} variant="outline" className="gap-1">
                            {a.key}:{a.value}
                          </Badge>
                        ))}
                        {(p.annotations || []).length === 0 && (
                          <span className="text-muted-foreground">No annotations</span>
                        )}
                        {(p.annotations || []).length > 6 && (
                          <span className="text-muted-foreground">
                            +{(p.annotations || []).length - 6} more
                          </span>
                        )}
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setManageLATPool(p)
                            setAnnotationsInitial(
                              new Set((p.annotations || []).map((x: any) => x.id as string))
                            )
                            setManageAnnotationsOpen(true)
                          }}
                        >
                          <LinkIcon className="mr-2 h-4 w-4" />
                          Manage
                        </Button>
                      </div>
                    </TableCell>

                    {/* Labels */}
                    <TableCell>
                      <div className="flex flex-wrap items-center gap-2">
                        {(p.labels || []).slice(0, 6).map((l: any) => (
                          <Badge key={l.id} variant="secondary" className="gap-1">
                            {l.key}:{l.value}
                          </Badge>
                        ))}
                        {(p.labels || []).length === 0 && (
                          <span className="text-muted-foreground">No labels</span>
                        )}
                        {(p.labels || []).length > 6 && (
                          <span className="text-muted-foreground">
                            +{(p.labels || []).length - 6} more
                          </span>
                        )}
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setManageLATPool(p)
                            setLabelsInitial(
                              new Set((p.labels || []).map((x: any) => x.id as string))
                            )
                            setManageLabelsOpen(true)
                          }}
                        >
                          <LinkIcon className="mr-2 h-4 w-4" />
                          Manage
                        </Button>
                      </div>
                    </TableCell>

                    {/* Taints */}
                    <TableCell>
                      <div className="flex flex-wrap items-center gap-2">
                        {(p.taints || []).slice(0, 6).map((t: any) => (
                          <Badge key={t.id} variant="outline" className="gap-1">
                            {t.key}:{t.value}
                            {t.effect ? <span className="ml-1">({t.effect})</span> : null}
                          </Badge>
                        ))}
                        {(p.taints || []).length === 0 && (
                          <span className="text-muted-foreground">No taints</span>
                        )}
                        {(p.taints || []).length > 6 && (
                          <span className="text-muted-foreground">
                            +{(p.taints || []).length - 6} more
                          </span>
                        )}
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setManageLATPool(p)
                            setTaintsInitial(
                              new Set((p.taints || []).map((x: any) => x.id as string))
                            )
                            setManageTaintsOpen(true)
                          }}
                        >
                          <LinkIcon className="mr-2 h-4 w-4" />
                          Manage
                        </Button>
                      </div>
                    </TableCell>

                    <TableCell className="text-right">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button size="icon" variant="ghost" className="h-8 w-8">
                            <Ellipsis className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem
                            onClick={() => {
                              setEditing(p)
                              setUpdateOpen(true)
                            }}
                          >
                            <Pencil className="mr-2 h-4 w-4" /> Edit
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            className="text-red-600 focus:text-red-600"
                            onClick={() => {
                              setDeleting(p)
                              setDeleteOpen(true)
                            }}
                          >
                            <Trash2 className="mr-2 h-4 w-4" /> Delete
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                )
              })}
              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={7} className="text-muted-foreground py-10 text-center">
                    No node pools found.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Edit dialog */}
      <Dialog open={updateOpen} onOpenChange={setUpdateOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Node Pool</DialogTitle>
          </DialogHeader>

          <Form {...updateForm}>
            <form className="space-y-4" onSubmit={updateForm.handleSubmit(onUpdateSubmit)}>
              <FormField
                control={updateForm.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="pool-name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="role"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role</FormLabel>
                    <Select
                      onValueChange={(v) =>
                        updateForm.setValue("role", v as Role, {
                          shouldDirty: true,
                          shouldValidate: true,
                        })
                      }
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select role" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="master">master</SelectItem>
                        <SelectItem value="worker">worker</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    setUpdateOpen(false)
                    setEditing(null)
                  }}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={updateForm.formState.isSubmitting}>
                  {updateForm.formState.isSubmitting ? "Saving…" : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete confirm */}
      <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Delete node pool</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This will permanently delete{" "}
            <span className="text-foreground font-medium">{deleting?.name}</span>.
          </p>
          <DialogFooter className="gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setDeleteOpen(false)
                setDeleting(null)
              }}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => deleteMut.mutate()}
              disabled={deleteMut.isPending}
            >
              {deleteMut.isPending ? "Deleting…" : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Manage Servers dialog */}
      <Dialog open={manageOpen} onOpenChange={setManageOpen}>
        <DialogContent className="sm:max-w-3xl">
          <DialogHeader>
            <DialogTitle>Manage Servers{managePool ? ` — ${managePool.name}` : ""}</DialogTitle>
          </DialogHeader>

          <div className="space-y-3">
            <div className="relative">
              <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
              <Input
                value={serverFilter}
                onChange={(e) => setServerFilter(e.target.value)}
                placeholder="Search by hostname, IP or role…"
                className="pl-8"
              />
            </div>

            <div className="max-h-[50vh] overflow-auto rounded border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[48px]"></TableHead>
                    <TableHead>Hostname</TableHead>
                    <TableHead>Private IP</TableHead>
                    <TableHead>Public IP</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {(serverQ.data ?? [])
                    .filter((s: any) => {
                      if (managePool?.role && !canAttachToPool(managePool.role as string, s.role)) {
                        return false
                      }
                      const q = serverFilter.trim().toLowerCase()
                      if (!q) return true
                      return (
                        (s.hostname ?? "").toLowerCase().includes(q) ||
                        (s.private_ip_address ?? "").toLowerCase().includes(q) ||
                        (s.public_ip_address ?? "").toLowerCase().includes(q) ||
                        (s.role ?? "").toLowerCase().includes(q)
                      )
                    })
                    .map((s: any) => {
                      const id = s.id as string
                      const checked = selected.has(id)
                      return (
                        <TableRow key={id}>
                          <TableCell className="text-center align-middle">
                            <input
                              type="checkbox"
                              className="h-4 w-4"
                              checked={checked}
                              onChange={(e) => {
                                const next = new Set(selected)
                                if (e.target.checked) next.add(id)
                                else next.delete(id)
                                setSelected(next)
                              }}
                            />
                          </TableCell>
                          <TableCell className="font-medium">{s.hostname || "—"}</TableCell>
                          <TableCell>{s.private_ip_address || "—"}</TableCell>
                          <TableCell>{s.public_ip_address || "—"}</TableCell>
                          <TableCell className="capitalize">{s.role || "—"}</TableCell>
                          <TableCell>
                            <StatusBadge status={s.status} />
                          </TableCell>
                        </TableRow>
                      )
                    })}
                  {(serverQ.data ?? []).length === 0 && (
                    <TableRow>
                      <TableCell colSpan={6} className="text-muted-foreground py-8 text-center">
                        {serverQ.isLoading ? "Loading servers…" : "No servers found."}
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>

            <div className="text-muted-foreground text-sm">
              Selected: <span className="text-foreground font-medium">{selected.size}</span>
            </div>
          </div>

          <DialogFooter className="gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setManageOpen(false)
                setManagePool(null)
                setSelected(new Set())
                setInitialSelected(new Set())
                setServerFilter("")
              }}
            >
              Close
            </Button>
            <Button
              onClick={async () => {
                if (!managePool) return
                const poolId = managePool.id as unknown as string

                // Compute diff
                const toAttach: string[] = []
                const toDetach: string[] = []
                for (const id of selected) if (!initialSelected.has(id)) toAttach.push(id)
                for (const id of initialSelected) if (!selected.has(id)) toDetach.push(id)

                try {
                  if (toAttach.length > 0) {
                    await nodePoolsApi.attachNodePoolServer(poolId, { server_ids: toAttach })
                  }
                  for (const id of toDetach) {
                    await nodePoolsApi.detachNodePoolServers(poolId, id)
                  }

                  await qc.invalidateQueries({ queryKey: ["node-pools"] })
                  await qc.invalidateQueries({ queryKey: ["servers"] })

                  toast.success("Servers updated for node pool.")
                  setManageOpen(false)
                  setManagePool(null)
                  setSelected(new Set())
                  setInitialSelected(new Set())
                  setServerFilter("")
                } catch (err: any) {
                  toast.error(err?.message ?? "Failed to update servers.")
                }
              }}
              disabled={serverQ.isLoading}
            >
              Save changes
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Manage Labels */}
      <ManageManyDialog
        open={manageLabelsOpen}
        onOpenChange={(o) => {
          setManageLabelsOpen(o)
          if (!o) setManageLATPool(null)
        }}
        title={`Manage Labels${manageLATPool ? ` — ${manageLATPool.name}` : ""}`}
        items={(labelQ.data ?? []).map((l: any) => ({
          id: l.id as string,
          key: l.key,
          value: l.value,
          name: `${l.key}:${l.value}`,
        }))}
        initialSelectedIds={labelsInitial}
        columns={[
          { header: "Key", render: (it) => <span className="font-medium">{it.key}</span> },
          { header: "Value", render: (it) => it.value ?? "—" },
        ]}
        onSave={async ({ toAttach, toDetach }) => {
          if (!manageLATPool) return
          const poolId = manageLATPool.id as unknown as string
          try {
            if (toAttach.length > 0) {
              await nodePoolsApi.attachNodePoolLabels(poolId, { label_ids: toAttach })
            }
            for (const id of toDetach) {
              await nodePoolsApi.detachNodePoolLabels(poolId, id)
            }
            await qc.invalidateQueries({ queryKey: ["node-pools"] })
            toast.success("Labels updated for node pool.")
          } catch (err: any) {
            toast.error(err?.message ?? "Failed to update labels.")
            throw err
          }
        }}
      />

      {/* Manage Annotations */}
      <ManageManyDialog
        open={manageAnnotationsOpen}
        onOpenChange={(o) => {
          setManageAnnotationsOpen(o)
          if (!o) setManageLATPool(null)
        }}
        title={`Manage Annotations${manageLATPool ? ` — ${manageLATPool.name}` : ""}`}
        items={(annotationQ.data ?? []).map((a: any) => ({
          id: a.id as string,
          key: a.key,
          value: a.value,
          name: `${a.key}:${a.value}`,
        }))}
        initialSelectedIds={annotationsInitial}
        columns={[
          { header: "Key", render: (it) => <span className="font-medium">{it.key}</span> },
          { header: "Value", render: (it) => it.value ?? "—" },
        ]}
        onSave={async ({ toAttach, toDetach }) => {
          if (!manageLATPool) return
          const poolId = manageLATPool.id as unknown as string
          try {
            if (toAttach.length > 0) {
              await nodePoolsApi.attachNodePoolAnnotations(poolId, { annotation_ids: toAttach })
            }
            for (const id of toDetach) {
              await nodePoolsApi.detachNodePoolAnnotations(poolId, id)
            }
            await qc.invalidateQueries({ queryKey: ["node-pools"] })
            toast.success("Annotations updated for node pool.")
          } catch (err: any) {
            toast.error(err?.message ?? "Failed to update annotations.")
            throw err
          }
        }}
      />

      {/* Manage Taints */}
      <ManageManyDialog
        open={manageTaintsOpen}
        onOpenChange={(o) => {
          setManageTaintsOpen(o)
          if (!o) setManageLATPool(null)
        }}
        title={`Manage Taints${manageLATPool ? ` — ${manageLATPool.name}` : ""}`}
        items={(taintQ.data ?? []).map((t: any) => ({
          id: t.id as string,
          key: t.key,
          value: t.value,
          effect: t.effect,
          name: `${t.key}:${t.value}`,
        }))}
        initialSelectedIds={taintsInitial}
        columns={[
          { header: "Key", render: (it) => <span className="font-medium">{it.key}</span> },
          { header: "Value", render: (it) => it.value ?? "—" },
          { header: "Effect", render: (it) => it.effect ?? "—" },
        ]}
        onSave={async ({ toAttach, toDetach }) => {
          if (!manageLATPool) return
          const poolId = manageLATPool.id as unknown as string
          try {
            if (toAttach.length > 0) {
              await nodePoolsApi.attachNodePoolTaints(poolId, { taint_ids: toAttach })
            }
            for (const id of toDetach) {
              await nodePoolsApi.detachNodePoolTaints(poolId, id)
            }
            await qc.invalidateQueries({ queryKey: ["node-pools"] })
            toast.success("Taints updated for node pool.")
          } catch (err: any) {
            toast.error(err?.message ?? "Failed to update taints.")
            throw err
          }
        }}
      />
    </div>
  )
}
