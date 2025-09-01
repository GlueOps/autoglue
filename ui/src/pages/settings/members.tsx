// src/pages/settings/members.tsx
import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { TrashIcon, UserPlus2 } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { api, ApiError } from "@/lib/api.ts"
import { EVT_ACTIVE_ORG_CHANGED, getActiveOrgId } from "@/lib/orgs-sync.ts"
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
import { Button } from "@/components/ui/button.tsx"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter as DialogFooterUI,
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
import { Separator } from "@/components/ui/separator.tsx"
import { Skeleton } from "@/components/ui/skeleton.tsx"

type Me = { id: string; email?: string; name?: string }

// Backend shape can vary; normalize to a safe shape for UI.
type MemberDTO = any
type Member = {
  userId: string
  email?: string
  name?: string
  role: string
  joinedAt?: string
}

function normalizeMember(m: MemberDTO): Member {
  const userId = m?.user_id ?? m?.UserID ?? m?.user?.id ?? m?.User?.ID ?? ""
  const email = m?.email ?? m?.Email ?? m?.user?.email ?? m?.User?.Email
  const name = m?.name ?? m?.Name ?? m?.user?.name ?? m?.User?.Name
  const role = m?.role ?? m?.Role ?? "member"
  const joinedAt = m?.created_at ?? m?.CreatedAt

  return { userId: String(userId), email, name, role: String(role), joinedAt }
}

const InviteSchema = z.object({
  email: z.email("Enter a valid email"),
  role: z.enum(["member", "admin"]),
})
type InviteValues = z.infer<typeof InviteSchema>

export const MemberManagement = () => {
  const [loading, setLoading] = useState(true)
  const [members, setMembers] = useState<Member[]>([])
  const [me, setMe] = useState<Me | null>(null)
  const [inviteOpen, setInviteOpen] = useState(false)
  const [inviting, setInviting] = useState(false)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  const activeOrgIdInitial = useMemo(() => getActiveOrgId(), [])

  const form = useForm<InviteValues>({
    resolver: zodResolver(InviteSchema),
    defaultValues: { email: "", role: "member" },
    mode: "onChange",
  })

  async function fetchMe() {
    try {
      const data = await api.get<Me>("/api/v1/auth/me")
      setMe(data)
    } catch {
      // non-blocking
    }
  }

  async function fetchMembers(orgId: string | null) {
    if (!orgId) {
      setMembers([])
      setLoading(false)
      return
    }
    setLoading(true)
    try {
      const data = await api.get<MemberDTO[]>("/api/v1/orgs/members")
      setMembers((data ?? []).map(normalizeMember))
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to load members"
      toast.error(msg)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void fetchMe()
    void fetchMembers(activeOrgIdInitial)
  }, [activeOrgIdInitial])

  // Refetch when active org changes (same tab or across tabs)
  useEffect(() => {
    const onActiveOrg = () => void fetchMembers(getActiveOrgId())
    const onStorage = (e: StorageEvent) => {
      if (e.key === "active_org_id") onActiveOrg()
    }
    window.addEventListener(EVT_ACTIVE_ORG_CHANGED, onActiveOrg as EventListener)
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener(EVT_ACTIVE_ORG_CHANGED, onActiveOrg as EventListener)
      window.removeEventListener("storage", onStorage)
    }
  }, [])

  async function onInvite(values: InviteValues) {
    const orgId = getActiveOrgId()
    if (!orgId) {
      toast.error("Select an organization first")
      return
    }
    try {
      setInviting(true)
      await api.post("/api/v1/orgs/invite", values)
      toast.success(`Invited ${values.email}`)
      setInviteOpen(false)
      form.reset({ email: "", role: "member" })
      // If you later expose pending invites, update that list; for now just refresh members.
      void fetchMembers(orgId)
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to invite member"
      toast.error(msg)
    } finally {
      setInviting(false)
    }
  }

  async function onRemove(userId: string) {
    const orgId = getActiveOrgId()
    if (!orgId) {
      toast.error("Select an organization first")
      return
    }
    try {
      setDeletingId(userId)
      await api.delete<void>(`/api/v1/orgs/members/${userId}`, {
        headers: { "X-Org-ID": orgId },
      })
      setMembers((prev) => prev.filter((m) => m.userId !== userId))
      toast.success("Member removed")
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to remove member"
      toast.error(msg)
    } finally {
      setDeletingId(null)
    }
  }

  const canManage = true // Server enforces admin; UI stays permissive.

  if (loading) {
    return (
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">Members</h1>
          <Button disabled>
            <UserPlus2 className="mr-2 h-4 w-4" />
            Invite
          </Button>
        </div>
        <Separator />
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-5 w-40" />
              </CardHeader>
              <CardContent className="space-y-2">
                <Skeleton className="h-4 w-56" />
                <Skeleton className="h-4 w-40" />
              </CardContent>
              <CardFooter>
                <Skeleton className="h-9 w-24" />
              </CardFooter>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  if (!getActiveOrgId()) {
    return (
      <div className="space-y-4 p-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold">Members</h1>
        </div>
        <Separator />
        <p className="text-muted-foreground text-sm">
          No organization selected. Choose an organization to manage its members.
        </p>
      </div>
    )
  }

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="text-2xl font-bold">Members</h1>

        <Dialog open={inviteOpen} onOpenChange={setInviteOpen}>
          <DialogTrigger asChild>
            <Button>
              <UserPlus2 className="mr-2 h-4 w-4" />
              Invite
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[520px]">
            <DialogHeader>
              <DialogTitle>Invite member</DialogTitle>
              <DialogDescription>Send an invite to join this organization.</DialogDescription>
            </DialogHeader>

            <Form {...form}>
              <form onSubmit={form.handleSubmit(onInvite)} className="grid gap-4 py-2">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email</FormLabel>
                      <FormControl>
                        <Input type="email" placeholder="jane@example.com" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="role"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Role</FormLabel>
                      <Select onValueChange={field.onChange} defaultValue={field.value}>
                        <FormControl>
                          <SelectTrigger className="w-[200px]">
                            <SelectValue placeholder="Select role" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="member">Member</SelectItem>
                          <SelectItem value="admin">Admin</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooterUI className="mt-2 flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                  <Button type="button" variant="outline" onClick={() => setInviteOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={!form.formState.isValid || inviting}>
                    {inviting ? "Sending…" : "Send invite"}
                  </Button>
                </DialogFooterUI>
              </form>
            </Form>
          </DialogContent>
        </Dialog>
      </div>

      <Separator />

      {members.length === 0 ? (
        <div className="text-muted-foreground text-sm">No members yet.</div>
      ) : (
        <div className="grid grid-cols-1 gap-4 pr-2 sm:grid-cols-2 lg:grid-cols-3">
          {members.map((m) => {
            const isSelf = me?.id && m.userId === me.id
            return (
              <Card key={m.userId} className="flex flex-col">
                <CardHeader>
                  <CardTitle className="text-base">{m.name || m.email || m.userId}</CardTitle>
                </CardHeader>
                <CardContent className="text-muted-foreground space-y-1 text-sm">
                  {m.email && <div>Email: {m.email}</div>}
                  <div>Role: {m.role}</div>
                  {m.joinedAt && <div>Joined: {new Date(m.joinedAt).toLocaleString()}</div>}
                </CardContent>
                <CardFooter className="mt-auto w-full flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                  <div />
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button
                        variant="destructive"
                        disabled={!canManage || isSelf || deletingId === m.userId}
                        className="ml-auto"
                      >
                        <TrashIcon className="mr-2 h-5 w-5" />
                        {deletingId === m.userId ? "Removing…" : "Remove"}
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Remove member?</AlertDialogTitle>
                        <AlertDialogDescription>
                          This will remove <b>{m.name || m.email || m.userId}</b> from the
                          organization.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter className="sm:justify-between">
                        <AlertDialogCancel disabled={deletingId === m.userId}>
                          Cancel
                        </AlertDialogCancel>
                        <AlertDialogAction asChild disabled={deletingId === m.userId}>
                          <Button variant="destructive" onClick={() => onRemove(m.userId)}>
                            Confirm remove
                          </Button>
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </CardFooter>
              </Card>
            )
          })}
        </div>
      )}
    </div>
  )
}
