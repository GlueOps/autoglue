import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateServerRequest, DtoUpdateServerRequest } from "@/sdk"
import { makeServersApi } from "@/sdkClient.ts"

const servers = makeServersApi()

export const serversApi = {
  listServers: () =>
    withRefresh(async () => {
      return await servers.listServers()
    }),
  createServer: (body: DtoCreateServerRequest) =>
    withRefresh(async () => {
      return await servers.createServer({ dtoCreateServerRequest: body })
    }),
  getServer: (id: string) =>
    withRefresh(async () => {
      return await servers.getServer({ id })
    }),
  updateServer: (id: string, body: DtoUpdateServerRequest) =>
    withRefresh(async () => {
      return await servers.updateServer({ id, dtoUpdateServerRequest: body })
    }),
  deleteServer: (id: string) =>
    withRefresh(async () => {
      await servers.deleteServer({ id })
    }),
}
