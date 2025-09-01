import { Navigate, Route, Routes } from "react-router-dom"

import { DashboardLayout } from "@/components/dashboard-layout.tsx"
import { ProtectedRoute } from "@/components/protected-route.tsx"
import { ForgotPassword } from "@/pages/auth/forgot-password.tsx"
import { Login } from "@/pages/auth/login.tsx"
import { Me } from "@/pages/auth/me.tsx"
import { Register } from "@/pages/auth/register.tsx"
import { ResetPassword } from "@/pages/auth/reset-password.tsx"
import { VerifyEmail } from "@/pages/auth/verify-email.tsx"
import {NotFoundPage} from "@/pages/error/not-found.tsx";
import {OrgManagement} from "@/pages/settings/orgs.tsx";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/auth/login" replace />} />
      {/* Public/auth branch */}
      <Route path="/auth">
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="forgot" element={<ForgotPassword />} />
        <Route path="reset" element={<ResetPassword />} />
        <Route path="verify" element={<VerifyEmail />} />

        <Route element={<ProtectedRoute />}>
            <Route element={<DashboardLayout />}>
                <Route path="me" element={<Me />} />
            </Route>
        </Route>
      </Route>

      <Route element={<ProtectedRoute />}>
        <Route element={<DashboardLayout />}>
          <Route path="/core">
              {/*
              <Route path="cluster" element={<ClusterListPage />} />
            <Route path="node-pools" element={<NodePoolsPage />} />
            <Route path="servers" element={<ServersPage />} />
            <Route path="taints" element={<TaintsPage />} />
            */}
          </Route>

          <Route path="/security">
              {/*<Route path="ssh" element={<SshKeysPage />} />*/}
          </Route>

          <Route path="/settings">
              <Route path="orgs" element={<OrgManagement />} />
              {/*<Route path="members" element={<MemberManagement />} />*/}
          </Route>

          <Route path="*" element={<NotFoundPage />} />
        </Route>
      </Route>
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}

export default App
