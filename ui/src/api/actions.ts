import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateActionRequest, DtoUpdateActionRequest } from "@/sdk"
import { makeActionsApi } from "@/sdkClient.ts"

const actions = makeActionsApi()
export const actionsApi = {
  listActions: () =>
    withRefresh(async () => {
      return await actions.listActions()
    }),
  createAction: (body: DtoCreateActionRequest) =>
    withRefresh(async () => {
      return await actions.createAction({
        createActionRequest: body,
      })
    }),
  updateAction: (id: string, body: DtoUpdateActionRequest) =>
    withRefresh(async () => {
      return await actions.updateAction({
        actionID: id,
        updateActionRequest: body,
      })
    }),
  deleteAction: (id: string) =>
    withRefresh(async () => {
      await actions.deleteAction({
        actionID: id,
      })
    }),
}
