import { orgStore } from "@/auth/org.ts"
import { authStore } from "@/auth/store.ts"
import {
  AnnotationsApi,
  ArcherAdminApi,
  AuthApi,
  Configuration,
  LabelsApi,
  MeApi,
  MeAPIKeysApi,
  MetaApi,
  NodePoolsApi,
  OrgsApi,
  ServersApi,
  SshApi,
  TaintsApi,
} from "@/sdk"

export const API_BASE = "/api/v1"

export function makeConfig() {
  return new Configuration({
    basePath: API_BASE,
    accessToken: async () => authStore.getAccessToken() ?? "",
    middleware: [
      {
        async pre(ctx) {
          const headers = new Headers(ctx.init?.headers ?? {})

          const hasBody =
            ctx.init?.body != null &&
            !(ctx.init.body instanceof FormData) &&
            !(ctx.init.body instanceof Blob)

          if (hasBody && !headers.has("Content-Type")) {
            headers.set("Content-Type", "application/json")
          }

          const token = authStore.getAccessToken()
          if (token) {
            headers.set("Authorization", `Bearer ${token}`)
          }

          const org = orgStore.get()
          if (org) {
            headers.set("X-Org-ID", org)
          }

          return {
            ...ctx,
            init: {
              ...ctx.init,
              headers,
            },
          }
        },
        async post(ctx) {
          return ctx.response
        },
      },
    ],
  })
}

function makeApiClient<T>(Ctor: new (cfg: Configuration) => T): T {
  return new Ctor(makeConfig())
}

export function makeAuthApi() {
  return makeApiClient(AuthApi)
}

export function makeMeApi() {
  return makeApiClient(MeApi)
}

export function makeMeKeysApi() {
  return makeApiClient(MeAPIKeysApi)
}

export function makeOrgsApi() {
  return makeApiClient(OrgsApi)
}

export function makeSshApi() {
  return makeApiClient(SshApi)
}

export function makeServersApi() {
  return makeApiClient(ServersApi)
}

export function makeTaintsApi() {
  return makeApiClient(TaintsApi)
}

export function makeLabelsApi() {
  return makeApiClient(LabelsApi)
}

export function makeAnnotationsApi() {
  return makeApiClient(AnnotationsApi)
}

export function makeArcherAdminApi() {
  return makeApiClient(ArcherAdminApi)
}

export function makeNodePoolApi() {
  return makeApiClient(NodePoolsApi)
}

export function makeMetaApi() {
  return makeApiClient(MetaApi)
}
