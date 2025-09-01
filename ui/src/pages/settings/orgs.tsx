import { useEffect, useRef, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { TrashIcon } from "lucide-react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
import { z } from "zod"

import { api, ApiError } from "@/lib/api.ts"
import {
  emitOrgsChanged,
  EVT_ACTIVE_ORG_CHANGED,
  EVT_ORGS_CHANGED,
  getActiveOrgId,
  setActiveOrgId as setActiveOrgIdLS,
} from "@/lib/orgs-sync.ts"
import { slugify } from "@/lib/utils.ts"
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
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"
import { Separator } from "@/components/ui/separator.tsx"
import { Skeleton } from "@/components/ui/skeleton.tsx"

type Organization = {
  id: string // confirm with your API; change to number if needed
  name: string
  slug: string
  created_at: string
}

const OrgSchema = z.object({
  name: z.string().min(2).max(100),
  slug: z
    .string()
    .min(2)
    .max(50)
    .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, "Use lowercase letters, numbers, and hyphens."),
})

type OrgFormValues = z.infer<typeof OrgSchema>

export const OrgManagement = () => {
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [loading, setLoading] = useState<boolean>(true)
  const [createOpen, setCreateOpen] = useState<boolean>(false)
  const [activeOrgId, setActiveOrgIdState] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<string | null>(null)
  const slugEditedRef = useRef(false)

  const form = useForm<OrgFormValues>({
    resolver: zodResolver(OrgSchema),
    mode: "onChange",
    defaultValues: {
      name: "",
      slug: "",
    },
  })

  // auto-generate slug from name unless user manually edited the slug
  const nameValue = form.watch("name")
  useEffect(() => {
    if (!slugEditedRef.current) {
      form.setValue("slug", slugify(nameValue || ""), { shouldValidate: true })
    }
  }, [nameValue, form])

  // fetch organizations
  const getOrgs = async () => {
    setLoading(true)
    try {
      const data = await api.get<Organization[]>("/api/v1/orgs")
      setOrganizations(data)
      setCreateOpen(data.length === 0)
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to load organizations"
      toast.error(msg)
    } finally {
      setLoading(false)
    }
  }

  // initial load + sync listeners
  useEffect(() => {
    // initialize active org from storage
    setActiveOrgIdState(getActiveOrgId())
    void getOrgs()

    // cross-tab sync for active org
    const onStorage = (e: StorageEvent) => {
      if (e.key === "active_org_id") setActiveOrgIdState(e.newValue)
    }
    window.addEventListener("storage", onStorage)

    // same-tab sync for active org (custom event)
    const onActive = (e: Event) => {
      const id = (e as CustomEvent<string | null>).detail ?? null
      setActiveOrgIdState(id)
    }
    window.addEventListener(EVT_ACTIVE_ORG_CHANGED, onActive as EventListener)

    // orgs list changes from elsewhere (custom event)
    const onOrgs = () => void getOrgs()
    window.addEventListener(EVT_ORGS_CHANGED, onOrgs)

    return () => {
      window.removeEventListener("storage", onStorage)
      window.removeEventListener(EVT_ACTIVE_ORG_CHANGED, onActive as EventListener)
      window.removeEventListener(EVT_ORGS_CHANGED, onOrgs)
    }
  }, [])

  async function onSubmit(values: OrgFormValues) {
    try {
      const newOrg = await api.post<Organization>("/api/v1/orgs", values)
      setOrganizations((prev) => [newOrg, ...prev])

      // set as current org and broadcast
      setActiveOrgIdLS(newOrg.id)
      setActiveOrgIdState(newOrg.id)
      emitOrgsChanged()

      toast.success(`Created ${newOrg.name}`)
      setCreateOpen(false)
      form.reset({ name: "", slug: "" })
      slugEditedRef.current = false
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to create organization"
      toast.error(msg)
    }
  }

  function handleSelectOrg(org: Organization) {
    setActiveOrgIdLS(org.id) // updates localStorage + emits event
    setActiveOrgIdState(org.id)
    toast.success(`Switched to ${org.name}`)
  }

  async function handleDeleteOrg(org: Organization) {
    try {
      setDeletingId(org.id)
      await api.delete<void>(`/api/v1/orgs/${org.id}`)

      setOrganizations((prev) => {
        const next = prev.filter((o) => o.id !== org.id)

        // if we deleted the active org, move to the first remaining org (or clear)
        if (activeOrgId === org.id) {
          const nextId = next[0]?.id ?? null
          setActiveOrgIdLS(nextId)
          setActiveOrgIdState(nextId)
        }

        return next
      })

      emitOrgsChanged()
      toast.success(`Deleted ${org.name}`)
    } catch (err) {
      const msg = err instanceof ApiError ? err.message : "Failed to delete organization"
      toast.error(msg)
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return (
      <div className="space-y-4 p-6">
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <h1 className="mb-4 text-2xl font-bold">Organizations</h1>
        </div>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-5 w-40" />
              </CardHeader>
              <CardContent>
                <Skeleton className="mb-2 h-4 w-24" />
                <Skeleton className="h-4 w-48" />
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

  return (
    <div className="space-y-4 p-6">
      <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
        <h1 className="mb-4 text-2xl font-bold">Organizations</h1>
        <Button onClick={() => setCreateOpen(true)}>New organization</Button>
      </div>
      <Separator />

      {organizations.length === 0 ? (
        <div className="text-muted-foreground text-sm">No organizations yet.</div>
      ) : (
        <div className="grid grid-cols-1 gap-4 pr-2 sm:grid-cols-2 lg:grid-cols-3">
          {organizations.map((org) => (
            <Card key={org.id} className="flex flex-col">
              <CardHeader>
                <CardTitle className="text-base">{org.name}</CardTitle>
              </CardHeader>
              <CardContent className="text-muted-foreground text-sm">
                <div>Slug: {org.slug}</div>
                <div className="mt-1">ID: {org.id}</div>
              </CardContent>
              <CardFooter className="mt-auto w-full flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Button onClick={() => handleSelectOrg(org)}>
                  {org.id === activeOrgId ? "Selected" : "Select"}
                </Button>

                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      variant="destructive"
                      className="ml-auto"
                      disabled={deletingId === org.id}
                    >
                      <TrashIcon className="mr-2 h-5 w-5" />
                      {deletingId === org.id ? "Deleting…" : "Delete"}
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Delete organization?</AlertDialogTitle>
                      <AlertDialogDescription>
                        This will permanently delete <b>{org.name}</b>. This action cannot be
                        undone.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter className="sm:justify-between">
                      <AlertDialogCancel disabled={deletingId === org.id}>Cancel</AlertDialogCancel>
                      <AlertDialogAction asChild disabled={deletingId === org.id}>
                        <Button variant="destructive" onClick={() => handleDeleteOrg(org)}>
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

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="sm:max-w-[480px]">
          <DialogHeader>
            <DialogTitle>Create organization</DialogTitle>
            <DialogDescription>Set a name and a URL-friendly slug.</DialogDescription>
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
                      <Input placeholder="Acme Inc" autoFocus {...field} />
                    </FormControl>
                    <FormDescription>This is your organization’s display name.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="slug"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Slug</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="acme-inc"
                        {...field}
                        onChange={(e) => {
                          slugEditedRef.current = true // user manually edited slug
                          field.onChange(e)
                        }}
                        onBlur={(e) => {
                          // normalize on blur
                          const normalized = slugify(e.target.value)
                          form.setValue("slug", normalized, { shouldValidate: true })
                          field.onBlur()
                        }}
                      />
                    </FormControl>
                    <FormDescription>Lowercase, numbers and hyphens only.</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <DialogFooter className="flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    form.reset()
                    setCreateOpen(false)
                    slugEditedRef.current = false
                  }}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={!form.formState.isValid || form.formState.isSubmitting}
                >
                  {form.formState.isSubmitting ? "Creating..." : "Create"}
                </Button>
              </DialogFooter>
            </form>
          </Form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
