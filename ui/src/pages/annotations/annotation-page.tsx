import { useMemo, useState } from "react"
import { annotationsApi } from "@/api/annotations.ts"
import { labelsApi } from "@/api/labels.ts"
import type { DtoLabelResponse } from "@/sdk"
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx"

const createAnnotationSchema = z.object({
  key: z.string().trim().min(1, "Key is required").max(120, "Max 120 chars"),
  value: z.string().trim().optional(),
})
type CreateAnnotationInput = z.input<typeof createAnnotationSchema>

const updateAnnotationSchema = createAnnotationSchema.partial()
type UpdateAnnotationValues = z.infer<typeof updateAnnotationSchema>

function AnnotationBadge({ t }: { t: Pick<DtoLabelResponse, "key" | "value"> }) {
  const label = `${t.key}${t.value ? `=${t.value}` : ""}`
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      <Tags className="mr-1 h-3 w-3" />
      {label}
    </Badge>
  )
}

export const AnnotationPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  const annotationQ = useQuery({
    queryKey: ["annotations"],
    queryFn: () => annotationsApi.listAnnotations(),
  })

  // --- Create

  const createForm = useForm<CreateAnnotationInput>({
    resolver: zodResolver(createAnnotationSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateAnnotationInput) => annotationsApi.createAnnotation(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["annotations"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Annotation Created Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while creating Annotation")
    },
  })

  const onCreateSubmit = (values: CreateAnnotationInput) => {
    createMut.mutate(values)
  }

  // --- Update
  const updateForm = useForm<UpdateAnnotationValues>({
    resolver: zodResolver(updateAnnotationSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateAnnotationValues }) =>
      annotationsApi.updateAnnotation(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["annotations"] })
      updateForm.reset()
      setUpdateOpen(false)
      toast.success("Annotation Updated Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while updating Annotation")
    },
  })

  const openEdit = (label: any) => {
    setEditingId(label.id)
    updateForm.reset({
      key: label.key,
      value: label.value,
    })
    setUpdateOpen(true)
  }

  // --- Delete ---

  const deleteMut = useMutation({
    mutationFn: (id: string) => annotationsApi.deleteAnnotation(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["annotations"] })
      setDeleteId(null)
      toast.success("Annotation Deleted Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while deleting Annotation")
    },
  })

  // --- Filter ---
  const filtered = useMemo(() => {
    const data = annotationQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((k: any) => {
          return k.key?.toLowerCase().includes(q) || k.value?.toLowerCase().includes(q)
        })
      : data
  }, [filter, annotationQ.data])

  if (annotationQ.isLoading) return <div className="p-6">Loading annotations…</div>
  if (annotationQ.error)
    return (
      <div className="p-6 text-red-500">
        Error loading annotations.<pre>{JSON.stringify(annotationQ, null, 2)}</pre>
      </div>
    )
  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Annotations</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search annotations"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Annotation
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create Label</DialogTitle>
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
                          <Input placeholder="environment" {...field} />
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
                        <FormLabel>Value</FormLabel>
                        <FormControl>
                          <Input placeholder="dev" {...field} />
                        </FormControl>
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
                <TableHead>Key</TableHead>
                <TableHead>Value</TableHead>
                <TableHead>Annotation</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((t) => (
                <TableRow key={t.id}>
                  <TableCell>{t.key}</TableCell>
                  <TableCell>{t.value}</TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <AnnotationBadge t={t} />
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
                  <TableCell colSpan={4} className="text-muted-foreground py-10 text-center">
                    <CircleSlash2 className="mx-auto mb-2 h-6 w-6 opacity-60" />
                    No labels match your search.
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
            <DialogTitle>Edit Annotation</DialogTitle>
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
            <DialogTitle>Delete annotation</DialogTitle>
          </DialogHeader>
          <p className="text-muted-foreground text-sm">
            This action cannot be undone. Are you sure you want to delete this annotation?
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
