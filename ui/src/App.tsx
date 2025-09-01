import { Navigate, Route, Routes } from "react-router-dom"

import { DashboardLayout } from "@/components/dashboard-layout.tsx"
import { ProtectedRoute } from "@/components/protected-route.tsx"
import { RequireAdmin } from "@/components/require-admin.tsx"
import { AdminUsersPage } from "@/pages/admin/users.tsx"
import { ForgotPassword } from "@/pages/auth/forgot-password.tsx"
import { Login } from "@/pages/auth/login.tsx"
import { Me } from "@/pages/auth/me.tsx"
import { Register } from "@/pages/auth/register.tsx"
import { ResetPassword } from "@/pages/auth/reset-password.tsx"
import { VerifyEmail } from "@/pages/auth/verify-email.tsx"
import { ServersPage } from "@/pages/core/servers-page.tsx"
import { Forbidden } from "@/pages/error/forbidden.tsx"
import { NotFoundPage } from "@/pages/error/not-found.tsx"
import { SshKeysPage } from "@/pages/security/ssh.tsx"
import { MemberManagement } from "@/pages/settings/members.tsx"
import { OrgManagement } from "@/pages/settings/orgs.tsx"

function App() {
  return (
    <Routes>
      <Route path="/403" element={<Forbidden />} />
      <Route path="/" element={<Navigate to="/auth/login" replace />} />
      {/* Public/auth branch */}
      <Route path="/auth">
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="forgot" element={<ForgotPassword />} />
        <Route path="reset" element={<ResetPassword />} />
        <Route path="verify" element={<VerifyEmail />} />
      </Route>

      <Route element={<ProtectedRoute />}>
        <Route element={<DashboardLayout />}>
          <Route element={<RequireAdmin />}>
            <Route path="/admin">
              <Route path="users" element={<AdminUsersPage />} />
            </Route>
          </Route>

          <Route path="/core">
            <Route path="servers" element={<ServersPage />} />
            {/*
              <Route path="cluster" element={<ClusterListPage />} />
            <Route path="node-pools" element={<NodePoolsPage />} />

            <Route path="taints" element={<TaintsPage />} />
            */}
          </Route>

          <Route path="/security">
            <Route path="ssh" element={<SshKeysPage />} />
          </Route>

          <Route path="/settings">
            <Route path="orgs" element={<OrgManagement />} />
            <Route path="members" element={<MemberManagement />} />
            <Route path="me" element={<Me />} />
          </Route>

          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}

export default App
