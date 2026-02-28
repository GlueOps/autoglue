import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateTaintRequest, DtoUpdateTaintRequest } from "@/sdk"
import { makeTaintsApi } from "@/sdkClient.ts"

const taints = makeTaintsApi()
export const taintsApi = {
  listTaints: () =>
    withRefresh(async () => {
      return await taints.listTaints()
    }),
  createTaint: (body: DtoCreateTaintRequest) =>
    withRefresh(async () => {
      return await taints.createTaint({ createTaintRequest: body })
    }),
  getTaint: (id: string) =>
    withRefresh(async () => {
      return await taints.getTaint({ id })
    }),
  deleteTaint: (id: string) =>
    withRefresh(async () => {
      await taints.deleteTaint({ id })
    }),
  updateTaint: (id: string, body: DtoUpdateTaintRequest) =>
    withRefresh(async () => {
      return await taints.updateTaint({ id, updateTaintRequest: body })
    }),
}
