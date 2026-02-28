import { authStore } from "@/auth/store.ts"
import type { DtoLogoutRequest } from "@/sdk"
import { makeAuthApi } from "@/sdkClient.ts"

export async function logoutEverywhere(): Promise<void> {
  const tokens = authStore.get()

  if (!tokens?.refresh_token) {
    authStore.logout()
    return
  }

  try {
    const body: DtoLogoutRequest = { refresh_token: tokens.refresh_token } as DtoLogoutRequest
    await makeAuthApi().logout({ logoutRequest: body })
  } catch (err) {
    console.warn("Logout API failed; clearing local state anyway", err)
  } finally {
    authStore.logout()
  }
}
