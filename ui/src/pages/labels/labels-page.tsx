import { useMemo, useState } from "react"
import { labelsApi } from "@/api/labels.ts"
import { taintsApi } from "@/api/taints.ts"
import type { DtoLabelResponse, DtoTaintResponse } from "@/sdk"
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

const createLabelSchema = z.object({
  key: z.string().trim().min(1, "Key is required").max(120, "Max 120 chars"),
  value: z.string().trim().optional(),
})
type CreateLabelInput = z.input<typeof createLabelSchema>

const updateLabelSchema = createLabelSchema.partial()
type UpdateLabelValues = z.infer<typeof updateLabelSchema>

function LabelBadge({ t }: { t: Pick<DtoLabelResponse, "key" | "value"> }) {
  const label = `${t.key}${t.value ? `=${t.value}` : ""}`
  return (
    <Badge variant="secondary" className="font-mono text-xs">
      <Tags className="mr-1 h-3 w-3" />
      {label}
    </Badge>
  )
}

export const LabelsPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [updateOpen, setUpdateOpen] = useState<boolean>(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<string | null>(null)

  const qc = useQueryClient()

  const labelsQ = useQuery({
    queryKey: ["labels"],
    queryFn: () => labelsApi.listLabels(),
  })

  // --- Create
  const createForm = useForm<CreateLabelInput>({
    resolver: zodResolver(createLabelSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  })

  const createMut = useMutation({
    mutationFn: (values: CreateLabelInput) => labelsApi.createLabel(values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["labels"] })
      createForm.reset()
      setCreateOpen(false)
      toast.success("Label Created Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while creating Label")
    },
  })

  const onCreateSubmit = (values: CreateLabelInput) => {
    createMut.mutate(values)
  }

  // --- Update
  const updateForm = useForm<UpdateLabelValues>({
    resolver: zodResolver(updateLabelSchema),
    defaultValues: {},
  })

  const updateMut = useMutation({
    mutationFn: ({ id, values }: { id: string; values: UpdateLabelValues }) =>
      labelsApi.updateLabel(id, values),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["labels"] })
      updateForm.reset()
      setUpdateOpen(false)
      toast.success("Label Updated Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while updating Label")
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
    mutationFn: (id: string) => labelsApi.deleteLabel(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["labels"] })
      setDeleteId(null)
      toast.success("Label Deleted Successfully.")
    },
    onError: (err) => {
      toast.error(err.message ?? "There was an error while deleting Label")
    },
  })

  // --- Filter ---
  const filtered = useMemo(() => {
    const data = labelsQ.data ?? []
    const q = filter.trim().toLowerCase()

    return q
      ? data.filter((k: any) => {
          return k.key?.toLowerCase().includes(q) || k.value?.toLowerCase().includes(q)
        })
      : data
  }, [filter, labelsQ.data])

  if (labelsQ.isLoading) return <div className="p-6">Loading labels…</div>
  if (labelsQ.error) return <div className="p-6 text-red-500">Error loading labels.</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Labels</h1>

        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <Search className="absolute top-2.5 left-2 h-4 w-4 opacity-60" />
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search labels"
              className="w-64 pl-8"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create Label
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
                <TableHead>Label</TableHead>
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
                      <LabelBadge t={t} />
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
      <pre>{JSON.stringify(labelsQ, null, 2)}</pre>
    </div>
  )
}
