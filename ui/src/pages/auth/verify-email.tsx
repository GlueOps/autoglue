import { useEffect, useState } from "react"
import { Link, useSearchParams } from "react-router-dom"
import { toast } from "sonner"

import { authStore } from "@/lib/auth.ts"
import { Button } from "@/components/ui/button.tsx"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx"

export function VerifyEmail() {
  const [params] = useSearchParams()
  const token = params.get("token")
  const [status, setStatus] = useState<"idle" | "ok" | "error">("idle")

  useEffect(() => {
    async function run() {
      if (!token) {
        setStatus("error")
        return
      }
      try {
        await authStore.verify(token)
        setStatus("ok")
      } catch (e: any) {
        toast.error(e.message || "Verification failed")
        setStatus("error")
      }
    }
    run()
  }, [token])

  return (
    <div className="mx-auto max-w-md">
      <Card>
        <CardHeader>
          <CardTitle>Email verification</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {status === "idle" && <p>Verifying…</p>}
          {status === "ok" && (
            <div>
              <p>Your email has been verified. You can now sign in.</p>
              <Button asChild className="mt-3">
                <Link to="/auth/login">Go to sign in</Link>
              </Button>
            </div>
          )}
          {status === "error" && (
            <div>
              <p>Verification failed. Please request a new verification email.</p>
              <Button asChild className="mt-3">
                <Link to="/auth/login">Back to sign in</Link>
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
