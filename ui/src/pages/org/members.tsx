import { useMemo, useState } from "react";
import { withRefresh } from "@/api/with-refresh.ts";
import { orgStore } from "@/auth/org.ts";
import { makeOrgsApi } from "@/sdkClient.ts";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { z } from "zod";



import { Button } from "@/components/ui/button.tsx";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table.tsx";





const addSchema = z.object({
  user_id: z.uuid("Invalid UUID"),
  role: z.enum(["owner", "admin", "member"]),
})
type AddValues = z.infer<typeof addSchema>

export const OrgMembers = () => {
  const api = makeOrgsApi()
  const qc = useQueryClient()
  const orgId = orgStore.get()
  const [updatingId, setUpdatingId] = useState<string | null>(null)

  const q = useQuery({
    enabled: !!orgId,
    queryKey: ["org:members", orgId],
    queryFn: () => withRefresh(() => api.listMembers({ id: orgId! })),
  })

  const ownersCount = useMemo(
    () => (q.data ?? []).filter((m) => m.role === "owner").length,
    [q.data]
  )

  const form = useForm<AddValues>({
    resolver: zodResolver(addSchema),
    defaultValues: {
      user_id: "",
      role: "member",
    },
  })

  const addMut = useMutation({
    mutationFn: (v: AddValues) => api.addOrUpdateMember({ id: orgId!, addOrUpdateMemberRequest: v }),
    onSuccess: () => {
      toast.success("Member added/updated")
      void qc.invalidateQueries({ queryKey: ["org:members", orgId] })
      form.reset({ user_id: "", role: "member" })
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed"),
  })

  const removeMut = useMutation({
    mutationFn: (userId: string) => api.removeMember({ id: orgId!, userId }),
    onSuccess: () => {
      toast.success("Member removed")
      void qc.invalidateQueries({ queryKey: ["org:members", orgId] })
    },
    onError: (e: any) => toast.error(e?.message ?? "Failed"),
  })

  const roleMut = useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: "owner" | "admin" | "member" }) =>
      api.addOrUpdateMember({ id: orgId!, addOrUpdateMemberRequest: { user_id: userId, role } }),
    onMutate: async ({ userId, role }) => {
      setUpdatingId(userId)
      // cancel queries and snapshot previous
      await qc.cancelQueries({ queryKey: ["org:members", orgId] })
      const prev = qc.getQueryData<any>(["org:members", orgId])
      // optimistic update
      qc.setQueryData(["org:members", orgId], (old: any[] = []) =>
        old.map((m) => (m.user_id === userId ? { ...m, role } : m))
      )
      return { prev }
    },
    onError: (e, _vars, ctx) => {
      if (ctx?.prev) qc.setQueryData(["org:members", orgId], ctx.prev)
      toast.error((e as any)?.message ?? "Failed to update role")
    },
    onSuccess: () => {
      toast.success("Role updated")
    },
    onSettled: () => {
      setUpdatingId(null)
      void qc.invalidateQueries({ queryKey: ["org:members", orgId] })
    },
  })

  const canDowngrade = (m: any) => !(m.role === "owner" && ownersCount <= 1)

  if (!orgId) return <p className="text-muted-foreground">Pick an organization.</p>
  if (q.isLoading) return <p>Loading...</p>
  if (q.error) return <p className="text-destructive">Failed to load members.</p>

  return (
    <Card>
      <CardHeader>
        <CardTitle>Members</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Add/Update */}
        <Form {...form}>
          <form
            className="grid grid-cols-1 items-end gap-3 md:grid-cols-12"
            onSubmit={form.handleSubmit((v) => addMut.mutate(v))}
          >
            <div className="md:col-span-6">
              <FormField
                control={form.control}
                name="user_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>User ID</FormLabel>
                    <FormControl>
                      <Input placeholder="UUID" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="md:col-span-4">
              <FormField
                control={form.control}
                name="role"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <SelectTrigger>
                        <SelectValue placeholder="Select role" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="member">member</SelectItem>
                        <SelectItem value="admin">admin</SelectItem>
                        <SelectItem value="owner">owner</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="md:col-span-2">
              <Button type="submit" className="w-full" disabled={addMut.isPending}>
                Save
              </Button>
            </div>
          </form>
        </Form>

        <div className="overflow-x-auto rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Id</TableHead>
                <TableHead>User</TableHead>
                <TableHead>Role</TableHead>
                <TableHead className="w-28"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {q.data?.map((m) => {
                const isRowPending = updatingId === m.user_id
                return (
                  <TableRow key={m.user_id} className="align-middle">
                    <TableCell className="font-mono text-xs">{m.user_id}</TableCell>
                    <TableCell>{m.email}</TableCell>

                    {/* Inline role select */}
                    <TableCell className="capitalize">
                      <div className="flex items-center gap-2">
                        <Select
                          value={m.role}
                          onValueChange={(next) => {
                            if (m.role === next) return
                            if (m.role === "owner" && next !== "owner" && !canDowngrade(m)) {
                              toast.error("You cannot demote the last owner.")
                              return
                            }
                            roleMut.mutate({
                              userId: m.user_id!,
                              role: next as any,
                            })
                          }}
                          disabled={isRowPending}
                        >
                          <SelectTrigger className="h-8 w-[140px]">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="member">member</SelectItem>
                            <SelectItem value="admin">admin</SelectItem>
                            <SelectItem value="owner">owner</SelectItem>
                          </SelectContent>
                        </Select>

                        {isRowPending && <Loader2 className="h-4 w-4 animate-spin" />}
                      </div>
                    </TableCell>

                    <TableCell className="text-right">
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => removeMut.mutate(m.user_id!)}
                        disabled={m.role === "owner" && ownersCount <= 1}
                        title={
                          m.role === "owner" && ownersCount <= 1
                            ? "Cannot remove the last owner"
                            : ""
                        }
                      >
                        Remove
                      </Button>
                    </TableCell>
                  </TableRow>
                )
              })}
              {q.data?.length === 0 && (
                <TableRow>
                  <TableCell colSpan={3} className="text-muted-foreground p-4">
                    No members.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  )
}
