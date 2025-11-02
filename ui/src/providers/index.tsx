import type { ReactNode } from "react"
import { ThemeProvider } from "@/providers/theme-provider.tsx"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"

import { Toaster } from "@/components/ui/sonner.tsx"

const queryClient = new QueryClient()

export const Providers = ({ children }: { children: ReactNode }) => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider defaultTheme="system" storageKey="dragon-theme">
        {children}
        <Toaster richColors expand position="top-center" />
      </ThemeProvider>
    </QueryClientProvider>
  )
}
