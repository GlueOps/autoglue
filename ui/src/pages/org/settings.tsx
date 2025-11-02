import { useEffect } from "react"
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
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"

const schema = z.object({
  name: z.string().min(1, "Required"),
  domain: z.string().optional(),
})

type Values = z.infer<typeof schema>

export const OrgSettings = () => {
  const api = makeOrgsApi()
  const qc = useQueryClient()
  const orgId = orgStore.get()

  const q = useQuery({
    enabled: !!orgId,
    queryKey: ["org", orgId],
    queryFn: () => withRefresh(() => api.getOrg({ id: orgId! })),
  })

  const form = useForm<Values>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: "",
      domain: "",
    },
  })

  useEffect(() => {
    if (q.data) {
      form.reset({
        name: q.data.name ?? "",
        domain: q.data.domain ?? "",
      })
    }
  }, [q.data])

  const updateMut = useMutation({
    mutationFn: (v: Partial<Values>) => api.updateOrg({ id: orgId!, body: v }),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["org", orgId] })
      toast.success("Organization updated")
    },
    onError: (e: any) => toast.error(e?.message ?? "Update failed"),
  })

  const deleteMut = useMutation({
    mutationFn: () => api.deleteOrg({ id: orgId! }),
    onSuccess: () => {
      toast.success("Organization deleted")
      orgStore.set("")
      void qc.invalidateQueries({ queryKey: ["orgs:mine"] })
    },
    onError: (e: any) => toast.error(e?.message ?? "Delete failed"),
  })

  if (!orgId) {
    return <p className="text-muted-foreground">Pick an organization.</p>
  }
  if (q.isLoading) return <p>Loading...</p>
  if (q.error) return <p className="text-destructive">Failed to load.</p>

  const onSubmit = (v: Values) => {
    const delta: Partial<Values> = {}
    if (v.name !== q.data?.name) delta.name = v.name
    const normDomain = v.domain?.trim() || undefined
    if ((normDomain ?? null) !== (q.data?.domain ?? null)) delta.domain = normDomain
    if (Object.keys(delta).length === 0) return
    updateMut.mutate(delta)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Organization Settings</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <Form {...form}>
          <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
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

            <div className="flex gap-2">
              <Button type="submit" disabled={updateMut.isPending}>
                Save
              </Button>
              <Button
                type="button"
                variant="destructive"
                onClick={() => deleteMut.mutate()}
                disabled={deleteMut.isPending}
              >
                Delete Org
              </Button>
            </div>
          </form>
        </Form>
      </CardContent>
    </Card>
  )
}
