import "./index.css"

import { StrictMode } from "react"
import App from "@/App.tsx"
import { createRoot } from "react-dom/client"
import { BrowserRouter } from "react-router-dom"

import { Toaster } from "@/components/ui/sonner.tsx"
import { ThemeProvider } from "@/components/theme-provider.tsx"

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <ThemeProvider defaultTheme="system" storageKey="dragon-theme">
        <App />
        <Toaster richColors position="top-right" />
      </ThemeProvider>
    </BrowserRouter>
  </StrictMode>
)
