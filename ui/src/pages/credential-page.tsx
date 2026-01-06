import { useMemo, useState } from "react"
import { credentialsApi } from "@/api/credentials"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  AlertTriangle,
  Eye,
  Loader2,
  MoreHorizontal,
  Pencil,
  Plus,
  Search,
  Trash2,
} from "lucide-react"
import { Controller, useForm } from "react-hook-form"
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
import { Switch } from "@/components/ui/switch"
import { Textarea } from "@/components/ui/textarea"

// -------------------- Constants --------------------

const AWS_ALLOWED_SERVICES = ["route53", "s3", "ec2", "iam", "rds", "dynamodb"] as const
type AwsSvc = (typeof AWS_ALLOWED_SERVICES)[number]

// -------------------- Schemas --------------------
// Zod v4 gotchas you hit:
// - .partial() cannot be used if the object contains refinements/effects (often true once you have transforms/refines).
// - .extend() cannot overwrite keys after refinements (requires .safeExtend()).
// Easiest fix: define CREATE and UPDATE schemas separately (no .partial(), no post-refinement .extend()).

const createCredentialSchema = z
  .object({
    credential_provider: z.enum(["aws", "cloudflare", "hetzner", "digitalocean", "generic"]),
    kind: z.enum(["aws_access_key", "api_token", "basic_auth", "oauth2"]),
    schema_version: z.number().default(1),
    name: z.string().min(1, "Name is required").max(100),
    scope_kind: z.enum(["provider", "service", "resource"]),
    scope_version: z.number().default(1),
    scope: z.any(),
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
    secret: z.any(),
  })
  .superRefine((val, ctx) => {
    // scope required unless provider scope
    if (val.scope_kind !== "provider" && !val.scope) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["scope"],
        message: `scope is required`,
      })
    }

    // AWS scope checks
    if (val.credential_provider === "aws") {
      if (val.scope_kind === "service") {
        const svc = (val.scope as any)?.service
        if (!AWS_ALLOWED_SERVICES.includes(svc)) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ["scope"],
            message: `For AWS service scope, "service" must be one of: ${AWS_ALLOWED_SERVICES.join(", ")}`,
          })
        }
      }
      if (val.scope_kind === "resource") {
        const arn = (val.scope as any)?.arn
        if (typeof arn !== "string" || !arn.startsWith("arn:")) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ["scope"],
            message: `For AWS resource scope, "arn" must start with "arn:"`,
          })
        }
      }
    }

    // secret requiredness by kind (create always validates)
    if (val.kind === "aws_access_key") {
      const sk = val.secret ?? {}
      const id = sk.access_key_id
      if (typeof id !== "string" || !/^[A-Z0-9]{20}$/.test(id)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `access_key_id must be 20 chars (A-Z0-9)`,
        })
      }
      if (typeof sk.secret_access_key !== "string" || sk.secret_access_key.length < 10) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `secret_access_key is required`,
        })
      }
    }

    if (val.kind === "api_token") {
      const token = (val.secret ?? {}).token
      if (!token) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `token is required`,
        })
      }
    }

    if (val.kind === "basic_auth") {
      const s = val.secret ?? {}
      if (!s.username || !s.password) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `username and password are required`,
        })
      }
    }

    if (val.kind === "oauth2") {
      const s = val.secret ?? {}
      if (!s.client_id || !s.client_secret || !s.refresh_token) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `client_id, client_secret, and refresh_token are required`,
        })
      }
    }
  })

type CreateCredentialValues = z.input<typeof createCredentialSchema>

// UPDATE schema: all fields optional, and validations are "patch-friendly".
const updateCredentialSchema = z
  .object({
    credential_provider: z
      .enum(["aws", "cloudflare", "hetzner", "digitalocean", "generic"])
      .optional(),
    kind: z.enum(["aws_access_key", "api_token", "basic_auth", "oauth2"]).optional(),
    schema_version: z.number().optional(),
    name: z.string().min(1, "Name is required").max(100).optional(),
    scope_kind: z.enum(["provider", "service", "resource"]).optional(),
    scope_version: z.number().optional(),
    scope: z.any().optional(),
    // allow "" so your form can keep empty strings; buildUpdateBody will drop them
    account_id: z.string().optional().or(z.literal("")),
    region: z.string().optional().or(z.literal("")),
    secret: z.any().optional(),
  })
  .superRefine((val, ctx) => {
    // If scope_kind is being changed to non-provider, require scope in the patch
    if (typeof val.scope_kind !== "undefined") {
      if (val.scope_kind !== "provider" && !val.scope) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["scope"],
          message: `scope is required`,
        })
      }
    }

    // AWS scope checks only if we have enough info
    if (val.credential_provider === "aws") {
      if (val.scope_kind === "service" && typeof val.scope !== "undefined") {
        const svc = (val.scope as any)?.service
        if (!AWS_ALLOWED_SERVICES.includes(svc)) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ["scope"],
            message: `For AWS service scope, "service" must be one of: ${AWS_ALLOWED_SERVICES.join(", ")}`,
          })
        }
      }
      if (val.scope_kind === "resource" && typeof val.scope !== "undefined") {
        const arn = (val.scope as any)?.arn
        if (typeof arn !== "string" || !arn.startsWith("arn:")) {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ["scope"],
            message: `For AWS resource scope, "arn" must start with "arn:"`,
          })
        }
      }
    }

    // Secret validation on update:
    // - only validate if rotating secret OR changing kind
    // - if rotating secret but kind is NOT provided, skip kind-specific checks (backend can validate)
    const rotatingSecret = typeof val.secret !== "undefined"
    const changingKind = typeof val.kind !== "undefined"
    if (!rotatingSecret && !changingKind) return
    if (!val.kind) return

    if (val.kind === "aws_access_key") {
      const sk = val.secret ?? {}
      const id = sk.access_key_id
      if (typeof id !== "string" || !/^[A-Z0-9]{20}$/.test(id)) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `access_key_id must be 20 chars (A-Z0-9)`,
        })
      }
      if (typeof sk.secret_access_key !== "string" || sk.secret_access_key.length < 10) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `secret_access_key is required`,
        })
      }
    }

    if (val.kind === "api_token") {
      const token = (val.secret ?? {}).token
      if (!token) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `token is required`,
        })
      }
    }

    if (val.kind === "basic_auth") {
      const s = val.secret ?? {}
      if (!s.username || !s.password) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `username and password are required`,
        })
      }
    }

    if (val.kind === "oauth2") {
      const s = val.secret ?? {}
      if (!s.client_id || !s.client_secret || !s.refresh_token) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["secret"],
          message: `client_id, client_secret, and refresh_token are required`,
        })
      }
    }
  })

type UpdateCredentialValues = z.input<typeof updateCredentialSchema>

// -------------------- Helpers --------------------

function pretty(obj: unknown) {
  try {
    if (obj == null) return ""
    if (typeof obj === "string") {
      try {
        return JSON.stringify(JSON.parse(obj), null, 2)
      } catch {
        return obj
      }
    }
    return JSON.stringify(obj, null, 2)
  } catch {
    return ""
  }
}

function extractErr(e: any): string {
  const raw = (e as any)?.body ?? (e as any)?.response ?? (e as any)?.message
  if (typeof raw === "string") return raw
  try {
    const msg = (e as any)?.response?.data?.message || (e as any)?.message
    if (msg) return String(msg)
  } catch {
    return "Unknown error"
  }
  return "Unknown error"
}

function isAwsServiceScope({
  credential_provider,
  scope_kind,
}: {
  credential_provider?: string
  scope_kind?: string
}) {
  return credential_provider === "aws" && scope_kind === "service"
}
function isAwsResourceScope({
  credential_provider,
  scope_kind,
}: {
  credential_provider?: string
  scope_kind?: string
}) {
  return credential_provider === "aws" && scope_kind === "resource"
}
function isProviderScope({ scope_kind }: { scope_kind?: string }) {
  return scope_kind === "provider"
}

function defaultCreateValues(): CreateCredentialValues {
  return {
    credential_provider: "aws",
    kind: "aws_access_key",
    schema_version: 1,
    name: "",
    scope_kind: "provider",
    scope_version: 1,
    scope: {},
    account_id: "",
    region: "",
    secret: {},
  }
}

// Build exact POST body as the SDK sends it
function buildCreateBody(v: CreateCredentialValues) {
  return {
    credential_provider: v.credential_provider,
    kind: v.kind,
    schema_version: v.schema_version ?? 1,
    name: v.name,
    scope_kind: v.scope_kind,
    scope_version: v.scope_version ?? 1,
    scope: v.scope ?? {},
    account_id: v.account_id,
    region: v.region,
    secret: v.secret ?? {},
  }
}

// Build exact PATCH body (only provided fields)
function buildUpdateBody(v: UpdateCredentialValues) {
  const body: any = {}
  const keys: (keyof UpdateCredentialValues)[] = [
    "name",
    "account_id",
    "region",
    "scope_kind",
    "scope_version",
    "scope",
    "secret",
    "credential_provider",
    "kind",
    "schema_version",
  ]
  for (const k of keys) {
    if (typeof v[k] !== "undefined" && v[k] !== "") body[k] = v[k]
  }
  return body
}

// -------------------- Page --------------------

export const CredentialPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [editOpen, setEditOpen] = useState<boolean>(false)
  const [revealOpen, setRevealOpen] = useState<boolean>(false)
  const [revealJson, setRevealJson] = useState<object | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)
  const [useRawSecretJSON, setUseRawSecretJSON] = useState<boolean>(false)
  const [useRawEditSecretJSON, setUseRawEditSecretJSON] = useState<boolean>(false)

  // Preview modals
  const [previewCreateOpen, setPreviewCreateOpen] = useState(false)
  const [previewCreateBody, setPreviewCreateBody] = useState<object | null>(null)
  const [previewUpdateOpen, setPreviewUpdateOpen] = useState(false)
  const [previewUpdateBody, setPreviewUpdateBody] = useState<object | null>(null)

  const qc = useQueryClient()

  // List
  const credentialQ = useQuery({
    queryKey: ["credentials"],
    queryFn: () => credentialsApi.listCredentials(),
  })

  // Create
  const createMutation = useMutation({
    mutationFn: (body: CreateCredentialValues) =>
      credentialsApi.createCredential(buildCreateBody(body) as any),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["credentials"] })
      toast.success("Credential created")
      setCreateOpen(false)
      createForm.reset(defaultCreateValues())
      setUseRawSecretJSON(false)
    },
    onError: (err: any) => {
      toast.error("Failed to create credential", { description: extractErr(err) })
    },
  })

  // Update
  const updateMutation = useMutation({
    mutationFn: (payload: { id: string; body: UpdateCredentialValues }) =>
      credentialsApi.updateCredential(payload.id, buildUpdateBody(payload.body)),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["credentials"] })
      toast.success("Credential updated")
      setEditOpen(false)
      setEditingId(null)
      setUseRawEditSecretJSON(false)
    },
    onError: (err: any) => {
      toast.error("Failed to update credential", { description: extractErr(err) })
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
      toast.error("Failed to delete credential", { description: extractErr(err) })
    },
  })

  // Reveal
  const revealMutation = useMutation({
    mutationFn: (id: string) => credentialsApi.revealCredential(id),
    onSuccess: (data) => {
      setRevealJson(data)
      setRevealOpen(true)
    },
    onError: (err: any) => {
      toast.error("Failed to reveal secret", { description: extractErr(err) })
    },
  })

  // ---------- Forms ----------

  const createForm = useForm<CreateCredentialValues>({
    resolver: zodResolver(createCredentialSchema),
    defaultValues: defaultCreateValues(),
    mode: "onBlur",
  })

  const editForm = useForm<UpdateCredentialValues>({
    resolver: zodResolver(updateCredentialSchema),
    defaultValues: {},
    mode: "onBlur",
  })

  function openEdit(row: any) {
    setEditingId(row.id)
    editForm.reset({
      // FIX: correct key (was "provider" in your original)
      credential_provider: row.credential_provider,
      kind: row.kind,
      schema_version: row.schema_version ?? 1,
      name: row.name,
      scope_kind: row.scope_kind,
      scope_version: row.scope_version ?? 1,
      account_id: row.account_id ?? "",
      region: row.region ?? "",
      scope: row.scope ?? (row.scope_kind === "provider" ? {} : undefined),
      secret: undefined, // keep existing unless user rotates
    } as any)
    setUseRawEditSecretJSON(false)
    setEditOpen(true)
  }

  // Derived lists
  const filtered = useMemo(() => {
    const items = credentialQ.data ?? []
    if (!filter.trim()) return items
    const f = filter.toLowerCase()
    return items.filter((c: any) =>
      [
        c.name,
        c.credential_provider,
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

  // -------------------- UI --------------------

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

  // Create form watchers
  const credential_provider = createForm.watch("credential_provider")
  const kind = createForm.watch("kind")
  const scopeKind = createForm.watch("scope_kind")

  const setCreateScope = (obj: any) =>
    createForm.setValue("scope", obj, { shouldDirty: true, shouldValidate: true })
  const setCreateSecret = (obj: any) =>
    createForm.setValue("secret", obj, { shouldDirty: true, shouldValidate: true })

  function ensureCreateDefaultsForSecret() {
    if (useRawSecretJSON) return

    if (credential_provider === "aws" && kind === "aws_access_key") {
      const s = createForm.getValues("secret") ?? {}
      setCreateSecret({
        access_key_id: s.access_key_id ?? "",
        secret_access_key: s.secret_access_key ?? "",
      })
    } else if (kind === "api_token") {
      const s = createForm.getValues("secret") ?? {}
      setCreateSecret({ token: s.token ?? "" })
    } else if (kind === "basic_auth") {
      const s = createForm.getValues("secret") ?? {}
      setCreateSecret({ username: s.username ?? "", password: s.password ?? "" })
    } else if (kind === "oauth2") {
      const s = createForm.getValues("secret") ?? {}
      setCreateSecret({
        client_id: s.client_id ?? "",
        client_secret: s.client_secret ?? "",
        refresh_token: s.refresh_token ?? "",
      })
    }
  }

  function onChangeCreateScopeKind(next: "provider" | "service" | "resource") {
    createForm.setValue("scope_kind", next, { shouldDirty: true, shouldValidate: true })
    if (next === "provider") setCreateScope({})
    if (next === "service") setCreateScope({ service: "route53" as AwsSvc })
    if (next === "resource") setCreateScope({ arn: "" })
  }

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <div>
          <h1 className="mb-1 text-2xl font-bold">Credentials</h1>
          <p className="text-muted-foreground text-sm">
            Store provider credentials. Secrets are encrypted server-side; revealing is a one-time
            read.
          </p>
        </div>

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
            <DialogContent className="sm:max-w-2xl">
              <DialogHeader>
                <DialogTitle>Create Credential</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form
                  onSubmit={createForm.handleSubmit((values) => {
                    const parsed = createCredentialSchema.safeParse(values)
                    if (!parsed.success) {
                      toast.error("Please fix validation errors")
                      return
                    }
                    createMutation.mutate(parsed.data)
                  })}
                  className="space-y-5 pt-2"
                >
                  <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                    <FormField
                      control={createForm.control}
                      name="credential_provider"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Provider</FormLabel>
                          <Select
                            onValueChange={(v) => {
                              field.onChange(v)
                              ensureCreateDefaultsForSecret()
                            }}
                            defaultValue={field.value}
                          >
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
                          <Select
                            onValueChange={(v) => {
                              field.onChange(v)
                              ensureCreateDefaultsForSecret()
                            }}
                            defaultValue={field.value}
                          >
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
                          <Select
                            onValueChange={(v: "provider" | "service" | "resource") => {
                              onChangeCreateScopeKind(v)
                            }}
                            defaultValue={field.value}
                          >
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

                  {/* Scope UI (create) */}
                  {!isProviderScope({ scope_kind: scopeKind }) && (
                    <>
                      {isAwsServiceScope({ credential_provider, scope_kind: scopeKind }) ? (
                        <FormItem>
                          <FormLabel>Service</FormLabel>
                          <Controller
                            control={createForm.control}
                            name="scope"
                            render={({ field }) => (
                              <Select
                                onValueChange={(svc) => field.onChange({ service: svc })}
                                value={(field.value as any)?.service ?? "route53"}
                              >
                                <FormControl>
                                  <SelectTrigger>
                                    <SelectValue placeholder="Select an AWS service" />
                                  </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                  {AWS_ALLOWED_SERVICES.map((s) => (
                                    <SelectItem key={s} value={s}>
                                      {s.toUpperCase()}
                                    </SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                            )}
                          />
                          <p className="text-muted-foreground mt-1 text-xs">
                            Must be one of: {AWS_ALLOWED_SERVICES.join(", ")}.
                          </p>
                        </FormItem>
                      ) : isAwsResourceScope({ credential_provider, scope_kind: scopeKind }) ? (
                        <FormItem>
                          <FormLabel>Resource ARN</FormLabel>
                          <Controller
                            control={createForm.control}
                            name="scope"
                            render={({ field }) => (
                              <Input
                                value={(field.value as any)?.arn ?? ""}
                                onChange={(e) => field.onChange({ arn: e.target.value })}
                                placeholder="arn:aws:service:region:account:resource"
                              />
                            )}
                          />
                          <FormMessage />
                        </FormItem>
                      ) : (
                        <FormField
                          control={createForm.control}
                          name="scope"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Scope (JSON)</FormLabel>
                              <Textarea
                                value={pretty(field.value ?? {})}
                                onChange={(e) => {
                                  try {
                                    const obj = JSON.parse(e.target.value)
                                    field.onChange(obj)
                                  } catch {
                                    field.onChange(e.target.value)
                                  }
                                }}
                                rows={3}
                                placeholder='{"service":"route53"} or {"arn":"arn:aws:..."}'
                                className="font-mono"
                              />
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      )}
                    </>
                  )}

                  {/* Secret UI (create) */}
                  <div className="flex items-center gap-2">
                    <Switch
                      checked={useRawSecretJSON}
                      onCheckedChange={(v) => {
                        setUseRawSecretJSON(v)
                        ensureCreateDefaultsForSecret()
                      }}
                      id="raw-secret-toggle"
                    />
                    <label htmlFor="raw-secret-toggle" className="text-sm">
                      Edit secret as raw JSON
                    </label>
                  </div>

                  {useRawSecretJSON ? (
                    <FormField
                      control={createForm.control}
                      name="secret"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Secret (JSON)</FormLabel>
                          <Textarea
                            value={pretty(field.value ?? {})}
                            onChange={(e) => {
                              try {
                                field.onChange(JSON.parse(e.target.value))
                              } catch {
                                field.onChange(e.target.value)
                              }
                            }}
                            rows={6}
                            placeholder={
                              kind === "aws_access_key"
                                ? '{"access_key_id":"...","secret_access_key":"..."}'
                                : kind === "api_token"
                                  ? '{"token":"..."}'
                                  : kind === "basic_auth"
                                    ? '{"username":"...","password":"..."}'
                                    : '{"client_id":"...","client_secret":"...","refresh_token":"..."}'
                            }
                            className="font-mono"
                          />
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  ) : (
                    <>
                      {credential_provider === "aws" && kind === "aws_access_key" && (
                        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                          <FormItem>
                            <FormLabel>Access Key ID</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  value={(field.value ?? {}).access_key_id ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      access_key_id: e.target.value.trim(),
                                    })
                                  }
                                  placeholder="AKIA..."
                                />
                              )}
                            />
                          </FormItem>
                          <FormItem>
                            <FormLabel>Secret Access Key</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  type="password"
                                  value={(field.value ?? {}).secret_access_key ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      secret_access_key: e.target.value,
                                    })
                                  }
                                  placeholder="•••••••••••••••"
                                />
                              )}
                            />
                          </FormItem>
                        </div>
                      )}

                      {kind === "api_token" && (
                        <FormItem>
                          <FormLabel>API Token</FormLabel>
                          <Controller
                            control={createForm.control}
                            name="secret"
                            render={({ field }) => (
                              <Input
                                value={(field.value ?? {}).token ?? ""}
                                onChange={(e) =>
                                  setCreateSecret({ ...(field.value ?? {}), token: e.target.value })
                                }
                                placeholder="token..."
                              />
                            )}
                          />
                        </FormItem>
                      )}

                      {kind === "basic_auth" && (
                        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                          <FormItem>
                            <FormLabel>Username</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  value={(field.value ?? {}).username ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      username: e.target.value,
                                    })
                                  }
                                />
                              )}
                            />
                          </FormItem>
                          <FormItem>
                            <FormLabel>Password</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  type="password"
                                  value={(field.value ?? {}).password ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      password: e.target.value,
                                    })
                                  }
                                />
                              )}
                            />
                          </FormItem>
                        </div>
                      )}

                      {kind === "oauth2" && (
                        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                          <FormItem>
                            <FormLabel>Client ID</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  value={(field.value ?? {}).client_id ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      client_id: e.target.value,
                                    })
                                  }
                                />
                              )}
                            />
                          </FormItem>
                          <FormItem>
                            <FormLabel>Client Secret</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  type="password"
                                  value={(field.value ?? {}).client_secret ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      client_secret: e.target.value,
                                    })
                                  }
                                  placeholder="••••••••••"
                                />
                              )}
                            />
                          </FormItem>
                          <FormItem>
                            <FormLabel>Refresh Token</FormLabel>
                            <Controller
                              control={createForm.control}
                              name="secret"
                              render={({ field }) => (
                                <Input
                                  value={(field.value ?? {}).refresh_token ?? ""}
                                  onChange={(e) =>
                                    setCreateSecret({
                                      ...(field.value ?? {}),
                                      refresh_token: e.target.value,
                                    })
                                  }
                                />
                              )}
                            />
                          </FormItem>
                        </div>
                      )}
                    </>
                  )}

                  <DialogFooter className="gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => {
                        const parsed = createCredentialSchema.safeParse(createForm.getValues())
                        if (!parsed.success) {
                          toast.error("Fix validation errors before previewing")
                          return
                        }
                        const body = buildCreateBody(parsed.data)
                        setPreviewCreateBody(body)
                        setPreviewCreateOpen(true)
                      }}
                    >
                      Preview request
                    </Button>

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
              <th className="w-[26%] px-4 py-2 text-left">Name</th>
              <th className="px-4 py-2 text-left">Provider</th>
              <th className="px-4 py-2 text-left">Kind</th>
              <th className="px-4 py-2 text-left">Scope</th>
              <th className="px-4 py-2 text-left">Account</th>
              <th className="px-4 py-2 text-left">Region</th>
              <th className="px-4 py-2 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.map((row: any) => (
              <tr key={row.id} className="border-t">
                <td className="px-4 py-2">
                  <div className="font-medium">{row.name}</div>
                  <div className="text-muted-foreground text-xs">
                    <span className="mr-1">id:</span>
                    <code className="bg-muted rounded px-1">{row.id.slice(0, 8)}…</code>
                  </div>
                </td>
                <td className="px-4 py-2">{row.credential_provider}</td>
                <td className="px-4 py-2">{row.kind}</td>
                <td className="px-4 py-2">
                  <Badge variant="secondary">{row.scope_kind}</Badge>
                </td>
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
                            recoverable.
                          </AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel disabled={deleteMutation.isPending}>
                            Cancel
                          </AlertDialogCancel>
                          <AlertDialogAction
                            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                            onClick={() => deleteMutation.mutate(row.id)}
                            disabled={deleteMutation.isPending}
                          >
                            {deleteMutation.isPending && (
                              <Loader2 className="mr-2 inline h-4 w-4 animate-spin" />
                            )}
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
                <td colSpan={7} className="px-4 py-12 text-center">
                  <div className="mx-auto max-w-md">
                    <div className="mb-2 flex items-center justify-center">
                      <AlertTriangle className="text-muted-foreground h-5 w-5" />
                    </div>
                    <p className="text-muted-foreground">No credentials match your search.</p>
                  </div>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Edit dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>Edit Credential</DialogTitle>
          </DialogHeader>

          <Form {...editForm}>
            <form
              onSubmit={editForm.handleSubmit((values) => {
                if (!editingId) return
                const parsed = updateCredentialSchema.safeParse(values)
                if (!parsed.success) {
                  toast.error("Please fix validation errors")
                  return
                }
                updateMutation.mutate({ id: editingId, body: parsed.data })
              })}
              className="space-y-5 pt-2"
            >
              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <FormField
                  control={editForm.control}
                  name="credential_provider"
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
                    <Textarea
                      value={pretty(
                        field.value ??
                          (editForm.getValues("scope_kind") === "provider" ? {} : undefined)
                      )}
                      onChange={(e) => {
                        try {
                          field.onChange(JSON.parse(e.target.value))
                        } catch {
                          field.onChange(e.target.value)
                        }
                      }}
                      rows={3}
                      className="font-mono"
                    />
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Rotate secret */}
              <div className="flex items-center gap-2">
                <Switch
                  checked={useRawEditSecretJSON}
                  onCheckedChange={setUseRawEditSecretJSON}
                  id="raw-edit-secret-toggle"
                />
                <label htmlFor="raw-edit-secret-toggle" className="text-sm">
                  Rotate secret with raw JSON (leave empty to keep existing)
                </label>
              </div>

              {useRawEditSecretJSON && (
                <FormField
                  control={editForm.control}
                  name="secret"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Rotate Secret (JSON)</FormLabel>
                      <Textarea
                        value={typeof field.value === "string" ? field.value : pretty(field.value)}
                        onChange={(e) => {
                          try {
                            field.onChange(JSON.parse(e.target.value))
                          } catch {
                            field.onChange(e.target.value)
                          }
                        }}
                        rows={6}
                        className="font-mono"
                        placeholder='{"access_key_id":"...","secret_access_key":"..."}'
                      />
                      <FormMessage />
                    </FormItem>
                  )}
                />
              )}

              <DialogFooter className="gap-2">
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => {
                    const parsed = updateCredentialSchema.safeParse(editForm.getValues())
                    if (!parsed.success) {
                      toast.error("Fix validation errors before previewing")
                      return
                    }
                    const body = buildUpdateBody(parsed.data)
                    setPreviewUpdateBody(body)
                    setPreviewUpdateOpen(true)
                  }}
                >
                  Preview request
                </Button>

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
            <DialogTitle className="flex items-center gap-2">
              <Eye className="h-4 w-4" /> Decrypted Secret
            </DialogTitle>
          </DialogHeader>
          <div className="bg-muted/40 rounded-lg border p-3">
            <pre className="max-h-[50vh] overflow-auto text-xs leading-relaxed">
              {pretty(revealJson ?? {})}
            </pre>
          </div>
          <div className="text-muted-foreground flex items-center gap-2 text-xs">
            <AlertTriangle className="h-3.5 w-3.5" />
            One-time read. Close this dialog to hide the secret.
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

      {/* Preview CREATE modal */}
      <Dialog open={previewCreateOpen} onOpenChange={setPreviewCreateOpen}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>Preview: POST /api/v1/credentials</DialogTitle>
          </DialogHeader>
          <div className="bg-muted/40 rounded-lg border p-3">
            <pre className="max-h-[50vh] overflow-auto text-xs leading-relaxed">
              {pretty(previewCreateBody ?? {})}
            </pre>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                navigator.clipboard.writeText(pretty(previewCreateBody ?? {}))
                toast.success("Copied body")
              }}
            >
              Copy body
            </Button>
            <Button onClick={() => setPreviewCreateOpen(false)}>Close</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Preview UPDATE modal */}
      <Dialog open={previewUpdateOpen} onOpenChange={setPreviewUpdateOpen}>
        <DialogContent className="sm:max-w-2xl">
          <DialogHeader>
            <DialogTitle>Preview: PATCH /api/v1/credentials/:id</DialogTitle>
          </DialogHeader>
          <div className="bg-muted/40 rounded-lg border p-3">
            <pre className="max-h-[50vh] overflow-auto text-xs leading-relaxed">
              {pretty(previewUpdateBody ?? {})}
            </pre>
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                navigator.clipboard.writeText(pretty(previewUpdateBody ?? {}))
                toast.success("Copied body")
              }}
            >
              Copy body
            </Button>
            <Button onClick={() => setPreviewUpdateOpen(false)}>Close</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
