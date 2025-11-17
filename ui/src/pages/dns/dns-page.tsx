import { useEffect, useMemo, useState } from "react"
import { credentialsApi } from "@/api/credentials"
import { dnsApi } from "@/api/dns"
import type {
  DtoCreateDomainRequest,
  DtoCreateRecordSetRequest,
  DtoCredentialOut,
  DtoDomainResponse,
  DtoRecordSetResponse,
  DtoUpdateDomainRequest,
  DtoUpdateRecordSetRequest,
} from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import {
  AlertTriangle,
  CheckCircle2,
  Circle,
  Loader2,
  MoreHorizontal,
  Pencil,
  Plus,
  Search,
  Trash2,
} from "lucide-react"
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
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card } from "@/components/ui/card"
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

// ---------- helpers ----------

const statusIcon = (s?: string) => {
  switch (s) {
    case "ready":
      return <CheckCircle2 className="h-4 w-4 text-emerald-600" />
    case "provisioning":
      return <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
    case "failed":
      return <AlertTriangle className="h-4 w-4 text-red-600" />
    default:
      return <Circle className="text-muted-foreground h-4 w-4" />
  }
}

const StatusBadge = ({ s }: { s?: string }) => (
  <Badge
    variant={s === "failed" ? "destructive" : s === "ready" ? "default" : "secondary"}
    className="gap-1"
    title={s}
  >
    {statusIcon(s)}
    <span className="capitalize">{s ?? "pending"}</span>
  </Badge>
)

const parseCommaList = (v: string) =>
  v
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean)

const joinCommaList = (arr?: string[] | null) => (arr && arr.length ? arr.join(",") : "")

const rrtypes = ["A", "AAAA", "CNAME", "TXT", "MX", "NS", "SRV", "CAA"]

const isR53 = (c: DtoCredentialOut) =>
  c.provider === "aws" &&
  c.scope_kind === "service" &&
  (() => {
    const s = (c as any).scope
    try {
      const obj = typeof s === "string" ? JSON.parse(s) : s || {}
      return obj?.service === "route53"
    } catch {
      return false
    }
  })()

const credLabel = (c: DtoCredentialOut) => {
  const bits = [c.name || "Unnamed", c.account_id, c.region].filter(Boolean)
  return bits.join(" · ")
}

// ---------- zod schemas ----------

const createDomainSchema = z.object({
  domain_name: z
    .string()
    .min(1, "Domain is required")
    .max(253)
    .transform((s) => s.trim().replace(/\.$/, "").toLowerCase()),
  credential_id: z.string().uuid("Pick a credential"),
  zone_id: z
    .string()
    .optional()
    .or(z.literal(""))
    .transform((v) => (v ? v.trim() : undefined)),
})
type CreateDomainValues = z.input<typeof createDomainSchema>

const updateDomainSchema = createDomainSchema.partial()
type UpdateDomainValues = z.infer<typeof updateDomainSchema>

const ttlSchema = z
  .union([
    z.number(),
    z
      .string()
      .regex(/^\d+$/)
      .transform((s) => Number(s)),
  ])
  .optional()
  .refine((v) => v === undefined || (v >= 1 && v <= 86400), {
    message: "TTL must be between 1 and 86400",
  })

const createRecordSchema = z
  .object({
    name: z
      .string()
      .min(1, "Name required")
      .max(253)
      .transform((s) => s.trim().replace(/\.$/, "").toLowerCase()),
    type: z.enum(rrtypes as [string, ...string[]]),
    ttl: ttlSchema,
    valuesCsv: z.string().optional(),
  })
  .superRefine((vals, ctx) => {
    const arr = parseCommaList(vals.valuesCsv ?? "")
    if (arr.length === 0) {
      ctx.addIssue({ code: "custom", message: "At least one value is required" })
    }
    if (vals.type === "CNAME" && arr.length !== 1) {
      ctx.addIssue({ code: "custom", message: "CNAME requires exactly one value" })
    }
  })
type CreateRecordValues = z.input<typeof createRecordSchema>

const updateRecordSchema = createRecordSchema.partial()
type UpdateRecordValues = z.input<typeof updateRecordSchema>

// ---------- main ----------

export const DnsPage = () => {
  const [filter, setFilter] = useState("")
  const [selected, setSelected] = useState<DtoDomainResponse | null>(null)

  const [createDomOpen, setCreateDomOpen] = useState(false)
  const [editDomOpen, setEditDomOpen] = useState(false)

  const [createRecOpen, setCreateRecOpen] = useState(false)
  const [editRecOpen, setEditRecOpen] = useState(false)
  const [editingRecord, setEditingRecord] = useState<DtoRecordSetResponse | null>(null)

  const qc = useQueryClient()

  // ---- queries ----

  const domainsQ = useQuery({
    queryKey: ["dns", "domains"],
    queryFn: () => dnsApi.listDomains(),
  })

  const recordsQ = useQuery({
    queryKey: ["dns", "records", selected?.id],
    queryFn: async () => {
      if (!selected) return []
      return await dnsApi.listRecordSetsByDomain(selected.id as string)
    },
    enabled: !!selected?.id,
  })

  const credentialQ = useQuery({
    queryKey: ["credentials", "r53"],
    queryFn: () => credentialsApi.listCredentials(),
  })

  const r53Credentials = useMemo(() => (credentialQ.data ?? []).filter(isR53), [credentialQ.data])

  useEffect(() => {
    const setSelectedDns = () => {
      if (!selected && domainsQ.data && domainsQ.data.length) {
        setSelected(domainsQ.data[0]!)
      }
    }
    setSelectedDns()
  }, [domainsQ.data, selected])

  const filteredDomains = useMemo(() => {
    const list: DtoDomainResponse[] = domainsQ.data ?? []
    if (!filter.trim()) return list
    const f = filter.toLowerCase()
    return list.filter((d) =>
      [d.domain_name, d.zone_id, d.status, d.domain_name]
        .filter(Boolean)
        .map((x) => String(x).toLowerCase())
        .some((s) => s.includes(f))
    )
  }, [domainsQ.data, filter])

  // ---- mutations: domains ----

  const createDomainForm = useForm<CreateDomainValues>({
    resolver: zodResolver(createDomainSchema),
    defaultValues: {
      domain_name: "",
      credential_id: "",
      zone_id: "",
    },
  })

  const createDomainMut = useMutation({
    mutationFn: (v: CreateDomainValues) =>
      dnsApi.createDomain(v as unknown as DtoCreateDomainRequest),
    onSuccess: async (d) => {
      toast.success("Domain created")
      setCreateDomOpen(false)
      createDomainForm.reset()
      await qc.invalidateQueries({ queryKey: ["dns", "domains"] })
      setSelected(d as DtoDomainResponse)
    },
    onError: (e: any) =>
      toast.error("Failed to create domain", { description: e?.message ?? "Unknown error" }),
  })

  const editDomainForm = useForm<UpdateDomainValues>({
    resolver: zodResolver(updateDomainSchema),
  })

  const openEditDomain = (d: DtoDomainResponse) => {
    setSelected(d)
    editDomainForm.reset({
      domain_name: d.domain_name,
      credential_id: d.credential_id,
      zone_id: d.zone_id || "",
    })
    setEditDomOpen(true)
  }

  const updateDomainMut = useMutation({
    mutationFn: (vals: UpdateDomainValues) => {
      if (!selected) throw new Error("No domain selected")
      return dnsApi.updateDomain(selected.id!, vals as unknown as DtoUpdateDomainRequest)
    },
    onSuccess: async () => {
      toast.success("Domain updated")
      setEditDomOpen(false)
      await qc.invalidateQueries({ queryKey: ["dns", "domains"] })
      await qc.invalidateQueries({ queryKey: ["dns", "records", selected?.id] })
    },
    onError: (e: any) =>
      toast.error("Failed to update domain", { description: e?.message ?? "Unknown error" }),
  })

  const deleteDomainMut = useMutation({
    mutationFn: (id: string) => dnsApi.deleteDomain(id),
    onSuccess: async () => {
      toast.success("Domain deleted")
      await qc.invalidateQueries({ queryKey: ["dns", "domains"] })
      setSelected(null)
    },
    onError: (e: any) =>
      toast.error("Failed to delete domain", { description: e?.message ?? "Unknown error" }),
  })

  // ---- mutations: record sets ----

  const createRecForm = useForm<CreateRecordValues>({
    resolver: zodResolver(createRecordSchema),
    defaultValues: {
      name: "",
      type: "A",
      ttl: 300,
      valuesCsv: "",
    },
  })

  const explainError = (e: any) => {
    const msg: string = e?.response?.data?.error || e?.message || "Unknown error"
    if (msg.includes("ownership_conflict")) {
      return "Ownership conflict: this (name,type) exists but isn’t owned by autoglue."
    }
    if (msg.includes("already_exists")) {
      return "A record with this (name,type) already exists. Use Edit instead."
    }
    return msg
  }

  const createRecordMut = useMutation({
    mutationFn: async (vals: CreateRecordValues) => {
      if (!selected) throw new Error("No domain selected")
      const body: DtoCreateRecordSetRequest = {
        name: vals.name,
        type: vals.type,
        // omit ttl when empty/undefined
        ...(vals.ttl ? { ttl: vals.ttl as unknown as number } : {}),
        values: parseCommaList(vals.valuesCsv ?? ""),
      }
      return dnsApi.createRecordSetsByDomain(selected.id!, body)
    },
    onSuccess: async () => {
      toast.success("Record set created")
      setCreateRecOpen(false)
      createRecForm.reset()
      await qc.invalidateQueries({ queryKey: ["dns", "records", selected?.id] })
    },
    onError: (e: any) =>
      toast.error("Failed to create record set", { description: explainError(e) }),
  })

  const editRecForm = useForm<UpdateRecordValues>({
    resolver: zodResolver(updateRecordSchema),
  })

  const openEditRecord = (r: DtoRecordSetResponse) => {
    setEditingRecord(r)
    const values = (r.values as any) || []
    editRecForm.reset({
      name: r.name,
      type: r.type,
      ttl: r.ttl ? Number(r.ttl) : undefined,
      valuesCsv: joinCommaList(values),
    })
    setEditRecOpen(true)
  }

  const updateRecordMut = useMutation({
    mutationFn: async (vals: UpdateRecordValues) => {
      if (!editingRecord) throw new Error("No record selected")
      const body: DtoUpdateRecordSetRequest = {}
      if (vals.name !== undefined) body.name = vals.name
      if (vals.type !== undefined) body.type = vals.type
      if (vals.ttl !== undefined && vals.ttl !== null) {
        // if blank string came through it would have been filtered; when undefined, omit
        body.ttl = vals.ttl as unknown as number | undefined
      }
      if (vals.valuesCsv !== undefined) {
        body.values = parseCommaList(vals.valuesCsv)
      }
      return dnsApi.updateRecordSetsByDomain(editingRecord.id!, body)
    },
    onSuccess: async () => {
      toast.success("Record set updated")
      setEditRecOpen(false)
      setEditingRecord(null)
      await qc.invalidateQueries({ queryKey: ["dns", "records", selected?.id] })
    },
    onError: (e: any) =>
      toast.error("Failed to update record set", { description: explainError(e) }),
  })

  const deleteRecordMut = useMutation({
    mutationFn: (id: string) => dnsApi.deleteRecordSetsByDomain(id),
    onSuccess: async () => {
      toast.success("Record set deleted")
      await qc.invalidateQueries({ queryKey: ["dns", "records", selected?.id] })
    },
    onError: (e: any) =>
      toast.error("Failed to delete record set", { description: e?.message ?? "Unknown error" }),
  })

  // ---------- UI ----------

  return (
    <div className="space-y-5 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="text-2xl font-bold">DNS</h1>
        <div className="flex items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search domains…"
              className="w-64 pl-8"
            />
          </div>
          <Dialog open={createDomOpen} onOpenChange={setCreateDomOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateDomOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Add Domain
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Add Domain</DialogTitle>
              </DialogHeader>
              <Form {...createDomainForm}>
                <form
                  className="space-y-4 pt-2"
                  onSubmit={createDomainForm.handleSubmit((v) => createDomainMut.mutate(v))}
                >
                  <FormField
                    control={createDomainForm.control}
                    name="domain_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Domain</FormLabel>
                        <FormControl>
                          <Input {...field} placeholder="example.com" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* CREDENTIAL SELECT (Create) */}
                  <FormField
                    control={createDomainForm.control}
                    name="credential_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Route53 Credential</FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                          disabled={credentialQ.isLoading || (r53Credentials?.length ?? 0) === 0}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue
                                placeholder={
                                  credentialQ.isLoading
                                    ? "Loading…"
                                    : (r53Credentials?.length ?? 0) === 0
                                      ? "No Route53 credentials found"
                                      : "Select credential"
                                }
                              />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {(r53Credentials ?? []).map((c) => (
                              <SelectItem key={c.id} value={c.id!}>
                                {credLabel(c)}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {credentialQ.error && (
                          <p className="text-destructive text-xs">Failed to load credentials.</p>
                        )}
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createDomainForm.control}
                    name="zone_id"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Zone ID (optional)</FormLabel>
                        <FormControl>
                          <Input {...field} placeholder="/hostedzone/Z123…" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateDomOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createDomainMut.isPending}>
                      {createDomainMut.isPending && (
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

      {/* domains panel */}
      <div>
        <Card className="p-3 md:col-span-5">
          <div className="mb-2 flex items-center justify-between">
            <div className="text-sm font-semibold">Domains</div>
            {domainsQ.isFetching && <Loader2 className="h-4 w-4 animate-spin" />}
          </div>
          <div className="max-h-[60vh] overflow-auto rounded-md border">
            <table className="min-w-full text-sm">
              <thead className="bg-muted/40 text-xs tracking-wide uppercase">
                <tr>
                  <th className="px-3 py-2 text-left">Domain</th>
                  <th className="px-3 py-2 text-left">Zone</th>
                  <th className="px-3 py-2 text-left">Status</th>
                  <th className="px-3 py-2 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {(filteredDomains ?? []).map((d) => (
                  <tr
                    key={d.id}
                    className={`hover:bg-muted/30 border-t ${
                      selected?.id === d.id ? "bg-muted/40" : ""
                    }`}
                    onClick={() => setSelected(d)}
                  >
                    <td className="cursor-pointer px-3 py-2 font-medium">{d.domain_name}</td>
                    <td className="px-3 py-2">{d.zone_id || "—"}</td>
                    <td className="px-3 py-2">
                      <StatusBadge s={d.status} />
                    </td>
                    <td className="px-3 py-2">
                      <div className="flex items-center justify-end gap-2">
                        <Button size="icon" variant="ghost" onClick={() => openEditDomain(d)}>
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button
                              size="icon"
                              variant="ghost"
                              onClick={(e) => e.stopPropagation()}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Delete “{d.domain_name}”?</AlertDialogTitle>
                              <AlertDialogDescription>
                                This deletes the domain metadata. External DNS records are not
                                touched.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancel</AlertDialogCancel>
                              <AlertDialogAction
                                className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                onClick={() => deleteDomainMut.mutate(d.id!)}
                              >
                                Delete
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      </div>
                    </td>
                  </tr>
                ))}
                {(!filteredDomains || filteredDomains.length === 0) && (
                  <tr>
                    <td colSpan={4} className="text-muted-foreground px-3 py-8 text-center">
                      No domains yet.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </Card>
      </div>

      <div>
        {/* records panel */}
        <Card className="p-3 md:col-span-7">
          <div className="mb-2 flex items-center justify-between">
            <div className="text-sm font-semibold">
              Records {selected ? `— ${selected.domain_name}` : ""}
            </div>
            <div className="flex items-center gap-2">
              <StatusBadge s={selected?.status} />
              <Dialog open={createRecOpen} onOpenChange={setCreateRecOpen}>
                <DialogTrigger asChild>
                  <Button disabled={!selected}>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Record
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-xl">
                  <DialogHeader>
                    <DialogTitle>Add Record</DialogTitle>
                  </DialogHeader>
                  <Form {...createRecForm}>
                    <form
                      className="space-y-4 pt-2"
                      onSubmit={createRecForm.handleSubmit((v) => createRecordMut.mutate(v))}
                    >
                      <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                        <FormField
                          control={createRecForm.control}
                          name="name"
                          render={({ field }) => (
                            <FormItem className="md:col-span-1">
                              <FormLabel>Name</FormLabel>
                              <FormControl>
                                <Input {...field} placeholder="endpoint (or @)" />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={createRecForm.control}
                          name="type"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Type</FormLabel>
                              <Select
                                onValueChange={field.onChange}
                                defaultValue={field.value as string}
                              >
                                <FormControl>
                                  <SelectTrigger>
                                    <SelectValue />
                                  </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                  {rrtypes.map((t) => (
                                    <SelectItem key={t} value={t}>
                                      {t}
                                    </SelectItem>
                                  ))}
                                </SelectContent>
                              </Select>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={createRecForm.control}
                          name="ttl"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>TTL (sec, optional)</FormLabel>
                              <FormControl>
                                <Input
                                  type="number"
                                  value={field.value as number | undefined}
                                  onChange={(e) =>
                                    field.onChange(
                                      e.target.value === "" ? undefined : Number(e.target.value)
                                    )
                                  }
                                  placeholder="300"
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <FormField
                        control={createRecForm.control}
                        name="valuesCsv"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Values (comma-separated)</FormLabel>
                            <FormControl>
                              <Textarea
                                {...field}
                                rows={3}
                                placeholder='e.g. 10.0.30.1, 10.0.30.2 or "v=spf1 ~all"'
                                className="font-mono"
                              />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <DialogFooter className="gap-2">
                        <Button
                          variant="outline"
                          type="button"
                          onClick={() => setCreateRecOpen(false)}
                        >
                          Cancel
                        </Button>
                        <Button type="submit" disabled={createRecordMut.isPending}>
                          {createRecordMut.isPending && (
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

          <div className="overflow-x-auto rounded-md border">
            {recordsQ.isLoading && (
              <div className="flex items-center gap-2 p-4">
                <Loader2 className="h-4 w-4 animate-spin" /> Loading records…
              </div>
            )}
            {!recordsQ.isLoading && (
              <table className="min-w-full text-sm">
                <thead className="bg-muted/40 text-xs tracking-wide uppercase">
                  <tr>
                    <th className="px-3 py-2 text-left">Name</th>
                    <th className="px-3 py-2 text-left">Type</th>
                    <th className="px-3 py-2 text-left">TTL</th>
                    <th className="px-3 py-2 text-left">Values</th>
                    <th className="px-3 py-2 text-left">Status</th>
                    <th className="px-3 py-2 text-left">Owner</th>
                    <th className="px-3 py-2 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {(recordsQ.data ?? []).map((r) => {
                    const values = (r.values as any) || []
                    return (
                      <tr key={r.id} className="border-t">
                        <td className="px-3 py-2 font-medium">{r.name || "@"}</td>
                        <td className="px-3 py-2">{r.type}</td>
                        <td className="px-3 py-2">{r.ttl ?? "—"}</td>
                        <td className="px-3 py-2">
                          <div className="max-w-[420px] truncate" title={(values || []).join(", ")}>
                            {(values || []).join(", ")}
                          </div>
                        </td>
                        <td className="px-3 py-2">
                          <StatusBadge s={r.status} />
                        </td>
                        <td className="px-3 py-2">{r.owner}</td>
                        <td className="px-3 py-2">
                          <div className="flex items-center justify-end gap-2">
                            <Button size="icon" variant="ghost" onClick={() => openEditRecord(r)}>
                              <Pencil className="h-4 w-4" />
                            </Button>

                            <AlertDialog>
                              <AlertDialogTrigger asChild>
                                <Button size="icon" variant="ghost">
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </AlertDialogTrigger>
                              <AlertDialogContent>
                                <AlertDialogHeader>
                                  <AlertDialogTitle>
                                    Delete “{r.name || "@"} {r.type}”?
                                  </AlertDialogTitle>
                                  <AlertDialogDescription>
                                    This removes the record set from your project. Your worker does
                                    not delete it from the DNS provider right now.
                                  </AlertDialogDescription>
                                </AlertDialogHeader>
                                <AlertDialogFooter>
                                  <AlertDialogCancel>Cancel</AlertDialogCancel>
                                  <AlertDialogAction
                                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                    onClick={() => deleteRecordMut.mutate(r.id!)}
                                  >
                                    Delete
                                  </AlertDialogAction>
                                </AlertDialogFooter>
                              </AlertDialogContent>
                            </AlertDialog>

                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => openEditRecord(r)}>
                                  Edit
                                </DropdownMenuItem>
                                <DropdownMenuItem
                                  className="text-destructive"
                                  onClick={() => deleteRecordMut.mutate(r.id!)}
                                >
                                  Delete
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </td>
                      </tr>
                    )
                  })}
                  {(!recordsQ.data || recordsQ.data.length === 0) && (
                    <tr>
                      <td colSpan={7} className="text-muted-foreground px-3 py-8 text-center">
                        {selected
                          ? "No records yet — add one."
                          : "Select a domain to view records."}
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            )}
          </div>
        </Card>
      </div>

      {/* Edit Domain Dialog */}
      <Dialog open={editDomOpen} onOpenChange={setEditDomOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Domain</DialogTitle>
          </DialogHeader>
          <Form {...editDomainForm}>
            <form
              className="space-y-4 pt-2"
              onSubmit={editDomainForm.handleSubmit((v) => updateDomainMut.mutate(v))}
            >
              <FormField
                control={editDomainForm.control}
                name="domain_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Domain</FormLabel>
                    <FormControl>
                      <Input {...field} placeholder="example.com" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* CREDENTIAL SELECT (Edit) */}
              <FormField
                control={editDomainForm.control}
                name="credential_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Route53 Credential</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      value={field.value ?? ""}
                      disabled={credentialQ.isLoading || (r53Credentials?.length ?? 0) === 0}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={
                              credentialQ.isLoading
                                ? "Loading…"
                                : (r53Credentials?.length ?? 0) === 0
                                  ? "No Route53 credentials found"
                                  : "Select credential"
                            }
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {(r53Credentials ?? []).map((c) => (
                          <SelectItem key={c.id} value={c.id!}>
                            {credLabel(c)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {credentialQ.error && (
                      <p className="text-destructive text-xs">Failed to load credentials.</p>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={editDomainForm.control}
                name="zone_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Zone ID</FormLabel>
                    <FormControl>
                      <Input {...field} placeholder="/hostedzone/Z123…" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter className="gap-2">
                <Button variant="outline" type="button" onClick={() => setEditDomOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateDomainMut.isPending}>
                  {updateDomainMut.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save Changes
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Edit Record Dialog */}
      <Dialog open={editRecOpen} onOpenChange={setEditRecOpen}>
        <DialogContent className="sm:max-w-xl">
          <DialogHeader>
            <DialogTitle>Edit Record</DialogTitle>
          </DialogHeader>
          <Form {...editRecForm}>
            <form
              className="space-y-4 pt-2"
              onSubmit={editRecForm.handleSubmit((v) => updateRecordMut.mutate(v))}
            >
              <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
                <FormField
                  control={editRecForm.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem className="md:col-span-1">
                      <FormLabel>Name</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={editRecForm.control}
                  name="type"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Type</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value as string}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {rrtypes.map((t) => (
                            <SelectItem key={t} value={t}>
                              {t}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={editRecForm.control}
                  name="ttl"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>TTL (sec, optional)</FormLabel>
                      <FormControl>
                        <Input
                          type="number"
                          value={field.value as number | undefined}
                          onChange={(e) =>
                            field.onChange(
                              e.target.value === "" ? undefined : Number(e.target.value)
                            )
                          }
                          placeholder="300"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={editRecForm.control}
                name="valuesCsv"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Values (comma-separated)</FormLabel>
                    <FormControl>
                      <Textarea {...field} rows={3} className="font-mono" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button variant="outline" type="button" onClick={() => setEditRecOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateRecordMut.isPending}>
                  {updateRecordMut.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save Changes
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
