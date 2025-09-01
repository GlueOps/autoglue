import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Link, useLocation, useNavigate } from "react-router-dom"
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
  email: z.email(),
  password: z.string().min(6),
})

type FormValues = z.infer<typeof schema>

export function Login() {
  const navigate = useNavigate()
  const location = useLocation()
  const form = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
      password: "",
    },
  })

  async function onSubmit(values: FormValues) {
    try {
      await authStore.login(values.email, values.password)
      toast.success("Welcome back!")
      const to = location.state?.from?.pathname ?? "/settings/me"
      navigate(to, { replace: true })
    } catch (e: any) {
      toast.error(e.message || "Login failed")
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card>
        <CardHeader>
          <CardTitle>Sign in</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                name="email"
                control={form.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input placeholder="you@example.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                name="password"
                control={form.control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder="••••••••" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full">
                Sign in
              </Button>
            </form>
          </Form>
          <div className="mt-4 flex justify-between text-sm">
            <Link to="/auth/forgot" className="underline">
              Forgot password?
            </Link>
            <Link to="/auth/register" className="underline">
              Create an account
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
