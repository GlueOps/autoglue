import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useEffect, useRef, useState } from "react";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { api, ApiError } from "@/lib/api"; // <-- import ApiError
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { slugify } from "@/lib/utils";
import { toast } from "sonner";
import {
    Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import {
    Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
    AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent,
    AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger
} from "@/components/ui/alert-dialog";
import { TrashIcon } from "lucide-react";

type Organization = {
    id: string;       // confirm with your API; change to number if needed
    name: string;
    slug: string;
    created_at: string;
};

const OrgSchema = z.object({
    name: z.string().min(2).max(100),
    slug: z.string()
        .min(2).max(50)
        .regex(/^[a-z0-9]+(?:-[a-z0-9]+)*$/, "Use lowercase letters, numbers, and hyphens."),
});
type OrgFormValues = z.infer<typeof OrgSchema>;

export const OrgManagement = () => {
    const [organizations, setOrganizations] = useState<Organization[]>([]);
    const [loading, setLoading] = useState(true);
    const [createOpen, setCreateOpen] = useState(false);
    const slugEditedRef = useRef(false);
    const [activeOrgId, setActiveOrgId] = useState<string | null>(null);
    const [deletingId, setDeletingId] = useState<string | null>(null);

    // initialize active org from localStorage once
    useEffect(() => {
        setActiveOrgId(localStorage.getItem("active_org_id"));
    }, []);

    // keep active org in sync across tabs
    useEffect(() => {
        const onStorage = (e: StorageEvent) => {
            if (e.key === "active_org_id") setActiveOrgId(e.newValue);
        };
        window.addEventListener("storage", onStorage);
        return () => window.removeEventListener("storage", onStorage);
    }, []);

    const form = useForm<OrgFormValues>({
        resolver: zodResolver(OrgSchema),
        mode: "onChange",
        defaultValues: { name: "", slug: "" },
    });

    // auto-generate slug from name unless user edited slug manually
    const nameValue = form.watch("name");
    useEffect(() => {
        if (!slugEditedRef.current) {
            form.setValue("slug", slugify(nameValue || ""), { shouldValidate: true });
        }
    }, [nameValue, form]);

    // fetch orgs once
    const getOrgs = async () => {
        setLoading(true);
        try {
            const data = await api.get<Organization[]>("/api/v1/orgs");
            setOrganizations(data);
            setCreateOpen(data.length === 0);
        } catch (err) {
            const msg = err instanceof ApiError ? err.message : "Failed to load organizations";
            toast.error(msg);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        void getOrgs();
    }, []);

    async function onSubmit(values: OrgFormValues) {
        try {
            const newOrg = await api.post<Organization>("/api/v1/orgs", values);
            setOrganizations(prev => [newOrg, ...prev]);
            localStorage.setItem("active_org_id", newOrg.id);
            setActiveOrgId(newOrg.id);
            toast.success(`Created ${newOrg.name}`);
            setCreateOpen(false);
            form.reset({ name: "", slug: "" });
            slugEditedRef.current = false;
        } catch (err) {
            const msg = err instanceof ApiError ? err.message : "Failed to create organization";
            toast.error(msg);
        }
    }

    function handleSelectOrg(org: Organization) {
        localStorage.setItem("active_org_id", org.id);
        setActiveOrgId(org.id);
        toast.success(`Switched to ${org.name}`);
    }

    async function handleDeleteOrg(org: Organization) {
        try {
            setDeletingId(org.id);
            await api.delete<void>(`/api/v1/orgs/${org.id}`); // <-- correct path
            setOrganizations(prev => {
                const next = prev.filter(o => o.id !== org.id); // <-- fix shadow bug
                if (activeOrgId === org.id) {
                    const nextId = next[0]?.id ?? null;
                    if (nextId) localStorage.setItem("active_org_id", nextId);
                    else localStorage.removeItem("active_org_id");
                    setActiveOrgId(nextId);
                }
                return next;
            });
            toast.success(`Deleted ${org.name}`);
        } catch (err) {
            const msg = err instanceof ApiError ? err.message : "Failed to delete organization";
            toast.error(msg);
        } finally {
            setDeletingId(null);
        }
    }

    if (loading) {
        return (
            <div className="p-6 space-y-4">
                <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                    <h1 className="text-2xl font-bold mb-4">Organizations</h1>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {[...Array(6)].map((_, i) => (
                        <Card key={i}>
                            <CardHeader><Skeleton className="h-5 w-40" /></CardHeader>
                            <CardContent>
                                <Skeleton className="h-4 w-24 mb-2" />
                                <Skeleton className="h-4 w-48" />
                            </CardContent>
                            <CardFooter><Skeleton className="h-9 w-24" /></CardFooter>
                        </Card>
                    ))}
                </div>
            </div>
        );
    }

    return (
        <div className="p-6 space-y-4">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <h1 className="text-2xl font-bold mb-4">Organizations</h1>
                <Button onClick={() => setCreateOpen(true)}>New organization</Button>
            </div>
            <Separator />

            {organizations.length === 0 ? (
                <div className="text-sm text-muted-foreground">No organizations yet.</div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 pr-2">
                    {organizations.map(org => (
                        <Card key={org.id} className="flex flex-col">
                            <CardHeader><CardTitle className="text-base">{org.name}</CardTitle></CardHeader>
                            <CardContent className="text-sm text-muted-foreground">
                                <div>Slug: {org.slug}</div>
                                <div className="mt-1">ID: {org.id}</div>
                                <div className="mt-1">Created: {new Date(org.created_at).toUTCString()}</div>
                            </CardContent>
                            <CardFooter className="mt-auto w-full flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-between">
                                <Button onClick={() => handleSelectOrg(org)}>
                                    {org.id === activeOrgId ? "Selected" : "Select"}
                                </Button>

                                <AlertDialog>
                                    <AlertDialogTrigger asChild>
                                        <Button variant="destructive" className="ml-auto">
                                            <TrashIcon className="h-5 w-5 mr-2" />
                                            Delete
                                        </Button>
                                    </AlertDialogTrigger>
                                    <AlertDialogContent>
                                        <AlertDialogHeader>
                                            <AlertDialogTitle>Delete organization?</AlertDialogTitle>
                                            <AlertDialogDescription>
                                                This will permanently delete <b>{org.name}</b>. This action cannot be undone.
                                            </AlertDialogDescription>
                                        </AlertDialogHeader>
                                        <AlertDialogFooter className="sm:justify-between">
                                            <AlertDialogCancel disabled={deletingId === org.id}>Cancel</AlertDialogCancel>
                                            <AlertDialogAction asChild disabled={deletingId === org.id}>
                                                <Button variant="destructive" onClick={() => handleDeleteOrg(org)}>
                                                    {deletingId === org.id ? "Deleting…" : "Delete"}
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
                                        <FormControl><Input placeholder="Acme Inc" autoFocus {...field} /></FormControl>
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
                                                onChange={(e) => { slugEditedRef.current = true; field.onChange(e); }}
                                                onBlur={(e) => {
                                                    const normalized = slugify(e.target.value);
                                                    form.setValue("slug", normalized, { shouldValidate: true });
                                                    field.onBlur();
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
                                    onClick={() => { form.reset(); setCreateOpen(false); }}
                                >
                                    Cancel
                                </Button>
                                <Button type="submit" disabled={!form.formState.isValid || form.formState.isSubmitting}>
                                    {form.formState.isSubmitting ? "Creating..." : "Create"}
                                </Button>
                            </DialogFooter>
                        </form>
                    </Form>
                </DialogContent>
            </Dialog>
        </div>
    );
};
