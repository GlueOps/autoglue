import { useMemo, useState } from "react"
import { credentialsApi } from "@/api/credentials"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Eye, Loader2, MoreHorizontal, Pencil, Plus, Search, Trash2 } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
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
import { Textarea } from "@/components/ui/textarea"

// ---------- Schemas ----------

const jsonTransform = z
  .string()
  .min(2, "JSON required")
  .refine((v) => {
    try {
      JSON.parse(v)
      return true
    } catch {
      return false
    }
  }, "Invalid JSON")
  .transform((v) => JSON.parse(v))

const createCredentialSchema = z.object({
  provider: z.enum(["aws", "cloudflare", "hetzner", "digitalocean", "generic"]),
  kind: z.enum(["aws_access_key", "api_token", "basic_auth", "oauth2"]),
  schema_version: z.number().default(1),
  name: z.string().min(1, "Name is required").max(100),
  scope_kind: z.enum(["provider", "service", "resource"]),
  scope_version: z.number().default(1),
  scope: jsonTransform,
  account_id: z
    .string()
    .optional()
    .or(z.literal(""))
    .transform((v) => (v ? v : undefined)),
  region: z
    .string()
    .optional()
    .or(z.literal(""))
    .transform((v) => (v ? v : undefined)),
  // Secrets are always JSON — makes rotate easy on update form too
  secret: jsonTransform,
})

type CreateCredentialInput = z.input<typeof createCredentialSchema>
type CreateCredentialValues = z.infer<typeof createCredentialSchema>

const updateCredentialSchema = createCredentialSchema.partial().extend({
  // allow rotating secret independently
  secret: jsonTransform.optional(),
  name: z.string().min(1, "Name is required").max(100).optional(),
})

// ---------- Helpers ----------

function pretty(obj: unknown) {
  try {
    return JSON.stringify(obj, null, 2)
  } catch {
    return ""
  }
}

function toFormDefaults<T extends Record<string, any>>(initial: Partial<T>) {
  return {
    schema_version: 1,
    scope_version: 1,
    ...initial,
  } as any
}

// ---------- Page ----------

export const CredentialPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [editOpen, setEditOpen] = useState<boolean>(false)
  const [revealOpen, setRevealOpen] = useState<boolean>(false)
  const [revealJson, setRevealJson] = useState<object | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  // List
  const credentialQ = useQuery({
    queryKey: ["credentials"],
    queryFn: () => credentialsApi.listCredentials(),
  })

  // Create
  const createMutation = useMutation({
    mutationFn: (body: CreateCredentialValues) =>
      credentialsApi.createCredential({
        provider: body.provider,
        kind: body.kind,
        schema_version: body.schema_version ?? 1,
        name: body.name,
        scope_kind: body.scope_kind,
        scope_version: body.scope_version ?? 1,
        scope: body.scope,
        account_id: body.account_id,
        region: body.region,
        secret: body.secret,
      }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["credentials"] })
      toast.success("Credential created")
      setCreateOpen(false)
      createForm.reset(createDefaults) // clear JSON textareas etc
    },
    onError: (err: any) => {
      toast.error("Failed to create credential", {
        description: err?.message ?? "Unknown error",
      })
    },
  })

  // Update
  const updateMutation = useMutation({
    mutationFn: (payload: { id: string; body: z.infer<typeof updateCredentialSchema> }) =>
      credentialsApi.updateCredential(payload.id, payload.body),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["credentials"] })
      toast.success("Credential updated")
      setEditOpen(false)
      setEditingId(null)
    },
    onError: (err: any) => {
      toast.error("Failed to update credential", {
        description: err?.message ?? "Unknown error",
      })
    },
  })

  // Delete
  const deleteMutation = useMutation({
    mutationFn: (id: string) => credentialsApi.deleteCredential(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["credentials"] })
      toast.success("Credential deleted")
    },
    onError: (err: any) => {
      toast.error("Failed to delete credential", {
        description: err?.message ?? "Unknown error",
      })
    },
  })

  // Reveal (one-time read)
  const revealMutation = useMutation({
    mutationFn: (id: string) => credentialsApi.revealCredential(id),
    onSuccess: (data) => {
      setRevealJson(data)
      setRevealOpen(true)
    },
    onError: (err: any) => {
      toast.error("Failed to reveal secret", {
        description: err?.message ?? "Unknown error",
      })
    },
  })

  // ---------- Forms ----------

  const createDefaults: CreateCredentialInput = toFormDefaults<CreateCredentialInput>({
    provider: "aws",
    kind: "aws_access_key",
    schema_version: 1,
    scope_kind: "provider",
    scope_version: 1,
    name: "",
    // IMPORTANT: default valid JSON strings so zod.transform succeeds
    scope: "{}" as any,
    secret: "{}" as any,
    account_id: "",
    region: "",
  })

  const createForm = useForm<CreateCredentialInput>({
    resolver: zodResolver(createCredentialSchema),
    defaultValues: createDefaults,
    mode: "onBlur",
  })

  const editForm = useForm<z.input<typeof updateCredentialSchema>>({
    resolver: zodResolver(updateCredentialSchema),
    defaultValues: {
      // populated on open
    },
    mode: "onBlur",
  })

  function openEdit(row: any) {
    setEditingId(row.id)
    editForm.reset({
      provider: row.provider,
      kind: row.kind,
      schema_version: row.schema_version ?? 1,
      name: row.name,
      scope_kind: row.scope_kind,
      scope_version: row.scope_version ?? 1,
      account_id: row.account_id ?? "",
      region: row.region ?? "",
      // show JSON in textareas
      scope: pretty(row.scope ?? {}),
      // secret is optional on update; leave empty to avoid rotate
      secret: undefined,
    } as any)
    setEditOpen(true)
  }

  const filtered = useMemo(() => {
    const items = credentialQ.data ?? []
    if (!filter.trim()) return items
    const f = filter.toLowerCase()
    return items.filter((c: any) =>
      [
        c.name,
        c.provider,
        c.kind,
        c.scope_kind,
        c.account_id,
        c.region,
        JSON.stringify(c.scope ?? {}),
      ]
        .filter(Boolean)
        .map((x: any) => String(x).toLowerCase())
        .some((s: string) => s.includes(f))
    )
  }, [credentialQ.data, filter])

  // ---------- UI ----------

  if (credentialQ.isLoading)
    return (
      <div className="flex items-center gap-2 p-6">
        <Loader2 className="h-4 w-4 animate-spin" /> Loading credentials…
      </div>
    )

  if (credentialQ.error)
    return (
      <div className="p-6 text-red-500">
        Error loading credentials.
        <pre className="mt-2 text-xs opacity-80">{JSON.stringify(credentialQ.error, null, 2)}</pre>
      </div>
    )

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-1 text-2xl font-bold">Credentials</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search by name, provider, kind, scope…"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Credential
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-xl">
              <DialogHeader>
                <DialogTitle>Create Credential</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form
                  onSubmit={createForm.handleSubmit((values) =>
                    createMutation.mutate(values as CreateCredentialValues)
                  )}
                  className="space-y-4 pt-2"
                >
                  <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                    <FormField
                      control={createForm.control}
                      name="provider"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Provider</FormLabel>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="aws">AWS</SelectItem>
                              <SelectItem value="cloudflare">Cloudflare</SelectItem>
                              <SelectItem value="hetzner">Hetzner</SelectItem>
                              <SelectItem value="digitalocean">DigitalOcean</SelectItem>
                              <SelectItem value="generic">Generic</SelectItem>
                            </SelectContent>
                          </Select>
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
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="aws_access_key">AWS Access Key</SelectItem>
                              <SelectItem value="api_token">API Token</SelectItem>
                              <SelectItem value="basic_auth">Basic Auth</SelectItem>
                              <SelectItem value="oauth2">OAuth2</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createForm.control}
                      name="scope_kind"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Scope Kind</FormLabel>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="provider">Provider</SelectItem>
                              <SelectItem value="service">Service</SelectItem>
                              <SelectItem value="resource">Resource</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createForm.control}
                      name="name"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Name</FormLabel>
                          <Input {...field} placeholder="My AWS Key" />
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createForm.control}
                      name="account_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Account ID (optional)</FormLabel>
                          <Input {...field} placeholder="e.g. 123456789012" />
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={createForm.control}
                      name="region"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Region (optional)</FormLabel>
                          <Input {...field} placeholder="e.g. us-east-1" />
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>

                  <FormField
                    control={createForm.control}
                    name="scope"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Scope (JSON)</FormLabel>
                        <Textarea
                          {...field}
                          rows={3}
                          placeholder='e.g. {"service":"s3"} or {"arn":"..."}'
                          className="font-mono"
                        />
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="secret"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Secret (JSON)</FormLabel>
                        <Textarea
                          {...field}
                          rows={6}
                          placeholder='{"access_key_id":"...","secret_access_key":"..."}'
                          className="font-mono"
                        />
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button variant="outline" type="button" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createMutation.isPending}>
                      {createMutation.isPending && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Create
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto rounded-xl border">
        <table className="min-w-full text-sm">
          <thead className="bg-muted/40 text-xs tracking-wide uppercase">
            <tr>
              <th className="w-[28%] px-4 py-2 text-left">Name</th>
              <th className="px-4 py-2 text-left">Provider</th>
              <th className="px-4 py-2 text-left">Kind</th>
              <th className="px-4 py-2 text-left">Scope Kind</th>
              <th className="px-4 py-2 text-left">Account</th>
              <th className="px-4 py-2 text-left">Region</th>
              <th className="px-4 py-2 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((row: any) => (
              <tr key={row.id} className="border-t">
                <td className="px-4 py-2 font-medium">{row.name}</td>
                <td className="px-4 py-2">{row.provider}</td>
                <td className="px-4 py-2">{row.kind}</td>
                <td className="px-4 py-2">{row.scope_kind}</td>
                <td className="px-4 py-2">{row.account_id ?? "—"}</td>
                <td className="px-4 py-2">{row.region ?? "—"}</td>
                <td className="px-4 py-2">
                  <div className="flex items-center justify-end gap-2">
                    <Button
                      size="icon"
                      variant="ghost"
                      title="Reveal secret (one-time read)"
                      onClick={() => revealMutation.mutate(row.id)}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                    <Button size="icon" variant="ghost" title="Edit" onClick={() => openEdit(row)}>
                      <Pencil className="h-4 w-4" />
                    </Button>

                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button size="icon" variant="ghost" title="Delete">
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Delete “{row.name}”?</AlertDialogTitle>
                          <AlertDialogDescription>
                            This will permanently remove the credential metadata. Secrets are not
                            recoverable from the service.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancel</AlertDialogCancel>
                          <AlertDialogAction
                            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                            onClick={() => deleteMutation.mutate(row.id)}
                          >
                            Delete
                          </AlertDialogAction>
                        </AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>

                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button size="icon" variant="ghost">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => openEdit(row)}>Edit</DropdownMenuItem>
                        <DropdownMenuItem onClick={() => revealMutation.mutate(row.id)}>
                          Reveal secret
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          className="text-destructive"
                          onClick={() => deleteMutation.mutate(row.id)}
                        >
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </td>
              </tr>
            ))}
            {filtered.length === 0 && (
              <tr>
                <td colSpan={7} className="text-muted-foreground px-4 py-10 text-center">
                  No credentials match your search.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Edit dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>Edit Credential</DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form
              onSubmit={editForm.handleSubmit((values) => {
                if (!editingId) return
                // Convert stringified JSON fields to objects via schema
                const parsed = updateCredentialSchema.safeParse(values)
                if (!parsed.success) {
                  toast.error("Please fix validation errors")
                  return
                }
                updateMutation.mutate({ id: editingId, body: parsed.data })
              })}
              className="space-y-4 pt-2"
            >
              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <FormField
                  control={editForm.control}
                  name="provider"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Provider</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="aws">AWS</SelectItem>
                          <SelectItem value="cloudflare">Cloudflare</SelectItem>
                          <SelectItem value="hetzner">Hetzner</SelectItem>
                          <SelectItem value="digitalocean">DigitalOcean</SelectItem>
                          <SelectItem value="generic">Generic</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
                  name="kind"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Kind</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="aws_access_key">AWS Access Key</SelectItem>
                          <SelectItem value="api_token">API Token</SelectItem>
                          <SelectItem value="basic_auth">Basic Auth</SelectItem>
                          <SelectItem value="oauth2">OAuth2</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
                  name="scope_kind"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Scope Kind</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="provider">Provider</SelectItem>
                          <SelectItem value="service">Service</SelectItem>
                          <SelectItem value="resource">Resource</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Name</FormLabel>
                      <Input {...field} />
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
                  name="account_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Account ID</FormLabel>
                      <Input {...field} placeholder="optional" />
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={editForm.control}
                  name="region"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Region</FormLabel>
                      <Input {...field} placeholder="optional" />
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={editForm.control}
                name="scope"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Scope (JSON)</FormLabel>
                    <Textarea {...field} rows={3} className="font-mono" />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editForm.control}
                name="secret"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Rotate Secret (JSON, optional)</FormLabel>
                    <Textarea
                      {...field}
                      rows={6}
                      className="font-mono"
                      placeholder="Leave empty to keep existing secret"
                    />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button variant="outline" type="button" onClick={() => setEditOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateMutation.isPending}>
                  {updateMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save changes
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Reveal modal */}
      <Dialog open={revealOpen} onOpenChange={setRevealOpen}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>Decrypted Secret</DialogTitle>
          </DialogHeader>
          <div className="bg-muted/40 rounded-lg border p-3">
            <pre className="max-h-[50vh] overflow-auto text-xs leading-relaxed">
              {pretty(revealJson ?? {})}
            </pre>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                navigator.clipboard.writeText(pretty(revealJson ?? {}))
                toast.success("Copied to clipboard")
              }}
            >
              Copy
            </Button>
            <Button onClick={() => setRevealOpen(false)}>Close</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
