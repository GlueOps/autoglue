import { withRefresh } from "@/api/with-refresh.ts";
import type { HandlersCreateUserKeyRequest, HandlersMeResponse, HandlersUpdateMeRequest, HandlersUserAPIKeyOut, ModelsUser } from "@/sdk";
import { makeMeApi, makeMeKeysApi } from "@/sdkClient.ts";





const me = makeMeApi()
const keys = makeMeKeysApi()

export type MeResponse = HandlersMeResponse & ModelsUser

export const meApi = {
  getMe: () =>
    withRefresh(async (): Promise<MeResponse> => {
      return await me.getMe()
    }),

  updateMe: (body: HandlersUpdateMeRequest) =>
    withRefresh(async (): Promise<ModelsUser> => {
      return await me.updateMe({ updateMeRequest: body })
    }),

  listKeys: () =>
    withRefresh(async (): Promise<HandlersUserAPIKeyOut[]> => {
      return await keys.listUserAPIKeys()
    }),

  createKey: (body: HandlersCreateUserKeyRequest) =>
    withRefresh(async (): Promise<HandlersUserAPIKeyOut> => {
      return await keys.createUserAPIKey({ createUserAPIKeyRequest: body })
    }),

  deleteKey: (id: string) =>
    withRefresh(async (): Promise<boolean> => {
      await keys.deleteUserAPIKey({ id })
      return true
    }),
}
