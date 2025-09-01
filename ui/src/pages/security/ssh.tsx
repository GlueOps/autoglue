import { useEffect, useMemo, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { CloudDownload, Copy, Plus, Trash } from "lucide-react"
import { useForm } from "react-hook-form"
import { z } from "zod"

import { api, API_BASE_URL } from "@/lib/api.ts"
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
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx"
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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx"

type SshKey = {
  id: string
  name: string
  public_keys: string
  fingerprint: string
  created_at: string
}

type Part = "public" | "private" | "both"

const CreateKeySchema = z.object({
  name: z.string().min(1, "Name is required").max(100, "Max 100 characters"),
  comment: z.string().trim().max(100, "Max 100 characters").default(""),
  bits: z.enum(["2048", "3072", "4096"]),
})

type CreateKeyInput = z.input<typeof CreateKeySchema>
type CreateKeyOutput = z.output<typeof CreateKeySchema>

function filenameFromDisposition(disposition?: string, fallback = "download.bin") {
  if (!disposition) return fallback
  const star = /filename\*=UTF-8''([^;]+)/i.exec(disposition)
  if (star?.[1]) return decodeURIComponent(star[1])
  const basic = /filename="?([^"]+)"?/i.exec(disposition)
  return basic?.[1] ?? fallback
}

function truncateMiddle(str: string, keep = 24) {
  if (!str || str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}

function getKeyType(publicKey: string) {
  return publicKey?.split(/\s+/)?.[0] ?? "ssh-key"
}

async function copyToClipboard(text: string) {
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    const el = document.createElement("textarea")
    el.value = text
    el.setAttribute("readonly", "")
    el.style.position = "absolute"
    el.style.left = "-9999px"
    document.body.appendChild(el)
    el.select()
    document.execCommand("copy")
    document.body.removeChild(el)
  }
}

export const SshKeysPage = () => {
  const [sshKeys, setSSHKeys] = useState<SshKey[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState("")
  const [createOpen, setCreateOpen] = useState(false)

  const hasOrg = useMemo(() => !!localStorage.getItem("active_org_id"), [])

  async function fetchSshKeys() {
    setLoading(true)
    setError(null)
    try {
      if (!hasOrg) {
        setSSHKeys([])
        setError("Select an organization first.")
        return
      }
      // api wrapper returns the parsed body directly
      const data = await api.get<SshKey[]>("/api/v1/ssh")
      setSSHKeys(data ?? [])
    } catch (err) {
      console.error(err)
      setError("Failed to fetch SSH keys")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchSshKeys()
    // re-fetch if active org changes in another tab
    const onStorage = (e: StorageEvent) => {
      if (e.key === "active_org_id") fetchSshKeys()
    }
    window.addEventListener("storage", onStorage)
    return () => window.removeEventListener("storage", onStorage)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const filtered = sshKeys.filter((k) => {
    const hay = `${k.name} ${k.public_keys} ${k.fingerprint}`.toLowerCase()
    return hay.includes(filter.toLowerCase())
  })

  // Use raw fetch for download so we can read headers and blob
  async function downloadKeyPair(id: string, part: Part = "both") {
    const token = localStorage.getItem("access_token")
    const orgId = localStorage.getItem("active_org_id")
    const url = `${API_BASE_URL}/api/v1/ssh/${encodeURIComponent(id)}/download?part=${encodeURIComponent(part)}`

    try {
      const res = await fetch(url, {
        method: "GET",
        headers: {
          ...(token ? { Authorization: `Bearer ${token}` } : {}),
          ...(orgId ? { "X-Org-ID": orgId } : {}),
        },
      })

      if (!res.ok) {
        const msg = await res.text().catch(() => "")
        throw new Error(msg || `HTTP ${res.status}`)
      }

      const blob = await res.blob()
      const fallback =
        part === "both"
          ? `ssh_key_${id}.zip`
          : part === "public"
            ? `id_rsa_${id}.pub`
            : `id_rsa_${id}.pem`
      const filename = filenameFromDisposition(
        res.headers.get("content-disposition") ?? undefined,
        fallback
      )

      const objectUrl = URL.createObjectURL(blob)
      const a = document.createElement("a")
      a.href = objectUrl
      a.download = filename
      document.body.appendChild(a)
      a.click()
      a.remove()
      URL.revokeObjectURL(objectUrl)
    } catch (e) {
      console.error(e)
      alert(e instanceof Error ? e.message : "Download failed")
    }
  }

  async function deleteKeyPair(id: string) {
    try {
      await api.delete<void>(`/api/v1/ssh/${encodeURIComponent(id)}`)
      await fetchSshKeys()
    } catch (e) {
      console.error(e)
      alert("Failed to delete key")
    }
  }

  const form = useForm<CreateKeyInput, any, CreateKeyOutput>({
    resolver: zodResolver(CreateKeySchema),
    defaultValues: { name: "", comment: "deploy@autoglue", bits: "4096" },
  })

  async function onSubmit(values: CreateKeyInput) {
    try {
      await api.post<SshKey>("/api/v1/ssh", {
        bits: Number(values.bits),
        comment: values.comment?.trim() ?? "",
        name: values.name.trim(),
        download: "none",
      })
      setCreateOpen(false)
      form.reset()
      await fetchSshKeys()
    } catch (e) {
      console.error(e)
      alert("Failed to create key")
    }
  }

  if (loading) return <div className="p-6">Loading SSH Keys…</div>
  if (error) return <div className="p-6 text-red-500">{error}</div>

  return (
    <TooltipProvider>
      <div className="space-y-4 p-6">
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
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                    name="bits"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Key size</FormLabel>
                        <FormControl>
                          <select
                            className="bg-background w-full rounded-md border px-3 py-2 text-sm"
                            value={field.value}
                            onChange={field.onChange}
                          >
                            <option value="2048">2048</option>
                            <option value="3072">3072</option>
                            <option value="4096">4096</option>
                          </select>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <DialogFooter className="gap-2">
                    <Button type="button" variant="outline" onClick={() => setCreateOpen(false)}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={form.formState.isSubmitting}>
                      {form.formState.isSubmitting ? "Creating…" : "Create"}
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
                  <TableHead className="min-w-[360px]">Public Key</TableHead>
                  <TableHead>Fingerprint</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="w-[160px] text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filtered.map((sshKey) => {
                  const keyType = getKeyType(sshKey.public_keys)
                  const truncated = truncateMiddle(sshKey.public_keys, 18)
                  return (
                    <TableRow key={sshKey.id}>
                      <TableCell className="align-top">{sshKey.name}</TableCell>

                      <TableCell className="align-top">
                        <div className="flex items-start gap-2">
                          <Badge variant="secondary" className="whitespace-nowrap">
                            {keyType}
                          </Badge>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <code className="font-mono text-sm break-all md:max-w-[48ch] md:truncate md:break-normal">
                                {truncated}
                              </code>
                            </TooltipTrigger>
                            <TooltipContent className="max-w-[70vw]">
                              <div className="max-w-full">
                                <p className="font-mono text-xs break-all">{sshKey.public_keys}</p>
                              </div>
                            </TooltipContent>
                          </Tooltip>
                        </div>
                      </TableCell>

                      <TableCell className="align-top">
                        <code className="font-mono text-sm">{sshKey.fingerprint}</code>
                      </TableCell>

                      <TableCell className="align-top">
                        {new Date(sshKey.created_at).toLocaleString(undefined, {
                          year: "numeric",
                          month: "short",
                          day: "2-digit",
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                      </TableCell>

                      <TableCell className="align-top">
                        <div className="flex justify-end gap-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => copyToClipboard(sshKey.public_keys)}
                            title="Copy public key"
                          >
                            <Copy className="mr-2 h-4 w-4" />
                            Copy
                          </Button>

                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="outline" size="sm">
                                <CloudDownload className="mr-2 h-4 w-4" />
                                Download
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => downloadKeyPair(sshKey.id, "both")}>
                                Public + Private (.zip)
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => downloadKeyPair(sshKey.id, "public")}
                              >
                                Public only (.pub)
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                onClick={() => downloadKeyPair(sshKey.id, "private")}
                              >
                                Private only (.pem)
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>

                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => deleteKeyPair(sshKey.id)}
                          >
                            <Trash className="mr-2 h-4 w-4" />
                            Delete
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          </div>
        </div>
      </div>
    </TooltipProvider>
  )
}
