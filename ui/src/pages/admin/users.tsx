import { useEffect, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { PencilIcon, PlusIcon, TrashIcon } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { api, ApiError } from "@/lib/api"
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
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter as DialogFooterUI,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
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
import { Separator } from "@/components/ui/separator"

type User = {
  id: string
  name: string
  email: string
  role: "admin" | "user" | string
  email_verified: boolean
  created_at: string
  updated_at?: string
}

type ListRes = { users: User[]; page: number; page_size: number; total: number }

const CreateSchema = z.object({
  name: z.string().min(1, "Name required"),
  email: z.email("Enter a valid email"),
  role: z.enum(["user", "admin"]),
  password: z.string().min(8, "Min 8 characters"),
})
type CreateValues = z.infer<typeof CreateSchema>

const EditSchema = z.object({
  name: z.string().min(1, "Name required"),
  email: z.email("Enter a valid email"),
  role: z.enum(["user", "admin"]),
  password: z.string().min(8, "Min 8 characters").optional().or(z.literal("")),
})
type EditValues = z.infer<typeof EditSchema>

export function AdminUsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)

  const [createOpen, setCreateOpen] = useState(false)
  const [editOpen, setEditOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)

  const createForm = useForm<CreateValues>({
    resolver: zodResolver(CreateSchema),
    mode: "onChange",
    defaultValues: { name: "", email: "", role: "user", password: "" },
  })

  const editForm = useForm<EditValues>({
    resolver: zodResolver(EditSchema),
    mode: "onChange",
    defaultValues: { name: "", email: "", role: "user", password: "" },
  })

  async function fetchUsers() {
    setLoading(true)
    try {
      const res = await api.get<ListRes>("/api/v1/admin/users?page=1&page_size=100")
      setUsers(res.users ?? [])
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Failed to load users")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void fetchUsers()
  }, [])

  async function onCreate(values: CreateValues) {
    try {
      const newUser = await api.post<User>("/api/v1/admin/users", values)
      setUsers((prev) => [newUser, ...prev])
      setCreateOpen(false)
      createForm.reset({ name: "", email: "", role: "user", password: "" })
      toast.success(`Created ${newUser.email}`)
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Failed to create user")
    }
  }

  function openEdit(u: User) {
    setEditing(u)
    editForm.reset({
      name: u.name || "",
      email: u.email,
      role: (u.role as "user" | "admin") ?? "user",
      password: "",
    })
    setEditOpen(true)
  }

  async function onEdit(values: EditValues) {
    if (!editing) return
    const payload: Record<string, unknown> = {
      name: values.name,
      email: values.email,
      role: values.role,
    }
    if (values.password && values.password.length >= 8) {
      payload.password = values.password
    }
    try {
      const updated = await api.patch<User>(`/api/v1/admin/users/${editing.id}`, payload)
      setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)))
      setEditOpen(false)
      setEditing(null)
      toast.success(`Updated ${updated.email}`)
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Failed to update user")
    }
  }

  async function onDelete(id: string) {
    try {
      setDeletingId(id)
      await api.delete<void>(`/api/v1/admin/users/${id}`)
      setUsers((prev) => prev.filter((u) => u.id !== id))
      toast.success("User deleted")
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Failed to delete user")
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <div className="space-y-4 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Users</h1>
        <Button onClick={() => setCreateOpen(true)}>
          <PlusIcon className="mr-2 h-4 w-4" />
          New user
        </Button>
      </div>
      <Separator />

      {loading ? (
        <div className="text-muted-foreground text-sm">Loading…</div>
      ) : users.length === 0 ? (
        <div className="text-muted-foreground text-sm">No users yet.</div>
      ) : (
        <div className="grid grid-cols-1 gap-4 pr-2 md:grid-cols-2 lg:grid-cols-3">
          {users.map((u) => (
            <Card key={u.id} className="flex flex-col">
              <CardHeader>
                <CardTitle className="text-base">{u.name || u.email}</CardTitle>
              </CardHeader>
              <CardContent className="text-muted-foreground space-y-1 text-sm">
                <div>Email: {u.email}</div>
                <div>Role: {u.role}</div>
                <div>Verified: {u.email_verified ? "Yes" : "No"}</div>
                <div>Joined: {new Date(u.created_at).toLocaleString()}</div>
              </CardContent>
              <CardFooter className="mt-auto w-full flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Button variant="outline" onClick={() => openEdit(u)}>
                  <PencilIcon className="mr-2 h-4 w-4" /> Edit
                </Button>

                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button variant="destructive" disabled={deletingId === u.id}>
                      <TrashIcon className="mr-2 h-4 w-4" />
                      {deletingId === u.id ? "Deleting…" : "Delete"}
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Delete user?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This will permanently delete <b>{u.email}</b>.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter className="sm:justify-between">
                      <AlertDialogCancel disabled={deletingId === u.id}>Cancel</AlertDialogCancel>
                      <AlertDialogAction asChild disabled={deletingId === u.id}>
                        <Button variant="destructive" onClick={() => onDelete(u.id)}>
                          Confirm delete
                        </Button>
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}

      {/* Create dialog */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-[520px]">
          <DialogHeader>
            <DialogTitle>Create user</DialogTitle>
            <DialogDescription>Add a new user account.</DialogDescription>
          </DialogHeader>
          <Form {...createForm}>
            <form onSubmit={createForm.handleSubmit(onCreate)} className="grid gap-4 py-2">
              <FormField
                name="name"
                control={createForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input {...field} placeholder="Jane Doe" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                name="email"
                control={createForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" {...field} placeholder="jane@example.com" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                name="role"
                control={createForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className="w-[200px]">
                          <SelectValue placeholder="Select role" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="user">User</SelectItem>
                        <SelectItem value="admin">Admin</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                name="password"
                control={createForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="password" {...field} placeholder="••••••••" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooterUI className="mt-2 flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={!createForm.formState.isValid || createForm.formState.isSubmitting}
                >
                  {createForm.formState.isSubmitting ? "Creating…" : "Create"}
                </Button>
              </DialogFooterUI>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Edit dialog */}
      <Dialog open={editOpen} onOpenChange={setEditOpen}>
        <DialogContent className="sm:max-w-[520px]">
          <DialogHeader>
            <DialogTitle>Edit user</DialogTitle>
            <DialogDescription>
              Update user details. Leave password blank to keep it unchanged.
            </DialogDescription>
          </DialogHeader>
          <Form {...editForm}>
            <form onSubmit={editForm.handleSubmit(onEdit)} className="grid gap-4 py-2">
              <FormField
                name="name"
                control={editForm.control}
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
                name="email"
                control={editForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input type="email" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                name="role"
                control={editForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Role</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className="w-[200px]">
                          <SelectValue placeholder="Select role" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="user">User</SelectItem>
                        <SelectItem value="admin">Admin</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                name="password"
                control={editForm.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New password (optional)</FormLabel>
                    <FormControl>
                      <Input type="password" {...field} placeholder="Leave blank to keep" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooterUI className="mt-2 flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Button type="button" variant="outline" onClick={() => setEditOpen(false)}>
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={!editForm.formState.isValid || editForm.formState.isSubmitting}
                >
                  {editForm.formState.isSubmitting ? "Saving…" : "Save changes"}
                </Button>
              </DialogFooterUI>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
