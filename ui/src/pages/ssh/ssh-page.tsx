import { useMemo, useState } from "react"
import { sshApi } from "@/api/ssh.ts"
import type { DtoCreateSSHRequest, DtoSshRevealResponse } from "@/sdk"
import { zodResolver } from "@hookform/resolvers/zod"
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { Download, Eye, Loader2, Plus, Trash2 } from "lucide-react"
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
import { Textarea } from "@/components/ui/textarea.tsx"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx"

const createKeySchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(100, "Max 100 characters"),
  comment: z.string().trim().min(1, "Comment is required").max(100, "Max 100 characters"),
  bits: z.enum(["2048", "3072", "4096"]).optional(),
  type: z.enum(["rsa", "ed25519"]).optional(),
})

type CreateKeyInput = z.input<typeof createKeySchema>

function saveBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement("a")
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

function copy(text: string, label = "Copied") {
  navigator.clipboard
    .writeText(text)
    .then(() => toast.success(label))
    .catch(() => toast.error("Copy failed"))
}

function getKeyType(publicKey: string) {
  return publicKey?.split(/\s+/)?.[0] ?? "ssh-key"
}

export const SshPage = () => {
  const [filter, setFilter] = useState<string>("")
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [revealFor, setRevealFor] = useState<DtoSshRevealResponse | null>(null)
  const [deleteId, setDeleteId] = useState<string | null>(null)

  const qc = useQueryClient()

  const sshQ = useQuery({
    queryKey: ["ssh"],
    queryFn: () => sshApi.listSshKeys(),
  })

  const form = useForm<CreateKeyInput>({
    resolver: zodResolver(createKeySchema),
    defaultValues: {
      name: "",
      comment: "",
      type: "rsa",
      bits: "4096",
    },
  })

  const createMutation = useMutation({
    mutationFn: async (values: CreateKeyInput) => {
      const payload: DtoCreateSSHRequest = {
        name: values.name,
        comment: values.comment,
        // Only send bits for RSA
        bits: values.type === "rsa" && values.bits ? Number(values.bits) : undefined,
        // Only send type if present
        type: values.type,
      }
      return await sshApi.createSshKey(payload)
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: ["ssh"] })
      setCreateOpen(false)
      form.reset({ name: "", comment: "", type: "rsa", bits: "4096" })
      toast.success("SSH Key created")
    },
    onError: (e: any) => toast.error(e?.message ?? "SSH Key creation failed"),
  })

  const revealMutation = useMutation({
    mutationFn: (id: string) => sshApi.revealSshKeyById(id),
    onSuccess: (data) => setRevealFor(data),
    onError: (e: any) => toast.error(e?.message ?? "Failed to reveal key"),
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => sshApi.deleteSshKey(id),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["ssh"] })
      setDeleteId(null)
      toast.success("SSH Key deleted")
    },
    onError: (e: any) => toast.error(e?.message ?? "Delete failed"),
  })

  const filtered = useMemo(() => {
    const q = filter.trim().toLowerCase()
    if (!q) return sshQ.data ?? []
    return (sshQ.data ?? []).filter((k) => {
      return (
        k.name?.toLowerCase().includes(q) ||
        k.fingerprint?.toLowerCase().includes(q) ||
        k.public_key?.toLowerCase().includes(q)
      )
    })
  }, [filter, sshQ.data])

  if (sshQ.isLoading) return <div className="p-6">Loading SSH Keys…</div>
  if (sshQ.error) return <div className="p-6 text-red-500">Error Loading SSH Keys</div>

  return (
    <TooltipProvider>
      <div className="space-y-4">
        <div className="flex items-center justify-between gap-3">
          <h1 className="text-2xl font-bold">SSH Keys</h1>

          <div className="w-full max-w-sm">
            <Input
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Search by name, fingerprint or key"
            />
          </div>

          <Dialog open={createOpen} onOpenChange={setCreateOpen}>
            <DialogTrigger asChild>
              <Button onClick={() => setCreateOpen(true)}>
                <Plus className="mr-2 h-4 w-4" />
                Create New Keypair
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Create SSH Keypair</DialogTitle>
              </DialogHeader>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit((values) => createMutation.mutate(values))}
                  className="space-y-4"
                >
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Name</FormLabel>
                        <FormControl>
                          <Input placeholder="e.g., CI deploy key" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="comment"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Comment</FormLabel>
                        <FormControl>
                          <Input placeholder="e.g., deploy@autoglue" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="type"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Type</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value}
                            onValueChange={(v) => {
                              field.onChange(v)
                              if (v === "ed25519") {
                                // bits not applicable
                                form.setValue("bits", undefined)
                              } else {
                                form.setValue("bits", "4096")
                              }
                            }}
                          >
                            <SelectTrigger className="w-[180px]">
                              <SelectValue placeholder="Select a ssh key type" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="rsa">RSA</SelectItem>
                              <SelectItem value="ed25519">ED25519</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="bits"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Key size</FormLabel>
                        <FormControl>
                          <Select
                            value={field.value}
                            disabled={form.watch("type") === "ed25519"}
                            onValueChange={field.onChange}
                          >
                            <SelectTrigger className="w-[180px]">
                              <SelectValue placeholder="RSA only" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="2048">2048</SelectItem>
                              <SelectItem value="3072">3072</SelectItem>
                              <SelectItem value="4096">4096</SelectItem>
                            </SelectContent>
                          </Select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => setCreateOpen(false)}
                      disabled={createMutation.isPending}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" disabled={createMutation.isPending}>
                      {createMutation.isPending ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                          Creating…
                        </>
                      ) : (
                        "Create"
                      )}
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
                  <TableHead>Name</TableHead>
                  <TableHead>Public Key</TableHead>
                  <TableHead>Fingerprint</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-[160px] text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((k) => {
                  const keyType = getKeyType(k.public_key!)
                  const truncated = truncateMiddle(k.public_key!, 18)
                  return (
                    <TableRow key={k.id}>
                      <TableCell className="font-medium">{k.name || "—"}</TableCell>
                      <TableCell>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Badge variant="secondary" className="whitespace-nowrap">
                              {keyType}
                            </Badge>
                          </TooltipTrigger>
                          <TooltipContent className="max-w-[70vw]">
                            <div className="max-w-full">
                              <p className="font-mono text-xs break-all">{k.public_key}</p>
                            </div>
                          </TooltipContent>
                        </Tooltip>
                      </TableCell>
                      <TableCell className="font-mono text-xs">{k.fingerprint}</TableCell>
                      <TableCell>
                        {k.created_at
                          ? new Date(k.created_at).toLocaleString(undefined, {
                              year: "numeric",
                              month: "short",
                              day: "2-digit",
                              hour: "2-digit",
                              minute: "2-digit",
                            })
                          : "—"}
                      </TableCell>
                      <TableCell className="space-x-2 text-right">
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => copy(k.public_key ?? "", "Public key copied")}
                        >
                          Copy Pub
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => copy(k.fingerprint ?? "", "Fingerprint copied")}
                        >
                          Copy FP
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => revealMutation.mutate(k.id!)}
                        >
                          <Eye className="mr-1 h-4 w-4" />
                          Reveal
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={async () => {
                            try {
                              const { filename, blob } = await sshApi.downloadBlob(k.id!, "both")
                              saveBlob(blob, filename)
                            } catch (e: any) {
                              toast.error(e?.message ?? "Download failed")
                            }
                          }}
                        >
                          <Download className="mr-1 h-4 w-4" />
                          Download
                        </Button>

                        <Button size="sm" variant="destructive" onClick={() => setDeleteId(k.id!)}>
                          <Trash2 className="mr-1 h-4 w-4" />
                          Delete
                        </Button>
                      </TableCell>
                    </TableRow>
                  )
                })}
                {filtered.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-muted-foreground py-10 text-center">
                      No SSH Keys
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>
        </div>

        {/* Reveal modal */}
        <Dialog open={!!revealFor} onOpenChange={(o) => !o && setRevealFor(null)}>
          <DialogContent className="sm:max-w-2xl">
            <DialogHeader>
              <DialogTitle>Private Key (read-only)</DialogTitle>
            </DialogHeader>
            <div className="space-y-3">
              <div className="text-sm">
                <div className="font-medium">{revealFor?.name ?? "SSH key"}</div>
                <div className="text-muted-foreground font-mono text-xs">
                  {revealFor?.fingerprint}
                </div>
                <Textarea
                  readOnly
                  className="h-64 w-full rounded-md border p-3 font-mono text-xs"
                  value={revealFor?.private_key ?? ""}
                />
                <div className="flex justify-end">
                  <Button
                    onClick={() =>
                      revealFor?.private_key && copy(revealFor.private_key, "Private key copied")
                    }
                  >
                    Copy
                  </Button>
                </div>
              </div>
            </div>
          </DialogContent>
        </Dialog>

        {/* Delete confirm */}
        <Dialog open={!!deleteId} onOpenChange={(o) => !o && setDeleteId(null)}>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Delete SSH Key</DialogTitle>
            </DialogHeader>
            <p className="text-muted-foreground text-sm">
              This will permanently delete the keypair. This action cannot be undone.
            </p>
            <DialogFooter className="gap-2">
              <Button variant="outline" onClick={() => setDeleteId(null)}>
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={() => deleteId && deleteMutation.mutate(deleteId)}
                disabled={deleteMutation.isPending}
              >
                {deleteMutation.isPending ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Deleting…
                  </>
                ) : (
                  "Delete"
                )}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </TooltipProvider>
  )
}
