import { useMemo, useState } from "react"
import { taintsApi } from "@/api/taints.ts"
import type { DtoTaintResponse } from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { CircleSlash2, Pencil, Plus, Search, Tags } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { truncateMiddle } from "@/lib/utils.ts"
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select.tsx"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx"

const EFFECTS = ["NoSchedule", "PreferNoSchedule", "NoExecute"] as const

const createTaintSchema = z.object({
  key: z.string().trim().min(1, "Key is required").max(120, "Max 120 chars"),
  value: z.string().trim().optional(),
  effect: z.enum(EFFECTS),
})

type CreateTaintInput = z.input<typeof createTaintSchema>

const updateTaintSchema = createTaintSchema.partial()
type UpdateTaintValues = z.infer<typeof updateTaintSchema>

function TaintBadge({ t }: { t: Pick<DtoTaintResponse, "key" | "value" | "effect"> }) {
  const label = `${t.key}${t.value ? `=${t.value}` : ""}${t.effect ? `:${t.effect}` : ""}`
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      <Tags className="mr-1 h-3 w-3" />
      {label}
    </Badge>
  )
}

export const TaintsPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  const taintsQ = useQuery({
    queryKey: ["taints"],
    queryFn: () => taintsApi.listTaints(),
  })

  // --- Create
  const createForm = useForm<CreateTaintInput>({
    resolver: zodResolver(createTaintSchema),
    defaultValues: {
      key: "",
      value: "",
      effect: undefined,
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateTaintInput) => taintsApi.createTaint(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["taints"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Taint Created Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while creating Taint")
    },
  })

  const onCreateSubmit = (values: CreateTaintInput) => {
    createMut.mutate(values)
  }

  // --- Update
  const updateForm = useForm<UpdateTaintValues>({
    resolver: zodResolver(updateTaintSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateTaintValues }) =>
      taintsApi.updateTaint(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["taints"] })
      updateForm.reset()
      setUpdateOpen(false)
      toast.success("Taint Updated Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while updating Taint")
    },
  })

  const openEdit = (taint: any) => {
    setEditingId(taint.id)
    updateForm.reset({
      key: taint.key,
      value: taint.value,
      effect: taint.effect,
    })
    setUpdateOpen(true)
  }

  // --- Delete ---

  const deleteMut = useMutation({
    mutationFn: (id: string) => taintsApi.deleteTaint(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["taints"] })
      setDeleteId(null)
      toast.success("Taint Deleted Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while deleting Taint")
    },
  })

  const filtered = useMemo(() => {
    const data = taintsQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((k: any) => {
          return (
            k.key?.toLowerCase().includes(q) ||
            k.value?.toLowerCase().includes(q) ||
            k.effect?.toLowerCase().includes(q)
          )
        })
      : data
  }, [filter, taintsQ.data])

  if (taintsQ.isLoading) return <div className="p-6">Loading taints…</div>
  if (taintsQ.error) return <div className="p-6 text-red-500">Error loading taints.</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Taints</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search taints"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" /> Create Taint
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create taint</DialogTitle>
              </DialogHeader>

              <Form {...createForm}>
                <form className="space-y-4" onSubmit={createForm.handleSubmit(onCreateSubmit)}>
                  <FormField
                    control={createForm.control}
                    name="key"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Key</FormLabel>
                        <FormControl>
                          <Input placeholder="dedicated" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="value"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Value (optional)</FormLabel>
                        <FormControl>
                          <Input placeholder="gpu" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={createForm.control}
                    name="effect"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Effect</FormLabel>
                        <Select onValueChange={field.onChange} value={field.value}>
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select effect" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {EFFECTS.map((e) => (
                              <SelectItem key={e} value={e}>
                                {e}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createForm.formState.isSubmitting}>
                      {createForm.formState.isSubmitting ? "Creating…" : "Create"}
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
                <TableHead>Taint</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((t) => (
                <TableRow key={t.id}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <TaintBadge t={t} />
                      <code className="text-muted-foreground text-xs">
                        {truncateMiddle(t.id!, 6)}
                      </code>
                    </div>
                  </TableCell>

                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm" onClick={() => openEdit(t)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => setDeleteId(t.id!)}
                        disabled={deleteMut.isPending && deleteId === t.id}
                      >
                        {deleteMut.isPending && deleteId === t.id ? "Deleting…" : "Delete"}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}

              {filtered.length === 0 && (
                <TableRow>
                  <TableCell colSpan={3} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No taints match your search.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* Update dialog */}
      <Dialog open={updateOpen} onOpenChange={setUpdateOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit taint</DialogTitle>
          </DialogHeader>
          <Form {...updateForm}>
            <form
              className="space-y-4"
              onSubmit={updateForm.handleSubmit((values) => {
                if (!editingId) return
                updateMut.mutate({ id: editingId, values })
              })}
            >
              <FormField
                control={updateForm.control}
                name="key"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Key</FormLabel>
                    <FormControl>
                      <Input placeholder="dedicated" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="value"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Value (optional)</FormLabel>
                    <FormControl>
                      <Input placeholder="gpu" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={updateForm.control}
                name="effect"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Effect</FormLabel>
                    <Select onValueChange={field.onChange} value={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select effect" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {EFFECTS.map((e) => (
                          <SelectItem key={e} value={e}>
                            {e}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="gap-2">
                <Button type="button" variant="outline" onClick={() => setUpdateOpen(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={updateMut.isPending}>
                  {updateMut.isPending ? "Saving…" : "Save changes"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>

      {/* Delete confirm dialog */}
      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>Delete taint</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This action cannot be undone. Are you sure you want to delete this taint?
          </p>
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
