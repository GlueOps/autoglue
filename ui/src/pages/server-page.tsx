import { useMemo, useState } from "react"
import { serversApi } from "@/api/servers.ts"
import { sshApi } from "@/api/ssh.ts"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { formatDistanceToNow } from "date-fns"
import { Plus, Search } from "lucide-react"
import { useForm, useWatch } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { cn } from "@/lib/utils"
import { truncateMiddle } from "@/lib/utils.ts"
import { Badge } from "@/components/ui/badge.tsx"
import { Button } from "@/components/ui/button.tsx"
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
import { TooltipProvider } from "@/components/ui/tooltip.tsx"

const ROLE_OPTIONS = ["master", "worker", "bastion"] as const
type Role = (typeof ROLE_OPTIONS)[number]

const STATUS = ["pending", "provisioning", "ready", "failed"] as const
type Status = (typeof STATUS)[number]

const createServerSchema = z
  .object({
    hostname: z.string().trim().max(60, "Max 60 chars"),
    public_ip_address: z.string().trim().optional().or(z.literal("")),
    private_ip_address: z.string().trim().min(1, "Private IP address required"),
    role: z.enum(ROLE_OPTIONS),
    ssh_key_id: z.uuid("Pick a valid SSH key"),
    ssh_user: z.string().trim().min(1, "SSH user is required"),
    status: z.enum(STATUS).default("pending"),
  })
  .refine(
    (v) => v.role !== "bastion" || (v.public_ip_address && v.public_ip_address.trim() !== ""),
    { message: "Public IP required for bastion", path: ["public_ip_address"] }
  )
type CreateServerInput = z.input<typeof createServerSchema>

const updateServerSchema = createServerSchema.partial()
type UpdateServerValues = z.infer<typeof updateServerSchema>

function StatusBadge({ status }: { status: Status }) {
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
      {status}
    </Badge>
  )
}

export const ServerPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)

  const [statusFilter, setStatusFilter] = useState<Status | "">("")
  const [roleFilter, setRoleFilter] = useState<Role | "">("")
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  const serverQ = useQuery({
    queryKey: ["servers"],
    queryFn: () => serversApi.listServers(),
  })

  const sshQ = useQuery({
    queryKey: ["ssh_keys"],
    queryFn: () => sshApi.listSshKeys(),
  })

  // Map of ssh_key_id -> label
  const sshLabelById = useMemo(() => {
    const m = new Map<string, string>()
    for (const k of sshQ.data ?? []) {
      const name = k.name ? k.name : "Unnamed key"
      const fp = k.fingerprint ? truncateMiddle(k.fingerprint, 8) : ""
      m.set(k.id!, fp ? `${name} — ${fp}` : name)
    }
    return m
  }, [sshQ.data])

  // --- Create ---
  const createForm = useForm<CreateServerInput>({
    resolver: zodResolver(createServerSchema),
    defaultValues: {
      hostname: "",
      private_ip_address: "",
      public_ip_address: "",
      role: "worker",
      ssh_key_id: "" as unknown as string,
      ssh_user: "",
      status: "pending",
    },
    mode: "onChange",
  })

  const watchedRoleCreate = useWatch({
    control: createForm.control,
    name: "role",
  })
  const roleIsBastion = watchedRoleCreate === "bastion"

  const watchedPublicIpCreate = useWatch({
    control: createForm.control,
    name: "public_ip_address",
  })
  const pubCreate = watchedPublicIpCreate?.trim() ?? ""
  const needPubCreate = roleIsBastion && pubCreate === ""

  const createMut = useMutation({
    mutationFn: (values: CreateServerInput) => serversApi.createServer(values as any),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["servers"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Server created successfully")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to create server")
    },
  })

  // --- Update ---
  const updateForm = useForm<UpdateServerValues>({
    resolver: zodResolver(updateServerSchema),
    defaultValues: {},
    mode: "onChange",
  })

  const watchedRoleUpdate = useWatch({
    control: updateForm.control,
    name: "role",
  })

  const watchedPublicIpAddressUpdate = useWatch({
    control: updateForm.control,
    name: "public_ip_address",
  })

  const roleIsBastionU = watchedRoleUpdate === "bastion"

  const pubUpdate = watchedPublicIpAddressUpdate?.trim() ?? ""
  const needPubUpdate = roleIsBastionU && pubUpdate === ""

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateServerValues }) =>
      serversApi.updateServer(id, values as any),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["servers"] })
      setUpdateOpen(false)
      setEditingId(null)
      toast.success("Server updated successfully")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to update server")
    },
  })

  // --- Delete ---
  const deleteMut = useMutation({
    mutationFn: (id: string) => serversApi.deleteServer(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["servers"] })
      setDeleteId(null)
      toast.success("Server deleted successfully")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to delete server")
    },
  })

  const filtered = useMemo(() => {
    const data = serverQ.data ?? []
    const q = filter.trim().toLowerCase()

    const textFiltered = q
      ? data.filter((k: any) => {
          return (
            k.hostname?.toLowerCase().includes(q) ||
            k.public_ip_address?.toLowerCase().includes(q) ||
            k.private_ip_address?.toLowerCase().includes(q) ||
            k.role?.toLowerCase().includes(q) ||
            k.ssh_user?.toLowerCase().includes(q)
          )
        })
      : data

    const roleFiltered = roleFilter
      ? textFiltered.filter((k: any) => k.role === roleFilter)
      : textFiltered

    const statusFiltered = statusFilter
      ? roleFiltered.filter((k: any) => k.status === statusFilter)
      : roleFiltered

    return statusFiltered
  }, [filter, roleFilter, statusFilter, serverQ.data])

  const onCreateSubmit = (values: CreateServerInput) => {
    createMut.mutate(values)
  }

  const openEdit = (srv: any) => {
    setEditingId(srv.id)
    updateForm.reset({
      hostname: srv.hostname ?? "",
      public_ip_address: srv.public_ip_address ?? "",
      private_ip_address: srv.private_ip_address ?? "",
      role: (srv.role as Role) ?? "worker",
      ssh_key_id: srv.ssh_key_id ?? "",
      ssh_user: srv.ssh_user ?? "",
      status: (srv.status as Status) ?? "pending",
    })
    setUpdateOpen(true)
  }

  if (sshQ.data?.length === 0)
    return <div className="p-6">Please create an SSH key for your organization first.</div>
  if (serverQ.isLoading) return <div className="p-6">Loading servers…</div>
  if (serverQ.error) return <div className="p-6 text-red-500">Error loading servers.</div>

  return (
    <TooltipProvider>
      <div className="space-y-4 p-6">
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <h1 className="mb-4 text-2xl font-bold">Servers</h1>

          <div className="flex flex-wrap items-center gap-2">
            <div className="relative">
              <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
              <Input
                value={filter}
                onChange={(e) => setFilter(e.target.value)}
                placeholder="Search hostname, Public IP, Private IP, role, user…"
                className="w-64 pl-8"
              />
            </div>

            <Select
              value={roleFilter || "all"} // map "" -> "all" for the UI
              onValueChange={(v) => setRoleFilter(v === "all" ? "" : (v as Role))}
            >
              <SelectTrigger className="w-36">
                <SelectValue placeholder="Role (all)" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All roles</SelectItem>
                {ROLE_OPTIONS.map((r) => (
                  <SelectItem key={r} value={r}>
                    {r}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select
              value={statusFilter || "all"} // map "" -> "all" for the UI
              onValueChange={(v) => setStatusFilter(v === "all" ? "" : (v as Status))}
            >
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Status (all)" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All statuses</SelectItem> {/* sentinel */}
                {STATUS.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Dialog open={createOpen} onOpenChange={setCreateOpen}>
              <DialogTrigger asChild>
                <Button onClick={() => setCreateOpen(true)}>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Server
                </Button>
              </DialogTrigger>
              <DialogContent className="sm:max-w-lg">
                <DialogHeader>
                  <DialogTitle>Create server</DialogTitle>
                </DialogHeader>

                <Form {...createForm}>
                  <form className="space-y-4" onSubmit={createForm.handleSubmit(onCreateSubmit)}>
                    <FormField
                      control={createForm.control}
                      name="hostname"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Hostname</FormLabel>
                          <FormControl>
                            <Input placeholder="worker-01" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                      <FormField
                        control={createForm.control}
                        name="public_ip_address"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="flex items-center justify-between">
                              <span>Public IP Address</span>
                              <span
                                className={cn(
                                  "rounded-full px-2 py-0.5 text-xs",
                                  roleIsBastion
                                    ? "bg-amber-100 text-amber-900"
                                    : "bg-muted text-muted-foreground"
                                )}
                              >
                                {roleIsBastion ? "Required for bastion" : "Optional"}
                              </span>
                            </FormLabel>
                            <FormControl>
                              <Input
                                placeholder={
                                  roleIsBastion
                                    ? "Required for bastion (e.g. 34.12.56.78)"
                                    : "34.12.56.78"
                                }
                                aria-required={roleIsBastion}
                                aria-invalid={
                                  needPubCreate || !!createForm.formState.errors.public_ip_address
                                }
                                required={roleIsBastion}
                                {...field}
                                className={cn(
                                  needPubCreate &&
                                    "border-destructive focus-visible:ring-destructive"
                                )}
                              />
                            </FormControl>
                            {roleIsBastion && (
                              <div className="rounded-md border border-amber-200 bg-amber-50 p-2 text-xs text-amber-900">
                                Bastion nodes must have a{" "}
                                <span className="font-medium">Public IP</span>.
                              </div>
                            )}
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={createForm.control}
                        name="private_ip_address"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Private IP Address</FormLabel>
                            <FormControl>
                              <Input placeholder="192.168.10.1" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
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
                                <SelectItem value="bastion">
                                  bastion — requires Public IP
                                </SelectItem>
                              </SelectContent>
                            </Select>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                      <FormField
                        control={createForm.control}
                        name="ssh_user"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>SSH user</FormLabel>
                            <FormControl>
                              <Input placeholder="ubuntu" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>

                    <FormField
                      control={createForm.control}
                      name="ssh_key_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>SSH key</FormLabel>
                          <Select onValueChange={field.onChange} value={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue
                                  placeholder={
                                    sshQ.data?.length ? "Select SSH key" : "No SSH keys found"
                                  }
                                />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {sshQ.data!.map((k) => (
                                <SelectItem key={k.id} value={k.id!}>
                                  {k.name ? k.name : "Unnamed key"} —{" "}
                                  {truncateMiddle(k.fingerprint!, 8)}
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
                      name="status"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Initial status</FormLabel>
                          <Select onValueChange={field.onChange} value={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="pending" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {STATUS.map((s) => (
                                <SelectItem key={s} value={s}>
                                  {s}
                                </SelectItem>
                              ))}
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
                      <Button
                        type="submit"
                        disabled={
                          createMut.isPending ||
                          createForm.formState.isSubmitting ||
                          !createForm.formState.isValid
                        }
                      >
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
                  <TableHead>Hostname</TableHead>
                  <TableHead>IP address</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>SSH user</TableHead>
                  <TableHead>SSH key</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-[220px] text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                {filtered.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={8}
                      className="text-muted-foreground py-10 text-center text-sm"
                    >
                      No servers found.
                    </TableCell>
                  </TableRow>
                ) : (
                  filtered.map((k: any) => (
                    <TableRow key={k.id}>
                      <TableCell className="font-medium">{k.hostname}</TableCell>
                      <TableCell>
                        <div className="flex flex-col">
                          <span
                            className={cn(
                              "tabular-nums",
                              !k.public_ip_address && "text-muted-foreground"
                            )}
                          >
                            {k.public_ip_address || "—"}
                          </span>
                          <span className="text-muted-foreground tabular-nums">
                            {k.private_ip_address}
                          </span>
                        </div>
                      </TableCell>
                      <TableCell className="capitalize">
                        <span
                          className={cn(
                            k.role === "bastion" &&
                              "rounded bg-amber-50 px-2 py-0.5 dark:bg-amber-900"
                          )}
                        >
                          {k.role}
                        </span>
                      </TableCell>
                      <TableCell className="tabular-nums">{k.ssh_user}</TableCell>
                      <TableCell className="truncate">
                        {sshLabelById.get(k.ssh_key_id) ?? "—"}
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={(k.status ?? "pending") as Status} />
                      </TableCell>
                      <TableCell title={k.created_at}>
                        {k.created_at
                          ? `${formatDistanceToNow(new Date(k.created_at), { addSuffix: true })}`
                          : "—"}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button variant="outline" size="sm" onClick={() => openEdit(k)}>
                            Edit
                          </Button>
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => setDeleteId(k.id)}
                            disabled={deleteMut.isPending && deleteId === k.id}
                          >
                            {deleteMut.isPending && deleteId === k.id ? "Deleting…" : "Delete"}
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>

      {/* Update dialog */}
      <Dialog open={updateOpen} onOpenChange={setUpdateOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit server</DialogTitle>
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
                name="hostname"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Hostname</FormLabel>
                    <FormControl>
                      <Input placeholder="worker-01" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <FormField
                  control={updateForm.control}
                  name="public_ip_address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="flex items-center justify-between">
                        <span>Public IP Address</span>
                        <span
                          className={cn(
                            "rounded-full px-2 py-0.5 text-xs",
                            roleIsBastionU
                              ? "bg-amber-100 text-amber-900"
                              : "bg-muted text-muted-foreground"
                          )}
                        >
                          {roleIsBastionU ? "Required for bastion" : "Optional"}
                        </span>
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder={
                            roleIsBastionU
                              ? "Required for bastion (e.g. 34.12.56.78)"
                              : "34.12.56.78"
                          }
                          aria-required={roleIsBastionU}
                          aria-invalid={
                            needPubUpdate || !!updateForm.formState.errors.public_ip_address
                          }
                          required={roleIsBastionU}
                          {...field}
                          className={cn(
                            needPubUpdate && "border-destructive focus-visible:ring-destructive"
                          )}
                        />
                      </FormControl>
                      {roleIsBastionU && (
                        <div className="rounded-md border border-amber-200 bg-amber-50 p-2 text-xs text-amber-900">
                          Bastion nodes must have a <span className="font-medium">Public IP</span>.
                        </div>
                      )}
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={updateForm.control}
                  name="private_ip_address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Private IP Address</FormLabel>
                      <FormControl>
                        <Input placeholder="192.168.10.1" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
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
                          <SelectItem value="bastion">bastion — requires Public IP</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={updateForm.control}
                  name="ssh_user"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>SSH user</FormLabel>
                      <FormControl>
                        <Input placeholder="ubuntu" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={updateForm.control}
                name="ssh_key_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SSH key</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select SSH key" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {sshQ.data!.map((k) => (
                          <SelectItem key={k.id} value={k.id!}>
                            {k.name ? k.name : "Unnamed key"} — {truncateMiddle(k.fingerprint!, 8)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={updateForm.control}
                name="status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Status</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="pending" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {STATUS.map((s) => (
                          <SelectItem key={s} value={s}>
                            {s}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
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
            <DialogTitle>Delete server</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This action cannot be undone. Are you sure you want to delete this server?
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
    </TooltipProvider>
  )
}
