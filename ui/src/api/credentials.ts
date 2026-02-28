import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateCredentialRequest, DtoUpdateCredentialRequest } from "@/sdk"
import { makeCredentialsApi } from "@/sdkClient.ts"

const credentials = makeCredentialsApi()

export const credentialsApi = {
  listCredentials: () =>
    withRefresh(async () => {
      return await credentials.listCredentials()
    }),
  createCredential: async (body: DtoCreateCredentialRequest) =>
    withRefresh(async () => {
      return await credentials.createCredential({ createCredentialRequest: body })
    }),
  getCredential: async (id: string) =>
    withRefresh(async () => {
      return await credentials.getCredential({ id })
    }),
  deleteCredential: async (id: string) =>
    withRefresh(async () => {
      await credentials.deleteCredential({ id })
    }),
  updateCredential: async (id: string, body: DtoUpdateCredentialRequest) =>
    withRefresh(async () => {
      return await credentials.updateCredential({ id, updateCredentialRequest: body })
    }),
  revealCredential: async (id: string) =>
    withRefresh(async () => {
      return await credentials.revealCredential({ id })
    }),
}
