import { AppShell } from "@/layouts/app-shell.tsx"
import { Route, Routes } from "react-router-dom"

import { ProtectedRoute } from "@/components/protected-route.tsx"
import { AnnotationPage } from "@/pages/annotations/annotation-page.tsx"
import { Login } from "@/pages/auth/login.tsx"
import { JobsPage } from "@/pages/jobs/jobs-page.tsx"
import { LabelsPage } from "@/pages/labels/labels-page.tsx"
import { MePage } from "@/pages/me/me-page.tsx"
import { NodePoolsPage } from "@/pages/nodepools/node-pools-page.tsx"
import { OrgApiKeys } from "@/pages/org/api-keys.tsx"
import { OrgMembers } from "@/pages/org/members.tsx"
import { OrgSettings } from "@/pages/org/settings.tsx"
import { ServerPage } from "@/pages/servers/server-page.tsx"
import { SshPage } from "@/pages/ssh/ssh-page.tsx"
import { TaintsPage } from "@/pages/taints/taints-page.tsx"

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route element={<ProtectedRoute />}>
        <Route element={<AppShell />}>
          <Route path="/me" element={<MePage />} />

          <Route path="/org/settings" element={<OrgSettings />} />
          <Route path="/org/members" element={<OrgMembers />} />
          <Route path="/org/api-keys" element={<OrgApiKeys />} />

          <Route path="/ssh" element={<SshPage />} />
          <Route path="/servers" element={<ServerPage />} />
          <Route path="/taints" element={<TaintsPage />} />
          <Route path="/labels" element={<LabelsPage />} />
          <Route path="/annotations" element={<AnnotationPage />} />
          <Route path="/node-pools" element={<NodePoolsPage />} />

          <Route path="/admin/jobs" element={<JobsPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Login />} />
    </Routes>
  )
}
