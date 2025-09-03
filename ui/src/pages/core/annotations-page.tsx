import { useEffect, useMemo, useState } from "react"
import { z } from "zod"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import {
    LinkIcon,
    Pencil,
    Plus,
    RefreshCw,
    Search,
    Trash,
    UnlinkIcon,
    ServerIcon,
} from "lucide-react"

import { api, ApiError } from "@/lib/api"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
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
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table"

/* ----------------------------- Types & Schemas ---------------------------- */

type NodePoolBrief = {
    id: string
    name: string
}

type Annotation = {
    id: string
    name: string
    value: string
    node_pools?: NodePoolBrief[]
}

const CreateSchema = z.object({
    name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
    value: z.string().trim().min(1, "Value is required").max(512, "Max 512 chars"),
    node_pool_ids: z.array(z.string().uuid()).optional().default([]),
})
type CreateInput = z.input<typeof CreateSchema>
type CreateValues = z.output<typeof CreateSchema>

const UpdateSchema = z.object({
    name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
    value: z.string().trim().min(1, "Value is required").max(512, "Max 512 chars"),
})
type UpdateValues = z.output<typeof UpdateSchema>

const AttachPoolsSchema = z.object({
    node_pool_ids: z.array(z.string().uuid()).min(1, "Pick at least one node pool"),
})
type AttachPoolsValues = z.output<typeof AttachPoolsSchema>

/* --------------------------------- Utils --------------------------------- */

function truncateMiddle(str: string, keep = 12) {
    if (!str || str.length <= keep * 2 + 3) return str
    return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

/* --------------------------------- Page ---------------------------------- */

export const AnnotationsPage = () => {
    const [loading, setLoading] = useState<boolean>(true)
    const [err, setErr] = useState<string | null>(null)

    const [annotations, setAnnotations] = useState<Annotation[]>([])
    const [allPools, setAllPools] = useState<NodePoolBrief[]>([])

    const [q, setQ] = useState("")

    // Dialog state
    const [createOpen, setCreateOpen] = useState(false)
    const [editTarget, setEditTarget] = useState<Annotation | null>(null)
    const [managePoolsTarget, setManagePoolsTarget] = useState<Annotation | null>(null)

    // Attached pools (for manage dialog)
    const [attachedPools, setAttachedPools] = useState<NodePoolBrief[]>([])
    const [attachedLoading, setAttachedLoading] = useState(false)
    const [attachedErr, setAttachedErr] = useState<string | null>(null)

    /* ------------------------------- Data Load ------------------------------ */

    async function loadAll() {
        setLoading(true)
        setErr(null)
        try {
            const [ann, pools] = await Promise.all([
                api.get<Annotation[]>("/api/v1/annotations?include=node_pools"),
                api.get<NodePoolBrief[]>("/api/v1/node-pools"),
            ])
            setAnnotations(ann || [])
            setAllPools(pools || [])

            // keep dialog targets in sync
            if (editTarget) {
                const refreshed = (ann || []).find((a) => a.id === editTarget.id) || null
                setEditTarget(refreshed)
            }
            if (managePoolsTarget) {
                const refreshed = (ann || []).find((a) => a.id === managePoolsTarget.id) || null
                setManagePoolsTarget(refreshed)
                if (refreshed) {
                    void loadAttachedPools(refreshed.id)
                }
            }
        } catch (e) {
            console.error(e)
            const msg =
                e instanceof ApiError ? e.message : "Failed to load annotations / node pools"
            setErr(msg)
        } finally {
            setLoading(false)
        }
    }

    async function loadAttachedPools(annotationId: string) {
        setAttachedLoading(true)
        setAttachedErr(null)
        try {
            const data = await api.get<NodePoolBrief[]>(
                `/api/v1/annotations/${annotationId}/node_pools`
            )
            setAttachedPools(data || [])
        } catch (e) {
            console.error(e)
            const msg =
                e instanceof ApiError ? e.message : "Failed to load pools for annotation"
            setAttachedErr(msg)
        } finally {
            setAttachedLoading(false)
        }
    }

    useEffect(() => {
        void loadAll()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    /* -------------------------------- Filters ------------------------------- */

    const filtered = useMemo(() => {
        const needle = q.trim().toLowerCase()
        if (!needle) return annotations
        return annotations.filter(
            (a) =>
                a.name.toLowerCase().includes(needle) ||
                a.value.toLowerCase().includes(needle) ||
                (a.node_pools || []).some((p) => p.name.toLowerCase().includes(needle))
        )
    }, [annotations, q])

    /* ------------------------------ Mutations ------------------------------- */

    async function deleteAnnotation(id: string) {
        if (!confirm("Delete this annotation? This cannot be undone.")) return
        await api.delete<void>(`/api/v1/annotations/${id}`)
        await loadAll()
    }

    // Create
    const createForm = useForm<CreateInput, any, CreateValues>({
        resolver: zodResolver(CreateSchema),
        defaultValues: { name: "", value: "", node_pool_ids: [] },
    })

    const submitCreate = async (values: CreateValues) => {
        const payload: any = { name: values.name.trim(), value: values.value.trim() }
        if (values.node_pool_ids && values.node_pool_ids.length > 0) {
            payload.node_pool_ids = values.node_pool_ids
        }
        await api.post("/api/v1/annotations", payload)
        setCreateOpen(false)
        createForm.reset({ name: "", value: "", node_pool_ids: [] })
        await loadAll()
    }

    // Edit
    const editForm = useForm<UpdateValues>({
        resolver: zodResolver(UpdateSchema),
        defaultValues: { name: "", value: "" },
    })

    function openEdit(a: Annotation) {
        setEditTarget(a)
        editForm.reset({ name: a.name, value: a.value })
    }

    const submitEdit = async (values: UpdateValues) => {
        if (!editTarget) return
        await api.patch(`/api/v1/annotations/${editTarget.id}`, {
            name: values.name.trim(),
            value: values.value.trim(),
        })
        setEditTarget(null)
        await loadAll()
    }

    // Manage pools (attach/detach)
    const attachPoolsForm = useForm<AttachPoolsValues>({
        resolver: zodResolver(AttachPoolsSchema),
        defaultValues: { node_pool_ids: [] },
    })

    function openManagePools(a: Annotation) {
        setManagePoolsTarget(a)
        attachPoolsForm.reset({ node_pool_ids: [] })
        void loadAttachedPools(a.id)
    }

    const submitAttachPools = async (values: AttachPoolsValues) => {
        if (!managePoolsTarget) return
        await api.post(`/api/v1/annotations/${managePoolsTarget.id}/node_pools`, {
            node_pool_ids: values.node_pool_ids,
        })
        attachPoolsForm.reset({ node_pool_ids: [] })
        await loadAttachedPools(managePoolsTarget.id)
        await loadAll()
    }

    async function detachPool(poolId: string) {
        if (!managePoolsTarget) return
        if (!confirm("Detach this node pool from the annotation?")) return
        await api.delete(`/api/v1/annotations/${managePoolsTarget.id}/node_pools/${poolId}`)
        await loadAttachedPools(managePoolsTarget.id)
        await loadAll()
    }

    /* --------------------------------- Render -------------------------------- */

    if (loading) return <div className="p-6">Loading annotations…</div>
    if (err) return <div className="p-6 text-red-500">{err}</div>

    return (
        <div className="space-y-4 p-6">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <h1 className="mb-4 text-2xl font-bold">Annotations</h1>

                <div className="flex flex-wrap items-center gap-2">
                    <div className="relative">
                        <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
                        <Input
                            value={q}
                            onChange={(e) => setQ(e.target.value)}
                            placeholder="Search name, value, pool…"
                            className="w-72 pl-8"
                        />
                    </div>

                    <Button variant="outline" onClick={loadAll}>
                        <RefreshCw className="mr-2 h-4 w-4" /> Refresh
                    </Button>

                    <Dialog open={createOpen} onOpenChange={setCreateOpen}>
                        <DialogTrigger asChild>
                            <Button onClick={() => setCreateOpen(true)}>
                                <Plus className="mr-2 h-4 w-4" /> Create Annotation
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-lg">
                            <DialogHeader>
                                <DialogTitle>Create annotation</DialogTitle>
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
                                                    <Input placeholder="cluster-autoscaler.kubernetes.io/safe-to-evict" {...field} />
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
                                                    <Input placeholder="true" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />
                                    <FormField
                                        control={createForm.control}
                                        name="node_pool_ids"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Initial node pools (optional)</FormLabel>
                                                <div className="max-h-56 space-y-2 overflow-auto rounded-xl border p-2">
                                                    {allPools.length === 0 && (
                                                        <div className="text-muted-foreground p-2 text-sm">No node pools available</div>
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
                                <TableHead>Name</TableHead>
                                <TableHead>Value</TableHead>
                                <TableHead>Node Pools</TableHead>
                                <TableHead className="w-[180px] text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {filtered.map((a) => {
                                const pools = a.node_pools || []
                                return (
                                    <TableRow key={a.id}>
                                        <TableCell className="font-mono text-sm">{a.name}</TableCell>
                                        <TableCell className="font-mono text-sm">{a.value}</TableCell>

                                        <TableCell>
                                            <div className="mb-2 flex flex-wrap gap-2">
                                                {pools.slice(0, 6).map((p) => (
                                                    <Badge key={p.id} variant="secondary" className="gap-1">
                                                        <ServerIcon className="h-3 w-3" /> {p.name}
                                                    </Badge>
                                                ))}
                                                {pools.length === 0 && (
                                                    <span className="text-muted-foreground">No node pools</span>
                                                )}
                                                {pools.length > 6 && (
                                                    <span className="text-muted-foreground">+{pools.length - 6} more</span>
                                                )}
                                            </div>
                                            <Button variant="outline" size="sm" onClick={() => openManagePools(a)}>
                                                <LinkIcon className="mr-2 h-4 w-4" /> Manage Node Pools
                                            </Button>
                                        </TableCell>

                                        <TableCell>
                                            <div className="flex justify-end gap-2">
                                                <Button variant="outline" size="sm" onClick={() => openEdit(a)}>
                                                    <Pencil className="mr-2 h-4 w-4" /> Edit
                                                </Button>
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button variant="destructive" size="sm">
                                                            <Trash className="mr-2 h-4 w-4" /> Delete
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align="end">
                                                        <DropdownMenuItem onClick={() => deleteAnnotation(a.id)}>
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
                                    <TableCell colSpan={4} className="text-muted-foreground py-10 text-center">
                                        No annotations match your search.
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
                        <DialogTitle>Edit annotation</DialogTitle>
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
                                            <Input placeholder="example.com/some" {...field} />
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
                                        <FormLabel>Value</FormLabel>
                                        <FormControl>
                                            <Input placeholder="true" {...field} />
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

            {/* Manage node pools dialog */}
            <Dialog open={!!managePoolsTarget} onOpenChange={(o) => !o && setManagePoolsTarget(null)}>
                <DialogContent className="sm:max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>
                            Manage node pools for{" "}
                            <span className="font-mono">{managePoolsTarget?.name}</span>
                        </DialogTitle>
                    </DialogHeader>

                    {/* Attached pools list */}
                    <div className="space-y-3">
                        <div className="text-sm font-medium">Attached node pools</div>

                        {attachedLoading ? (
                            <div className="text-muted-foreground rounded-md border p-3 text-sm">Loading…</div>
                        ) : attachedErr ? (
                            <div className="rounded-md border p-3 text-sm text-red-500">{attachedErr}</div>
                        ) : (
                            <div className="overflow-hidden rounded-xl border">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Name</TableHead>
                                            <TableHead className="w-[120px] text-right">Detach</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {attachedPools.map((p) => (
                                            <TableRow key={p.id}>
                                                <TableCell className="font-medium">{p.name}</TableCell>
                                                <TableCell>
                                                    <div className="flex justify-end">
                                                        <Button
                                                            variant="destructive"
                                                            size="sm"
                                                            onClick={() => detachPool(p.id)}
                                                        >
                                                            <UnlinkIcon className="mr-2 h-4 w-4" /> Detach
                                                        </Button>
                                                    </div>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                        {attachedPools.length === 0 && (
                                            <TableRow>
                                                <TableCell colSpan={2} className="text-muted-foreground py-8 text-center">
                                                    No node pools attached yet.
                                                </TableCell>
                                            </TableRow>
                                        )}
                                    </TableBody>
                                </Table>
                            </div>
                        )}
                    </div>

                    {/* Attach pools */}
                    <div className="pt-4">
                        <Form {...attachPoolsForm}>
                            <form onSubmit={attachPoolsForm.handleSubmit(submitAttachPools)} className="space-y-3">
                                <FormField
                                    control={attachPoolsForm.control}
                                    name="node_pool_ids"
                                    render={({ field }) => (
                                        <FormItem>
                                            <FormLabel>Attach more node pools</FormLabel>
                                            <div className="grid max-h-64 grid-cols-1 gap-2 overflow-auto rounded-xl border p-2 md:grid-cols-2">
                                                {(() => {
                                                    const attachedIds = new Set(attachedPools.map((p) => p.id))
                                                    const attachable = allPools.filter((p) => !attachedIds.has(p.id))
                                                    if (attachable.length === 0) {
                                                        return (
                                                            <div className="text-muted-foreground p-2 text-sm">
                                                                No more node pools available to attach
                                                            </div>
                                                        )
                                                    }
                                                    return attachable.map((p) => {
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
                                                    })
                                                })()}
                                            </div>
                                            <FormMessage />
                                        </FormItem>
                                    )}
                                />

                                <DialogFooter className="gap-2">
                                    <Button type="submit" disabled={attachPoolsForm.formState.isSubmitting}>
                                        <LinkIcon className="mr-2 h-4 w-4" />{" "}
                                        {attachPoolsForm.formState.isSubmitting ? "Attaching…" : "Attach selected"}
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
