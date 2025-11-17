import { useMemo, useState } from "react"
import { loadBalancersApi } from "@/api/loadbalancers"
import type { DtoLoadBalancerResponse } from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { CircleSlash2, Network, Pencil, Plus, Search } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { truncateMiddle } from "@/lib/utils"
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

// --- schemas ---

const createLoadBalancerSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(120, "Max 120 chars"),
  kind: z.enum(["glueops", "public"]).default("public"),
  public_ip_address: z
    .string()
    .trim()
    .min(1, "Public IP/hostname is required")
    .max(255, "Max 255 chars"),
  private_ip_address: z
    .string()
    .trim()
    .min(1, "Private IP/hostname is required")
    .max(255, "Max 255 chars"),
})
type CreateLoadBalancerInput = z.input<typeof createLoadBalancerSchema>

const updateLoadBalancerSchema = createLoadBalancerSchema.partial()
type UpdateLoadBalancerValues = z.infer<typeof updateLoadBalancerSchema>

// --- badge ---

function LoadBalancerBadge({
  lb,
}: {
  lb: Pick<DtoLoadBalancerResponse, "name" | "kind" | "public_ip_address" | "private_ip_address">
}) {
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      <Network className="mr-1 h-3 w-3" />
      {lb.name} · {lb.kind} · {lb.public_ip_address} → {lb.private_ip_address}
    </Badge>
  )
}

export const LoadBalancersPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  const lbsQ = useQuery({
    queryKey: ["loadBalancers"],
    queryFn: () => loadBalancersApi.listLoadBalancers(),
  })

  // --- Create ---

  const createForm = useForm<CreateLoadBalancerInput>({
    resolver: zodResolver(createLoadBalancerSchema),
    defaultValues: {
      name: "",
      kind: "public",
      public_ip_address: "",
      private_ip_address: "",
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateLoadBalancerInput) => loadBalancersApi.createLoadBalancer(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["loadBalancers"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Load balancer created successfully.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while creating the load balancer")
    },
  })

  const onCreateSubmit = (values: CreateLoadBalancerInput) => {
    createMut.mutate(values)
  }

  // --- Update ---

  const updateForm = useForm<UpdateLoadBalancerValues>({
    resolver: zodResolver(updateLoadBalancerSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateLoadBalancerValues }) =>
      loadBalancersApi.updateLoadBalancer(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["loadBalancers"] })
      updateForm.reset()
      setUpdateOpen(false)
      toast.success("Load balancer updated successfully.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while updating the load balancer")
    },
  })

  const openEdit = (lb: DtoLoadBalancerResponse) => {
    setEditingId(lb.id!)
    updateForm.reset({
      name: lb.name ?? "",
      kind: (lb.kind as "public" | "glueops") ?? "public",
      public_ip_address: lb.public_ip_address ?? "",
      private_ip_address: lb.private_ip_address ?? "",
    })
    setUpdateOpen(true)
  }

  // --- Delete ---

  const deleteMut = useMutation({
    mutationFn: (id: string) => loadBalancersApi.deleteLoadBalancer(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["loadBalancers"] })
      setDeleteId(null)
      toast.success("Load balancer deleted successfully.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "There was an error while deleting the load balancer")
    },
  })

  // --- Filter ---

  const filtered = useMemo(() => {
    const data = lbsQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((lb: any) => {
          return (
            lb.name?.toLowerCase().includes(q) ||
            lb.kind?.toLowerCase().includes(q) ||
            lb.public_ip_address?.toLowerCase().includes(q) ||
            lb.private_ip_address?.toLowerCase().includes(q)
          )
        })
      : data
  }, [filter, lbsQ.data])

  if (lbsQ.isLoading) return <div className="p-6">Loading load balancers…</div>
  if (lbsQ.error) return <div className="p-6 text-red-500">Error loading load balancers.</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Load Balancers</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search load balancers"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Load Balancer
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create Load Balancer</DialogTitle>
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
                          <Input placeholder="apps-lb-01" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="kind"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Kind</FormLabel>
                        <FormControl>
                          <Select onValueChange={field.onChange} value={field.value ?? "public"}>
                            <SelectTrigger>
                              <SelectValue placeholder="Select kind" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="public">Public</SelectItem>
                              <SelectItem value="glueops">GlueOps</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="public_ip_address"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Public IP</FormLabel>
                        <FormControl>
                          <Input placeholder="1.2.3.4" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="private_ip_address"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Private IP</FormLabel>
                        <FormControl>
                          <Input placeholder="10.0.30.10" {...field} />
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
                <TableHead>Kind</TableHead>
                <TableHead>Public IP / Hostname</TableHead>
                <TableHead>Private IP / Hostname</TableHead>
                <TableHead>Summary</TableHead>
                <TableHead className="w-[220px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((lb: DtoLoadBalancerResponse) => (
                <TableRow key={lb.id}>
                  <TableCell>{lb.name}</TableCell>
                  <TableCell>{lb.kind}</TableCell>
                  <TableCell>{lb.public_ip_address}</TableCell>
                  <TableCell>{lb.private_ip_address}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <LoadBalancerBadge lb={lb} />
                      {lb.id && (
                        <code className="text-muted-foreground text-xs">
                          {truncateMiddle(lb.id, 6)}
                        </code>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => openEdit(lb)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => setDeleteId(lb.id!)}
                        disabled={deleteMut.isPending && deleteId === lb.id}
                      >
                        {deleteMut.isPending && deleteId === lb.id ? "Deleting…" : "Delete"}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={6} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No load balancers match your search.
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
            <DialogTitle>Edit Load Balancer</DialogTitle>
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
                      <Input placeholder="apps-lb-01" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="kind"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Kind</FormLabel>
                    <FormControl>
                      <Select onValueChange={field.onChange} value={field.value ?? "public"}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select kind" />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="public">Public</SelectItem>
                          <SelectItem value="glueops">GlueOps</SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="public_ip_address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Public IP / Hostname</FormLabel>
                    <FormControl>
                      <Input placeholder="1.2.3.4 or apps.example.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="private_ip_address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Private IP / Hostname</FormLabel>
                    <FormControl>
                      <Input placeholder="10.0.30.10" {...field} />
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

      {/* Delete confirm dialog */}
      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Delete load balancer</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This action cannot be undone. Are you sure you want to delete this load balancer?
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
