import { useState } from "react"
import { withRefresh } from "@/api/with-refresh.ts"
import { orgStore } from "@/auth/org.ts"
import { makeOrgsApi } from "@/sdkClient.ts"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { Button } from "@/components/ui/button.tsx"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, } from "@/components/ui/dialog.tsx"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, } from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow, } from "@/components/ui/table.tsx"

const createSchema = z.object({
  name: z.string(),
  expires_in_hours: z.number().min(1).max(43800),
})
type CreateValues = z.infer<typeof createSchema>

export const OrgApiKeys = () => {
  const api = makeOrgsApi()
  const qc = useQueryClient()
  const orgId = orgStore.get()

  const q = useQuery({
    enabled: !!orgId,
    queryKey: ["org:keys", orgId],
    queryFn: () => withRefresh(() => api.listOrgKeys({ id: orgId! })),
  })

  const form = useForm<CreateValues>({
    resolver: zodResolver(createSchema),
    defaultValues: {
      name: "",
      expires_in_hours: 720,
    },
  })

  const [showSecret, setShowSecret] = useState<{
    key?: string
    secret?: string
  } | null>(null)

  const createMut = useMutation({
    mutationFn: (v: CreateValues) => api.createOrgKey({ id: orgId!, handlersOrgKeyCreateReq: v }),
    onSuccess: (resp) => {
      void qc.invalidateQueries({ queryKey: ["org:keys", orgId] })
      setShowSecret({ key: resp.org_key, secret: resp.org_secret })
      toast.success("Key created")
      form.reset({ name: "", expires_in_hours: undefined })
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed to create key"),
  })

  const deleteMut = useMutation({
    mutationFn: (keyId: string) => api.deleteOrgKey({ id: orgId!, keyId }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["org:keys", orgId] })
      toast.success("Key deleted")
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed to delete key"),
  })

  if (!orgId) return <p className="text-muted-foreground">Pick an organization.</p>
  if (q.isLoading) return <p>Loading...</p>
  if (q.error) return <p className="text-destructive">Failed to load keys.</p>

  return (
    <Card>
      <CardHeader>
        <CardTitle>Org API Keys</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit((v) => createMut.mutate(v))}
            className="grid grid-cols-1 items-end gap-3 md:grid-cols-12"
          >
            <div className="md:col-span-6">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="automation-bot" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="md:col-span-4">
              <FormField
                control={form.control}
                name="expires_in_hours"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Expires In (hours)</FormLabel>
                    <FormControl>
                      <Input placeholder="e.g. 720" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="md:col-span-2">
              <Button type="submit" className="w-full" disabled={createMut.isPending}>
                Create
              </Button>
            </div>
          </form>
        </Form>

        <div className="overflow-x-auto rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Scope</TableHead>
                <TableHead>Created</TableHead>
                <TableHead>Expires</TableHead>
                <TableHead className="w-28" />
              </TableRow>
            </TableHeader>
            <TableBody>
              {q.data?.map((k) => (
                <TableRow key={k.id}>
                  <TableCell>{k.name ?? "-"}</TableCell>
                  <TableCell>{k.scope}</TableCell>
                  <TableCell>{new Date(k.created_at!).toLocaleString()}</TableCell>
                  <TableCell>
                    {k.expires_at ? new Date(k.expires_at).toLocaleString() : "-"}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="destructive" size="sm" onClick={() => deleteMut.mutate(k.id!)}>
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {q.data?.length === 0 && (
                <TableRow>
                  <TableCell colSpan={5} className="text-muted-foreground p-4">
                    No keys.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        {/* Show once dialog with key/secret */}
        <Dialog open={!!showSecret} onOpenChange={(o) => !o && setShowSecret(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Copy your credentials</DialogTitle>
            </DialogHeader>
            <div className="space-y-2">
              <div>
                <div className="text-muted-foreground mb-1 text-xs">Org Key</div>
                <Input
                  readOnly
                  value={showSecret?.key ?? ""}
                  onFocus={(e) => e.currentTarget.select()}
                />
              </div>
              <div>
                <div className="text-muted-foreground mb-1 text-xs">Org Secret</div>
                <Input
                  readOnly
                  value={showSecret?.secret ?? ""}
                  onFocus={(e) => e.currentTarget.select()}
                />
              </div>
              <p className="text-muted-foreground text-xs">
                This secret is shown once. Store it securely.
              </p>
            </div>
            <DialogFooter>
              <Button onClick={() => setShowSecret(null)}>Done</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </CardContent>
    </Card>
  )
}
