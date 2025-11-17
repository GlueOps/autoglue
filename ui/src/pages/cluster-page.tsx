;
// src/pages/ClustersPage.tsx

import { useEffect, useMemo, useState } from "react";
import { clustersApi } from "@/api/clusters";
import { dnsApi } from "@/api/dns";
import { loadBalancersApi } from "@/api/loadbalancers";
import { nodePoolsApi } from "@/api/node_pools";
import { serversApi } from "@/api/servers";
import type { DtoClusterResponse, DtoDomainResponse, DtoLoadBalancerResponse, DtoNodePoolResponse, DtoRecordSetResponse, DtoServerResponse } from "@/sdk";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { AlertCircle, CheckCircle2, CircleSlash2, FileCode2, Globe2, Loader2, MapPin, Pencil, Plus, Search, Server, Wrench } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";



import { truncateMiddle } from "@/lib/utils";
import { Badge } from "@/components/ui/badge.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";





// --- Schemas ---

const createClusterSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
  provider: z.string().trim().min(1, "Provider is required").max(120, "Max 120 chars"),
  region: z.string().trim().min(1, "Region is required").max(120, "Max 120 chars"),
})
type CreateClusterInput = z.input<typeof createClusterSchema>

const updateClusterSchema = createClusterSchema.partial()
type UpdateClusterValues = z.infer<typeof updateClusterSchema>

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

function ClusterSummary({ c }: { c: DtoClusterResponse }) {
  return (
    <div className="flex flex-col gap-1 text-xs text-muted-foreground">
      <div className="flex flex-wrap items-center gap-2">
        {c.provider && (
          <span className="inline-flex items-center gap-1">
            <Globe2 className="h-3 w-3" />
            {c.provider}
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

  // Config dialog state
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
    queryFn: () => clustersApi.listClusters(),
  })

  const lbsQ = useQuery({
    queryKey: ["load-balancers"],
    queryFn: () => loadBalancersApi.listLoadBalancers(),
  })

  const domainsQ = useQuery({
    queryKey: ["domains"],
    queryFn: () => dnsApi.listDomains(),
  })

  // record sets fetched per captain domain
  const recordSetsQ = useQuery({
    queryKey: ["record-sets", captainDomainId],
    enabled: !!captainDomainId,
    queryFn: () => dnsApi.listRecordSetsByDomain(captainDomainId),
  })

  const serversQ = useQuery({
    queryKey: ["servers"],
    queryFn: () => serversApi.listServers(),
  })

  const npQ = useQuery({
    queryKey: ["node-pools"],
    queryFn: () => nodePoolsApi.listNodePools(),
  })

  // --- Create ---

  const createForm = useForm<CreateClusterInput>({
    resolver: zodResolver(createClusterSchema),
    defaultValues: {
      name: "",
      provider: "",
      region: "",
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
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while creating the cluster")
    },
  })

  const onCreateSubmit = (values: CreateClusterInput) => {
    createMut.mutate(values)
  }

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
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while updating the cluster")
    },
  })

  const openEdit = (cluster: DtoClusterResponse) => {
    if (!cluster.id) return
    setEditingId(cluster.id)
    updateForm.reset({
      name: cluster.name ?? "",
      provider: cluster.provider ?? "",
      region: cluster.region ?? "",
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
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while deleting the cluster")
    },
  })

  // --- Filter ---

  const filtered = useMemo(() => {
    const data: DtoClusterResponse[] = clustersQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((c) => {
        return (
          c.name?.toLowerCase().includes(q) ||
          c.provider?.toLowerCase().includes(q) ||
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

    // Prefill IDs from current attachments
    if (configCluster.captain_domain?.id) {
      setCaptainDomainId(configCluster.captain_domain.id)
    }
    if (configCluster.control_plane_record_set?.id) {
      setRecordSetId(configCluster.control_plane_record_set.id)
    }
    if (configCluster.apps_load_balancer?.id) {
      setAppsLbId(configCluster.apps_load_balancer.id)
    }
    if (configCluster.glueops_load_balancer?.id) {
      setGlueopsLbId(configCluster.glueops_load_balancer.id)
    }
    if (configCluster.bastion_server?.id) {
      setBastionId(configCluster.bastion_server.id)
    }
  }, [configCluster])

  async function refreshConfigCluster() {
    if (!configCluster?.id) return
    try {
      const updated = await clustersApi.getCluster(configCluster.id)
      setConfigCluster(updated)
      await qc.invalidateQueries({ queryKey: ["clusters"] })
    } catch {
      // ignore
    }
  }

  async function handleAttachCaptain() {
    if (!configCluster?.id) return
    if (!captainDomainId) {
      toast.error("Domain is required")
      return
    }
    setBusyKey("captain")
    try {
      await clustersApi.attachCaptainDomain(configCluster.id, {
        domain_id: captainDomainId,
      })
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
    if (!recordSetId) {
      toast.error("Record set is required")
      return
    }
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
    if (!appsLbId) {
      toast.error("Load balancer is required")
      return
    }
    setBusyKey("apps-lb")
    try {
      await clustersApi.attachAppsLoadBalancer(configCluster.id, {
        load_balancer_id: appsLbId,
      })
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
    if (!glueopsLbId) {
      toast.error("Load balancer is required")
      return
    }
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
    if (!bastionId) {
      toast.error("Server is required")
      return
    }
    setBusyKey("bastion")
    try {
      await clustersApi.attachBastion(configCluster.id, {
        server_id: bastionId,
      })
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
    if (!nodePoolId) {
      toast.error("Node pool is required")
      return
    }
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
    if (!kubeconfigText.trim()) {
      toast.error("Kubeconfig is required")
      return
    }
    setBusyKey("kubeconfig")
    try {
      await clustersApi.setKubeconfig(configCluster.id, {
        kubeconfig: kubeconfigText,
      })
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
                <form className="space-y-4" onSubmit={createForm.handleSubmit(onCreateSubmit)}>
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
                    name="provider"
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
                <TableHead>Summary</TableHead>
                <TableHead className="w-[320px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((c: DtoClusterResponse) => (
                <TableRow key={c.id}>
                  <TableCell className="font-medium">{c.name}</TableCell>
                  <TableCell>{c.provider}</TableCell>
                  <TableCell>{c.region}</TableCell>
                  <TableCell>
                    <StatusBadge status={c.status} />
                    {c.last_error && (
                      <div className="mt-1 text-xs text-red-500">
                        {truncateMiddle(c.last_error, 80)}
                      </div>
                    )}
                  </TableCell>
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
                  <TableCell colSpan={6} className="text-muted-foreground py-10 text-center">
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
                name="provider"
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

      {/* Configure dialog (attachments + kubeconfig + node pools) */}
      <Dialog open={!!configCluster} onOpenChange={(open) => !open && setConfigCluster(null)}>
        <DialogContent className="max-h-[90vh] w-full max-w-3xl overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              Configure Cluster{configCluster?.name ? `: ${configCluster.name}` : ""}
            </DialogTitle>
          </DialogHeader>

          {configCluster && (
            <div className="space-y-6 py-2">
              {/* Kubeconfig */}
              <section className="space-y-2 rounded-xl border p-4">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <div className="flex items-center gap-2">
                      <FileCode2 className="h-4 w-4" />
                      <h3 className="font-semibold text-sm">Kubeconfig</h3>
                    </div>
                    <p className="text-muted-foreground text-xs">
                      Paste the kubeconfig for this cluster. It will be stored encrypted and never
                      returned by the API.
                    </p>
                  </div>
                </div>

                <Textarea
                  value={kubeconfigText}
                  onChange={(e) => setKubeconfigText(e.target.value)}
                  rows={6}
                  placeholder="apiVersion: v1&#10;clusters:&#10;  - cluster: ..."
                  className="font-mono text-xs"
                />

                <div className="flex flex-wrap gap-2">
                  <Button
                    size="sm"
                    onClick={handleSetKubeconfig}
                    disabled={isBusy("kubeconfig")}
                  >
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
                    <h3 className="font-semibold text-sm">Captain Domain</h3>
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
                          placeholder={
                            domainsQ.isLoading ? "Loading domains…" : "Select domain"
                          }
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

              {/* Control Plane Record Set (shown once we have a captainDomainId) */}
              {captainDomainId && (
                <section className="space-y-2 rounded-xl border p-4">
                  <div className="flex items-center justify-between gap-2">
                    <div>
                      <h3 className="font-semibold text-sm">Control Plane Record Set</h3>
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
                              recordSetsQ.isLoading
                                ? "Loading record sets…"
                                : "Select record set"
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
                        disabled={
                          isBusy("recordset") || !configCluster.control_plane_record_set
                        }
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
                    <h3 className="font-semibold text-sm">Apps Load Balancer</h3>
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
                    <h3 className="font-semibold text-sm">GlueOps / Control-plane Load Balancer</h3>
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
                      <h3 className="font-semibold text-sm">Bastion Server</h3>
                    </div>
                    <p className="text-muted-foreground text-xs">
                      SSH bastion used to reach the cluster nodes.
                    </p>
                  </div>
                  <div className="text-right text-xs">
                    <div className="font-mono">
                      {configCluster.bastion_server
                        ? configCluster.bastion_server.hostname ??
                        configCluster.bastion_server.id
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
                          placeholder={
                            serversQ.isLoading ? "Loading servers…" : "Select server"
                          }
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
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <h3 className="font-semibold text-sm">Node Pools</h3>
                    <p className="text-muted-foreground text-xs">
                      Attach node pools to this cluster. Each node pool may have its own labels,
                      taints, and backing servers.
                    </p>
                  </div>
                </div>

                <div className="flex flex-col gap-2 md:flex-row md:items-end">
                  <div className="flex-1">
                    <Label className="text-xs">Available Node Pools</Label>
                    <Select value={nodePoolId} onValueChange={(val) => setNodePoolId(val)}>
                      <SelectTrigger className="w-full">
                        <SelectValue
                          placeholder={
                            npQ.isLoading ? "Loading node pools…" : "Select node pool"
                          }
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
