import { withRefresh } from "@/api/with-refresh.ts"
import type {
  HandlersCreateUserKeyRequest,
  HandlersMeResponse,
  HandlersUpdateMeRequest,
  HandlersUserAPIKeyOut,
  ModelsUser,
} from "@/sdk"
import { makeMeApi, makeMeKeysApi } from "@/sdkClient.ts"

const me = makeMeApi()
const keys = makeMeKeysApi()

export const meApi = {
  getMe: () =>
    withRefresh(async (): Promise<HandlersMeResponse> => {
      return await me.getMe()
    }),

  updateMe: (body: HandlersUpdateMeRequest) =>
    withRefresh(async (): Promise<ModelsUser> => {
      return await me.updateMe({ handlersUpdateMeRequest: body })
    }),

  listKeys: () =>
    withRefresh(async (): Promise<HandlersUserAPIKeyOut[]> => {
      return await keys.listUserAPIKeys()
    }),

  createKey: (body: HandlersCreateUserKeyRequest) =>
    withRefresh(async (): Promise<HandlersUserAPIKeyOut> => {
      return await keys.createUserAPIKey({ handlersCreateUserKeyRequest: body })
    }),

  deleteKey: (id: string) =>
    withRefresh(async (): Promise<boolean> => {
      await keys.deleteUserAPIKey({ id })
      return true
    }),
}
