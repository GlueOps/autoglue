import { useEffect, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { PencilIcon, Plus, TrashIcon } from "lucide-react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { api } from "@/lib/api.ts"
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

type Label = {
  id: string
  key: string
  value: string
}

const CreateLabelSchema = z.object({
  key: z.string().min(2),
  value: z.string().min(2),
})
type CreateLabelValues = z.infer<typeof CreateLabelSchema>

export const LabelsPage = () => {
  const [labels, setLabels] = useState<Label[]>([])
  const [loading, setLoading] = useState<boolean>(false)
  const [err, setErr] = useState<string | null>(null)

  const [createOpen, setCreateOpen] = useState(false)

  async function loadAll() {
    setLoading(true)
    setErr(null)
    try {
      const labelData = await api.get<Label[]>("/api/v1/labels")
      console.log(JSON.stringify(labelData))
      setLabels(labelData)
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void loadAll()
  }, [])

  const createForm = useForm<CreateLabelValues>({
    resolver: zodResolver(CreateLabelSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  })

  const submitCreate = async (values: CreateLabelValues) => {
    const payload: Record<string, any> = {
      key: values.key,
      value: values.value,
    }
    await api.post<Label>("/api/v1/labels", payload)
    setCreateOpen(false)
    createForm.reset()
    await loadAll()
  }

  if (loading) return <div className="p-6">Loading servers…</div>
  if (err) return <div className="p-6 text-red-500">{err}</div>

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Labels</h1>

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
              <form onSubmit={createForm.handleSubmit(submitCreate)} className="space-y-4">
                <FormField
                  control={createForm.control}
                  name="key"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Key</FormLabel>
                      <FormControl>
                        <Input placeholder="app.kubernetes.io/managed-by" {...field} />
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
                        <Input placeholder="GlueOps" {...field} />
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

      <div className="bg-background overflow-hidden rounded-2xl border shadow-sm">
        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Key</TableHead>
                <TableHead>Values</TableHead>
                <TableHead className="w-[180px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {labels.map((l) => (
                <TableRow key={l.id}>
                  <TableCell>{l.key}</TableCell>
                  <TableCell>{l.value}</TableCell>
                  <TableCell>
                    <div className="flex justify-end gap-2">
                      <Button variant="outline" size="sm">
                        <PencilIcon className="mr-2 h-4 w-4" />
                        Edit
                      </Button>
                      <Button variant="destructive" size="sm">
                        <TrashIcon className="mr-2 h-4 w-4" />
                        Delete
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  )
}
