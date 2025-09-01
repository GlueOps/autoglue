import { useEffect, useState } from "react"
import { useNavigate } from "react-router-dom"
import { toast } from "sonner"

import { authStore } from "@/lib/auth.ts"
import { Button } from "@/components/ui/button.tsx"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card.tsx"

export function Me() {
  const [me, setMe] = useState<any>(null)
  const navigate = useNavigate()

  useEffect(() => {
    ;(async () => {
      try {
        const data = await authStore.me()
        setMe(data)
      } catch (e: any) {
        toast.error(e.message || "Failed to load profile")
      }
    })()
  }, [])

  async function handleLogout() {
    await authStore.logout()
    navigate("/auth/login")
  }

  return (
    <div className="mx-auto max-w-xl">
      <Card>
        <CardHeader>
          <CardTitle>My Account</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {me ? (
            <pre className="bg-muted overflow-auto rounded p-3 text-sm">
              {JSON.stringify(me, null, 2)}
            </pre>
          ) : (
            <p>Loading…</p>
          )}
          <Button onClick={handleLogout}>Sign out</Button>
        </CardContent>
      </Card>
    </div>
  )
}
