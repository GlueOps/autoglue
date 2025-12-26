import { useEffect, useMemo, useState } from "react";
import { actionsApi } from "@/api/actions";
import { clustersApi } from "@/api/clusters";
import { dnsApi } from "@/api/dns";
import { loadBalancersApi } from "@/api/loadbalancers";
import { nodePoolsApi } from "@/api/node_pools";
import { serversApi } from "@/api/servers";
import type { DtoActionResponse, DtoClusterResponse, DtoClusterRunResponse, DtoDomainResponse, DtoLoadBalancerResponse, DtoNodePoolResponse, DtoRecordSetResponse, DtoServerResponse } from "@/sdk";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, CheckCircle2, CircleSlash2, FileCode2, Globe2, Loader2, MapPin, Pencil, Plus, Search, Server, Wrench } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";



import { truncateMiddle } from "@/lib/utils";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Textarea } from "@/components/ui/textarea";





// --- Schemas ---

const createClusterSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
  cluster_provider: z.string().trim().min(1, "Provider is required").max(120, "Max 120 chars"),
  region: z.string().trim().min(1, "Region is required").max(120, "Max 120 chars"),
  docker_image: z.string().trim().min(1, "Docker Image is required"),
  docker_tag: z.string().trim().min(1, "Docker Tag is required"),
})
type CreateClusterInput = z.input<typeof createClusterSchema>

const updateClusterSchema = createClusterSchema.partial()
type UpdateClusterValues = z.infer<typeof updateClusterSchema>

// --- Data normalization helpers (fixes rows.some is not a function) ---

function asArray<T>(res: any): T[] {
  if (Array.isArray(res)) return res as T[]
  if (Array.isArray(res?.data)) return res.data as T[]
  if (Array.isArray(res?.body)) return res.body as T[]
  if (Array.isArray(res?.result)) return res.result as T[]
  return []
}

function asObject<T>(res: any): T {
  // for get endpoints that might return {data: {...}}
  if (res?.data && typeof res.data === "object") return res.data as T
  return res as T
}

// --- UI helpers ---

function StatusBadge({ status }: { status?: string | null }) {
  const value = (status ?? "").toLowerCase()

  if (!value) {
    return (
      <Badge variant="outline" className="text-xs">
        unknown
      </Badge>
    )
  }

  if (value === "ready") {
    return (
      <Badge variant="default" className="flex items-center gap-1 text-xs">
        <CheckCircle2 className="h-3 w-3" />
        ready
      </Badge>
    )
  }

  if (value === "failed") {
    return (
      <Badge variant="destructive" className="flex items-center gap-1 text-xs">
        <AlertCircle className="h-3 w-3" />
        failed
      </Badge>
    )
  }

  if (value === "provisioning" || value === "pending" || value === "pre_pending") {
    return (
      <Badge variant="secondary" className="flex items-center gap-1 text-xs">
        <Loader2 className="h-3 w-3 animate-spin" />
        {value.replace("_", " ")}
      </Badge>
    )
  }

  if (value === "incomplete") {
    return (
      <Badge variant="outline" className="flex items-center gap-1 text-xs">
        <AlertCircle className="h-3 w-3" />
        incomplete
      </Badge>
    )
  }

  return (
    <Badge variant="outline" className="text-xs">
      {value}
    </Badge>
  )
}

function RunStatusBadge({ status }: { status?: string | null }) {
  const s = (status ?? "").toLowerCase()

  if (!s)
    return (
      <Badge variant="outline" className="text-xs">
        unknown
      </Badge>
    )

  if (s === "succeeded" || s === "success") {
    return (
      <Badge variant="default" className="flex items-center gap-1 text-xs">
        <CheckCircle2 className="h-3 w-3" />
        succeeded
      </Badge>
    )
  }

  if (s === "failed" || s === "error") {
    return (
      <Badge variant="destructive" className="flex items-center gap-1 text-xs">
        <AlertCircle className="h-3 w-3" />
        failed
      </Badge>
    )
  }

  if (s === "queued" || s === "running") {
    return (
      <Badge variant="secondary" className="flex items-center gap-1 text-xs">
        <Loader2 className="h-3 w-3 animate-spin" />
        {s}
      </Badge>
    )
  }

  return (
    <Badge variant="outline" className="text-xs">
      {s}
    </Badge>
  )
}

function fmtTime(v: any): string {
  if (!v) return "-"
  try {
    const d = v instanceof Date ? v : new Date(v)
    if (Number.isNaN(d.getTime())) return "-"
    return d.toLocaleString()
  } catch {
    return "-"
  }
}

function ClusterSummary({ c }: { c: DtoClusterResponse }) {
  return (
    <div className="text-muted-foreground flex flex-col gap-1 text-xs">
      <div className="flex flex-wrap items-center gap-2">
        {c.cluster_provider && (
          <span className="inline-flex items-center gap-1">
            <Globe2 className="h-3 w-3" />
            {c.cluster_provider}
          </span>
        )}
        {c.region && (
          <span className="inline-flex items-center gap-1">
            <MapPin className="h-3 w-3" />
            {c.region}
          </span>
        )}
      </div>
      <div className="flex flex-wrap items-center gap-2 font-mono">
        {c.random_token && (
          <span>
            token: <span className="ml-1">{truncateMiddle(c.random_token, 8)}</span>
          </span>
        )}
        {c.certificate_key && (
          <span>
            cert: <span className="ml-1">{truncateMiddle(c.certificate_key, 8)}</span>
          </span>
        )}
      </div>
    </div>
  )
}

export const ClustersPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  // Configure dialog state
  const [configCluster, setConfigCluster] = useState<DtoClusterResponse | null>(null)

  const [captainDomainId, setCaptainDomainId] = useState("")
  const [recordSetId, setRecordSetId] = useState("")
  const [appsLbId, setAppsLbId] = useState("")
  const [glueopsLbId, setGlueopsLbId] = useState("")
  const [bastionId, setBastionId] = useState("")
  const [nodePoolId, setNodePoolId] = useState("")
  const [kubeconfigText, setKubeconfigText] = useState("")
  const [busyKey, setBusyKey] = useState<string | null>(null)

  const isBusy = (k: string) => busyKey === k

  const qc = useQueryClient()

  // --- Queries ---

  const clustersQ = useQuery({
    queryKey: ["clusters"],
    queryFn: async () => asArray<DtoClusterResponse>(await clustersApi.listClusters()),
  })

  const lbsQ = useQuery({
    queryKey: ["load-balancers"],
    queryFn: async () =>
      asArray<DtoLoadBalancerResponse>(await loadBalancersApi.listLoadBalancers()),
  })

  const domainsQ = useQuery({
    queryKey: ["domains"],
    queryFn: async () => asArray<DtoDomainResponse>(await dnsApi.listDomains()),
  })

  const recordSetsQ = useQuery({
    queryKey: ["record-sets", captainDomainId],
    enabled: !!captainDomainId,
    queryFn: async () =>
      asArray<DtoRecordSetResponse>(await dnsApi.listRecordSetsByDomain(captainDomainId)),
  })

  const serversQ = useQuery({
    queryKey: ["servers"],
    queryFn: async () => asArray<DtoServerResponse>(await serversApi.listServers()),
  })

  const npQ = useQuery({
    queryKey: ["node-pools"],
    queryFn: async () => asArray<DtoNodePoolResponse>(await nodePoolsApi.listNodePools()),
  })

  const actionsQ = useQuery({
    queryKey: ["actions"],
    queryFn: async () => asArray<DtoActionResponse>(await actionsApi.listActions()),
  })

  const runsQ = useQuery({
    queryKey: ["cluster-runs", configCluster?.id],
    enabled: !!configCluster?.id,
    queryFn: async () =>
      asArray<DtoClusterRunResponse>(await clustersApi.listClusterRuns(configCluster!.id!)),
    refetchInterval: (data) => {
      // IMPORTANT: data might not be array if queryFn isn't normalizing. But it is here anyway.
      const rows = Array.isArray(data) ? data : []
      const active = rows.some((r: any) => {
        const s = String(r?.status ?? "").toLowerCase()
        return s === "queued" || s === "running"
      })
      return active ? 2000 : false
    },
  })

  const actionLabelByTarget = useMemo(() => {
    const m = new Map<string, string>()
    ;(actionsQ.data ?? []).forEach((a) => {
      if (a.make_target) m.set(a.make_target, a.label ?? a.make_target)
    })
    return m
  }, [actionsQ.data])

  const runDisplayName = (r: DtoClusterRunResponse) =>
    actionLabelByTarget.get(r.action ?? "") ?? r.action ?? "unknown"

  // --- Create ---

  const createForm = useForm<CreateClusterInput>({
    resolver: zodResolver(createClusterSchema),
    defaultValues: {
      name: "",
      cluster_provider: "",
      region: "",
      docker_image: "",
      docker_tag: "",
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateClusterInput) => clustersApi.createCluster(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["clusters"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Cluster created successfully.")
    },
    onError: (err: any) =>
      toast.error(err?.message ?? "There was an error while creating the cluster"),
  })

  // --- Update basic details ---

  const updateForm = useForm<UpdateClusterValues>({
    resolver: zodResolver(updateClusterSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateClusterValues }) =>
      clustersApi.updateCluster(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["clusters"] })
      updateForm.reset()
      setUpdateOpen(false)
      toast.success("Cluster updated successfully.")
    },
    onError: (err: any) =>
      toast.error(err?.message ?? "There was an error while updating the cluster"),
  })

  const openEdit = (cluster: DtoClusterResponse) => {
    if (!cluster.id) return
    setEditingId(cluster.id)
    updateForm.reset({
      name: cluster.name ?? "",
      cluster_provider: cluster.cluster_provider ?? "",
      region: cluster.region ?? "",
      docker_image: cluster.docker_image ?? "",
      docker_tag: cluster.docker_tag ?? "",
    })
    setUpdateOpen(true)
  }

  // --- Delete ---

  const deleteMut = useMutation({
    mutationFn: (id: string) => clustersApi.deleteCluster(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["clusters"] })
      setDeleteId(null)
      toast.success("Cluster deleted successfully.")
    },
    onError: (err: any) =>
      toast.error(err?.message ?? "There was an error while deleting the cluster"),
  })

  // --- Run Action ---

  const runActionMut = useMutation({
    mutationFn: ({ clusterID, actionID }: { clusterID: string; actionID: string }) =>
      clustersApi.runClusterAction(clusterID, actionID),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["cluster-runs", configCluster?.id] })
      toast.success("Action enqueued.")
    },
    onError: (err: any) => toast.error(err?.message ?? "Failed to enqueue action."),
  })

  async function handleRunAction(actionID: string) {
    if (!configCluster?.id) return
    setBusyKey(`run:${actionID}`)
    try {
      await runActionMut.mutateAsync({ clusterID: configCluster.id, actionID })
    } finally {
      setBusyKey(null)
    }
  }

  // --- Filter ---

  const filtered = useMemo(() => {
    const data: DtoClusterResponse[] = clustersQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((c) => {
          return (
            c.name?.toLowerCase().includes(q) ||
            c.cluster_provider?.toLowerCase().includes(q) ||
            c.region?.toLowerCase().includes(q) ||
            c.status?.toLowerCase().includes(q)
          )
        })
      : data
  }, [filter, clustersQ.data])

  // --- Config dialog helpers ---

  useEffect(() => {
    if (!configCluster) {
      setCaptainDomainId("")
      setRecordSetId("")
      setAppsLbId("")
      setGlueopsLbId("")
      setBastionId("")
      setNodePoolId("")
      setKubeconfigText("")
      return
    }

    if (configCluster.captain_domain?.id) setCaptainDomainId(configCluster.captain_domain.id)
    if (configCluster.control_plane_record_set?.id)
      setRecordSetId(configCluster.control_plane_record_set.id)
    if (configCluster.apps_load_balancer?.id) setAppsLbId(configCluster.apps_load_balancer.id)
    if (configCluster.glueops_load_balancer?.id)
      setGlueopsLbId(configCluster.glueops_load_balancer.id)
    if (configCluster.bastion_server?.id) setBastionId(configCluster.bastion_server.id)
  }, [configCluster])

  async function refreshConfigCluster() {
    if (!configCluster?.id) return
    try {
      const updatedRaw = await clustersApi.getCluster(configCluster.id)
      const updated = asObject<DtoClusterResponse>(updatedRaw)
      setConfigCluster(updated)
      await qc.invalidateQueries({ queryKey: ["clusters"] })
      await qc.invalidateQueries({ queryKey: ["cluster-runs", configCluster.id] })
    } catch {
      // ignore
    }
  }

  async function handleAttachCaptain() {
    if (!configCluster?.id) return
    if (!captainDomainId) return toast.error("Domain is required")
    setBusyKey("captain")
    try {
      await clustersApi.attachCaptainDomain(configCluster.id, { domain_id: captainDomainId })
      toast.success("Captain domain attached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach captain domain.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachCaptain() {
    if (!configCluster?.id) return
    setBusyKey("captain")
    try {
      await clustersApi.detachCaptainDomain(configCluster.id)
      toast.success("Captain domain detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach captain domain.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleAttachRecordSet() {
    if (!configCluster?.id) return
    if (!recordSetId) return toast.error("Record set is required")
    setBusyKey("recordset")
    try {
      await clustersApi.attachControlPlaneRecordSet(configCluster.id, {
        record_set_id: recordSetId,
      })
      toast.success("Control plane record set attached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach record set.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachRecordSet() {
    if (!configCluster?.id) return
    setBusyKey("recordset")
    try {
      await clustersApi.detachControlPlaneRecordSet(configCluster.id)
      toast.success("Control plane record set detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach record set.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleAttachAppsLb() {
    if (!configCluster?.id) return
    if (!appsLbId) return toast.error("Load balancer is required")
    setBusyKey("apps-lb")
    try {
      await clustersApi.attachAppsLoadBalancer(configCluster.id, { load_balancer_id: appsLbId })
      toast.success("Apps load balancer attached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach apps load balancer.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachAppsLb() {
    if (!configCluster?.id) return
    setBusyKey("apps-lb")
    try {
      await clustersApi.detachAppsLoadBalancer(configCluster.id)
      toast.success("Apps load balancer detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach apps load balancer.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleAttachGlueopsLb() {
    if (!configCluster?.id) return
    if (!glueopsLbId) return toast.error("Load balancer is required")
    setBusyKey("glueops-lb")
    try {
      await clustersApi.attachGlueOpsLoadBalancer(configCluster.id, {
        load_balancer_id: glueopsLbId,
      })
      toast.success("GlueOps load balancer attached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach GlueOps load balancer.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachGlueopsLb() {
    if (!configCluster?.id) return
    setBusyKey("glueops-lb")
    try {
      await clustersApi.detachGlueOpsLoadBalancer(configCluster.id)
      toast.success("GlueOps load balancer detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach GlueOps load balancer.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleAttachBastion() {
    if (!configCluster?.id) return
    if (!bastionId) return toast.error("Server is required")
    setBusyKey("bastion")
    try {
      await clustersApi.attachBastion(configCluster.id, { server_id: bastionId })
      toast.success("Bastion server attached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach bastion server.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachBastion() {
    if (!configCluster?.id) return
    setBusyKey("bastion")
    try {
      await clustersApi.detachBastion(configCluster.id)
      toast.success("Bastion server detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach bastion server.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleAttachNodePool() {
    if (!configCluster?.id) return
    if (!nodePoolId) return toast.error("Node pool is required")
    setBusyKey("nodepool")
    try {
      await clustersApi.attachNodePool(configCluster.id, nodePoolId)
      toast.success("Node pool attached.")
      setNodePoolId("")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to attach node pool.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleDetachNodePool(npId: string) {
    if (!configCluster?.id) return
    setBusyKey("nodepool")
    try {
      await clustersApi.detachNodePool(configCluster.id, npId)
      toast.success("Node pool detached.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to detach node pool.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleSetKubeconfig() {
    if (!configCluster?.id) return
    if (!kubeconfigText.trim()) return toast.error("Kubeconfig is required")
    setBusyKey("kubeconfig")
    try {
      await clustersApi.setKubeconfig(configCluster.id, { kubeconfig: kubeconfigText })
      toast.success("Kubeconfig updated.")
      setKubeconfigText("")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to set kubeconfig.")
    } finally {
      setBusyKey(null)
    }
  }

  async function handleClearKubeconfig() {
    if (!configCluster?.id) return
    setBusyKey("kubeconfig")
    try {
      await clustersApi.clearKubeconfig(configCluster.id)
      toast.success("Kubeconfig cleared.")
      await refreshConfigCluster()
    } catch (err: any) {
      toast.error(err?.message ?? "Failed to clear kubeconfig.")
    } finally {
      setBusyKey(null)
    }
  }

  if (clustersQ.isLoading) return <div className="p-6">Loading clusters…</div>
  if (clustersQ.error) return <div className="p-6 text-red-500">Error loading clusters.</div>

  const allLbs: DtoLoadBalancerResponse[] = lbsQ.data ?? []
  const appsLbs = allLbs.filter((lb) => lb.kind === "public")
  const glueopsLbs = allLbs.filter((lb) => lb.kind === "glueops")

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Clusters</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search clusters"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Cluster
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create Cluster</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form
                  className="space-y-4"
                  onSubmit={createForm.handleSubmit((v) => createMut.mutate(v))}
                >
                  <FormField
                    control={createForm.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Name</FormLabel>
                        <FormControl>
                          <Input placeholder="prod-cluster-eu-west-1" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="cluster_provider"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Provider</FormLabel>
                        <FormControl>
                          <Input placeholder="aws / hetzner / baremetal" {...field} />
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
                    name="docker_image"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Docker Image</FormLabel>
                        <FormControl>
                          <Input placeholder="ghcr.io/glueops/gluekube" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="docker_tag"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Docker Tag</FormLabel>
                        <FormControl>
                          <Input placeholder="v1.33" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createMut.isPending}>
                      {createMut.isPending ? "Creating…" : "Create"}
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
                <TableHead>Provider</TableHead>
                <TableHead>Region</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Docker</TableHead>
                <TableHead>Summary</TableHead>
                <TableHead className="w-[320px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((c: DtoClusterResponse) => (
                <TableRow key={c.id}>
                  <TableCell className="font-medium">{c.name}</TableCell>
                  <TableCell>{c.cluster_provider}</TableCell>
                  <TableCell>{c.region}</TableCell>
                  <TableCell>
                    <StatusBadge status={c.status} />
                    {c.last_error && (
                      <div className="mt-1 text-xs text-red-500">
                        {truncateMiddle(c.last_error, 80)}
                      </div>
                    )}
                  </TableCell>
                  <TableCell>{(c.docker_image ?? "") + ":" + (c.docker_tag ?? "")}</TableCell>
                  <TableCell>
                    <ClusterSummary c={c} />
                    {c.id && (
                      <code className="text-muted-foreground mt-1 block text-xs">
                        {truncateMiddle(c.id, 6)}
                      </code>
                    )}
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-wrap justify-end gap-2">
                      <Button variant="ghost" size="sm" onClick={() => setConfigCluster(c)}>
                        <Wrench className="mr-1 h-4 w-4" /> Configure
                      </Button>
                      <Button variant="outline" size="sm" onClick={() => openEdit(c)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => c.id && setDeleteId(c.id)}
                        disabled={deleteMut.isPending && deleteId === c.id}
                      >
                        {deleteMut.isPending && deleteId === c.id ? "Deleting…" : "Delete"}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={7} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No clusters match your search.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Update dialog */}
      <Dialog open={updateOpen} onOpenChange={setUpdateOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Cluster</DialogTitle>
          </DialogHeader>

          <Form {...updateForm}>
            <form
              className="space-y-4"
              onSubmit={updateForm.handleSubmit((values) => {
                if (!editingId) return
                updateMut.mutate({ id: editingId, values })
              })}
            >
              <FormField
                control={updateForm.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="prod-cluster-eu-west-1" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="cluster_provider"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Provider</FormLabel>
                    <FormControl>
                      <Input placeholder="aws / hetzner / baremetal" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
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
                control={updateForm.control}
                name="docker_image"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Docker Image</FormLabel>
                    <FormControl>
                      <Input placeholder="ghcr.io/glueops/gluekube" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="docker_tag"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Docker Tag</FormLabel>
                    <FormControl>
                      <Input placeholder="v1.33" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button type="button" variant="outline" onClick={() => setUpdateOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateMut.isPending}>
                  {updateMut.isPending ? "Saving…" : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Configure dialog (attachments + kubeconfig + node pools + actions/runs) */}
      <Dialog open={!!configCluster} onOpenChange={(open) => !open && setConfigCluster(null)}>
        <DialogContent className="max-h-[90vh] w-full max-w-3xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              Configure Cluster{configCluster?.name ? `: ${configCluster.name}` : ""}
            </DialogTitle>
          </DialogHeader>

          {configCluster && (
            <div className="space-y-6 py-2">
              {/* Cluster Actions */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <div className="flex items-center gap-2">
                      <Wrench className="h-4 w-4" />
                      <h3 className="text-sm font-semibold">Cluster Actions</h3>
                    </div>
                    <p className="text-muted-foreground text-xs">
                      Run admin-configured actions on this cluster. Actions are executed
                      asynchronously.
                    </p>
                  </div>

                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => runsQ.refetch()}
                    disabled={runsQ.isFetching || !configCluster?.id}
                  >
                    {runsQ.isFetching ? "Refreshing…" : "Refresh runs"}
                  </Button>
                </div>

                <div className="space-y-2">
                  {actionsQ.isLoading ? (
                    <p className="text-muted-foreground text-xs">Loading actions…</p>
                  ) : (actionsQ.data ?? []).length === 0 ? (
                    <p className="text-muted-foreground text-xs">
                      No actions configured yet. Create actions in Admin → Actions.
                    </p>
                  ) : (
                    <div className="divide-border rounded-md border">
                      {(actionsQ.data ?? []).map((a: DtoActionResponse) => (
                        <div
                          key={a.id}
                          className="flex items-center justify-between gap-3 px-3 py-2"
                        >
                          <div className="flex min-w-0 flex-col">
                            <div className="flex items-center gap-2">
                              <span className="text-sm font-medium">{a.label}</span>
                              {a.make_target && (
                                <code className="text-muted-foreground text-xs">
                                  {a.make_target}
                                </code>
                              )}
                            </div>
                            {a.description && (
                              <p className="text-muted-foreground line-clamp-2 text-xs">
                                {a.description}
                              </p>
                            )}
                          </div>

                          <Button
                            size="sm"
                            onClick={() => a.id && handleRunAction(a.id)}
                            disabled={!a.id || isBusy(`run:${a.id}`)}
                          >
                            {a.id && isBusy(`run:${a.id}`) ? "Enqueueing…" : "Run"}
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                <div className="mt-3 space-y-1">
                  <Label className="text-xs">Recent Runs</Label>

                  {runsQ.isLoading ? (
                    <p className="text-muted-foreground text-xs">Loading runs…</p>
                  ) : (runsQ.data ?? []).length === 0 ? (
                    <p className="text-muted-foreground text-xs">No runs yet for this cluster.</p>
                  ) : (
                    <div className="overflow-x-auto rounded-md border">
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>Action</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead>Created</TableHead>
                            <TableHead>Finished</TableHead>
                            <TableHead>Error</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {(runsQ.data ?? []).slice(0, 20).map((r) => (
                            <TableRow key={r.id}>
                              <TableCell className="min-w-[220px]">
                                <div className="flex flex-col">
                                  <span className="text-sm font-medium">{runDisplayName(r)}</span>
                                  {r.id && (
                                    <code className="text-muted-foreground text-xs">
                                      {truncateMiddle(r.id, 8)}
                                    </code>
                                  )}
                                </div>
                              </TableCell>
                              <TableCell>
                                <RunStatusBadge status={r.status} />
                              </TableCell>
                              <TableCell className="text-xs">
                                {fmtTime((r as any).created_at)}
                              </TableCell>
                              <TableCell className="text-xs">
                                {fmtTime((r as any).finished_at)}
                              </TableCell>
                              <TableCell className="text-xs">
                                {r.error ? truncateMiddle(r.error, 80) : "-"}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </div>
                  )}
                </div>
              </section>

              {/* Kubeconfig */}
              <section className="space-y-2 rounded-xl border p-4">
                <div>
                  <div className="flex items-center gap-2">
                    <FileCode2 className="h-4 w-4" />
                    <h3 className="text-sm font-semibold">Kubeconfig</h3>
                  </div>
                  <p className="text-muted-foreground text-xs">
                    Paste the kubeconfig for this cluster. It will be stored encrypted and never
                    returned by the API.
                  </p>
                </div>

                <Textarea
                  value={kubeconfigText}
                  onChange={(e) => setKubeconfigText(e.target.value)}
                  rows={6}
                  placeholder={"apiVersion: v1\nclusters:\n  - cluster: ..."}
                  className="font-mono text-xs"
                />

                <div className="flex flex-wrap gap-2">
                  <Button size="sm" onClick={handleSetKubeconfig} disabled={isBusy("kubeconfig")}>
                    {isBusy("kubeconfig") ? "Saving…" : "Save kubeconfig"}
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={handleClearKubeconfig}
                    disabled={isBusy("kubeconfig")}
                  >
                    Clear kubeconfig
                  </Button>
                </div>
              </section>

              {/* Captain Domain */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <h3 className="text-sm font-semibold">Captain Domain</h3>
                    <p className="text-muted-foreground text-xs">
                      Domain used for the AutoGlue captain endpoint.
                    </p>
                  </div>
                  <div className="text-right text-xs">
                    <div className="font-mono">
                      {configCluster.captain_domain
                        ? configCluster.captain_domain.domain_name
                        : "Not attached"}
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">Domain</Label>
                    <Select
                      value={captainDomainId}
                      onValueChange={(val) => {
                        setCaptainDomainId(val)
                        setRecordSetId("")
                      }}
                    >
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={domainsQ.isLoading ? "Loading domains…" : "Select domain"}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {(domainsQ.data ?? []).map((d: DtoDomainResponse) => (
                          <SelectItem key={d.id!} value={d.id!}>
                            {d.domain_name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <p className="text-muted-foreground mt-1 text-[11px]">
                      Once a domain is attached, control plane record sets for that domain will be
                      available below.
                    </p>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleAttachCaptain}
                      disabled={isBusy("captain") || domainsQ.isLoading}
                    >
                      {isBusy("captain") ? "Attaching…" : "Attach"}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleDetachCaptain}
                      disabled={isBusy("captain") || !configCluster.captain_domain}
                    >
                      Detach
                    </Button>
                  </div>
                </div>
              </section>

              {/* Control Plane Record Set */}
              {captainDomainId && (
                <section className="space-y-2 rounded-xl border p-4">
                  <div className="flex items-center justify-between gap-2">
                    <div>
                      <h3 className="text-sm font-semibold">Control Plane Record Set</h3>
                      <p className="text-muted-foreground text-xs">
                        DNS record set used for the cluster control plane endpoint.
                      </p>
                    </div>
                    <div className="text-right text-xs">
                      <div className="font-mono">
                        {configCluster.control_plane_record_set
                          ? configCluster.control_plane_record_set.name
                          : "Not attached"}
                      </div>
                    </div>
                  </div>

                  <div className="flex flex-col gap-2 md:flex-row md:items-end">
                    <div className="flex-1">
                      <Label className="text-xs">Record Set</Label>
                      <Select value={recordSetId} onValueChange={(val) => setRecordSetId(val)}>
                        <SelectTrigger className="w-full">
                          <SelectValue
                            placeholder={
                              recordSetsQ.isLoading ? "Loading record sets…" : "Select record set"
                            }
                          />
                        </SelectTrigger>
                        <SelectContent>
                          {(recordSetsQ.data ?? []).map((rs: DtoRecordSetResponse) => (
                            <SelectItem key={rs.id!} value={rs.id!}>
                              {rs.name} · {rs.type}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        onClick={handleAttachRecordSet}
                        disabled={isBusy("recordset") || recordSetsQ.isLoading}
                      >
                        {isBusy("recordset") ? "Attaching…" : "Attach"}
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={handleDetachRecordSet}
                        disabled={isBusy("recordset") || !configCluster.control_plane_record_set}
                      >
                        Detach
                      </Button>
                    </div>
                  </div>
                </section>
              )}

              {/* Apps Load Balancer */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <h3 className="text-sm font-semibold">Apps Load Balancer</h3>
                    <p className="text-muted-foreground text-xs">
                      Frontend load balancer for application traffic.
                    </p>
                  </div>
                  <div className="text-right text-xs">
                    <div className="font-mono">
                      {configCluster.apps_load_balancer
                        ? configCluster.apps_load_balancer.name
                        : "Not attached"}
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">Apps Load Balancer</Label>
                    <Select value={appsLbId} onValueChange={(val) => setAppsLbId(val)}>
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={
                            lbsQ.isLoading ? "Loading load balancers…" : "Select apps LB"
                          }
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {appsLbs.map((lb) => (
                          <SelectItem key={lb.id!} value={lb.id!}>
                            {lb.name} · {lb.public_ip_address}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleAttachAppsLb}
                      disabled={isBusy("apps-lb") || lbsQ.isLoading}
                    >
                      {isBusy("apps-lb") ? "Attaching…" : "Attach"}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleDetachAppsLb}
                      disabled={isBusy("apps-lb") || !configCluster.apps_load_balancer}
                    >
                      Detach
                    </Button>
                  </div>
                </div>
              </section>

              {/* GlueOps Load Balancer */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <h3 className="text-sm font-semibold">GlueOps / Control-plane Load Balancer</h3>
                    <p className="text-muted-foreground text-xs">
                      Load balancer for GlueOps/control-plane traffic.
                    </p>
                  </div>
                  <div className="text-right text-xs">
                    <div className="font-mono">
                      {configCluster.glueops_load_balancer
                        ? configCluster.glueops_load_balancer.name
                        : "Not attached"}
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">GlueOps Load Balancer</Label>
                    <Select value={glueopsLbId} onValueChange={(val) => setGlueopsLbId(val)}>
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={
                            lbsQ.isLoading ? "Loading load balancers…" : "Select GlueOps LB"
                          }
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {glueopsLbs.map((lb) => (
                          <SelectItem key={lb.id!} value={lb.id!}>
                            {lb.name} · {lb.private_ip_address}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleAttachGlueopsLb}
                      disabled={isBusy("glueops-lb") || lbsQ.isLoading}
                    >
                      {isBusy("glueops-lb") ? "Attaching…" : "Attach"}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleDetachGlueopsLb}
                      disabled={isBusy("glueops-lb") || !configCluster.glueops_load_balancer}
                    >
                      Detach
                    </Button>
                  </div>
                </div>
              </section>

              {/* Bastion Server */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <div className="flex items-center gap-2">
                      <Server className="h-4 w-4" />
                      <h3 className="text-sm font-semibold">Bastion Server</h3>
                    </div>
                    <p className="text-muted-foreground text-xs">
                      SSH bastion used to reach the cluster nodes.
                    </p>
                  </div>
                  <div className="text-right text-xs">
                    <div className="font-mono">
                      {configCluster.bastion_server
                        ? (configCluster.bastion_server.hostname ?? configCluster.bastion_server.id)
                        : "Not attached"}
                    </div>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">Bastion Server</Label>
                    <Select value={bastionId} onValueChange={(val) => setBastionId(val)}>
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={serversQ.isLoading ? "Loading servers…" : "Select server"}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {(serversQ.data ?? []).map((s: DtoServerResponse) => (
                          <SelectItem key={s.id!} value={s.id!}>
                            {s.hostname ?? s.id} · {s.private_ip_address}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleAttachBastion}
                      disabled={isBusy("bastion") || serversQ.isLoading}
                    >
                      {isBusy("bastion") ? "Attaching…" : "Attach"}
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={handleDetachBastion}
                      disabled={isBusy("bastion") || !configCluster.bastion_server}
                    >
                      Detach
                    </Button>
                  </div>
                </div>
              </section>

              {/* Node Pools */}
              <section className="space-y-2 rounded-xl border p-4">
                <div>
                  <h3 className="text-sm font-semibold">Node Pools</h3>
                  <p className="text-muted-foreground text-xs">
                    Attach node pools to this cluster. Each node pool may have its own labels,
                    taints, and backing servers.
                  </p>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">Available Node Pools</Label>
                    <Select value={nodePoolId} onValueChange={(val) => setNodePoolId(val)}>
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={npQ.isLoading ? "Loading node pools…" : "Select node pool"}
                        />
                      </SelectTrigger>
                      <SelectContent>
                        {(npQ.data ?? []).map((np: DtoNodePoolResponse) => (
                          <SelectItem key={np.id!} value={np.id!}>
                            {np.name} · {np.role}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleAttachNodePool}
                      disabled={isBusy("nodepool") || npQ.isLoading}
                    >
                      {isBusy("nodepool") ? "Attaching…" : "Attach"}
                    </Button>
                  </div>
                </div>

                <div className="mt-3 space-y-1">
                  <Label className="text-xs">Attached Node Pools</Label>
                  {configCluster.node_pools && configCluster.node_pools.length > 0 ? (
                    <div className="divide-border mt-1 rounded-md border">
                      {configCluster.node_pools.map((np: DtoNodePoolResponse) => (
                        <div
                          key={np.id}
                          className="flex items-center justify-between gap-3 px-3 py-2 text-xs"
                        >
                          <div className="flex flex-col">
                            <span className="font-medium">{np.name}</span>
                            <span className="text-muted-foreground">
                              role: {np.role} · servers: {np.servers?.length ?? 0}
                            </span>
                          </div>
                          <Button
                            size="sm"
                            variant="outline"
                            onClick={() => np.id && handleDetachNodePool(np.id)}
                            disabled={isBusy("nodepool")}
                          >
                            Detach
                          </Button>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted-foreground mt-1 text-xs">
                      No node pools attached to this cluster yet.
                    </p>
                  )}
                </div>
              </section>

              <DialogFooter className="mt-2">
                <Button variant="outline" onClick={() => setConfigCluster(null)}>
                  Close
                </Button>
              </DialogFooter>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Delete confirm dialog */}
      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Delete cluster</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This action cannot be undone. Are you sure you want to delete this cluster?
          </p>
          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => deleteId && deleteMut.mutate(deleteId)}
              disabled={deleteMut.isPending}
            >
              {deleteMut.isPending ? "Deleting…" : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
