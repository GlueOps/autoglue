import { Outlet } from "react-router-dom"

import { SidebarProvider } from "@/components/ui/sidebar.tsx"
import { DashboardSidebar } from "@/components/dashboard-sidebar.tsx"
import { Footer } from "@/components/footer.tsx"

export function DashboardLayout() {
  return (
    <div className="flex h-screen">
      <SidebarProvider>
        <DashboardSidebar />
        <div className="flex flex-1 flex-col">
          <main className="flex-1 overflow-auto p-4">
            <Outlet />
          </main>
          <Footer />
        </div>
      </SidebarProvider>
    </div>
  )
}
