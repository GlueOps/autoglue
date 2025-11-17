import { useMemo, useState } from "react"
import { meApi } from "@/api/me.ts"
import { withRefresh } from "@/api/with-refresh.ts"
import { makeOrgsApi } from "@/sdkClient.ts"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
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
} from "@/components/ui/alert-dialog.tsx"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, } from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import { Label } from "@/components/ui/label.tsx"
import { Separator } from "@/components/ui/separator.tsx"
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow, } from "@/components/ui/table.tsx"

const orgsApi = makeOrgsApi()

const orgApi = {
  create: (body: { name: string; domain?: string }) =>
    withRefresh(async () => orgsApi.createOrg({ handlersOrgCreateReq: body })), // POST /orgs
}

const profileSchema = z.object({
  display_name: z.string().min(2, "Too short").max(100, "Too long"),
})

type ProfileForm = z.infer<typeof profileSchema>

const createKeySchema = z.object({
  name: z.string().min(2, "Too short").max(100, "Too long"),
  expires_in_hours: z.number().min(1).max(43800),
})

type CreateKeyForm = z.infer<typeof createKeySchema>

const createOrgSchema = z.object({
  name: z.string().min(2, "Too short").max(100, "Too long"),
  domain: z
    .string()
    .trim()
    .toLowerCase()
    .optional()
    .or(z.literal(""))
    .refine((v) => !v || /^[a-z0-9.-]+\.[a-z]{2,}$/i.test(v), "Invalid domain (e.g. example.com)"),
})

type CreateOrgForm = z.infer<typeof createOrgSchema>

export const MePage = () => {
  const qc = useQueryClient()

  const meQ = useQuery({
    queryKey: ["me"],
    queryFn: () => meApi.getMe(),
  })

  const form = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      display_name: "",
    },
    values: meQ.data ? { display_name: meQ.data.display_name ?? "" } : undefined,
  })

  const updateMut = useMutation({
    mutationFn: (values: ProfileForm) => meApi.updateMe(values),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["me"] })
      toast.success("Profile updated")
    },
    onError: (e: any) => toast.error(e?.message ?? "Update failed"),
  })

  const keysQ = useQuery({
    queryKey: ["me", "api-keys"],
    queryFn: () => meApi.listKeys(),
  })

  const [createOpen, setCreateOpen] = useState(false)
  const [justCreated, setJustCreated] = useState<ReturnType<typeof Object> | null>(null)

  const createForm = useForm<CreateKeyForm>({
    resolver: zodResolver(createKeySchema),
    defaultValues: {
      name: "",
      expires_in_hours: 720,
    },
  })

  const createMut = useMutation({
    mutationFn: (v: CreateKeyForm) =>
      meApi.createKey({
        name: v.name,
        expires_in_hours: v.expires_in_hours,
      } as CreateKeyForm),
    onSuccess: (resp: any) => {
      setJustCreated(resp)
      setCreateOpen(false)
      void qc.invalidateQueries({ queryKey: ["me", "api-keys"] })
      toast.success("API key created")
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed to create key"),
  })

  const [deleteId, setDeleteId] = useState<string | null>(null)

  const delMut = useMutation({
    mutationFn: (id: string) => meApi.deleteKey(id),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["me", "api-keys"] })
      setDeleteId(null)
      toast.success("Key deleted")
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed to delete key"),
  })

  const primaryEmail = useMemo(
    () => meQ.data?.emails?.find((e) => e.is_primary)?.email ?? meQ.data?.primary_email ?? "",
    [meQ.data]
  )

  // --- Create Org dialog + mutation ---
  const [orgOpen, setOrgOpen] = useState<boolean>(false)

  const orgForm = useForm<CreateOrgForm>({
    resolver: zodResolver(createOrgSchema),
    defaultValues: {
      name: "",
      domain: "",
    },
  })

  const orgCreateMut = useMutation({
    mutationFn: (v: CreateOrgForm) =>
      orgApi.create({
        name: v.name.trim(),
        domain: v.domain?.trim() ? v.domain.trim().toLowerCase() : undefined,
      }),
    onSuccess: () => {
      setOrgOpen(false)
      orgForm.reset()
      void qc.invalidateQueries({ queryKey: ["me"] })
      toast.success("Organization created")
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed to create organization"),
  })

  if (meQ.isLoading) return <div className="p-6">Loading…</div>
  if (meQ.error) return <div className="text-destructive p-6">Failed to load profile</div>

  return (
    <div className="space-y-6 p-6">
      {/* Profile */}
      <Card>
        <CardHeader>
          <CardTitle>Profile</CardTitle>
          <CardDescription>Manage your personal information.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-6 md:grid-cols-2">
            <div className="space-y-4">
              <div>
                <Label>Email</Label>
                <div className="text-muted-foreground mt-1 text-sm">{primaryEmail || "—"}</div>
              </div>

              <div>
                <Label>ID</Label>
                <div className="text-muted-foreground mt-1 text-sm">{meQ.data?.id || "—"}</div>
                <div className="text-muted-foreground mt-1 text-sm">
                  Share this ID with the organization owner of the Organization to join
                </div>
              </div>
              <Form {...form}>
                <form
                  className="space-y-4"
                  onSubmit={form.handleSubmit((v) => updateMut.mutate(v))}
                >
                  <FormField
                    control={form.control}
                    name="display_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Display name</FormLabel>
                        <FormControl>
                          <Input placeholder="Your name" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <Button type="submit" disabled={updateMut.isPending}>
                    Save
                  </Button>
                </form>
              </Form>
            </div>

            {/* Organizations + Create Org */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Organizations</Label>
                <Dialog open={orgOpen} onOpenChange={setOrgOpen}>
                  <DialogTrigger asChild>
                    <Button size="sm">New Organization</Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Create organization</DialogTitle>
                      <DialogDescription>
                        Give it a name, and optionally assign your company domain.
                      </DialogDescription>
                    </DialogHeader>

                    <Form {...orgForm}>
                      <form
                        className="space-y-4"
                        onSubmit={orgForm.handleSubmit((v) => orgCreateMut.mutate(v))}
                      >
                        <FormField
                          control={orgForm.control}
                          name="name"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Name</FormLabel>
                              <FormControl>
                                <Input placeholder="Acme Inc." {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <FormField
                          control={orgForm.control}
                          name="domain"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Domain (optional)</FormLabel>
                              <FormControl>
                                <Input placeholder="acme.com" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <DialogFooter>
                          <DialogClose asChild>
                            <Button type="button" variant="outline">
                              Cancel
                            </Button>
                          </DialogClose>
                          <Button type="submit" disabled={orgCreateMut.isPending}>
                            Create
                          </Button>
                        </DialogFooter>
                      </form>
                    </Form>
                  </DialogContent>
                </Dialog>
              </div>

              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Name</TableHead>
                      <TableHead>Domain</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {meQ.data?.organizations?.map((o) => (
                      <TableRow key={o.id}>
                        <TableCell>{o.name}</TableCell>
                        <TableCell>{(o as any).domain ?? "—"}</TableCell>
                      </TableRow>
                    ))}
                    {(!meQ.data?.organizations || meQ.data.organizations.length === 0) && (
                      <TableRow>
                        <TableCell colSpan={2} className="text-muted-foreground">
                          No organizations
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Separator />

      {/* API Keys (unchanged) */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle>User API Keys</CardTitle>
            <CardDescription>Personal keys for API access.</CardDescription>
          </div>
          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button>New Key</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Create API Key</DialogTitle>
                <DialogDescription>Give it a label and expiry.</DialogDescription>
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
                        <FormLabel>Label</FormLabel>
                        <FormControl>
                          <Input placeholder="CI script, local dev, ..." {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="expires_in_hours"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Expires in hours</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            inputMode="numeric"
                            step={1}
                            min={1}
                            placeholder="e.g. 720"
                            {...field}
                            onChange={(e) =>
                              field.onChange(e.target.value === "" ? "" : Number(e.target.value))
                            }
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter>
                    <DialogClose asChild>
                      <Button type="button" variant="outline">
                        Cancel
                      </Button>
                    </DialogClose>
                    <Button type="submit" disabled={createMut.isPending}>
                      Create
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </DialogContent>
          </Dialog>
        </CardHeader>

        <CardContent>
          <div className="overflow-x-auto rounded-md border">
            <Table className="text-sm">
              <TableCaption>Your user-scoped API keys.</TableCaption>
              <TableHeader>
                <TableRow>
                  <TableHead>Label</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Expires</TableHead>
                  <TableHead>Last used</TableHead>
                  <TableHead className="w-24" />
                </TableRow>
              </TableHeader>
              <TableBody>
                {keysQ.data?.map((k) => (
                  <TableRow key={k.id}>
                    <TableCell>{k.name ?? "—"}</TableCell>
                    <TableCell>{new Date(k.created_at!).toLocaleString()}</TableCell>
                    <TableCell>
                      {k.expires_at ? new Date(k.expires_at).toLocaleString() : "—"}
                    </TableCell>
                    <TableCell>
                      {k.last_used_at ? new Date(k.last_used_at).toLocaleString() : "—"}
                    </TableCell>
                    <TableCell className="text-right">
                      <AlertDialog
                        open={deleteId === k.id}
                        onOpenChange={(o) => !o && setDeleteId(null)}
                      >
                        <AlertDialogTrigger asChild>
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => setDeleteId(k.id!)}
                          >
                            Delete
                          </Button>
                        </AlertDialogTrigger>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Delete this key?</AlertDialogTitle>
                            <AlertDialogDescription>
                              This action cannot be undone. Requests using this key will stop
                              working.
                            </AlertDialogDescription>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction onClick={() => delMut.mutate(k.id!)}>
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialog>
                    </TableCell>
                  </TableRow>
                ))}
                {(!keysQ.data || keysQ.data.length === 0) && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-muted-foreground">
                      No API keys yet.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      {/* Plaintext key shown once */}
      <Dialog open={!!justCreated} onOpenChange={(o) => !o && setJustCreated(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Copy your new API key</DialogTitle>
            <DialogDescription>This is only shown once. Store it securely.</DialogDescription>
          </DialogHeader>
          <div className="rounded-md border p-3 font-mono text-sm break-all">
            {(justCreated as any)?.plain ?? "—"}
          </div>
          <div className="flex justify-end gap-2">
            <Button
              variant="outline"
              onClick={() => {
                const val = (justCreated as any)?.plain
                if (val) {
                  navigator.clipboard.writeText(val)
                  toast.success("Copied")
                }
              }}
            >
              Copy
            </Button>
            <Button onClick={() => setJustCreated(null)}>Done</Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
