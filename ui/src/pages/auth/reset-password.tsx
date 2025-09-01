import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Link, useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"
import { z } from "zod"

import { authStore } from "@/lib/auth.ts"
import { Button } from "@/components/ui/button.tsx"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx"
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx"
import { Input } from "@/components/ui/input.tsx"

const schema = z.object({
  new_password: z.string().min(6),
})

type FormValues = z.infer<typeof schema>

export function ResetPassword() {
  const [params] = useSearchParams()
  const token = params.get("token")
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { new_password: "" },
  })
  const navigate = useNavigate()

  async function onSubmit(values: FormValues) {
    if (!token) {
      toast.error("Missing token")
      return
    }
    try {
      await authStore.reset(token, values.new_password)
      toast.success("Password updated. Please sign in.")
      navigate("/auth/login")
    } catch (e: any) {
      toast.error(e.message || "Reset failed")
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card>
        <CardHeader>
          <CardTitle>Reset password</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                name="new_password"
                control={form.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>New password</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder="••••••••" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full">
                Update password
              </Button>
            </form>
          </Form>
          <div className="mt-4 text-sm">
            <Link to="/auth/login" className="underline">
              Back to sign in
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
