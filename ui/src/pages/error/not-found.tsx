import { useNavigate } from "react-router-dom"

import { Button } from "@/components/ui/button.tsx"

export const NotFoundPage = () => {
  const navigate = useNavigate()

  return (
    <div className="bg-background text-foreground flex min-h-screen flex-col items-center justify-center">
      <h1 className="mb-4 text-6xl font-bold">404</h1>
      <p className="mb-8 text-2xl">Oops! Page not found</p>
      <Button onClick={() => navigate("/dashboard")}>Go back to Dashboard</Button>
    </div>
  )
}
