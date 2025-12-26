import { AppShell } from "@/layouts/app-shell.tsx"
import { Route, Routes } from "react-router-dom"

import { ProtectedRoute } from "@/components/protected-route.tsx"
import { ActionsPage } from "@/pages/actions-page.tsx"
import { AnnotationPage } from "@/pages/annotation-page.tsx"
import { ClustersPage } from "@/pages/cluster-page"
import { CredentialPage } from "@/pages/credential-page.tsx"
import { DnsPage } from "@/pages/dns-page.tsx"
import { DocsPage } from "@/pages/docs-page.tsx"
import { JobsPage } from "@/pages/jobs-page.tsx"
import { LabelsPage } from "@/pages/labels-page.tsx"
import { LoadBalancersPage } from "@/pages/load-balancers-page"
import { Login } from "@/pages/login.tsx"
import { MePage } from "@/pages/me-page.tsx"
import { NodePoolsPage } from "@/pages/node-pools-page.tsx"
import { OrgApiKeys } from "@/pages/org/api-keys.tsx"
import { OrgMembers } from "@/pages/org/members.tsx"
import { OrgSettings } from "@/pages/org/settings.tsx"
import { ServerPage } from "@/pages/server-page.tsx"
import { SshPage } from "@/pages/ssh-page.tsx"
import { TaintsPage } from "@/pages/taints-page.tsx"

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/docs" element={<DocsPage />} />

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
          <Route path="/credentials" element={<CredentialPage />} />
          <Route path="/dns" element={<DnsPage />} />
          <Route path="/load-balancers" element={<LoadBalancersPage />} />
          <Route path="/clusters" element={<ClustersPage />} />

          <Route path="/admin/jobs" element={<JobsPage />} />
          <Route path="/admin/actions" element={<ActionsPage />} />
        </Route>
      </Route>
      <Route path="*" element={<Login />} />
    </Routes>
  )
}
