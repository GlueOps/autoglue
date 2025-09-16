import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import {
    LinkIcon,
    PencilIcon,
    Plus,
    RefreshCcw,
    Server as ServerIcon,
    TrashIcon,
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
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import { Textarea } from "@/components/ui/textarea.tsx"
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table.tsx"

// ---- Types ----
export type NodePoolBrief = {
    id: string
    name: string
}

export type ServerBrief = {
    id: string
    hostname: string
    ip: string
    role: string
    status: string
}

export type Cluster = {
    id: string
    name: string
    provider: string
    region: string
    status: string
    cluster_load_balancer?: string
    control_load_balancer?: string
    node_pools?: NodePoolBrief[]
    bastion_server?: ServerBrief | null
}

// ---- Schemas ----
const CreateClusterSchema = z.object({
    name: z.string().trim().min(2, "Name is too short"),
    provider: z.string().trim().min(2, "Provider is too short"),
    region: z.string().trim().min(1, "Region is required"),
    node_pool_ids: z.array(z.string().uuid()).optional().default([]),
    bastion_server_id: z.string().uuid().optional(),
    cluster_load_balancer: z.string().optional(),
    control_load_balancer: z.string().optional(),
    kubeconfig: z.string().optional(),
})
export type CreateClusterValues = z.infer<typeof CreateClusterSchema>

const UpdateClusterSchema = z
    .object({
        name: z.string().trim().min(2, "Name is too short").optional(),
        provider: z.string().trim().min(2, "Provider is too short").optional(),
        region: z.string().trim().min(1, "Region is required").optional(),
        status: z.string().trim().min(1, "Status is required").optional(),
        bastion_server_id: z.string().uuid().or(z.literal("")).optional(),
        cluster_load_balancer: z.string().optional(),
        control_load_balancer: z.string().optional(),
        kubeconfig: z.string().optional(),
    })
    .refine(
        (v) =>
            v.name !== undefined ||
            v.provider !== undefined ||
            v.region !== undefined ||
            v.status !== undefined ||
            v.bastion_server_id !== undefined ||
            v.cluster_load_balancer !== undefined ||
            v.control_load_balancer !== undefined ||
            v.kubeconfig !== undefined,
        { message: "Provide at least one change", path: ["name"] }
    )
export type UpdateClusterValues = z.infer<typeof UpdateClusterSchema>

const AttachPoolsSchema = z.object({
    node_pool_ids: z.array(z.string().uuid()).min(1, "Pick at least one node pool"),
})
export type AttachPoolsValues = z.infer<typeof AttachPoolsSchema>

const SetBastionSchema = z.object({
    server_id: z.string().uuid({ message: "Enter a valid Server UUID" }),
})
export type SetBastionValues = z.infer<typeof SetBastionSchema>

// ---- Utils ----
function truncateMiddle(str: string, keep = 8) {
    if (!str || str.length <= keep * 2 + 3) return str
    return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

// ---- Component ----
export  function ClustersPage() {
    const [clusters, setClusters] = useState<Cluster[]>([])
    const [allPools, setAllPools] = useState<NodePoolBrief[]>([])
    const [bastionCandidates, setBastionCandidates] = useState<ServerBrief[]>([])

    const [loading, setLoading] = useState(false)
    const [err, setErr] = useState<string | null>(null)

    const [q, setQ] = useState("")

    // dialogs
    const [createOpen, setCreateOpen] = useState(false)
    const [editOpen, setEditOpen] = useState(false)
    const [editTarget, setEditTarget] = useState<Cluster | null>(null)

    const [managePoolsTarget, setManagePoolsTarget] = useState<Cluster | null>(null)
    const [manageBastionTarget, setManageBastionTarget] = useState<Cluster | null>(null)

    async function loadAll() {
        setLoading(true)
        setErr(null)
        try {
            const url = `/api/v1/clusters?include=node_pools,bastion${q ? `&q=${encodeURIComponent(q)}` : ""}`
            const [clustersRaw, poolsRaw, serversRaw] = await Promise.all([
                api.get<any[]>(url),
                api.get<NodePoolBrief[]>("/api/v1/node-pools"),
                // Best-effort; if this endpoint doesn't exist, we'll just fall back to manual input
                api.get<ServerBrief[]>("/api/v1/servers?role=bastion").catch(() => [] as any),
            ])

            const normalized: Cluster[] = (clustersRaw || []).map((c) => ({
                id: c.id,
                name: c.name,
                provider: c.provider,
                region: c.region,
                status: c.status,
                cluster_load_balancer: c.cluster_load_balancer,
                control_load_balancer: c.control_load_balancer,
                node_pools: c.node_pools ?? [],
                bastion_server: c.bastion_server ?? null,
            }))

            setClusters(normalized)
            setAllPools(poolsRaw || [])
            setBastionCandidates(serversRaw || [])

            // keep dialogs in sync after refresh
            if (managePoolsTarget) {
                const refreshed = normalized.find((x) => x.id === managePoolsTarget.id) || null
                setManagePoolsTarget(refreshed)
            }
            if (manageBastionTarget) {
                const refreshed = normalized.find((x) => x.id === manageBastionTarget.id) || null
                setManageBastionTarget(refreshed)
            }
            if (editTarget) {
                const refreshed = normalized.find((x) => x.id === editTarget.id) || null
                setEditTarget(refreshed)
            }
        } catch (e) {
            console.error(e)
            const msg = e instanceof ApiError ? e.message : "Failed to load clusters"
            setErr(msg)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        void loadAll()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    // ---- CREATE ----
    const createForm = useForm<CreateClusterValues>({
        resolver: zodResolver(CreateClusterSchema),
        defaultValues: {
            name: "",
            provider: "",
            region: "",
            node_pool_ids: [],
            bastion_server_id: undefined,
            cluster_load_balancer: "",
            control_load_balancer: "",
            kubeconfig: "",
        },
    })

    const submitCreate = async (values: CreateClusterValues) => {
        const payload: any = {
            name: values.name.trim(),
            provider: values.provider.trim(),
            region: values.region.trim(),
            node_pool_ids: values.node_pool_ids || [],
        }
        if (values.bastion_server_id) payload.bastion_server_id = values.bastion_server_id
        if (values.cluster_load_balancer) payload.cluster_load_balancer = values.cluster_load_balancer
        if (values.control_load_balancer) payload.control_load_balancer = values.control_load_balancer
        if (values.kubeconfig && values.kubeconfig.trim()) payload.kubeconfig = values.kubeconfig.trim()

        await api.post<Cluster>("/api/v1/clusters", payload)
        setCreateOpen(false)
        createForm.reset()
        await loadAll()
    }

    // ---- EDIT ----
    const editForm = useForm<UpdateClusterValues>({
        resolver: zodResolver(UpdateClusterSchema),
        defaultValues: {
            name: undefined,
            provider: undefined,
            region: undefined,
            status: undefined,
            bastion_server_id: undefined,
            cluster_load_balancer: undefined,
            control_load_balancer: undefined,
            kubeconfig: undefined,
        },
    })

    function openEdit(c: Cluster) {
        setEditTarget(c)
        editForm.reset({
            name: undefined,
            provider: undefined,
            region: undefined,
            status: undefined,
            bastion_server_id: undefined,
            cluster_load_balancer: undefined,
            control_load_balancer: undefined,
        })
        setEditOpen(true)
    }

    const submitEdit = async (values: UpdateClusterValues) => {
        if (!editTarget) return
        const payload: any = {}
        if (values.name !== undefined) payload.name = values.name.trim()
        if (values.provider !== undefined) payload.provider = values.provider.trim()
        if (values.region !== undefined) payload.region = values.region.trim()
        if (values.status !== undefined) payload.status = values.status.trim()
        if (values.bastion_server_id !== undefined) payload.bastion_server_id = values.bastion_server_id || ""
        if (values.cluster_load_balancer !== undefined)
            payload.cluster_load_balancer = values.cluster_load_balancer
        if (values.control_load_balancer !== undefined)
            payload.control_load_balancer = values.control_load_balancer
        if (values.kubeconfig !== undefined && values.kubeconfig.trim())
            payload.kubeconfig = values.kubeconfig.trim()

        await api.patch<Cluster>(`/api/v1/clusters/${editTarget.id}`, payload)
        setEditOpen(false)
        setEditTarget(null)
        await loadAll()
    }

    // ---- DELETE ----
    async function deleteCluster(id: string) {
        if (!confirm("Delete this cluster? This cannot be undone.")) return
        await api.delete(`/api/v1/clusters/${id}`)
        await loadAll()
    }

    // ---- MANAGE NODE POOLS ----
    const attachForm = useForm<AttachPoolsValues>({
        resolver: zodResolver(AttachPoolsSchema),
        defaultValues: { node_pool_ids: [] },
    })

    function openManagePools(c: Cluster) {
        setManagePoolsTarget(c)
        attachForm.reset({ node_pool_ids: [] })
    }

    const submitAttachPools = async (values: AttachPoolsValues) => {
        if (!managePoolsTarget) return
        await api.post(`/api/v1/clusters/${managePoolsTarget.id}/node_pools`, {
            node_pool_ids: values.node_pool_ids,
        })
        attachForm.reset({ node_pool_ids: [] })
        await loadAll()
    }

    async function detachPool(clusterId: string, poolId: string) {
        if (!confirm("Detach selected node pool?")) return
        await api.delete(`/api/v1/clusters/${clusterId}/node_pools/${poolId}`)
        await loadAll()
    }

    const attachablePools = useMemo(() => {
        if (!managePoolsTarget) return [] as NodePoolBrief[]
        const attached = new Set((managePoolsTarget.node_pools || []).map((p) => p.id))
        return allPools.filter((p) => !attached.has(p.id))
    }, [managePoolsTarget, allPools])

    // ---- MANAGE BASTION ----
    const setBastionForm = useForm<SetBastionValues>({
        resolver: zodResolver(SetBastionSchema),
        defaultValues: { server_id: "" },
    })

    function openManageBastion(c: Cluster) {
        setManageBastionTarget(c)
        setBastionForm.reset({ server_id: "" })
    }

    const submitSetBastion = async (values: SetBastionValues) => {
        if (!manageBastionTarget) return
        await api.post(`/api/v1/clusters/${manageBastionTarget.id}/bastion`, {
            server_id: values.server_id,
        })
        await loadAll()
    }

    async function clearBastion() {
        if (!manageBastionTarget) return
        if (!confirm("Clear bastion for this cluster?")) return
        await api.delete(`/api/v1/clusters/${manageBastionTarget.id}/bastion`)
        await loadAll()
    }

    // ---- UI ----
    if (loading) return <div className="p-6">Loading clusters…</div>
    if (err) return <div className="p-6 text-red-500">{err}</div>

    return (
        <div className="space-y-4 p-6">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <h1 className="mb-4 text-2xl font-bold">Clusters</h1>

                <div className="flex flex-1 items-center gap-2 md:justify-end">
                    <Input
                        className="max-w-xs"
                        placeholder="Filter by name…"
                        value={q}
                        onChange={(e) => setQ(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === "Enter") void loadAll()
                        }}
                    />
                    <Button variant="outline" onClick={() => void loadAll()}>
                        <RefreshCcw className="mr-2 h-4 w-4" />
                        Apply
                    </Button>

                    <Dialog open={createOpen} onOpenChange={setCreateOpen}>
                        <DialogTrigger asChild>
                            <Button onClick={() => setCreateOpen(true)}>
                                <Plus className="mr-2 h-4 w-4" />
                                Create Cluster
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="sm:max-w-2xl">
                            <DialogHeader>
                                <DialogTitle>Create Cluster</DialogTitle>
                            </DialogHeader>

                            <Form {...createForm}>
                                <form onSubmit={createForm.handleSubmit(submitCreate)} className="grid gap-4 md:grid-cols-2">
                                    <FormField
                                        control={createForm.control}
                                        name="name"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Name</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="my-eks-prod" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={createForm.control}
                                        name="provider"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Provider</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="aws|gcp|azure|onprem" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={createForm.control}
                                        name="region"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Region</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="eu-west-1" {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={createForm.control}
                                        name="bastion_server_id"
                                        render={({ field }) => (
                                            <FormItem>
                                                <FormLabel>Bastion server (UUID, optional)</FormLabel>
                                                <FormControl>
                                                    <Input placeholder="paste server UUID" {...field} />
                                                </FormControl>
                                                {bastionCandidates.length > 0 && (
                                                    <div className="text-xs text-muted-foreground">
                                                        Suggestions:
                                                        <div className="mt-1 flex flex-wrap gap-2">
                                                            {bastionCandidates.slice(0, 6).map((s) => (
                                                                <Button
                                                                    key={s.id}
                                                                    type="button"
                                                                    size="sm"
                                                                    variant={field.value === s.id ? "default" : "outline"}
                                                                    onClick={() => field.onChange(s.id)}
                                                                    className="font-normal"
                                                                >
                                                                    <ServerIcon className="mr-1 h-3 w-3" /> {s.hostname || truncateMiddle(s.id, 6)}
                                                                </Button>
                                                            ))}
                                                        </div>
                                                    </div>
                                                )}
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={createForm.control}
                                        name="cluster_load_balancer"
                                        render={({ field }) => (
                                            <FormItem className="md:col-span-2">
                                                <FormLabel>Cluster Load Balancer (optional)</FormLabel>
                                                <FormControl>
                                                    <Textarea placeholder="e.g. JSON or URL or ARN" rows={2} {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <FormField
                                        control={createForm.control}
                                        name="control_load_balancer"
                                        render={({ field }) => (
                                            <FormItem className="md:col-span-2">
                                                <FormLabel>Control Load Balancer (optional)</FormLabel>
                                                <FormControl>
                                                    <Textarea placeholder="e.g. JSON or URL or ARN" rows={2} {...field} />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    {/* Node pools */}
                                    <FormField
                                        control={createForm.control}
                                        name="node_pool_ids"
                                        render={({ field }) => (
                                            <FormItem className="md:col-span-2">
                                                <FormLabel>Attach node pools (optional)</FormLabel>
                                                <div className="grid max-h-64 grid-cols-1 gap-2 overflow-auto rounded-xl border p-2 md:grid-cols-2">
                                                    {allPools.length === 0 && (
                                                        <div className="text-muted-foreground p-2 text-sm">No node pools available</div>
                                                    )}
                                                    {allPools.map((p) => {
                                                        const checked = field.value?.includes(p.id) || false
                                                        return (
                                                            <label key={p.id} className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1">
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
                                                                    <div className="text-muted-foreground text-xs">{truncateMiddle(p.id, 8)}</div>
                                                                </div>
                                                            </label>
                                                        )
                                                    })}
                                                </div>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    {/* KubeConfig */}
                                    <FormField
                                        control={createForm.control}
                                        name="kubeconfig"
                                        render={({field}) => (
                                            <FormItem className='md:colspan-2'>
                                                <FormLabel>Kubeconfig (optional)</FormLabel>
                                                <FormControl>
                                                    <Textarea
                                                        placeholder="Paste full kubeconfig YAML here. It will be encrypted and never returned by the API."
                                                        rows={8}
                                                        className="font-mono"
                                                        {...field}
                                                    />
                                                </FormControl>
                                                <FormMessage />
                                            </FormItem>
                                        )}
                                    />

                                    <DialogFooter className="md:col-span-2 gap-2">
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

            {/* Table */}
            <div className="bg-background overflow-hidden rounded-2xl border shadow-sm">
                <div className="overflow-x-auto">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Name</TableHead>
                                <TableHead>Provider / Region</TableHead>
                                <TableHead>Status</TableHead>
                                <TableHead>Node Pools</TableHead>
                                <TableHead>Bastion</TableHead>
                                <TableHead className="w-[360px] text-right">Actions</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {clusters.map((c) => (
                                <TableRow key={c.id}>
                                    <TableCell className="font-medium">{c.name}</TableCell>
                                    <TableCell>
                                        <div className="text-sm font-medium">{c.provider || "—"}</div>
                                        <div className="text-muted-foreground text-xs">{c.region || "—"}</div>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant={c.status === "ready" ? "default" : c.status === "error" ? "destructive" : "secondary"}>
                                            {c.status || "unknown"}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex max-w-[280px] flex-wrap gap-2">
                                            {(c.node_pools || []).slice(0, 4).map((p) => (
                                                <Badge key={p.id} variant="secondary">{p.name}</Badge>
                                            ))}
                                            {(c.node_pools || []).length === 0 && (
                                                <span className="text-muted-foreground">No pools</span>
                                            )}
                                            {(c.node_pools || []).length > 4 && (
                                                <span className="text-muted-foreground">+{(c.node_pools || []).length - 4} more</span>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell>
                                        {c.bastion_server ? (
                                            <div className="leading-tight">
                                                <div className="text-sm font-medium">{c.bastion_server.hostname || truncateMiddle(c.bastion_server.id, 6)}</div>
                                                <div className="text-muted-foreground text-xs">{c.bastion_server.ip}</div>
                                            </div>
                                        ) : (
                                            <span className="text-muted-foreground">None</span>
                                        )}
                                    </TableCell>
                                    <TableCell>
                                        <div className="flex justify-end gap-2">
                                            <Button variant="outline" size="sm" onClick={() => openManagePools(c)}>
                                                <LinkIcon className="mr-2 h-4 w-4" /> Manage pools
                                            </Button>
                                            <Button variant="outline" size="sm" onClick={() => openManageBastion(c)}>
                                                <ServerIcon className="mr-2 h-4 w-4" /> Bastion
                                            </Button>
                                            <Button variant="outline" size="sm" onClick={() => openEdit(c)}>
                                                <PencilIcon className="mr-2 h-4 w-4" /> Edit
                                            </Button>
                                            <Button variant="destructive" size="sm" onClick={() => deleteCluster(c.id)}>
                                                <TrashIcon className="mr-2 h-4 w-4" /> Delete
                                            </Button>
                                        </div>
                                    </TableCell>
                                </TableRow>
                            ))}

                            {clusters.length === 0 && (
                                <TableRow>
                                    <TableCell colSpan={6} className="text-muted-foreground py-10 text-center">
                                        No clusters yet.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </div>
            </div>

            {/* Edit cluster */}
            <Dialog open={editOpen} onOpenChange={(o) => !o && setEditOpen(false)}>
                <DialogContent className="sm:max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>
                            Edit Cluster
                            {editTarget ? (
                                <span className="text-muted-foreground ml-2 font-mono text-sm">({editTarget.name})</span>
                            ) : null}
                        </DialogTitle>
                    </DialogHeader>

                    <Form {...editForm}>
                        <form onSubmit={editForm.handleSubmit(submitEdit)} className="grid gap-4 md:grid-cols-2">
                            <FormField
                                control={editForm.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>New name (optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder={editTarget?.name || "e.g. my-eks-prod"} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={editForm.control}
                                name="provider"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>New provider (optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder={editTarget?.provider || "aws"} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={editForm.control}
                                name="region"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>New region (optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder={editTarget?.region || "eu-west-1"} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={editForm.control}
                                name="status"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>New status (optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder={editTarget?.status || "pending|ready|error"} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={editForm.control}
                                name="bastion_server_id"
                                render={({ field }) => (
                                    <FormItem className="md:col-span-2">
                                        <FormLabel>Replace/clear bastion (optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder="paste new server UUID or leave blank to clear" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={editForm.control}
                                name="cluster_load_balancer"
                                render={({ field }) => (
                                    <FormItem className="md:col-span-2">
                                        <FormLabel>Cluster Load Balancer (optional)</FormLabel>
                                        <FormControl>
                                            <Textarea placeholder={editTarget?.cluster_load_balancer || "e.g. JSON or URL or ARN"} rows={2} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <FormField
                                control={editForm.control}
                                name="control_load_balancer"
                                render={({ field }) => (
                                    <FormItem className="md:col-span-2">
                                        <FormLabel>Control Load Balancer (optional)</FormLabel>
                                        <FormControl>
                                            <Textarea placeholder={editTarget?.control_load_balancer || "e.g. JSON or URL or ARN"} rows={2} {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            {/* KubeConfig */}
                            <FormField
                                control={editForm.control}
                                name="kubeconfig"
                                render={({field}) => (
                                    <FormItem className='md:colspan-2'>
                                        <FormLabel>Replace Kubeconfig (optional)</FormLabel>
                                        <FormControl>
                                            <Textarea
                                                placeholder="Paste NEW kubeconfig YAML to replace the stored one. Leave empty for no change."
                                                rows={8}
                                                className="font-mono"
                                                {...field}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <DialogFooter className="md:col-span-2 gap-2">
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

            {/* Manage node pools */}
            <Dialog open={!!managePoolsTarget} onOpenChange={(o) => !o && setManagePoolsTarget(null)}>
                <DialogContent className="sm:max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>
                            Manage node pools for <span className="font-mono">{managePoolsTarget?.name}</span>
                        </DialogTitle>
                    </DialogHeader>

                    {/* Attached */}
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
                                    {(managePoolsTarget?.node_pools || []).map((p) => (
                                        <TableRow key={p.id}>
                                            <TableCell className="font-medium">{p.name}</TableCell>
                                            <TableCell>
                                                <div className="flex justify-end">
                                                    <Button variant="destructive" size="sm" onClick={() => detachPool(managePoolsTarget!.id, p.id)}>
                                                        <UnlinkIcon className="mr-2 h-4 w-4" /> Detach
                                                    </Button>
                                                </div>
                                            </TableCell>
                                        </TableRow>
                                    ))}

                                    {(managePoolsTarget?.node_pools || []).length === 0 && (
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
                                                    <div className="text-muted-foreground p-2 text-sm">No more node pools available to attach</div>
                                                )}
                                                {attachablePools.map((p) => {
                                                    const checked = field.value?.includes(p.id) || false
                                                    return (
                                                        <label key={p.id} className="hover:bg-accent flex cursor-pointer items-start gap-2 rounded p-1">
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
                                                                <div className="text-muted-foreground text-xs">{truncateMiddle(p.id, 8)}</div>
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

            {/* Manage bastion */}
            <Dialog open={!!manageBastionTarget} onOpenChange={(o) => !o && setManageBastionTarget(null)}>
                <DialogContent className="sm:max-w-lg">
                    <DialogHeader>
                        <DialogTitle>
                            Manage bastion for <span className="font-mono">{manageBastionTarget?.name}</span>
                        </DialogTitle>
                    </DialogHeader>

                    <div className="space-y-2">
                        <div className="text-sm font-medium">Current</div>
                        <div className="rounded-xl border p-3 text-sm">
                            {manageBastionTarget?.bastion_server ? (
                                <div>
                                    <div className="font-medium">{manageBastionTarget.bastion_server.hostname}</div>
                                    <div className="text-muted-foreground">{manageBastionTarget.bastion_server.ip}</div>
                                </div>
                            ) : (
                                <div className="text-muted-foreground">None</div>
                            )}
                        </div>
                    </div>

                    <Form {...setBastionForm}>
                        <form onSubmit={setBastionForm.handleSubmit(submitSetBastion)} className="space-y-4">
                            <FormField
                                control={setBastionForm.control}
                                name="server_id"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>New bastion server (UUID)</FormLabel>
                                        <FormControl>
                                            <Input placeholder="paste server UUID" {...field} />
                                        </FormControl>
                                        {bastionCandidates.length > 0 && (
                                            <div className="text-xs text-muted-foreground">
                                                Suggestions:
                                                <div className="mt-1 flex flex-wrap gap-2">
                                                    {bastionCandidates.slice(0, 8).map((s) => (
                                                        <Button
                                                            key={s.id}
                                                            type="button"
                                                            size="sm"
                                                            variant={field.value === s.id ? "default" : "outline"}
                                                            onClick={() => field.onChange(s.id)}
                                                            className="font-normal"
                                                        >
                                                            <ServerIcon className="mr-1 h-3 w-3" /> {s.hostname || truncateMiddle(s.id, 6)}
                                                        </Button>
                                                    ))}
                                                </div>
                                            </div>
                                        )}
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />

                            <DialogFooter className="gap-2">
                                <Button type="button" variant="secondary" onClick={() => void clearBastion()}>
                                    Clear bastion
                                </Button>
                                <Button type="submit" disabled={setBastionForm.formState.isSubmitting}>
                                    {setBastionForm.formState.isSubmitting ? "Saving…" : "Set bastion"}
                                </Button>
                            </DialogFooter>
                        </form>
                    </Form>
                </DialogContent>
            </Dialog>
            <pre>{JSON.stringify(clusters, null, 2)}</pre>
        </div>
    )
}
