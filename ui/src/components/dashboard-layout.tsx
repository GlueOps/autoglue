import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import {Outlet} from "react-router-dom";
import {Footer} from "@/components/footer.tsx";
import {DashboardSidebar} from "@/components/dashboard-sidebar.tsx";

export function DashboardLayout() {
  return (
      <div className="flex h-screen">
          <SidebarProvider>
              <DashboardSidebar />
              <div className="flex flex-col flex-1">
                  <main className="flex-1 p-4 overflow-auto">
                      <Outlet />
                  </main>
                  <Footer />
              </div>
          </SidebarProvider>
      </div>
  )
}
