// src/pages/core/servers-page.tsx
import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { Pencil, Plus, RefreshCw, Search, Trash } from "lucide-react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { api } from "@/lib/api.ts"
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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx"

type Server = {
  id: string
  hostname?: string | null
  ip_address: string
  role: string
  ssh_key_id: string
  ssh_user: string
  status: "pending" | "provisioning" | "ready" | "failed" | string
  organization_id: string
  created_at: string
  updated_at: string
}

type SshKey = {
  id: string
  name?: string | null
  public_keys: string
  fingerprint: string
  created_at: string
  updated_at: string
  organization_id: string
}

const STATUS = ["pending", "provisioning", "ready", "failed"] as const
type Status = (typeof STATUS)[number]

const ROLE_OPTIONS = ["master", "worker", "bastion"] as const
type Role = (typeof ROLE_OPTIONS)[number]

const CreateServerSchema = z.object({
  hostname: z.string().trim().max(120, "Max 120 chars").optional(),
  ip_address: z.string().trim().min(1, "IP address is required"),
  role: z.enum(ROLE_OPTIONS),
  ssh_key_id: z.string().uuid("Pick a valid SSH key"),
  ssh_user: z.string().trim().min(1, "SSH user is required"),
  status: z.enum(STATUS).default("pending"),
})
type CreateServerInput = z.input<typeof CreateServerSchema>
type CreateServerValues = z.output<typeof CreateServerSchema>

const UpdateServerSchema = CreateServerSchema.partial()
type UpdateServerValues = z.infer<typeof UpdateServerSchema>

function StatusBadge({ status }: { status: string }) {
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

function truncateMiddle(str: string, keep = 16) {
  if (!str || str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

export const ServersPage = () => {
  const [servers, setServers] = useState<Server[]>([])
  const [sshKeys, setSshKeys] = useState<SshKey[]>([])
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)

  const [q, setQ] = useState("")
  const [statusFilter, setStatusFilter] = useState<Status | "">("")
  const [roleFilter, setRoleFilter] = useState<Role | "">("")

  const [createOpen, setCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Server | null>(null)

  function buildServersURL() {
    const qs = new URLSearchParams()
    if (statusFilter) qs.set("status", statusFilter)
    if (roleFilter) qs.set("role", roleFilter)
    const suffix = qs.toString() ? `?${qs.toString()}` : ""
    return `/api/v1/servers${suffix}`
  }

  async function loadAll() {
    setLoading(true)
    setErr(null)
    try {
      // NOTE: your api.get<T> returns parsed JSON directly
      const [srv, keys] = await Promise.all([
        api.get<Server[]>(buildServersURL()),
        api.get<SshKey[]>("/api/v1/ssh"),
      ])
      setServers(srv ?? [])
      setSshKeys(keys ?? [])
    } catch (e) {
      console.error(e)
      setErr("Failed to load servers or SSH keys")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadAll()
    const onStorage = (e: StorageEvent) => {
      if (e.key === "active_org_id") loadAll()
    }
    window.addEventListener("storage", onStorage)
    return () => window.removeEventListener("storage", onStorage)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    loadAll()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [statusFilter, roleFilter])

  const keyById = useMemo(() => {
    const m = new Map<string, SshKey>()
    sshKeys.forEach((k) => m.set(k.id, k))
    return m
  }, [sshKeys])

  const filtered = useMemo(() => {
    const needle = q.trim().toLowerCase()
    if (!needle) return servers
    return servers.filter(
      (s) =>
        (s.hostname ?? "").toLowerCase().includes(needle) ||
        s.ip_address.toLowerCase().includes(needle) ||
        s.role.toLowerCase().includes(needle) ||
        s.ssh_user.toLowerCase().includes(needle)
    )
  }, [servers, q])

  async function deleteServer(id: string) {
    if (!confirm("Delete this server? This cannot be undone.")) return
    await api.delete<void>(`/api/v1/servers/${encodeURIComponent(id)}`)
    await loadAll()
  }

  const createForm = useForm<CreateServerInput, any, CreateServerValues>({
    resolver: zodResolver(CreateServerSchema),
    defaultValues: {
      hostname: "",
      ip_address: "",
      role: "worker",
      ssh_key_id: "",
      ssh_user: "ubuntu",
      status: "pending",
    },
  })

  const submitCreate = async (values: CreateServerValues) => {
    const payload: Record<string, any> = {
      ip_address: values.ip_address.trim(),
      role: values.role,
      ssh_key_id: values.ssh_key_id,
      ssh_user: values.ssh_user.trim(),
      status: values.status,
    }
    if (values.hostname && values.hostname.trim()) {
      payload.hostname = values.hostname.trim()
    }

    await api.post<Server>("/api/v1/servers", payload)
    setCreateOpen(false)
    createForm.reset()
    await loadAll()
  }

  const editForm = useForm<UpdateServerValues>({
    resolver: zodResolver(UpdateServerSchema),
    defaultValues: {},
  })

  function openEdit(s: Server) {
    setEditTarget(s)
    editForm.reset({
      hostname: s.hostname ?? "",
      ip_address: s.ip_address,
      role: (ROLE_OPTIONS.includes(s.role as Role) ? (s.role as Role) : "worker") as any,
      ssh_key_id: s.ssh_key_id,
      ssh_user: s.ssh_user,
      status: (STATUS.includes(s.status as Status) ? (s.status as Status) : "pending") as any,
    })
  }

  const submitEdit = async (values: UpdateServerValues) => {
    if (!editTarget) return
    const payload: Record<string, any> = {}
    if (values.hostname !== undefined) payload.hostname = values.hostname?.trim() || ""
    if (values.ip_address !== undefined) payload.ip_address = values.ip_address.trim()
    if (values.role !== undefined) payload.role = values.role
    if (values.ssh_key_id !== undefined) payload.ssh_key_id = values.ssh_key_id
    if (values.ssh_user !== undefined) payload.ssh_user = values.ssh_user.trim()
    if (values.status !== undefined) payload.status = values.status

    await api.patch<Server>(`/api/v1/servers/${encodeURIComponent(editTarget.id)}`, payload)
    setEditTarget(null)
    await loadAll()
  }

  if (loading) return <div className="p-6">Loading servers…</div>
  if (err) return <div className="p-6 text-red-500">{err}</div>

  return (
    <TooltipProvider>
      <div className="space-y-4 p-6">
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <h1 className="mb-4 text-2xl font-bold">Servers</h1>

          <div className="flex flex-wrap items-center gap-2">
            <div className="relative">
              <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
              <Input
                value={q}
                onChange={(e) => setQ(e.target.value)}
                placeholder="Search hostname, IP, role, user…"
                className="w-64 pl-8"
              />
            </div>

            <Select value={roleFilter} onValueChange={(v) => setRoleFilter(v as Role | "")}>
              <SelectTrigger className="w-36">
                <SelectValue placeholder="Role (all)" />
              </SelectTrigger>
              <SelectContent>
                {/* <SelectItem value="">All roles</SelectItem> */}
                {ROLE_OPTIONS.map((r) => (
                  <SelectItem key={r} value={r}>
                    {r}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Select value={statusFilter} onValueChange={(v) => setStatusFilter(v as Status | "")}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Status (all)" />
              </SelectTrigger>
              <SelectContent>
                {/* <SelectItem value="">All status</SelectItem> */}
                {STATUS.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <Button variant="outline" onClick={loadAll}>
              <RefreshCw className="mr-2 h-4 w-4" />
              Refresh
            </Button>

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
                  <form onSubmit={createForm.handleSubmit(submitCreate)} className="space-y-4">
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

                    <FormField
                      control={createForm.control}
                      name="ip_address"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>IP address</FormLabel>
                          <FormControl>
                            <Input placeholder="10.0.1.23" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                      <FormField
                        control={createForm.control}
                        name="role"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Role</FormLabel>
                            <Select onValueChange={field.onChange} value={field.value}>
                              <FormControl>
                                <SelectTrigger>
                                  <SelectValue placeholder="Select role" />
                                </SelectTrigger>
                              </FormControl>
                              <SelectContent>
                                {ROLE_OPTIONS.map((r) => (
                                  <SelectItem key={r} value={r}>
                                    {r}
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
                                    sshKeys.length ? "Select SSH key" : "No SSH keys found"
                                  }
                                />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {sshKeys.map((k) => (
                                <SelectItem key={k.id} value={k.id}>
                                  {k.name ? k.name : "Unnamed key"} —{" "}
                                  {truncateMiddle(k.fingerprint, 8)}
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
                  <TableHead>Hostname</TableHead>
                  <TableHead>IP address</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>SSH user</TableHead>
                  <TableHead>SSH key</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-[180px] text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((s) => {
                  const key = keyById.get(s.ssh_key_id)
                  return (
                    <TableRow key={s.id}>
                      <TableCell className="font-medium">{s.hostname || "—"}</TableCell>
                      <TableCell>
                        <code className="font-mono text-sm">{s.ip_address}</code>
                      </TableCell>
                      <TableCell className="capitalize">{s.role}</TableCell>
                      <TableCell>
                        <code className="font-mono text-sm">{s.ssh_user}</code>
                      </TableCell>
                      <TableCell>
                        {key ? (
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <span className="inline-flex items-center gap-2">
                                <Badge variant="secondary">{key.name || "SSH key"}</Badge>
                                <code className="font-mono text-xs">
                                  {truncateMiddle(key.fingerprint, 8)}
                                </code>
                              </span>
                            </TooltipTrigger>
                            <TooltipContent className="max-w-[70vw]">
                              <p className="font-mono text-xs break-all">{key.public_keys}</p>
                            </TooltipContent>
                          </Tooltip>
                        ) : (
                          <span className="text-muted-foreground">Unknown</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <StatusBadge status={s.status} />
                      </TableCell>
                      <TableCell>{new Date(s.created_at).toLocaleString()}</TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-2">
                          <Button variant="outline" size="sm" onClick={() => openEdit(s)}>
                            <Pencil className="mr-2 h-4 w-4" /> Edit
                          </Button>

                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="destructive" size="sm">
                                <Trash className="mr-2 h-4 w-4" /> Delete
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => deleteServer(s.id)}>
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
                    <TableCell colSpan={8} className="text-muted-foreground py-10 text-center">
                      No servers match your filters.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>

      {/* Edit dialog */}
      <Dialog open={!!editTarget} onOpenChange={(o) => !o && setEditTarget(null)}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit server</DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form onSubmit={editForm.handleSubmit(submitEdit)} className="space-y-4">
              <FormField
                control={editForm.control}
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

              <FormField
                control={editForm.control}
                name="ip_address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>IP address</FormLabel>
                    <FormControl>
                      <Input placeholder="10.0.1.23" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <FormField
                  control={editForm.control}
                  name="role"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Role</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value as any}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select role" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {ROLE_OPTIONS.map((r) => (
                            <SelectItem key={r} value={r}>
                              {r}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
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
                control={editForm.control}
                name="ssh_key_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SSH key</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value as any}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={sshKeys.length ? "Select SSH key" : "No SSH keys found"}
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {sshKeys.map((k) => (
                          <SelectItem key={k.id} value={k.id}>
                            {k.name ? k.name : "SSH key"} — {truncateMiddle(k.fingerprint, 8)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="status"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Status</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value as any}>
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
    </TooltipProvider>
  )
}
