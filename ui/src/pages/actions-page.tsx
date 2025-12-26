import { useMemo, useState } from "react"
import { actionsApi } from "@/api/actions.ts"
import type { DtoActionResponse } from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { AlertCircle, CircleSlash2, Loader2, Pencil, Plus, Search, Trash2 } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

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
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx"
import { Textarea } from "@/components/ui/textarea.tsx"

const createActionSchema = z.object({
  label: z.string().trim().min(1, "Label is required").max(255, "Max 255 chars"),
  description: z.string().trim().min(1, "Description is required"),
  make_target: z
    .string()
    .trim()
    .min(1, "Make target is required")
    .max(255, "Max 255 chars")
    // keep client-side fairly strict to avoid surprises; server should also validate
    .regex(/^[a-zA-Z0-9][a-zA-Z0-9._-]{0,63}$/, "Invalid make target (allowed: a-z A-Z 0-9 . _ -)"),
})
type CreateActionInput = z.input<typeof createActionSchema>

const updateActionSchema = createActionSchema.partial()
type UpdateActionInput = z.input<typeof updateActionSchema>

function TargetBadge({ target }: { target?: string | null }) {
  if (!target) {
    return (
      <Badge variant="outline" className="text-xs">
        —
      </Badge>
    )
  }
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      {target}
    </Badge>
  )
}

export const ActionsPage = () => {
  const qc = useQueryClient()

  const [filter, setFilter] = useState("")
  const [createOpen, setCreateOpen] = useState(false)
  const [updateOpen, setUpdateOpen] = useState(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editing, setEditing] = useState<DtoActionResponse | null>(null)

  const actionsQ = useQuery({
    queryKey: ["admin-actions"],
    queryFn: () => actionsApi.listActions(),
  })

  const filtered = useMemo(() => {
    const data: DtoActionResponse[] = actionsQ.data ?? []
    const q = filter.trim().toLowerCase()
    if (!q) return data

    return data.filter((a) => {
      return (
        (a.label ?? "").toLowerCase().includes(q) ||
        (a.description ?? "").toLowerCase().includes(q) ||
        (a.make_target ?? "").toLowerCase().includes(q)
      )
    })
  }, [filter, actionsQ.data])

  const createForm = useForm<CreateActionInput>({
    resolver: zodResolver(createActionSchema),
    defaultValues: {
      label: "",
      description: "",
      make_target: "",
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateActionInput) => actionsApi.createAction(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["admin-actions"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Action created.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to create action.")
    },
  })

  const updateForm = useForm<UpdateActionInput>({
    resolver: zodResolver(updateActionSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateActionInput }) =>
      actionsApi.updateAction(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["admin-actions"] })
      updateForm.reset()
      setUpdateOpen(false)
      setEditing(null)
      toast.success("Action updated.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to update action.")
    },
  })

  const openEdit = (a: DtoActionResponse) => {
    if (!a.id) return
    setEditing(a)
    updateForm.reset({
      label: a.label ?? "",
      description: a.description ?? "",
      make_target: a.make_target ?? "",
    })
    setUpdateOpen(true)
  }

  const deleteMut = useMutation({
    mutationFn: (id: string) => actionsApi.deleteAction(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["admin-actions"] })
      setDeleteId(null)
      toast.success("Action deleted.")
    },
    onError: (err: any) => {
      toast.error(err?.message ?? "Failed to delete action.")
    },
  })

  if (actionsQ.isLoading) return <div className="p-6">Loading actions…</div>
  if (actionsQ.error) return <div className="p-6 text-red-500">Error loading actions.</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="text-2xl font-bold">Admin Actions</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search actions"
              className="w-72 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Action
              </Button>
            </DialogTrigger>

            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create Action</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form
                  className="space-y-4"
                  onSubmit={createForm.handleSubmit((v) => createMut.mutate(v))}
                >
                  <FormField
                    control={createForm.control}
                    name="label"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Label</FormLabel>
                        <FormControl>
                          <Input placeholder="Setup" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="make_target"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Make Target</FormLabel>
                        <FormControl>
                          <Input placeholder="setup" className="font-mono" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Description</FormLabel>
                        <FormControl>
                          <Textarea
                            rows={4}
                            placeholder="Runs prepare, ping-servers, then make setup on the bastion."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createMut.isPending}>
                      {createMut.isPending ? "Creating…" : "Create"}
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
                <TableHead>Label</TableHead>
                <TableHead>Make Target</TableHead>
                <TableHead>Description</TableHead>
                <TableHead className="w-[260px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((a) => (
                <TableRow key={a.id}>
                  <TableCell className="font-medium">{a.label}</TableCell>
                  <TableCell>
                    <TargetBadge target={a.make_target} />
                  </TableCell>
                  <TableCell className="text-muted-foreground max-w-[680px] truncate">
                    {a.description}
                  </TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => openEdit(a)}>
                        <Pencil className="mr-2 h-4 w-4" />
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => a.id && setDeleteId(a.id)}
                        disabled={deleteMut.isPending && deleteId === a.id}
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        {deleteMut.isPending && deleteId === a.id ? "Deleting…" : "Delete"}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={4} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No actions match your search.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Update dialog */}
      <Dialog
        open={updateOpen}
        onOpenChange={(open) => {
          setUpdateOpen(open)
          if (!open) setEditing(null)
        }}
      >
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Action</DialogTitle>
          </DialogHeader>

          {editing ? (
            <Form {...updateForm}>
              <form
                className="space-y-4"
                onSubmit={updateForm.handleSubmit((values) => {
                  if (!editing.id) return
                  updateMut.mutate({ id: editing.id, values })
                })}
              >
                <FormField
                  control={updateForm.control}
                  name="label"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Label</FormLabel>
                      <FormControl>
                        <Input {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={updateForm.control}
                  name="make_target"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Make Target</FormLabel>
                      <FormControl>
                        <Input className="font-mono" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={updateForm.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea rows={4} {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <DialogFooter className="gap-2">
                  <Button type="button" variant="outline" onClick={() => setUpdateOpen(false)}>
                    Cancel
                  </Button>
                  <Button type="submit" disabled={updateMut.isPending}>
                    {updateMut.isPending ? (
                      <span className="inline-flex items-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        Saving…
                      </span>
                    ) : (
                      "Save changes"
                    )}
                  </Button>
                </DialogFooter>
              </form>
            </Form>
          ) : (
            <div className="text-muted-foreground text-sm">No action selected.</div>
          )}
        </DialogContent>
      </Dialog>

      {/* Delete confirm dialog */}
      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Delete action</DialogTitle>
          </DialogHeader>

          <div className="flex items-start gap-3">
            <AlertCircle className="mt-0.5 h-5 w-5 text-red-500" />
            <p className="text-muted-foreground text-sm">
              This action cannot be undone. Are you sure you want to delete it?
            </p>
          </div>

          <DialogFooter className="gap-2">
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => deleteId && deleteMut.mutate(deleteId)}
              disabled={deleteMut.isPending}
            >
              {deleteMut.isPending ? "Deleting…" : "Delete"}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
