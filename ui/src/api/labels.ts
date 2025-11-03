import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateLabelRequest, DtoUpdateLabelRequest } from "@/sdk"
import { makeLabelsApi } from "@/sdkClient.ts"


const labels = makeLabelsApi()

export const labelsApi = {
  listLabels: () =>
    withRefresh(async () => {
      return await labels.listLabels()
    }),
  createLabel: (body: DtoCreateLabelRequest) =>
    withRefresh(async () => {
      return await labels.createLabel({ body })
    }),
  getLabel: (id: string) =>
    withRefresh(async () => {
      return await labels.getLabel({ id })
    }),
  deleteLabel: (id: string) =>
    withRefresh(async () => {
      await labels.deleteLabel({ id })
    }),
  updateLabel: (id: string, body: DtoUpdateLabelRequest) =>
    withRefresh(async () => {
      return await labels.updateLabel({ id, body })
    }),
}
