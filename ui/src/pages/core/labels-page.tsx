import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { LinkIcon, PencilIcon, Plus, TrashIcon, UnlinkIcon } from "lucide-react"
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

type NodePoolBrief = {
  id: string
  name: string
}

type Label = {
  id: string
  key: string
  value: string
  node_pools?: NodePoolBrief[] // normalized from API's "node_groups"
}

const CreateLabelSchema = z.object({
  key: z.string().trim().min(2, "Key is too short"),
  value: z.string().trim().min(2, "Value is too short"),
})
type CreateLabelValues = z.infer<typeof CreateLabelSchema>

const UpdateLabelSchema = z
  .object({
    key: z.string().trim().min(2, "Key is too short").optional(),
    value: z.string().trim().min(2, "Value is too short").optional(),
  })
  .refine((v) => v.key !== undefined || v.value !== undefined, {
    message: "Provide a new key or value",
    path: ["key"],
  })
type UpdateLabelValues = z.infer<typeof UpdateLabelSchema>

const AttachPoolsSchema = z.object({
  node_pool_ids: z.array(z.string().uuid()).min(1, "Pick at least one node pool"),
})
type AttachPoolsValues = z.infer<typeof AttachPoolsSchema>

function truncateMiddle(str: string, keep = 8) {
  if (!str || str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

export const LabelsPage = () => {
  const [labels, setLabels] = useState<Label[]>([])
  const [allPools, setAllPools] = useState<NodePoolBrief[]>([])
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)

  const [createOpen, setCreateOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Label | null>(null)

  const [manageTarget, setManageTarget] = useState<Label | null>(null)

  async function loadAll() {
    setLoading(true)
    setErr(null)
    try {
      // include=node_pools -> backend returns "node_groups" field; normalize it.
      const [labelsRaw, poolsRaw] = await Promise.all([
        api.get<any[]>("/api/v1/labels?include=node_pools"),
        api.get<NodePoolBrief[]>("/api/v1/node-pools"),
      ])

      const normalized: Label[] = (labelsRaw || []).map((l) => ({
        id: l.id,
        key: l.key,
        value: l.value,
        node_pools: l.node_pools ?? l.node_groups ?? [], // support either
      }))

      setLabels(normalized)
      setAllPools(poolsRaw || [])

      if (manageTarget) {
        const refreshed = normalized.find((x) => x.id === manageTarget.id) || null
        setManageTarget(refreshed)
      }
      if (editTarget) {
        const refreshed = normalized.find((x) => x.id === editTarget.id) || null
        setEditTarget(refreshed)
      }
    } catch (e) {
      console.error(e)
      const msg = e instanceof ApiError ? e.message : "Failed to load labels/pools"
      setErr(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadAll()
  }, [])

  // CREATE
  const createForm = useForm<CreateLabelValues>({
    resolver: zodResolver(CreateLabelSchema),
    defaultValues: { key: "", value: "" },
  })

  const submitCreate = async (values: CreateLabelValues) => {
    await api.post<Label>("/api/v1/labels", {
      key: values.key.trim(),
      value: values.value.trim(),
    })
    setCreateOpen(false)
    createForm.reset()
    await loadAll()
  }

  // EDIT
  const editForm = useForm<UpdateLabelValues>({
    resolver: zodResolver(UpdateLabelSchema),
    defaultValues: { key: undefined, value: undefined },
  })

  function openEdit(l: Label) {
    setEditTarget(l)
    editForm.reset({ key: undefined, value: undefined })
    setEditOpen(true)
  }

  const submitEdit = async (values: UpdateLabelValues) => {
    if (!editTarget) return
    const payload: Partial<Label> = {}
    if (values.key !== undefined) payload.key = values.key.trim()
    if (values.value !== undefined) payload.value = values.value.trim()
    await api.patch<Label>(`/api/v1/labels/${editTarget.id}`, payload)
    setEditOpen(false)
    setEditTarget(null)
    await loadAll()
  }

  // DELETE
  async function deleteLabel(id: string) {
    if (!confirm("Delete this label? This cannot be undone.")) return
    await api.delete(`/api/v1/labels/${id}`)
    await loadAll()
  }

  // MANAGE NODE POOLS (attach/detach)
  const attachForm = useForm<AttachPoolsValues>({
    resolver: zodResolver(AttachPoolsSchema),
    defaultValues: { node_pool_ids: [] },
  })

  function openManage(l: Label) {
    setManageTarget(l)
    attachForm.reset({ node_pool_ids: [] })
  }

  const submitAttachPools = async (values: AttachPoolsValues) => {
    if (!manageTarget) return
    await api.post<Label>(`/api/v1/labels/${manageTarget.id}/node_pools`, {
      node_pool_ids: values.node_pool_ids,
    })
    attachForm.reset({ node_pool_ids: [] })
    await loadAll()
  }

  async function detachPool(poolId: string) {
    if (!manageTarget) return
    if (!confirm("Detach this label from the selected node pool?")) return
    await api.delete(`/api/v1/labels/${manageTarget.id}/node_pools/${poolId}`)
    await loadAll()
  }

  const attachablePools = useMemo(() => {
    if (!manageTarget) return [] as NodePoolBrief[]
    const attached = new Set((manageTarget.node_pools || []).map((p) => p.id))
    return allPools.filter((p) => !attached.has(p.id))
  }, [manageTarget, allPools])

  if (loading) return <div className="p-6">Loading labels…</div>
  if (err) return <div className="p-6 text-red-500">{err}</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Labels</h1>

        <Dialog open={createOpen} onOpenChange={setCreateOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => setCreateOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Create Label
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-lg">
            <DialogHeader>
              <DialogTitle>Create Label</DialogTitle>
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
                        <Input placeholder="app.kubernetes.io/managed-by" {...field} />
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
                      <FormLabel>Value</FormLabel>
                      <FormControl>
                        <Input placeholder="GlueOps" {...field} />
                      </FormControl>
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

      <div className="bg-background overflow-hidden rounded-2xl border shadow-sm">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Key</TableHead>
                <TableHead>Value</TableHead>
                <TableHead>Node Pools</TableHead>
                <TableHead className="w-[260px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {labels.map((l) => (
                <TableRow key={l.id}>
                  <TableCell className="font-medium">{l.key}</TableCell>
                  <TableCell>{l.value}</TableCell>
                  <TableCell>
                    <div className="flex flex-wrap gap-2">
                      {(l.node_pools || []).slice(0, 6).map((p) => (
                        <Badge key={p.id} variant="secondary">
                          {p.name}
                        </Badge>
                      ))}
                      {(l.node_pools || []).length === 0 && (
                        <span className="text-muted-foreground">No pools</span>
                      )}
                      {(l.node_pools || []).length > 6 && (
                        <span className="text-muted-foreground">
                          +{(l.node_pools || []).length - 6} more
                        </span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => openManage(l)}>
                        <LinkIcon className="mr-2 h-4 w-4" />
                        Manage node pools
                      </Button>
                      <Button variant="outline" size="sm" onClick={() => openEdit(l)}>
                        <PencilIcon className="mr-2 h-4 w-4" />
                        Edit
                      </Button>
                      <Button variant="destructive" size="sm" onClick={() => deleteLabel(l.id)}>
                        <TrashIcon className="mr-2 h-4 w-4" />
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {labels.length === 0 && (
                <TableRow>
                  <TableCell colSpan={4} className="text-muted-foreground py-10 text-center">
                    No labels yet.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Edit label */}
      <Dialog open={editOpen} onOpenChange={(o) => !o && setEditOpen(false)}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>
              Edit Label{" "}
              {editTarget ? (
                <span className="text-muted-foreground ml-2 font-mono text-sm">
                  ({editTarget.key} = {editTarget.value})
                </span>
              ) : null}
            </DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form onSubmit={editForm.handleSubmit(submitEdit)} className="space-y-4">
              <FormField
                control={editForm.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New key (optional)</FormLabel>
                    <FormControl>
                      <Input placeholder={editTarget?.key || "e.g. app"} {...field} />
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
                    <FormLabel>New value (optional)</FormLabel>
                    <FormControl>
                      <Input placeholder={editTarget?.value || "e.g. GlueOps"} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button type="button" variant="outline" onClick={() => setEditOpen(false)}>
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

      {/* Manage node pools for a label */}
      <Dialog open={!!manageTarget} onOpenChange={(o) => !o && setManageTarget(null)}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>
              Manage node pools for{" "}
              <span className="font-mono">
                {manageTarget ? `${manageTarget.key}=${manageTarget.value}` : ""}
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
                  {(manageTarget?.node_pools || []).map((p) => (
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

                  {(manageTarget?.node_pools || []).length === 0 && (
                    <TableRow>
                      <TableCell colSpan={2} className="text-muted-foreground py-8 text-center">
                        No pools attached yet.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          </div>

          {/* Attach more */}
          <div className="pt-4">
            <Form {...attachForm}>
              <form onSubmit={attachForm.handleSubmit(submitAttachPools)} className="space-y-3">
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
                    <LinkIcon className="mr-2 h-4 w-4" />
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
